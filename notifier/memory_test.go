/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-23 21:35:06
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-23 21:35:06
 * @FilePath: \go-risk\notifier\memory_test.go
 * @Description: 内存通知器单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package notifier

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryPublishSubscribe(t *testing.T) {
	m := NewMemory()
	var got []byte
	err := m.Subscribe(context.Background(), TopicBan, func(p []byte) {
		got = p
	})
	assert.NoError(t, err)

	payload := []byte(`{"key":"ip"}`)
	err = m.Publish(context.Background(), TopicBan, payload)
	assert.NoError(t, err)
	assert.Equal(t, string(payload), string(got))
}

func TestMemoryPublishNoSubscribers(t *testing.T) {
	m := NewMemory()
	err := m.Publish(context.Background(), "risk.none", []byte(`{}`))
	assert.NoError(t, err)
}

func TestMemoryPublishMultipleSubscribers(t *testing.T) {
	m := NewMemory()
	count := 0
	for i := 0; i < 3; i++ {
		err := m.Subscribe(context.Background(), TopicUnban, func([]byte) {
			count++
		})
		assert.NoError(t, err)
	}
	err := m.Publish(context.Background(), TopicUnban, []byte(`{}`))
	assert.NoError(t, err)
	assert.Equal(t, 3, count)
}