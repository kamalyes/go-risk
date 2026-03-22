/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-22 13:38:05
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-22 13:38:05
 * @FilePath: \go-risk\bouncer\extauthz_test.go
 * @Description: ext_authz 授权服务单测（放行、拦截映射、来源地址解析）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bouncer

import (
	"context"
	"testing"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/kamalyes/go-risk/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestExtAuthzAllow(t *testing.T) {
	svc := New(&stubEngine{verdict: core.Decision{Verdict: core.Allow}})

	resp, err := svc.AuthorizationServer().Check(context.Background(), authorizationRequest(map[string]string{"user-agent": "curl"}, "/api"))
	require.NoError(t, err)

	assert.Equal(t, int32(codes.OK), resp.GetStatus().GetCode())
	assert.NotNil(t, resp.GetOkResponse())
	assert.Len(t, resp.GetOkResponse().GetHeaders(), 1)
}

func TestExtAuthzDeny(t *testing.T) {
	svc := New(&stubEngine{verdict: core.Decision{Verdict: core.Ban, BanScope: []string{"ip"}}})

	resp, err := svc.AuthorizationServer().Check(context.Background(), authorizationRequest(nil, "/admin"))
	require.NoError(t, err)

	assert.Equal(t, int32(codes.PermissionDenied), resp.GetStatus().GetCode())
	denied := resp.GetDeniedResponse()
	require.NotNil(t, denied)
	assert.Equal(t, int32(403), int32(denied.GetStatus().GetCode()))
}

func TestExtractAuthzSourceIP(t *testing.T) {
	req := &authv3.CheckRequest{
		Attributes: &authv3.AttributeContext{
			Request: &authv3.AttributeContext_Request{
				Http: &authv3.AttributeContext_HttpRequest{
					Method:  "GET",
					Path:    "/api",
					Headers: map[string]string{"user-agent": "curl"},
				},
			},
			Source: &authv3.AttributeContext_Peer{
				Address: &corev3.Address{
					Address: &corev3.Address_SocketAddress{
						SocketAddress: &corev3.SocketAddress{Address: "1.2.3.10"},
					},
				},
			},
		},
	}

	rc := extractAuthz(req)

	assert.Equal(t, "1.2.3.10", rc.Subject.IP)
	assert.Equal(t, "GET", rc.Method)
	assert.Equal(t, "/api", rc.Path)
	assert.Equal(t, "curl", rc.UserAgent)
	assert.NotEmpty(t, rc.TraceID)
}

// authorizationRequest 构造仅含 HTTP 请求的 ext_authz 授权请求
func authorizationRequest(headers map[string]string, path string) *authv3.CheckRequest {
	return &authv3.CheckRequest{
		Attributes: &authv3.AttributeContext{
			Request: &authv3.AttributeContext_Request{
				Http: &authv3.AttributeContext_HttpRequest{
					Method:  "GET",
					Path:    path,
					Headers: headers,
				},
			},
		},
	}
}
