/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-19 10:08:51
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-19 10:08:51
 * @FilePath: \go-risk\bouncer\constants.go
 * @Description: 决策端点路径、响应头与封禁时长默认常量
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bouncer

import "time"

// 决策端点路径与响应头常量，统一维护避免散落魔法字符串
const (
	// DefaultForwardAuthPath forwardAuth 决策端点默认挂载路径，Traefik forwardAuth 指向该地址
	DefaultForwardAuthPath = "/"
	// RiskIDHeader 决策回写的链路标识响应头，与 nethttp 适配器保持一致
	RiskIDHeader = "X-Risk-Id"
	// RetryAfterHeader 挑战或限速的等待时间响应头，拦截时透传给边缘网关
	RetryAfterHeader = "Retry-After"
)

// defaultBanDuration 封禁默认时长，决策结果不含封禁时长时兜底
const defaultBanDuration = time.Hour
