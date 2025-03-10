// Copyright 2020 gorse Project Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package base

import (
	"encoding/gob"
	"strconv"
)

// Index keeps the mapping between names (string) and indices (integer).
type Index interface {
	Len() int32
	Add(name string)
	ToNumber(name string) int32
	ToName(index int32) string
	GetNames() []string
}

func init() {
	gob.Register(&MapIndex{})
	gob.Register(&DirectIndex{})
}

// MapIndex manages the map between sparse Names and dense indices. A sparse ID is
// a user ID or item ID. The dense index is the internal user index or item index
// optimized for faster parameter access and less memory usage.
type MapIndex struct {
	Numbers map[string]int32 // sparse ID -> dense index
	Names   []string         // dense index -> sparse ID
}

// NotId represents an ID doesn't exist.
const NotId = int32(-1)

// NewMapIndex creates a MapIndex.
func NewMapIndex() *MapIndex {
	set := new(MapIndex)
	set.Numbers = make(map[string]int32)
	set.Names = make([]string, 0)
	return set
}

// Len returns the number of indexed Names.
func (idx *MapIndex) Len() int32 {
	if idx == nil {
		return 0
	}
	return int32(len(idx.Names))
}

// Add adds a new ID to the indexer.
func (idx *MapIndex) Add(name string) {
	if _, exist := idx.Numbers[name]; !exist {
		idx.Numbers[name] = int32(len(idx.Names))
		idx.Names = append(idx.Names, name)
	}
}

// ToNumber converts a sparse ID to a dense index.
func (idx *MapIndex) ToNumber(name string) int32 {
	if denseId, exist := idx.Numbers[name]; exist {
		return denseId
	}
	return NotId
}

// ToName converts a dense index to a sparse ID.
func (idx *MapIndex) ToName(index int32) string {
	return idx.Names[index]
}

// GetNames returns all names in current index.
func (idx *MapIndex) GetNames() []string {
	return idx.Names
}

// DirectIndex means that the name and its index is the same. For example,
// the index of "1" is 1, vice versa.
type DirectIndex struct {
	Limit int32
}

// NewDirectIndex create a direct mapping index.
func NewDirectIndex() *DirectIndex {
	return &DirectIndex{Limit: 0}
}

// Len returns the number of names in current index.
func (idx *DirectIndex) Len() int32 {
	return idx.Limit
}

// Add a name to current index.
func (idx *DirectIndex) Add(s string) {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	if int32(i) >= idx.Limit {
		idx.Limit = int32(i + 1)
	}
}

// ToNumber converts a name to corresponding index.
func (idx *DirectIndex) ToNumber(name string) int32 {
	i, err := strconv.Atoi(name)
	if err != nil {
		panic(err)
	}
	if int32(i) >= idx.Limit {
		return NotId
	}
	return int32(i)
}

// ToName converts a index to corresponding name.
func (idx *DirectIndex) ToName(index int32) string {
	if index >= idx.Limit {
		panic("index out of range")
	}
	return strconv.Itoa(int(index))
}

// GetNames returns all names in current index.
func (idx *DirectIndex) GetNames() []string {
	names := make([]string, idx.Limit)
	for i := range names {
		names[i] = strconv.Itoa(i)
	}
	return names
}
