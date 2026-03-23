/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-23 10:11:22
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-23 10:11:22
 * @FilePath: \go-risk\banlist\constants.go
 * @Description: 封禁名单拦截原因常量
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package banlist

// 封禁名单命中的原因标识，统一维护避免散落魔法字符串
const (
	// ReasonID 命中原因 ID，用于审计与热规则统计聚合
	ReasonID = "banlist"
	// ReasonGroup 命中原因所属分组
	ReasonGroup = "banlist"
)
