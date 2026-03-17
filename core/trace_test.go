/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-17 11:58:53
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-17 12:07:22
 * @FilePath: \go-risk\core\trace_test.go
 * @Description: 全链路追踪标识单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTraceIDFormat(t *testing.T) {
	id := NewTraceID()
	assert.Len(t, id, 32)
	assert.Regexp(t, "^[0-9a-f]{32}$", id)
}

func TestNewTraceIDUnique(t *testing.T) {
	assert.NotEqual(t, NewTraceID(), NewTraceID())
}
