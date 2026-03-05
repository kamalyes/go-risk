/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-05 19:26:11
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-05 19:26:11
 * @FilePath: \go-risk\adapters\gin\gin.go
 * @Description: Gin 框架风控适配器（薄中间件）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package gin

import (
	"bytes"
	"io"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamalyes/go-risk/core"
)

const maxBodyRead = 1 << 20 // 最多读取 1MB 请求体

// Middleware 返回 Gin 中间件，对每个请求执行风控决策；非放行结论直接中断
func Middleware(e core.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		rc := extract(c.Request)
		decision := e.Evaluate(c.Request.Context(), rc)
		if !allow(decision) {
			writeDecision(c, decision)
			return
		}
		c.Next()
	}
}

// extract 从 Gin 请求提取归一化风控上下文。
func extract(r *http.Request) *core.RiskContext {
	headers := make(map[string]string, len(r.Header))
	order := make([]string, 0, len(r.Header))
	for k, v := range r.Header {
		headers[k] = v[0]
		order = append(order, k)
	}
	body, _ := readBody(r)
	return &core.RiskContext{
		Subject: core.Subject{
			IP:         clientIP(r),
			TenantID:   r.Header.Get("X-Tenant-Id"),
			UserID:     r.Header.Get("X-User-Id"),
			PlatformID: r.Header.Get("X-Platform-Id"),
		},
		Method:      r.Method,
		Path:        r.URL.Path,
		UserAgent:   r.UserAgent(),
		Headers:     headers,
		HeaderOrder: order,
		Query:       r.URL.RawQuery,
		Body:        body,
		BodySize:    r.ContentLength,
	}
}

func readBody(r *http.Request) (string, error) {
	b, err := io.ReadAll(io.LimitReader(r.Body, maxBodyRead))
	if err != nil {
		return "", err
	}
	// 恢复请求体，避免下游 handler 读取为空
	r.Body = io.NopCloser(bytes.NewReader(b))
	return string(b), nil
}

func allow(d core.Decision) bool {
	return d.Verdict == core.Allow || d.Verdict == core.Observe
}

func writeDecision(c *gin.Context, d core.Decision) {
	switch d.Verdict {
	case core.Challenge, core.Throttle:
		if d.RetryAfter > 0 {
			c.Header("Retry-After", d.RetryAfter.String())
		}
		c.AbortWithStatus(http.StatusTooManyRequests)
	case core.Ban:
		c.AbortWithStatus(http.StatusForbidden)
	default:
		c.AbortWithStatus(http.StatusForbidden)
	}
}

func clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
