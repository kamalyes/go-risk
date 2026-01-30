/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-30 15:26:02
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-30 15:26:02
 * @FilePath: \go-risk\store\store_test.go
 * @Description: 内存计数后端单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package store

import (
	"context"
	"testing"
	"time"

	"github.com/kamalyes/go-risk/core"
)

func TestMemoryIncrBatch(t *testing.T) {
	m := NewMemory()
	ctx := context.Background()
	if err := m.IncrBatch(ctx, []core.CounterOp{{Key: "ip:1.2.3.5", Delta: 1}}); err != nil {
		t.Fatal(err)
	}
	c, err := m.WindowCount(ctx, "ip:1.2.3.5", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if c != 1 {
		t.Fatalf("count = %d, want 1", c)
	}
}

func TestMemoryWindowReset(t *testing.T) {
	m := NewMemory()
	ctx := context.Background()
	_ = m.IncrBatch(ctx, []core.CounterOp{{Key: "k", Delta: 1}})
	time.Sleep(1100 * time.Millisecond)
	c, _ := m.WindowCount(ctx, "k", time.Second)
	if c != 0 {
		t.Fatalf("want reset to 0, got %d", c)
	}
}