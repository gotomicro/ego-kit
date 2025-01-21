// Copyright 2021 ecodeclub
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

package iox

import (
	"io"
	"sync"
)

// MultipleBytes 是一个实现了 io.Reader 和 io.Writer 接口的结构体
// 它可以安全地在多个 goroutine 之间共享
type MultipleBytes struct {
	data  [][]byte
	idx1  int // 第几个切片
	idx2  int // data[idx1] 中的下标
	mutex sync.RWMutex
}

// NewMultipleBytes 创建一个新的 MultipleBytes 实例
// capacity 参数用于预分配内部缓冲区的容量
func NewMultipleBytes(capacity int) *MultipleBytes {
	return &MultipleBytes{
		data: make([][]byte, 0, 1),
	}
}

// Read 实现 io.Reader 接口
// 从当前位置读取数据到 p 中，如果没有数据可读返回 io.EOF
func (m *MultipleBytes) Read(p []byte) (n int, err error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// 如果没有数据或者已经读完了所有数据
	if len(m.data) == 0 || (m.idx1 >= len(m.data)) {
		return 0, io.EOF
	}

	totalRead := 0
	for m.idx1 < len(m.data) {
		currentSlice := m.data[m.idx1]
		remaining := len(currentSlice) - m.idx2
		if remaining <= 0 {
			m.idx1++
			m.idx2 = 0
			continue
		}

		toRead := len(p) - totalRead
		if toRead <= 0 {
			break
		}

		if remaining > toRead {
			n = copy(p[totalRead:], currentSlice[m.idx2:m.idx2+toRead])
			m.idx2 += n
		} else {
			n = copy(p[totalRead:], currentSlice[m.idx2:])
			m.idx1++
			m.idx2 = 0
		}
		totalRead += n
	}

	return totalRead, nil
}

// Write 实现 io.Writer 接口
// 将 p 中的数据写入到内部缓冲区
func (m *MultipleBytes) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 创建新的切片来存储数据
	newSlice := make([]byte, len(p))
	copy(newSlice, p)
	m.data = append(m.data, newSlice)

	return len(p), nil
}

// Len 返回当前缓冲区中的数据总长度
func (m *MultipleBytes) Len() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	total := 0
	for _, slice := range m.data {
		total += len(slice)
	}
	return total
}

// Cap 返回当前缓冲区的容量
func (m *MultipleBytes) Cap() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	total := 0
	for _, slice := range m.data {
		total += cap(slice)
	}
	return total
}

// Reset 重置读取位置到开始处
func (m *MultipleBytes) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.idx1 = 0
	m.idx2 = 0
}

// Clear 清空缓冲区并重置读取位置
func (m *MultipleBytes) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data = m.data[:0]
	m.idx1 = 0
	m.idx2 = 0
}

// Bytes 返回内部缓冲区的副本
func (m *MultipleBytes) Bytes() []byte {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	total := 0
	for _, slice := range m.data {
		total += len(slice)
	}

	result := make([]byte, total)
	pos := 0
	for _, slice := range m.data {
		pos += copy(result[pos:], slice)
	}

	return result
}
