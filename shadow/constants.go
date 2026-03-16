/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-16 09:16:33
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-16 09:16:33
 * @FilePath: \go-risk\shadow\constants.go
 * @Description: 影子模式默认常量定义
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package shadow

const (
	// defaultScoreDelta 判异阈值，为 0 时仅按 Verdict 判异
	defaultScoreDelta = 0
	// defaultBufferSize 差异样本异步队列默认容量
	defaultBufferSize = 256
)