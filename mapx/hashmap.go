// Copyright 2021 gotomicro
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

package mapx

import "github.com/gotomicro/ekit/syncx"

type node[T HashKey, ValType any] struct {
	key   HashKey
	value ValType
	next  *node[T, ValType]
}

func (m *MyHashMap[T, ValType]) NewNode(key HashKey, val ValType) *node[T, ValType] {
	new_node := m.nodePool.Get()
	new_node.value = val
	new_node.key = key
	return new_node
}

type HashKey interface {
	Code() uint64
	Equals(key any) bool
}

type MyHashMap[T HashKey, ValType any] struct {
	hashmap  map[uint64]*node[T, ValType]
	nodePool *syncx.Pool[*node[T, ValType]]
}

func (m *MyHashMap[T, ValType]) Put(key T, val ValType) error {
	hash := key.Code()
	root, ok := m.hashmap[hash]
	if !ok {
		hash = key.Code()
		new_node := m.NewNode(key, val)
		m.hashmap[hash] = new_node
		return nil
	}
	pre := root
	for root != nil {
		if root.key.Equals(key) {
			root.value = val
			return nil
		}
		pre = root
		root = root.next
	}
	new_node := m.NewNode(key, val)
	pre.next = new_node
	return nil
}

func (m *MyHashMap[T, ValType]) Get(key T) (ValType, bool) {
	hash := key.Code()
	root, ok := m.hashmap[hash]
	var val ValType
	if !ok {
		return val, false
	}
	for root != nil {
		if root.key.Equals(key) {
			return root.value, true
		}
		root = root.next
	}
	return val, false
}

func NewHashMap[T HashKey, ValType any](size int) *MyHashMap[T, ValType] {
	return &MyHashMap[T, ValType]{
		nodePool: syncx.NewPool[*node[T, ValType]](func() *node[T, ValType] {
			return &node[T, ValType]{}
		}),
		hashmap: make(map[uint64]*node[T, ValType], size),
	}
}

type HashMap[T HashKey, ValType any] interface {
	Put(key T, val ValType) error
	Get(key T) (ValType, bool)
}
