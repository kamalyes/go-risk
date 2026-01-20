/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-20 10:08:27
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-20 10:08:27
 * @FilePath: \go-risk\transform\transform.go
 * @Description: 请求解码与规范化转换链
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package transform

// Transform 字符串转换函数，用于命中规则前的解码/规范化，防止编码绕过。
type Transform func(string) string

// Chain 顺序应用多个转换的转换链。
type Chain struct {
	transforms []Transform
}

// NewChain 创建转换链，按传入顺序执行。
func NewChain(transforms ...Transform) *Chain {
	return &Chain{transforms: transforms}
}

// Apply 依次应用所有转换并返回结果。
func (c *Chain) Apply(s string) string {
	for _, t := range c.transforms {
		s = t(s)
	}
	return s
}