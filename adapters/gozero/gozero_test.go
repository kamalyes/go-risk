/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-08 19:31:06
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-08 19:31:06
 * @FilePath: \go-risk\adapters\gozero\gozero_test.go
 * @Description: go-zero 适配器闭环单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package gozero

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kamalyes/go-risk/engine"
	"github.com/stretchr/testify/assert"
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
	h(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}