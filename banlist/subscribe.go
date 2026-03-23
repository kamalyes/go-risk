/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-23 10:35:11
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-23 10:35:11
 * @FilePath: \go-risk\banlist\subscribe.go
 * @Description: 封禁广播消费端，订阅封禁与解禁事件维护本地名单
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package banlist

import (
	"context"
	"encoding/json"

	"github.com/kamalyes/go-risk/core"
	"github.com/kamalyes/go-risk/notifier"
)

// Subscribe 订阅封禁与解禁主题并维护本地名单，封禁广播在此落地为强制拦截
func Subscribe(ctx context.Context, n core.Notifier, list *List) error {
	if n == nil {
		return nil
	}
	if err := n.Subscribe(ctx, notifier.TopicBan, func(payload []byte) {
		var entries []core.BanEntry
		if json.Unmarshal(payload, &entries) == nil {
			list.Apply(entries)
		}
	}); err != nil {
		return err
	}
	return n.Subscribe(ctx, notifier.TopicUnban, func(payload []byte) {
		var entries []core.BanEntry
		if json.Unmarshal(payload, &entries) == nil {
			list.Remove(entries)
		}
	})
}
