/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-16 09:20:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-16 09:20:18
 * @FilePath: \go-risk\shadow\shadow.go
 * @Description: 影子模式：主引擎决策放行、备引擎并行观测、决策差异异步落样本
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package shadow

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/kamalyes/go-risk/core"
)

// Divergence 主备引擎决策差异样本，供离线回放与规则自迭代
type Divergence struct {
	TraceID   string
	Subject   core.Subject
	Method    string
	Path      string
	Primary   core.Decision
	Shadow    core.Decision
	Timestamp time.Time
}

// Engine 影子引擎，实现 core.Engine：主引擎决策决定放行，备引擎并行观测。
// 决策差异通过异步队列投递给 Sink，队列满时丢弃，保证主链路不被阻塞。
type Engine struct {
	primary    core.Engine
	shadow     core.Engine
	sink       func(Divergence)
	scoreDelta int
	bufferSize int

	samples chan Divergence
	wg      sync.WaitGroup
	once    sync.Once
}

// Option 影子引擎配置项
type Option func(*Engine)

// New 构建影子引擎，primary 为放行依据，shadow 为并行观测的候选引擎
func New(primary, shadow core.Engine, opts ...Option) *Engine {
	e := &Engine{
		primary:    primary,
		shadow:     shadow,
		scoreDelta: defaultScoreDelta,
		bufferSize: defaultBufferSize,
	}
	for _, opt := range opts {
		opt(e)
	}
	if e.sink != nil {
		e.samples = make(chan Divergence, e.bufferSize)
		e.wg.Add(1)
		go e.loop()
	}
	return e
}

// WithSink 注入差异样本消费者；缺省不落样且不启动后台协程
func WithSink(fn func(Divergence)) Option {
	return func(e *Engine) {
		e.sink = fn
	}
}

// WithScoreDelta 在主备评分差绝对值超过阈值时也落样，缺省仅按 Verdict 判异
func WithScoreDelta(n int) Option {
	return func(e *Engine) {
		if n > 0 {
			e.scoreDelta = n
		}
	}
}

// WithBufferSize 调整差异样本异步队列容量，缺省 256
func WithBufferSize(n int) Option {
	return func(e *Engine) {
		if n > 0 {
			e.bufferSize = n
		}
	}
}

// Evaluate 同步评估主备引擎，返回主引擎决策；备引擎仅在配置落样时参与差异判定
func (e *Engine) Evaluate(ctx context.Context, rc *core.RiskContext) core.Decision {
	pd := e.primary.Evaluate(ctx, rc)
	if e.sink == nil {
		return pd
	}
	sd := e.shadow.Evaluate(ctx, rc)
	if e.divergent(pd, sd) {
		e.record(pd, sd, rc)
	}
	return pd
}

// MarkResult 回写结果到主备引擎，两份行为模型并行累积
func (e *Engine) MarkResult(ctx context.Context, rc *core.RiskContext, statusCode int, ok bool) {
	e.primary.MarkResult(ctx, rc, statusCode, ok)
	e.shadow.MarkResult(ctx, rc, statusCode, ok)
}

// Close 关闭样本队列并等待消费者退出，随后释放主备引擎资源
func (e *Engine) Close() error {
	e.once.Do(func() {
		if e.samples != nil {
			close(e.samples)
		}
	})
	e.wg.Wait()
	return errors.Join(e.primary.Close(), e.shadow.Close())
}

// divergent 判定主备决策是否构成差异样本：Verdict 不同，或评分差超阈值
func (e *Engine) divergent(p, s core.Decision) bool {
	if p.Verdict != s.Verdict {
		return true
	}
	if e.scoreDelta <= 0 {
		return false
	}
	d := p.Score - s.Score
	if d < 0 {
		d = -d
	}
	return d >= e.scoreDelta
}

// record 非阻塞投递样本，队列满直接丢弃以免拖慢主链路
func (e *Engine) record(pd, sd core.Decision, rc *core.RiskContext) {
	select {
	case e.samples <- Divergence{
		TraceID:   rc.TraceID,
		Subject:   rc.Subject,
		Method:    rc.Method,
		Path:      rc.Path,
		Primary:   pd,
		Shadow:    sd,
		Timestamp: time.Now(),
	}:
	default:
	}
}

// loop 样本消费者，随队列关闭退出
func (e *Engine) loop() {
	defer e.wg.Done()
	for d := range e.samples {
		e.sink(d)
	}
}