/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-01 10:07:22
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-01 10:07:22
 * @FilePath: \go-risk\semantic\engine.go
 * @Description: SQL/XSS/命令注入语义检测引擎
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package semantic

import (
	"html"
	"net/url"
	"strings"
	"unicode"
)

// Detector 语义检测引擎，基于词法结构与 Token 比对识别注入，零规则库降低误报
type Detector struct{}

// New 创建语义检测引擎
func New() Detector { return Detector{} }

// Scan 对输入做 URL/HTML 解码与归一化后，依次做 SQL/XSS/命令注入语义分析
func (Detector) Scan(value string) (bool, string) {
	if value == "" {
		return false, ""
	}
	decoded, err := url.QueryUnescape(value)
	if err != nil {
		decoded = value
	}
	v := strings.ToLower(html.UnescapeString(decoded))
	if r, ok := scanSQL(v); ok {
		return true, r
	}
	if r, ok := scanXSS(v); ok {
		return true, r
	}
	if r, ok := scanCommand(v); ok {
		return true, r
	}
	return false, ""
}

// sqlKeywords SQL 注入高危关键字（全小写）
var sqlKeywords = []string{
	"select", "union", "insert", "update", "delete", "drop", "alter",
	"from", "where", "sleep", "benchmark", "execute",
	"information_schema", "concat", "substr", "load_file", "outfile", "group_concat",
}

// scanSQL 识别 SQL 注释终结符、关键字 Token 与恒真条件
func scanSQL(v string) (string, bool) {
	if strings.Contains(v, "--") || strings.Contains(v, "/*") || strings.Contains(v, "#") {
		return "sql-comment", true
	}
	for _, kw := range sqlKeywords {
		if tokenContains(v, kw) {
			return "sql-keyword", true
		}
	}
	if tautology(v) {
		return "sql-tautology", true
	}
	return "", false
}

// xssTags XSS 危险标签，命中即视为注入
var xssTags = []string{"<script", "</script", "<iframe", "<img", "<svg", "<math", "<object", "<embed"}

// xssHandlers XSS 危险事件属性，后接 = 才命中
var xssHandlers = []string{
	"onerror", "onload", "onclick", "onmouseover", "onfocus", "onblur",
	"onchange", "oninput", "onsubmit", "onkeydown", "onkeyup",
}

// scanXSS 识别协议、危险标签与事件属性注入
func scanXSS(v string) (string, bool) {
	if strings.Contains(v, "javascript:") {
		return "xss-protocol", true
	}
	for _, t := range xssTags {
		if strings.Contains(v, t) {
			return "xss-tag", true
		}
	}
	for _, h := range xssHandlers {
		if strings.Contains(v, h+"=") {
			return "xss-handler", true
		}
	}
	return "", false
}

// cmdSeparators 命令注入的分隔符/替换特性
var cmdSeparators = []string{";", "&&", "||", "|", "`", "$(", "${"}

// cmdExecs 常见命令执行载体（全小写）
var cmdExecs = []string{
	"wget", "curl", "cat", "whoami", "id", "echo", "nc", "bash", "sh",
	"ping", "ls", "rm", "chmod", "python", "perl", "powershell",
}

// scanCommand 先判分隔符再判命令 Token，避免把普通文本含分号误报为注入
func scanCommand(v string) (string, bool) {
	hasSep := false
	for _, sep := range cmdSeparators {
		if strings.Contains(v, sep) {
			hasSep = true
			break
		}
	}
	if !hasSep {
		return "", false
	}
	for _, ex := range cmdExecs {
		if tokenContains(v, ex) {
			return "cmd-injection", true
		}
	}
	return "", false
}

// tokens 按词法边界切分，将 = 作为独立 Token 保留，供恒真条件识别
func tokens(s string) []string {
	var out []string
	var cur strings.Builder
	for _, r := range s {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_':
			cur.WriteRune(r)
		case r == '=':
			if cur.Len() > 0 {
				out = append(out, cur.String())
				cur.Reset()
			}
			out = append(out, "=")
		default:
			if cur.Len() > 0 {
				out = append(out, cur.String())
				cur.Reset()
			}
		}
	}
	if cur.Len() > 0 {
		out = append(out, cur.String())
	}
	return out
}

// tokenContains 判断 kw 是否作为完整 Token 出现
func tokenContains(s, kw string) bool {
	for _, tok := range tokens(s) {
		if tok == kw {
			return true
		}
	}
	return false
}

// tautology 识别恒真条件：数字相等（1=1）或字符串相等（'a'='a'）
func tautology(s string) bool {
	toks := tokens(s)
	for i := 1; i+1 < len(toks); i++ {
		if toks[i] == "=" && toks[i-1] == toks[i+1] && isNumeric(toks[i-1]) {
			return true
		}
	}
	return quoteTautology(s)
}

// quoteTautology 扫描形如 'x'='x' / "x"="x" 的相邻相等字符串字面量
func quoteTautology(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] != '\'' && s[i] != '"' {
			continue
		}
		q := s[i]
		j := i + 1
		for j < len(s) && s[j] != q {
			j++
		}
		if j >= len(s) || s[i+1:j] == "" {
			continue
		}
		literal := s[i+1 : j]
		k := j + 1
		for k < len(s) {
			if s[k] == '=' {
				m := k + 1
				for m < len(s) && s[m] == ' ' {
					m++
				}
				if m < len(s) && s[m] == q {
					n := m + 1
					for n < len(s) && s[n] != q {
						n++
					}
					if n < len(s) && s[m+1:n] == literal {
						return true
					}
				}
			}
			k++
		}
	}
	return false
}

// isNumeric 判断 Token 是否为纯数字字面量
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
