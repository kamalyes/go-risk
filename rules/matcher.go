/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-10 19:32:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-10 19:32:08
 * @FilePath: \go-risk\rules\matcher.go
 * @Description: 规则匹配器
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rules

import (
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/kamalyes/go-risk/core"
	"github.com/kamalyes/go-risk/transform"
)

// transformRegistry 内置转换名称映射，规则通过名称引用，避免 core 依赖 transform 实现
var transformRegistry = map[string]transform.Transform{
	"lower":         transform.Lower(),
	"urldecode":     transform.URLDecode(),
	"normalizepath": transform.NormalizePath(),
}

// regexpCache 正则预编译缓存，避免热路径重复编译
var regexpCache sync.Map

// Match 判断规则是否命中，返回 (是否命中, 命中的目标值)
// 命中前先按规则声明的转换链做解码/规范化，防止编码绕过
func Match(rc *core.RiskContext, rule core.Rule) (bool, string) {
	if !rule.Enabled {
		return false, ""
	}
	chain := buildChain(rule.Transforms)
	for _, target := range rule.Targets {
		v := extractTarget(rc, target)
		if v == "" {
			continue
		}
		if chain != nil {
			v = chain.Apply(v)
		}
		if matchOp(v, rule.Op, rule.Pattern) {
			return true, v
		}
	}
	return false, ""
}

// buildChain 按名称解析转换链；未命中名称则忽略
func buildChain(names []string) *transform.Chain {
	ts := make([]transform.Transform, 0, len(names))
	for _, n := range names {
		if t, ok := transformRegistry[n]; ok {
			ts = append(ts, t)
		}
	}
	if len(ts) == 0 {
		return nil
	}
	return transform.NewChain(ts...)
}

// extractTarget 按规则目标从请求上下文提取待匹配值
func extractTarget(rc *core.RiskContext, target core.Target) string {
	switch target.Collection {
	case "path":
		return rc.Path
	case "query":
		return rc.Query
	case "body":
		return rc.Body
	case "header":
		return rc.Headers[target.Field]
	case "ua":
		return rc.UserAgent
	case "arg":
		return argFromQuery(rc.Query, target.Field)
	}
	return ""
}

func argFromQuery(query, field string) string {
	vals, _ := url.ParseQuery(query)
	return vals.Get(field)
}

// matchOp 执行匹配操作符
func matchOp(v string, op core.Op, pattern string) bool {
	switch op {
	case core.OpContains:
		return strings.Contains(v, pattern)
	case core.OpEquals:
		return v == pattern
	case core.OpStartsWith:
		return strings.HasPrefix(v, pattern)
	case core.OpRegex:
		re, err := compileRegexp(pattern)
		if err != nil {
			return false
		}
		return re.MatchString(v)
	}
	return false
}

func compileRegexp(pattern string) (*regexp.Regexp, error) {
	if v, ok := regexpCache.Load(pattern); ok {
		return v.(*regexp.Regexp), nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	regexpCache.Store(pattern, re)
	return re, nil
}