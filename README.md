# go-risk 风控引擎 🛡️

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue.svg)](https://go.dev/)

**go-risk** 是一个框架无关的精准风控引擎，对齐阿里 Sentinel 的 Slot 责任链与 OWASP Coraza 的分阶段决策模型。核心库零 Web 框架依赖，gin / go-zero / go-rpc-gateway / 标准库各自通过一个薄适配层接入；处置分级（观察→限速→挑战→封禁）、规则热更新、解码规范化转换链、分布式计数与封禁广播后端全程可插拔。

## 🏗️ 决策流水线

请求进入后按 Slot 责任链依次处理，处置类结论短路返回：

```
提取 → 白名单 → 黑名单 → 规则匹配(WAF/蜜罐/情报) → 评分 → 分级处置 → 事后回写
```

- 白名单 / 黑名单命中即短路，99% 正常流量与已知攻击在最早阶段拦截
- 规则与引擎分离，版本化快照 + 灰度下发
- 处置分级：`Allow` / `Observe` 观察 / `Throttle` 限速 / `Challenge` 挑战 / `Ban` 封禁

## 📦 模块分层

| 包 | 职责 |
| --- | --- |
| [core](./core) | 核心契约：请求主体、决策、规则、Slot / Engine / CounterStore / Notifier 接口，零依赖 |
| [engine](./engine) | 引擎编排：链式装配 + SlotChain 流水线 + 生命周期管理 |
| [transform](./transform) | 解码 / 规范化转换链（urlDecode / normalizePath / lower），防编码绕过 |
| [store](./store) | 计数后端（内存分片原子计数，Redis 待接入） |
| [notifier](./notifier) | 封禁广播（内存实现，NATS 待接入） |
| [fingerprint](./fingerprint) | 设备指纹（接口占位，HTTP/TLS 指纹待实现） |
| [rules](./rules) | 内置规则集（WAF/蜜罐/威胁情报，待实现） |
| [adapters/nethttp](./adapters/nethttp) | 标准库 net/http 适配器 |

## ✨ 核心特性

### 🎯 决策流水线

- **Slot 责任链**：可插拔、可裁剪，处置类结论短路返回
- **规则引擎分离**：规则快照版本化，支持热更新与灰度

### 🛡️ 精准防绕过

- **转换链**：规则命中前先做 URL 解码、路径规范化、小写归一，杜绝编码绕过

### 🏢 分布式可插拔

- **计数后端**：内存分片原子计数（零锁）为默认，扩展时注入 Redis 实现
- **封禁广播**：内存为默认，跨实例封禁通过注入 NATS 实现

## 📦 安装

```bash
go get github.com/kamalyes/go-risk
```

**系统要求**：Go 1.25+

## 🚀 快速开始

### 最小接入（标准库）

```go
package main

import (
	"net/http"

	"github.com/kamalyes/go-risk/adapters/nethttp"
	"github.com/kamalyes/go-risk/engine"
)

func main() {
	e := engine.New() // 默认内存后端 + 空 Slot 链，直通放行

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.ListenAndServe(":8080", nethttp.Wrap(e, mux))
}
```

### 注入自定义 Slot

```go
e := engine.New(
	engine.WithSlots(myWAFSlot, myScorerSlot),
)
```

## 🎯 处置分级

| Verdict | 语义 |
| --- | --- |
| Allow | 放行 |
| Observe | 观察（埋点记录，不拦截） |
| Throttle | 限速（动态收紧） |
| Challenge | 挑战（429 / 验证码 / 二次验证） |
| Ban | 封禁 |

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE) 开源。

## 📌 Commit Emoji 图例

| Emoji | 类型 | 说明 |
| --- | --- | --- |
| ✨ | feat | 新增功能 |
| 🐛 | fix | Bug 修复 |
| ♻️ | refactor | 代码重构 |
| ✅ | test | 测试相关 |
| 📝 | docs | 文档更新 |
| 🔧 | chore | 配置 / 基础设施 |