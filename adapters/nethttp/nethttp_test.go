/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-31 09:55:31
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-17 12:08:53
 * @FilePath: \go-risk\adapters\nethttp\nethttp_test.go
 * @Description: net/http 适配器闭环单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package nethttp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kamalyes/go-risk/engine"
	"github.com/stretchr/testify/assert"
)

func TestWrapPassthrough(t *testing.T) {
	e := engine.New()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := Wrap(e, next)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestExtractBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api", strings.NewReader("hello"))
	rc := Extract(req)
	assert.Equal(t, "hello", rc.Body)
}

func TestExtractSubjectConfigured(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api?a=1", nil)
	req.Header.Set("X-Forwarded-For", "10.0.0.1")
	req.Header.Set("X-Tenant-Id", "t-100")
	req.Header.Set("X-User-Id", "u-200")
	req.Header.Set("X-Platform-Id", "p-300")
	rc := ExtractWith(req, SubjectAttributes{
		"tenant":   "X-Tenant-Id",
		"user":     "X-User-Id",
		"platform": "X-Platform-Id",
	})
	assert.Equal(t, "10.0.0.1", rc.Subject.IP)
	assert.Equal(t, "t-100", rc.Subject.Attributes["tenant"])
	assert.Equal(t, "u-200", rc.Subject.Attributes["user"])
	assert.Equal(t, "p-300", rc.Subject.Attributes["platform"])
}

func TestExtractSubjectDefaultEmpty(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	req.Header.Set("X-Tenant-Id", "t-100")
	req.Header.Set("X-User-Id", "u-200")
	req.Header.Set("X-Platform-Id", "p-300")
	rc := Extract(req)
	assert.Nil(t, rc.Subject.Attributes)
}

func TestExtractTraceIDPassthrough(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	req.Header.Set("X-Risk-Id", "risk-100")
	rc := Extract(req)
	assert.Equal(t, "risk-100", rc.TraceID)
}

func TestExtractTraceIDGenerated(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	rc := Extract(req)
	assert.NotEmpty(t, rc.TraceID)
}

func TestWrapSetsRiskIDHeader(t *testing.T) {
	e := engine.New()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := Wrap(e, next)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	assert.NotEmpty(t, rec.Header().Get("X-Risk-Id"))
}
