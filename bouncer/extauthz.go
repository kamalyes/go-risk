/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-21 10:26:33
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-21 10:26:33
 * @FilePath: \go-risk\bouncer\extauthz.go
 * @Description: Envoy ext_authz 授权服务，供 Envoy/Istio 网格在线决策
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bouncer

import (
	"context"
	"net/http"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/kamalyes/go-risk/core"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

// AuthorizationServer 返回 Envoy ext_authz 授权服务，注册进 gRPC 服务后供 Envoy/Istio 调用
func (s *Service) AuthorizationServer() authv3.AuthorizationServer {
	return &extAuthzServer{svc: s}
}

// extAuthzServer 实现 envoy.service.auth.v3.Authorization 的 Check 方法
type extAuthzServer struct {
	svc *Service
}

// Check 对单次授权请求做风控决策并翻译为 ext_authz 响应
func (s *extAuthzServer) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	rc := extractAuthz(req)
	d := s.svc.evaluate(ctx, rc)
	return authzResponse(rc, d), nil
}

// extractAuthz 从 ext_authz 请求属性提取归一化风控上下文
func extractAuthz(req *authv3.CheckRequest) *core.RiskContext {
	httpReq := req.GetAttributes().GetRequest().GetHttp()
	headers := httpReq.GetHeaders()
	rc := &core.RiskContext{
		TraceID:   core.NewTraceID(),
		Method:    httpReq.GetMethod(),
		Path:      httpReq.GetPath(),
		UserAgent: headers["user-agent"],
		Headers:   headers,
		Body:      httpReq.GetBody(),
		BodySize:  int64(len(httpReq.GetBody())),
	}
	if id := headers["x-risk-id"]; id != "" {
		rc.TraceID = id
	}
	rc.Subject.IP = sourceIP(req.GetAttributes())
	return rc
}

// sourceIP 解析 ext_authz 来源地址，优先 socket 地址回退管道地址
func sourceIP(attr *authv3.AttributeContext) string {
	addr := attr.GetSource().GetAddress()
	if sa := addr.GetSocketAddress(); sa != nil {
		return sa.GetAddress()
	}
	if pipe := addr.GetPipe(); pipe != nil {
		return pipe.GetPath()
	}
	return ""
}

// authzResponse 将处置结论翻译为 ext_authz 放行或拒绝响应
func authzResponse(rc *core.RiskContext, d core.Decision) *authv3.CheckResponse {
	if allowed(d) {
		return &authv3.CheckResponse{
			Status: &status.Status{Code: int32(codes.OK)},
			HttpResponse: &authv3.CheckResponse_OkResponse{
				OkResponse: &authv3.OkHttpResponse{Headers: []*corev3.HeaderValueOption{riskHeader(rc)}},
			},
		}
	}
	code := denyStatus(d)
	headers := []*corev3.HeaderValueOption{riskHeader(rc)}
	if d.RetryAfter > 0 {
		headers = append(headers, headerValue(RetryAfterHeader, d.RetryAfter.String()))
	}
	return &authv3.CheckResponse{
		Status: &status.Status{Code: int32(codes.PermissionDenied)},
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: &authv3.DeniedHttpResponse{
				Status:  &typev3.HttpStatus{Code: typev3.StatusCode(code)},
				Headers: headers,
				Body:    http.StatusText(code),
			},
		},
	}
}

// riskHeader 构造链路标识回写头
func riskHeader(rc *core.RiskContext) *corev3.HeaderValueOption {
	return headerValue(RiskIDHeader, rc.TraceID)
}

// headerValue 构造覆盖式响应头
func headerValue(key, value string) *corev3.HeaderValueOption {
	return &corev3.HeaderValueOption{
		Header:       &corev3.HeaderValue{Key: key, Value: value},
		AppendAction: corev3.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD,
	}
}
