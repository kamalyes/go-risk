/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-17 11:52:36
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-17 11:52:36
 * @FilePath: \go-risk\core\trace.go
 * @Description: 全链路追踪标识生成
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package core

import (
	"encoding/binary"
	"encoding/hex"
	"math/rand/v2"
)

// NewTraceID 生成 128 位随机十六进制链路追踪标识，串联请求从入口到决策与日志的全链路
func NewTraceID() string {
	var b [16]byte
	binary.LittleEndian.PutUint64(b[0:8], rand.Uint64())
	binary.LittleEndian.PutUint64(b[8:16], rand.Uint64())
	return hex.EncodeToString(b[:])
}