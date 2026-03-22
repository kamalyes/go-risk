/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-19 11:22:17
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-19 11:22:17
 * @FilePath: \go-risk\bouncer\interfaces.go
 * @Description: 封禁落地接口契约
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bouncer

import (
	"context"

	"github.com/kamalyes/go-risk/core"
)

// BanSink 封禁落地抽象，网络层（eBPF/nftables/NetworkPolicy 或远端风控中心）实现后接入
type BanSink interface {
	// Ban 广播封禁条目，多维度封禁时一次提交整批
	Ban(ctx context.Context, entries []core.BanEntry) error

	// Unban 广播解禁条目，白名单解除封禁
	Unban(ctx context.Context, entries []core.BanEntry) error
}
