/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-25 19:37:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-25 19:37:16
 * @FilePath: \go-risk\store\memory_test.go
 * @Description: 内存计数后端补充单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package store

import (
	"context"
	"testing"
	"time"

	"github.com/kamalyes/go-risk/core"
	"github.com/stretchr/testify/assert"
)

func TestMemoryIncrBatchReuseCounter(t *testing.T) {
	m := NewMemory()
	ctx := context.Background()
	err := m.IncrBatch(ctx, []core.CounterOp{{Key: "ip:1.2.3.6", Delta: 1}})
	assert.NoError(t, err)
	err = m.IncrBatch(ctx, []core.CounterOp{{Key: "ip:1.2.3.6", Delta: 2}})
	assert.NoError(t, err)
	c, err := m.WindowCount(ctx, "ip:1.2.3.6", time.Second)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), c)
}

func TestMemoryWindowCountMissingKey(t *testing.T) {
	m := NewMemory()
	c, err := m.WindowCount(context.Background(), "ip:missing", time.Second)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), c)
}

func TestMemoryLoadBanList(t *testing.T) {
	m := NewMemory()
	m.banMu.Lock()
	m.bans["ip:1.2.3.7"] = core.BanEntry{Key: "ip:1.2.3.7", Scope: "ip", ExpireAt: time.Now().Add(time.Minute)}
	m.banMu.Unlock()

	list, err := m.LoadBanList(context.Background())
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "ip:1.2.3.7", list[0].Key)
}
