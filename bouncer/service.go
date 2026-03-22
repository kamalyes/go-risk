/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-20 09:41:27
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-20 09:41:27
 * @FilePath: \go-risk\bouncer\service.go
 * @Description: 决策服务门面，装配引擎、封禁落地与共享决策流程
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bouncer

import (
	"context"
	"log/slog"
	"time"

	"github.com/kamalyes/go-risk/adapters/nethttp"
	"github.com/kamalyes/go-risk/core"
)

// Service bouncer 决策服务，面向边缘网关把引擎决策翻译为放行或拦截
type Service struct {
	engine      core.Engine
	sink        BanSink
	logger      *slog.Logger
	banDuration time.Duration
	attrs       nethttp.SubjectAttributes
}

// Option 决策服务配置项
type Option func(*Service)

// New 创建决策服务，默认不广播封禁也不注入扩展身份属性
func New(engine core.Engine, opts ...Option) *Service {
	s := &Service{engine: engine, banDuration: defaultBanDuration}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// WithBanSink 注入封禁落地，封禁结论时广播到网络层
func WithBanSink(sink BanSink) Option {
	return func(s *Service) { s.sink = sink }
}

// WithLogger 注入结构化日志，广播失败等旁路错误以单行 KV 输出
func WithLogger(l *slog.Logger) Option {
	return func(s *Service) { s.logger = l }
}

// WithBanDuration 覆盖封禁默认时长
func WithBanDuration(d time.Duration) Option {
	return func(s *Service) {
		if d > 0 {
			s.banDuration = d
		}
	}
}

// WithSubjectAttributes 声明扩展身份属性来源，从请求头提取租户或用户等维度
func WithSubjectAttributes(attrs nethttp.SubjectAttributes) Option {
	return func(s *Service) { s.attrs = attrs }
}

// evaluate 执行决策，封禁结论时广播到网络层落地，广播失败不阻断本次拦截
func (s *Service) evaluate(ctx context.Context, rc *core.RiskContext) core.Decision {
	d := s.engine.Evaluate(ctx, rc)
	if d.Verdict == core.Ban && s.sink != nil {
		if err := s.sink.Ban(ctx, banEntries(rc, d, s.banDuration)); err != nil && s.logger != nil {
			s.logger.Error("广播封禁失败", "trace_id", rc.TraceID, "error", err)
		}
	}
	return d
}
