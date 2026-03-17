/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-06 19:15:37
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-17 12:25:16
 * @FilePath: \go-risk\adapters\gin\gin_test.go
 * @Description: Gin 适配器闭环单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package gin

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kamalyes/go-risk/adapters/nethttp"
	"github.com/kamalyes/go-risk/engine"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewarePassthrough(t *testing.T) {
	gin.SetMode(gin.TestMode)
	e := engine.New()
	r := gin.New()
	r.Use(Middleware(e))
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "ok", rec.Body.String())
}

func TestExtractSubject(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api?a=1", nil)
	req.Header.Set("X-Forwarded-For", "10.0.0.1")
	req.Header.Set("X-Tenant-Id", "t-100")
	req.Header.Set("X-User-Id", "u-200")
	req.Header.Set("X-Platform-Id", "p-300")
	rc := nethttp.ExtractWith(req, nethttp.SubjectAttributes{
		"tenant":   "X-Tenant-Id",
		"user":     "X-User-Id",
		"platform": "X-Platform-Id",
	})
	assert.Equal(t, "10.0.0.1", rc.Subject.IP)
	assert.Equal(t, "t-100", rc.Subject.Attributes["tenant"])
	assert.Equal(t, "u-200", rc.Subject.Attributes["user"])
	assert.Equal(t, "p-300", rc.Subject.Attributes["platform"])
	assert.Equal(t, "GET", rc.Method)
	assert.Equal(t, "/api", rc.Path)
	assert.Equal(t, "a=1", rc.Query)
}

func TestExtractBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api", strings.NewReader("hello"))
	rc := nethttp.ExtractWith(req, nil)
	assert.Equal(t, "hello", rc.Body)
	assert.Equal(t, int64(5), rc.BodySize)
}

func TestMiddlewareSetsTraceHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	e := engine.New()
	r := gin.New()
	r.Use(Middleware(e))
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.NotEmpty(t, rec.Header().Get(nethttp.RiskIDHeader))
}
