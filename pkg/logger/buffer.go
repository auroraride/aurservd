// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-11-29
// Based on shiguangju by liasica, magicrolan@qq.com.

package logger

import (
    "bytes"
    "sync"
)

var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

// NewBuffer 从池中获取新 bytes.Buffer
func NewBuffer() *bytes.Buffer {
    return bufferPool.Get().(*bytes.Buffer)
}

// PutBuffer 将 Buffer放入池中
func PutBuffer(buf *bytes.Buffer) {
    // See https://golang.org/issue/23199
    const maxSize = 1 << 16
    if buf != nil && buf.Cap() < maxSize { // 对于大Buffer直接丢弃
        buf.Reset()
        bufferPool.Put(buf)
    }
}
