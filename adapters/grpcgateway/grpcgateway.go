/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-09 20:19:33
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-17 12:26:50
 * @FilePath: \go-risk\adapters\grpcgateway\grpcgateway.go
 * @Description: grpc-gateway 风控适配器（HTTP 中间件 + gRPC 一元拦截器）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package grpcgateway

import (
	"context"
	"net/http"

	"github.com/kamalyes/go-risk/adapters/nethttp"
	"github.com/kamalyes/go-risk/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// Middleware 返回 grpc-gateway 的 HTTP 层中间件（包装 runtime.ServeMux）。
func Middleware(e core.Engine) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return nethttp.Wrap(e, next)
	}
}

// UnaryServerInterceptor 返回 gRPC 服务端一元拦截器，不注入扩展身份属性
func UnaryServerInterceptor(e core.Engine) grpc.UnaryServerInterceptor {
	return UnaryServerInterceptorWith(e, nil)
}

// UnaryServerInterceptorWith 支持自定义扩展身份属性来源 metadata key 的一元拦截器
func UnaryServerInterceptorWith(e core.Engine, h nethttp.SubjectAttributes) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		rc := extract(ctx, info.FullMethod, h)
		grpc.SetHeader(ctx, metadata.Pairs("x-risk-id", rc.TraceID))
		decision := e.Evaluate(ctx, rc)
		if !allow(decision) {
			return nil, decisionError(decision)
		}
		return handler(ctx, req)
	}
}

// extract 从 gRPC 上下文（metadata + peer）提取归一化风控上下文。
func extract(ctx context.Context, method string, h nethttp.SubjectAttributes) *core.RiskContext {
	rc := &core.RiskContext{
		TraceID: core.NewTraceID(),
		Method:  method,
		Path:    method,
		Headers: map[string]string{},
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return rc
	}
	headers := make(map[string]string, len(md))
	order := make([]string, 0, len(md))
	for k, v := range md {
		headers[k] = v[0]
		order = append(order, k)
	}
	rc.Headers = headers
	rc.HeaderOrder = order
	if tid := first(md, "x-risk-id"); tid != "" {
		rc.TraceID = tid
	}
	if len(h) > 0 {
		rc.Subject.Attributes = make(map[string]string, len(h))
		for attr, key := range h {
			if v := first(md, key); v != "" {
				rc.Subject.Attributes[attr] = v
			}
		}
	}
	rc.UserAgent = first(md, "user-agent")
	if ip := first(md, "x-forwarded-for"); ip != "" {
		rc.Subject.IP = ip
	}
	if rc.Subject.IP == "" {
		if p, ok := peer.FromContext(ctx); ok {
			rc.Subject.IP = p.Addr.String()
		}
	}
	return rc
}

func first(md metadata.MD, key string) string {
	if v := md.Get(key); len(v) > 0 {
		return v[0]
	}
	return ""
}

func allow(d core.Decision) bool {
	return d.Verdict == core.Allow || d.Verdict == core.Observe
}

func decisionError(d core.Decision) error {
	switch d.Verdict {
	case core.Challenge, core.Throttle:
		return status.Error(codes.ResourceExhausted, "too many requests")
	case core.Ban:
		return status.Error(codes.PermissionDenied, "forbidden")
	default:
		return status.Error(codes.PermissionDenied, "rejected")
	}
}
