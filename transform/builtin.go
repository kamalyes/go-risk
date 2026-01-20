/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-20 10:08:27
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-20 15:35:52
 * @FilePath: \go-risk\transform\builtin.go
 * @Description: 内置转换实现
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package transform

import (
	"net/url"
	"path"
	"strings"
)

// Lower 返回转为小写的转换。
func Lower() Transform {
	return strings.ToLower
}

// URLDecode 返回 URL 解码转换，只解码百分号编码，避免字面加号被误转空格。
func URLDecode() Transform {
	return func(s string) string {
		decoded, err := url.PathUnescape(s)
		if err != nil {
			return s
		}
		return decoded
	}
}

// NormalizePath 返回路径规范化转换，折叠冗余分隔符与相对路径。
func NormalizePath() Transform {
	return func(s string) string {
		return path.Clean(s)
	}
}
