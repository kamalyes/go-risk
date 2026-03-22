/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-20 15:19:38
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-20 15:19:38
 * @FilePath: \go-risk\bouncer\forwardauth.go
 * @Description: forwardAuth HTTP 决策端点，供 Traefik forwardAuth 中间件调用
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bouncer

import (
	"net/http"

	"github.com/kamalyes/go-risk/adapters/nethttp"
	"github.com/kamalyes/go-risk/core"
)

// Handler 返回 forwardAuth 决策端点，Traefik 依返回码放行或拦截
func (s *Service) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rc := nethttp.ExtractWith(r, s.attrs)
		w.Header().Set(RiskIDHeader, rc.TraceID)
		d := s.evaluate(r.Context(), rc)
		if !allowed(d) {
			writeDeny(w, d)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

// writeDeny 依处置结论回写拦截状态码，挑战或限速附带 Retry-After
func writeDeny(w http.ResponseWriter, d core.Decision) {
	if d.RetryAfter > 0 {
		w.Header().Set(RetryAfterHeader, d.RetryAfter.String())
	}
	http.Error(w, http.StatusText(denyStatus(d)), denyStatus(d))
}
