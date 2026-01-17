/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-17 11:28:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-17 11:28:16
 * @FilePath: \go-risk\core\slot.go
 * @Description: 风控 Slot 链抽象
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package core

import "context"

// Slot 风控流水线中的一个处理单元，按序串联，可插拔、可裁剪。
// 对齐 Sentinel 的 ProcessorSlot 责任链模式。
type Slot interface {
	// Name 返回 Slot 名称，用于日志与排障。
	Name() string

	// Evaluate 处理请求并返回决策；Verdict=Allow 表示交由下一 Slot 继续，
	// 其余 Verdict 表示中断（短路）并作为最终处置。
	Evaluate(ctx context.Context, rc *RiskContext) Decision
}