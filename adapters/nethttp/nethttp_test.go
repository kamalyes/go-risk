/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-31 09:55:31
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-17 11:18:11
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
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestExtractBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api", strings.NewReader("hello"))
	rc := Extract(req)
	if rc.Body != "hello" {
		t.Fatalf("body = %q, want hello", rc.Body)
	}
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
	if rc.Subject.IP != "10.0.0.1" {
		t.Fatalf("ip = %q, want 10.0.0.1", rc.Subject.IP)
	}
	if rc.Subject.Attributes["tenant"] != "t-100" {
		t.Fatalf("tenant = %q, want t-100", rc.Subject.Attributes["tenant"])
	}
	if rc.Subject.Attributes["user"] != "u-200" {
		t.Fatalf("user = %q, want u-200", rc.Subject.Attributes["user"])
	}
	if rc.Subject.Attributes["platform"] != "p-300" {
		t.Fatalf("platform = %q, want p-300", rc.Subject.Attributes["platform"])
	}
}

func TestExtractSubjectDefaultEmpty(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	req.Header.Set("X-Tenant-Id", "t-100")
	req.Header.Set("X-User-Id", "u-200")
	req.Header.Set("X-Platform-Id", "p-300")
	rc := Extract(req)
	if rc.Subject.Attributes != nil {
		t.Fatalf("attributes = %v, want nil", rc.Subject.Attributes)
	}
}
