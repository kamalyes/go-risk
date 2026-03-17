/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-10 19:37:22
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-17 11:26:52
 * @FilePath: \go-risk\adapters\grpcgateway\grpcgateway_test.go
 * @Description: grpc-gateway 适配器闭环单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package grpcgateway

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kamalyes/go-risk/adapters/nethttp"
	"github.com/kamalyes/go-risk/engine"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestMiddlewarePassthrough(t *testing.T) {
	e := engine.New()
	mw := Middleware(e)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := mw(next)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUnaryInterceptorPassthrough(t *testing.T) {
	e := engine.New()
	interceptor := UnaryServerInterceptor(e)
	ctx := context.Background()
	info := &grpc.UnaryServerInfo{FullMethod: "/demo.Service/Hello"}
	handler := func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	}
	resp, err := interceptor(ctx, "req", info, handler)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp)
}

func TestExtractFromMetadata(t *testing.T) {
	md := metadata.Pairs(
		"tenant-id", "t-100",
		"user-id", "u-200",
		"platform-id", "p-300",
		"user-agent", "unit-test",
		"x-forwarded-for", "10.0.0.5",
	)
	ctx := metadata.NewIncomingContext(context.Background(), md)
	rc := extract(ctx, "/demo.Service/Hello", nethttp.SubjectAttributes{
		"tenant":   "tenant-id",
		"user":     "user-id",
		"platform": "platform-id",
	})
	assert.Equal(t, "t-100", rc.Subject.Attributes["tenant"])
	assert.Equal(t, "u-200", rc.Subject.Attributes["user"])
	assert.Equal(t, "p-300", rc.Subject.Attributes["platform"])
	assert.Equal(t, "unit-test", rc.UserAgent)
	assert.Equal(t, "10.0.0.5", rc.Subject.IP)
	assert.Equal(t, "/demo.Service/Hello", rc.Path)
}
