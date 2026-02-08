/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-08 20:15:33
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-08 20:15:33
 * @FilePath: \go-risk\fingerprint\http.go
 * @Description: HTTP 层被动指纹与机器人概率识别
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/kamalyes/go-risk/core"
)

// botUAPatterns 已知自动化工具 User-Agent 特征，命中即视作机器人信号
var botUAPatterns = []string{
	"curl/",           // curl
	"python-requests", // python requests
	"go-http-client",  // Go 默认 client
	"wget/",           // wget
	"sqlmap",          // SQL 注入工具
	"nmap",            // 端口扫描
	"nikto",           // Web 扫描
	"dirbuster",       // 目录爆破
	"masscan",         // 端口扫描
	"zgrab",           // 指纹采集
	"nuclei",          // 漏洞扫描
	"axios/",          // node axios
}

// HTTPFingerprint HTTP 层被动指纹器：基于 header 顺序、存在性与 UA 生成稳定指纹，
// 并输出机器人概率（0~1）零第三方依赖，任意网关卡均可使用
type HTTPFingerprint struct{}

// NewHTTP 创建 HTTP 指纹器
func NewHTTP() *HTTPFingerprint { return &HTTPFingerprint{} }

// Identify 计算指纹与机器人概率
func (h *HTTPFingerprint) Identify(rc *core.RiskContext) Result {
	return Result{
		Hash:     h.hash(rc),
		BotScore: h.botScore(rc),
	}
}

// hash 基于 header 顺序、UA 与关键 Accept 头生成稳定指纹
func (h *HTTPFingerprint) hash(rc *core.RiskContext) string {
	var sb strings.Builder
	sb.WriteString(strings.Join(rc.HeaderOrder, ","))
	sb.WriteString("|")
	sb.WriteString(rc.UserAgent)
	sb.WriteString("|")
	sb.WriteString(rc.Headers["Accept-Language"])
	sb.WriteString("|")
	sb.WriteString(rc.Headers["Accept-Encoding"])
	sum := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(sum[:])
}

// botScore 多信号累加计算机器人概率
func (h *HTTPFingerprint) botScore(rc *core.RiskContext) float64 {
	score := 0.0
	ua := strings.ToLower(rc.UserAgent)

	if ua == "" {
		score += 0.3
	}
	if rc.Headers["Accept-Language"] == "" {
		score += 0.15
	}
	for _, p := range botUAPatterns {
		if strings.Contains(ua, p) {
			score += 0.35
			break
		}
	}
	if len(rc.HeaderOrder) < 3 {
		score += 0.15
	}
	if score > 1 {
		score = 1
	}
	return score
}