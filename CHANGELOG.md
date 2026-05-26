# CHANGELOG

## v0.1.0 — 2026-05-26

初始版本发布。

### 新增功能

- **Cookie 认证模式**：通过浏览器导出的 Cookie 模拟登录，调用妙响 Studio 全部 API
- **API Key 认证模式**：支持官方开放接口，通过 `--api-key` 参数指定
- **`gen` 命令**：提交 AI 歌曲生成任务
  - 同步模式：`--output` 指定输出路径后阻塞等待任务完成并自动下载音频
  - 异步模式：无 `--output` 时立即返回任务 ID（JSON 格式）
  - 支持 `--lyrics`（歌词）、`--style`（风格）、`--singer`（歌手）、`--model`（模型ID）、`--instrumental`（纯音乐模式）参数
- **`job` 命令**：查询单个任务状态和作品详情
  - `--id` 指定任务 ID
  - `--output` 指定后进入同步等待并下载模式
- **`jobs` 命令**：列出最近 50 个任务，显示状态、标题、VID、时长
- **`singers` 命令**：获取 AI 歌手列表（SingerID、名称、描述、音域）
- **`tags` 命令**：获取风格标签列表（TagID、中英文名称）
- **自动轮询**：任务提交后自动轮询状态，显示队列剩余数量和当前状态
- **自动下载**：任务完成后自动通过 VID 获取播放 URL 并下载音频文件

### 技术实现

- **API 逆向**：通过 BrowserOS + CDP 捕获妙响平台前端 JS 中的全部 API 端点定义
- **TaskType=46**：发现正确的 AI 歌曲创作任务类型（TaskType=6 无法完成）
- **响应结构解析**：正确处理 `baseResp` / `data` 嵌套响应格式
- **VID 播放 URL**：通过 `/studio_api/video/get-vid-play-info` 获取 MP3 下载链接
- **HTTP Cookie Jar**：使用 Go 标准库 `cookiejar` 管理 Cookie，自动附加到所有请求
- **Cobra CLI**：使用 spf13/cobra 构建命令行界面

### Cookie 文件

- 自动搜索路径：可执行文件同目录 → 当前工作目录
- 支持 `--cookies` 参数显式指定路径
- 格式：`name1=value1; name2=value2; ...`

### 已知限制

- Cookie 模式：Cookie 具有时效性，过期后需重新获取
- 每次任务生成两首变体，同步模式默认下载第一首
- 任务队列等待时间取决于平台负载，通常 3-10 分钟
