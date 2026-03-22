/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-19 09:15:32
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-19 09:15:32
 * @FilePath: \go-risk\bouncer\doc.go
 * @Description: bouncer 决策服务，面向边缘网关暴露风控决策端点
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

// Package bouncer 提供面向边缘网关的风控决策服务，对标 CrowdSec bouncer 的治理组件定位。
// 它把 go-risk 引擎的决策结果翻译为放行或拦截，以 HTTP forwardAuth 与 gRPC ext_authz
// 两种形式接入 Traefik、Envoy/Istio 等入口，并对封禁结论广播到网络层落地（eBPF/nftables/NetworkPolicy）。
package bouncer
