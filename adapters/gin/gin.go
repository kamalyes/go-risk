/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-05 19:26:11
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-17 12:23:08
 * @FilePath: \go-risk\adapters\gin\gin.go
 * @Description: Gin 框架风控适配器（薄中间件）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamalyes/go-risk/adapters/nethttp"
	"github.com/kamalyes/go-risk/core"
)

// Middleware 返回 Gin 中间件，对每个请求执行风控决策；非放行结论直接中断。
func Middleware(e core.Engine) gin.HandlerFunc {
	return MiddlewareWith(e, nil)
}

// MiddlewareWith 支持自定义扩展身份属性来源 header 的 Gin 中间件
func MiddlewareWith(e core.Engine, h nethttp.SubjectAttributes) gin.HandlerFunc {
	return func(c *gin.Context) {
		rc := nethttp.ExtractWith(c.Request, h)
		c.Header(nethttp.RiskIDHeader, rc.TraceID)
		decision := e.Evaluate(c.Request.Context(), rc)
		if !allow(decision) {
			writeDecision(c, decision)
			return
		}
		c.Next()
	}
}

func allow(d core.Decision) bool {
	return d.Verdict == core.Allow || d.Verdict == core.Observe
}

func writeDecision(c *gin.Context, d core.Decision) {
	switch d.Verdict {
	case core.Challenge, core.Throttle:
		if d.RetryAfter > 0 {
			c.Header("Retry-After", d.RetryAfter.String())
		}
		c.AbortWithStatus(http.StatusTooManyRequests)
	case core.Ban:
		c.AbortWithStatus(http.StatusForbidden)
	default:
		c.AbortWithStatus(http.StatusForbidden)
	}
}
