/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-08 22:07:51
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-08 22:07:51
 * @FilePath: \go-risk\fingerprint\tls.go
 * @Description: TLS/JA3 指纹透传识别
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package fingerprint

import "github.com/kamalyes/go-risk/core"

// tlsHeaderNames 常用 TLS 指纹透传头，由 TLS 终结层（LB/Envoy/ingress）计算后注入
var tlsHeaderNames = []string{"X-Tls-Fingerprint", "X-Ja3", "X-Ja4"}

// TLSFingerprint 基于 TLS 终结层透传的 JA3/JA4 指纹
// 网关若被前端 LB 终结 TLS，则无法直取 ClientHello，需依赖上游透传
type TLSFingerprint struct{}

// NewTLS 创建 TLS 指纹器
func NewTLS() *TLSFingerprint { return &TLSFingerprint{} }

// Identify 从透传头读取 TLS 指纹；未透传时返回空结果
func (t *TLSFingerprint) Identify(rc *core.RiskContext) Result {
	for _, name := range tlsHeaderNames {
		if v := rc.Headers[name]; v != "" {
			return Result{Hash: v}
		}
	}
	return Result{}
}