/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-28 15:38:05
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-17 11:28:09
 * @FilePath: \go-risk\observe\logger.go
 * @Description: 单行 KV 结构化决策日志
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package observe

import (
	"context"
	"log/slog"
	"sort"
	"strings"

	"github.com/kamalyes/go-risk/core"
)

// Log 输出单行 KV 结构化决策日志，trace 等信息随 ctx 透传
func Log(ctx context.Context, logger *slog.Logger, rc *core.RiskContext, d core.Decision) {
	if logger == nil {
		return
	}
	logger.InfoContext(ctx, "risk",
		"trace_id", rc.TraceID,
		"subject", subject(rc),
		"attributes", attributesKV(rc.Subject.Attributes),
		"verdict", d.Verdict.String(),
		"score", d.Score,
		"reasons", reasons(d.Reasons),
	)
}

// subject 返回风控主体标识，优先指纹其次 IP
func subject(rc *core.RiskContext) string {
	if rc.Subject.Fingerprint != "" {
		return rc.Subject.Fingerprint
	}
	return rc.Subject.IP
}

// attributesKV 将扩展身份属性排序后拼为 k=v 单行，保证日志输出稳定
func attributesKV(attrs map[string]string) string {
	if len(attrs) == 0 {
		return ""
	}
	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+attrs[k])
	}
	return strings.Join(parts, ",")
}

// reasons 将命中原因链拼接为单行描述
func reasons(rs []core.Reason) string {
	if len(rs) == 0 {
		return ""
	}
	ids := make([]string, 0, len(rs))
	for _, r := range rs {
		ids = append(ids, r.ID)
	}
	return strings.Join(ids, ",")
}
