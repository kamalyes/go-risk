/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-27 13:06:28
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-17 11:16:39
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

// SubjectAttributes 扩展身份属性来源映射，属性名 -> 来源 header 名
// 通用框架不预设任何业务维度，接入方按需声明要注入的属性及其来源
type SubjectAttributes map[string]string

// Wrap 包装 handler，为每个请求执行风控决策；非放行结论直接中断。
func Wrap(e core.Engine, next http.Handler) http.Handler {
	return WrapWith(e, next, nil)
}

// WrapWith 支持自定义扩展身份属性来源 header 的包装
func WrapWith(e core.Engine, next http.Handler, h SubjectAttributes) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rc := ExtractWith(r, h)
		decision := e.Evaluate(r.Context(), rc)
		if !allow(decision) {
			writeDecision(w, decision)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Extract 从 net/http 请求提取归一化风控上下文，不注入任何扩展身份属性
func Extract(r *http.Request) *core.RiskContext {
	return ExtractWith(r, nil)
}

// ExtractWith 按映射提取扩展身份属性到风控主体
func ExtractWith(r *http.Request, h SubjectAttributes) *core.RiskContext {
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
			Attributes: attributes(r, h),
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

// attributes 按映射读取非空 header 值，构成扩展身份属性；无映射时返回 nil 避免分配
func attributes(r *http.Request, m SubjectAttributes) map[string]string {
	if len(m) == 0 {
		return nil
	}
	attrs := make(map[string]string, len(m))
	for attr, key := range m {
		if v := r.Header.Get(key); v != "" {
			attrs[attr] = v
		}
	}
	return attrs
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
