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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultipleBytes_ReadWrite(t *testing.T) {
	testCases := []struct {
		name     string
		write    []byte
		readSize int
		wantRead []byte
		wantErr  error
	}{
		{
			name:     "empty read",
			write:    []byte{},
			readSize: 1,
			wantRead: []byte{},
			wantErr:  io.EOF,
		},
		{
			name:     "single byte",
			write:    []byte{1},
			readSize: 1,
			wantRead: []byte{1},
			wantErr:  nil,
		},
		{
			name:     "multiple bytes",
			write:    []byte{1, 2, 3, 4, 5},
			readSize: 3,
			wantRead: []byte{1, 2, 3},
			wantErr:  nil,
		},
		{
			name:     "read more than available",
			write:    []byte{1, 2},
			readSize: 4,
			wantRead: []byte{1, 2},
			wantErr:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mb := NewMultipleBytes(len(tc.write))
			n, err := mb.Write(tc.write)
			assert.Equal(t, len(tc.write), n)
			assert.Nil(t, err)

			read := make([]byte, tc.readSize)
			n, err = mb.Read(read)
			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.wantRead, read[:n])
			}
		})
	}
}

func TestMultipleBytes_Reset(t *testing.T) {
	mb := NewMultipleBytes(4)
	data := []byte{1, 2, 3, 4}

	// 写入数据
	n, err := mb.Write(data)
	assert.Equal(t, len(data), n)
	assert.Nil(t, err)

	// 第一次读取
	read := make([]byte, 2)
	n, err = mb.Read(read)
	assert.Equal(t, 2, n)
	assert.Nil(t, err)
	assert.Equal(t, []byte{1, 2}, read)

	// 重置
	mb.Reset()

	// 重置后再次读取
	read = make([]byte, 4)
	n, err = mb.Read(read)
	assert.Equal(t, 4, n)
	assert.Nil(t, err)
	assert.Equal(t, data, read)
}

func TestMultipleBytes_Clear(t *testing.T) {
	mb := NewMultipleBytes(4)
	data := []byte{1, 2, 3, 4}

	// 写入数据
	n, err := mb.Write(data)
	assert.Equal(t, len(data), n)
	assert.Nil(t, err)

	// 清空
	mb.Clear()

	// 清空后读取
	read := make([]byte, 1)
	n, err = mb.Read(read)
	assert.Equal(t, 0, n)
	assert.Equal(t, io.EOF, err)

	// 验证长度
	assert.Equal(t, 0, mb.Len())
}

func TestMultipleBytes_Bytes(t *testing.T) {
	mb := NewMultipleBytes(4)
	data := []byte{1, 2, 3, 4}

	// 写入数据
	n, err := mb.Write(data)
	assert.Equal(t, len(data), n)
	assert.Nil(t, err)

	// 获取副本
	copy := mb.Bytes()
	assert.Equal(t, data, copy)

	// 修改副本不影响原数据
	copy[0] = 5
	original := mb.Bytes()
	assert.Equal(t, data, original)
}
