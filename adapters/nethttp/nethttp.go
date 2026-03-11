/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-27 13:06:28
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-11 19:53:26
 * @FilePath: \go-risk\adapters\nethttp\nethttp.go
 * @Description: 标准库 net/http 适配器
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package nethttp

import (
	"bytes"
	"io"
	"net"
	"net/http"

	"github.com/kamalyes/go-risk/core"
)

const maxBodyRead = 1 << 20 // 最多读取 1MB 请求体

// Wrap 包装 handler，为每个请求执行风控决策；非放行结论直接中断。
func Wrap(e core.Engine, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rc := Extract(r)
		decision := e.Evaluate(r.Context(), rc)
		if !allow(decision) {
			writeDecision(w, decision)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Extract 从 net/http 请求提取归一化风控上下文
func Extract(r *http.Request) *core.RiskContext {
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

func writeDecision(w http.ResponseWriter, d core.Decision) {
	switch d.Verdict {
	case core.Challenge, core.Throttle:
		if d.RetryAfter > 0 {
			w.Header().Set("Retry-After", d.RetryAfter.String())
		}
		http.Error(w, "too many requests", http.StatusTooManyRequests)
	case core.Ban:
		http.Error(w, "forbidden", http.StatusForbidden)
	default:
		http.Error(w, "rejected", http.StatusForbidden)
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
