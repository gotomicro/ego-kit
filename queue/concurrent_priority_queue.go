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

package queue

import (
	"sync"

	"github.com/gotomicro/ekit"
	"github.com/gotomicro/ekit/internal/queue"
)

type ConcurrentPriorityQueue[T any] struct {
	pg queue.PriorityQueue[T]
	m  sync.RWMutex
}

func (c *ConcurrentPriorityQueue[T]) Len() int {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.pg.Len()
}

func (c *ConcurrentPriorityQueue[T]) Cap() int {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.pg.Cap()
}

func (c *ConcurrentPriorityQueue[T]) Peek() (T, error) {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.pg.Peek()
}

func (c *ConcurrentPriorityQueue[T]) Enqueue(t T) error {
	c.m.Lock()
	defer c.m.Unlock()
	return c.pg.Enqueue(t)
}

func (c *ConcurrentPriorityQueue[T]) Dequeue() (T, error) {
	c.m.Lock()
	defer c.m.Unlock()
	return c.pg.Dequeue()
}

// NewConcurrentPriorityQueue 创建优先队列 capacity <= 0 时，为无界队列
func NewConcurrentPriorityQueue[T any](capacity int, compare ekit.Comparator[T]) *ConcurrentPriorityQueue[T] {
	return &ConcurrentPriorityQueue[T]{
		pg: *queue.NewPriorityQueue[T](capacity, compare),
	}
}
