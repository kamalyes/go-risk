/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-07 20:28:52
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-07 20:28:52
 * @FilePath: \go-risk\adapters\gozero\gozero.go
 * @Description: go-zero REST 框架风控适配器（薄中间件）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package gozero

import (
	"net/http"

	"github.com/kamalyes/go-risk/adapters/nethttp"
	"github.com/kamalyes/go-risk/core"
	"github.com/zeromicro/go-zero/rest"
)

// Middleware 返回 go-zero REST 中间件，复用标准库适配完成风控决策。
func Middleware(e core.Engine) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return nethttp.Wrap(e, next).ServeHTTP
	}
}