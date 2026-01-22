/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-22 10:12:31
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-22 10:12:31
 * @FilePath: \go-risk\notifier\memory.go
 * @Description: 内存通知器发布订阅实现
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package notifier

import "context"

// Publish 向订阅者广播消息。
func (m *Memory) Publish(ctx context.Context, topic string, payload []byte) error {
	m.mu.RLock()
	handlers := m.subs[topic]
	m.mu.RUnlock()
	for _, h := range handlers {
		h(payload)
	}
	return nil
}

// Subscribe 订阅指定主题，handler 在消息到达时回调。
func (m *Memory) Subscribe(ctx context.Context, topic string, handler func([]byte)) error {
	m.mu.Lock()
	m.subs[topic] = append(m.subs[topic], handler)
	m.mu.Unlock()
	return nil
}
