# go-risk

框架无关的精准风控库，对齐阿里 Sentinel / OWASP Coraza 的设计精华。核心库零 Web 框架依赖，各框架只需薄适配层接入。

## 特性

- 决策流水线（Slot 链）可插拔、可裁剪
- 规则与引擎分离，支持热更新与灰度
- 请求解码 / 规范化转换链，防编码绕过
- 分布式计数与封禁广播（后端可插拔）

## 接入

```go
e := engine.New()
handler := nethttp.Wrap(e, mux)
```

## 模块

- `core`：核心契约（主体、决策、规则、接口）
- `engine`：风控引擎编排
- `transform`：解码 / 规范化转换链
- `store`：计数后端（内存实现）
- `notifier`：封禁广播（内存实现）
- `fingerprint`：设备指纹（占位）
- `rules`：内置规则集（占位）
- `adapters/nethttp`：标准库适配器