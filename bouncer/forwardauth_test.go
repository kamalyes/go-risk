/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-22 11:23:19
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-22 11:23:19
 * @FilePath: \go-risk\bouncer\forwardauth_test.go
 * @Description: forwardAuth HTTP 决策端点单测（放行、拦截、限速、封禁广播）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bouncer

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kamalyes/go-risk/core"
	"github.com/stretchr/testify/assert"
)

func TestHandlerAllow(t *testing.T) {
	svc := New(&stubEngine{verdict: core.Decision{Verdict: core.Allow}})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://auth/", nil)
	svc.Handler().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, rec.Header().Get(RiskIDHeader))
}

func TestHandlerBan(t *testing.T) {
	svc := New(&stubEngine{verdict: core.Decision{Verdict: core.Ban, BanScope: []string{"ip"}}})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://auth/", nil)
	req.RemoteAddr = "1.2.3.9:80"
	svc.Handler().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.NotEmpty(t, rec.Header().Get(RiskIDHeader))
}

func TestHandlerThrottleRetryAfter(t *testing.T) {
	svc := New(&stubEngine{verdict: core.Decision{Verdict: core.Throttle, RetryAfter: time.Minute}})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://auth/", nil)
	svc.Handler().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
	assert.Equal(t, time.Minute.String(), rec.Header().Get(RetryAfterHeader))
}

func TestHandlerBanBroadcast(t *testing.T) {
	sink := &stubSink{}
	svc := New(
		&stubEngine{verdict: core.Decision{Verdict: core.Ban, BanScope: []string{"ip", "fingerprint"}}},
		WithBanSink(sink),
	)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://auth/", nil)
	req.RemoteAddr = "1.2.3.9:80"
	svc.Handler().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Len(t, sink.bans, 1)
	assert.Len(t, sink.bans[0], 1)
	assert.Equal(t, "ip", sink.bans[0][0].Scope)
}

func TestHandlerAllowNoBroadcast(t *testing.T) {
	sink := &stubSink{}
	svc := New(&stubEngine{verdict: core.Decision{Verdict: core.Allow}}, WithBanSink(sink))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://auth/", nil)
	svc.Handler().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, sink.bans)
}
