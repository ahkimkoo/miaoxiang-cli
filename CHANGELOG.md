# CHANGELOG

## v0.2.0 — 2026-05-26

### 新增

- **`models` 命令**：查询可用模型列表（ModelType、ModelVersion），不再需要猜 model ID
- **全量下载**：同步模式下载全部生成的歌曲，不再只下第一首
- **歌曲标题命名**：下载文件以歌曲标题命名，同名自动加后缀（`-2`、`-3`）
- **JSON 下载结果**：下载完成后输出 JSON 格式的文件路径、标题、时长
- **跨平台发布**：提供 macOS（amd64/arm64）、Linux（amd64/arm64）、Windows（amd64）预编译二进制
- **OpenClaw Skill**：`releases/miaoxiang-cli/SKILL.md`，支持 OpenClaw 自动发现和加载

### 修复

- **singers 命令返回空**：API 返回字段为 `AISingers`（复数），代码误解析为 `Singers`，导致 JSON 反序列化为空数组
- **纯音乐模式要求歌词**：`--lyrics` 不再强制必填，`--instrumental` 模式无需歌词
- **HTTP 错误静默吞掉**：`doRequest` 新增 HTTP 状态码检查，非 2xx 返回明确错误信息
- **JSON 解析失败无提示**：反序列化失败时输出原始响应体到 stderr，便于调试

### 变更

- `--model` help 更新为 "模型类型"，提示用 `models` 命令查看
- `--singer` help 改为 "AI歌手名称"
- `--style` help 补充 "用 tags 命令查看"
- `gen` 命令 Long description 新增 4 个典型示例
- 进度日志统一输出到 stderr，JSON 结果输出到 stdout
- `AISinger` 结构体补全 `Gender`、`DemoURL`、`CoverURL`、`ArtistID`、`MatchPercent` 字段

## v0.1.0 — 2026-05-26

### 新增

- **Cookie 认证模式**：通过浏览器导出的 Cookie 调用妙响 Studio API
- **API Key 认证模式**：支持 `--api-key` 参数
- **`gen` 命令**：提交 AI 歌曲生成任务，支持同步/异步模式
- **`job` 命令**：查询单个任务状态和作品详情
- **`jobs` 命令**：列出最近 50 个任务
- **`singers` 命令**：获取 AI 歌手列表
- **`tags` 命令**：获取风格标签列表
- **自动轮询**：任务提交后自动轮询状态，显示队列数量
- **自动下载**：任务完成后通过 VID 获取播放 URL 下载音频

### 技术实现

- API 逆向：通过 BrowserOS + CDP 捕获妙响平台前端 API 端点
- TaskType=46：发现正确的 AI 歌曲创作任务类型
- VID 播放 URL：通过 `/studio_api/video/get-vid-play-info` 获取 MP3 下载链接
- Cobra CLI：使用 spf13/cobra 构建命令行界面

### 已知限制

- Cookie 具有时效性，过期后需重新获取
- 同步模式仅下载第一首歌曲变体（v0.2.0 已修复）
