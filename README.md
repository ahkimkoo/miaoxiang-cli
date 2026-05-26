# 妙响 CLI (miaoxiang-cli)

[抖音妙响音乐生成平台](https://music.douyin.com/studio/create) 命令行工具，支持 AI 歌曲生成、任务查询、作品下载，提供 Cookie 和 API Key 两种认证模式。

## 特性

- 两种认证模式：Cookie（模拟浏览器）和 API Key（官方开放接口）
- 同步/异步两种任务提交模式
- 自动轮询任务状态并下载完成的音频文件
- 查询歌手列表和风格标签
- 管理全部历史任务
- 零外部运行时依赖，单二进制文件即可运行

## 安装

### 从源码编译

需要 Go 1.21 或更高版本。

```bash
git clone https://github.com/your-org/miaoxiang-cli.git
cd miaoxiang-cli
go build -o miaoxiang-cli .
```

编译后的 `miaoxiang-cli` 是一个独立的二进制文件，可以直接复制到任意位置使用。

### 安装到系统 PATH

```bash
sudo cp miaoxiang-cli /usr/local/bin/
```

## 认证

CLI 支持两种认证方式，Cookie 模式为默认推荐方式。

### Cookie 模式（推荐）

Cookie 模式模拟浏览器登录行为，可以调用完整的 Studio API。

#### 获取 Cookie

**方法一：浏览器开发者工具**

1. 使用 Chrome/Edge 打开 https://music.douyin.com/studio/create
2. 确保已登录抖音账号
3. 按 F12 打开开发者工具
4. 进入 Application 面板 → Storage → Cookies → https://music.douyin.com
5. 全选所有 cookie，复制为 `name1=value1; name2=value2; ...` 格式

**方法二：使用 BrowserOS 自动化提取**

```python
# 通过 CDP 连接 BrowserOS 实例，提取完整 cookies（含 httpOnly）
```

#### 配置 Cookie

将 cookie 字符串保存到 `cookies.txt` 文件中：

```bash
# 方式1：放在可执行文件同目录
cp cookies.txt /usr/local/bin/cookies.txt

# 方式2：放在当前工作目录
cp cookies.txt ./cookies.txt

# 方式3：通过 --cookies 参数指定路径
miaoxiang-cli --cookies /path/to/cookies.txt gen ...
```

cookie 文件格式为分号分隔的键值对：

```
passport_csrf_token=abc123; passport_csrf_token_default=abc123; sessionid=xyz789; ...
```

### API Key 模式

API Key 模式使用妙响平台官方开放接口，适用于服务器端集成。

```bash
miaoxiang-cli --api-key "YOUR_API_KEY" gen --lyrics "xxx" --style "流行"
```

> API Key 获取：登录妙响平台后在设置页面获取。

## 使用

### 提交歌曲生成任务

**同步模式**（阻塞等待，完成后自动下载）：

```bash
miaoxiang-cli gen \
  --lyrics "春天的花开，夏天的风吹" \
  --style "流行" \
  --output my_song.wav
```

输出：
```
任务已创建: 7644058967982934824
等待任务 7644058967982934824 完成...
  状态: 处理中, 队列: 54
  状态: 处理中, 队列: 49
  ...
  状态: 已完成, 队列: 0
作品: 花迎夏风 (VID: v0d003g10004d8aigl2ljht7vn7nnbv0)
下载音频到: my_song.wav
下载完成: my_song.wav (203200)
```

**异步模式**（返回任务 ID，不等待完成）：

```bash
miaoxiang-cli gen \
  --lyrics "春天的花开，夏天的风吹" \
  --style "流行"
```

输出 JSON：
```json
{
  "Status": "created",
  "TaskID": "7644060013820152576",
  "TaskType": 46
}
```

### 查询单个任务

**查询状态**：

```bash
miaoxiang-cli job --id 7644060013820152576
```

输出 JSON：
```json
{
  "QueueTaskCount": "50",
  "Status": 2,
  "StatusText": "处理中",
  "TaskID": "7644060013820152576",
  "Works": [
    {
      "Duration": "0",
      "Status": 2,
      "Title": "",
      "VID": "",
      "WorkID": "7644060046762216255"
    },
    {
      "Duration": "0",
      "Status": 2,
      "Title": "",
      "VID": "",
      "WorkID": "7644060147132074767"
    }
  ]
}
```

**查询并等待下载**：

```bash
miaoxiang-cli job --id 7644058967982934824 --output my_song.wav
```

### 列出所有任务

```bash
miaoxiang-cli jobs
```

输出：
```
[处理中] 7644060013820152576 |  | VID:  | 0
[已完成] 7644058967982934824 | 花迎夏风 | VID: v0d003g10004d8aigl2ljht7vn7nnbv0 | 203200
[已完成] 7644058967982934824 | 风随花来 | VID: v02003g10004d8aigjaljht81b8bc500 | 238320
```

### 获取 AI 歌手列表

```bash
miaoxiang-cli singers
```

输出 JSON 格式的歌手列表，包含 SingerID、Name、Description、PitchRange 等字段。

歌手 ID 可在 `gen` 命令中通过 `--singer` 参数指定。

### 获取风格标签

```bash
miaoxiang-cli tags
```

输出 JSON 格式的风格标签列表，包含 TagID、TagName、TagNameZh、TagNameEn 等字段。

常见风格：
- 轻音乐/原声 (Easy Listening/Acoustic)
- 电影/史诗 (Cinematic/Epic)
- 世界/异域 (World/Exotic)
- 流行 (Pop)
- 电子 (Electronic)
- R&B

### 获取 AI 配置

通过 Go API 获取支持的模型列表：

```go
client, _ := NewClient("", "cookies.txt")
config, _ := client.GetConfig()
for _, m := range config.Data.SupportModels {
    fmt.Printf("Model %d: %s (v%s)\n", m.ModelID, m.DisplayName, m.ModelVersion)
}
```

## 命令参考

### 全局参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--api-key string` | API Key | 空（使用 Cookie 模式） |
| `--cookies string` | Cookie 文件路径 | 可执行文件同目录或当前目录的 cookies.txt |
| `-h, --help` | 显示帮助信息 | - |

### `gen` 命令 — 提交歌曲生成任务

| 参数 | 必填 | 说明 | 默认值 |
|------|------|------|--------|
| `--lyrics string` | 是 | 歌词内容 | - |
| `--style string` | 是 | 音乐风格 | - |
| `--output string` | 否 | 输出文件路径。设置此参数后进入同步模式，等待任务完成并下载 | 无（异步模式） |
| `--singer string` | 否 | 歌手 ID（通过 singers 命令获取） | 无 |
| `--model int` | 否 | 模型 ID | 7 (Sway v5.5) |
| `--instrumental` | 否 | 纯音乐模式（无人声） | false |
| `--interval int` | 否 | 轮询间隔，单位秒 | 15 |

### `job` 命令 — 查询单个任务

| 参数 | 必填 | 说明 | 默认值 |
|------|------|------|--------|
| `--id string` | 是 | 任务 ID | - |
| `--output string` | 否 | 输出文件路径。设置后等待任务完成并下载 | 无 |
| `--interval int` | 否 | 轮询间隔，单位秒 | 15 |

### `jobs` 命令 — 列出所有任务

无额外参数。列出最近 50 个任务。

### `singers` 命令 — 获取 AI 歌手列表

无额外参数。

### `tags` 命令 — 获取风格标签

无额外参数。

## 任务状态码

| 状态码 | 说明 | 行为 |
|--------|------|------|
| 1 | 排队中 | 继续轮询 |
| 2 | 处理中 | 继续轮询 |
| 3 | 已完成 | 获取作品详情，下载音频 |
| 5 | 失败 | 报错退出 |

## API 端点

Cookie 模式下调用的妙响 Studio API 端点：

| 方法 | 端点 | 说明 |
|------|------|------|
| POST | `/studio_api/create-studio-task` | 创建 AI 歌曲任务（TaskType=46） |
| GET | `/studio_api/get-studio-task` | 查询任务状态 |
| GET | `/studio_api/assets/work` | 获取作品详情 |
| GET | `/studio_api/assets/work-list` | 列出所有作品 |
| POST | `/studio_api/video/get-vid-play-info` | 获取视频播放 URL（含 MP3） |
| POST | `/studio_api/aigc/ai-singers` | 获取 AI 歌手列表 |
| GET | `/studio_api/tag/list` | 获取风格标签 |
| GET | `/studio_api/aigc/config` | 获取 AI 创作配置 |

API Key 模式下调用的端点：

| 方法 | 端点 | 说明 |
|------|------|------|
| POST | `/studio_api/create-song-task-with-apikey` | 创建歌曲任务 |
| GET | `/studio_api/assets/get-song-work-with-apikey` | 查询作品（含 AudioURL） |

## Go API

CLI 核心功能封装为 `Client` 结构体，可以直接在 Go 项目中使用：

```go
package main

import (
    "fmt"
    "miaoxiang-cli"  // 替换为实际模块路径
)

func main() {
    // 初始化客户端
    client, err := NewClient("", "cookies.txt")
    if err != nil {
        panic(err)
    }

    // 创建任务
    task, err := client.CreateTask("夏日微风", "", 7, nil, false)
    if err != nil {
        panic(err)
    }
    fmt.Println("TaskID:", task.TaskID)

    // 查询状态
    status, err := client.GetTaskStatus(task.TaskID)
    if err != nil {
        panic(err)
    }
    fmt.Println("Status:", TaskStatusText(status.Status))

    // 自动等待并下载
    err = client.AutoDownload(task.TaskID, "output.wav", 0)
    if err != nil {
        panic(err)
    }
}
```

### Client 方法

| 方法 | 说明 | 认证模式 |
|------|------|----------|
| `NewClient(apiKey, cookieFile)` | 创建客户端 | 两种 |
| `CreateTask(prompt, lyrics, modelID, tags, instrumental)` | 创建任务 | Cookie |
| `GetTaskStatus(taskID)` | 查询任务状态 | Cookie |
| `GetWork(taskID)` | 获取作品详情 | Cookie |
| `GetVidPlayInfo(vids)` | 获取播放 URL | Cookie |
| `ListJobs(page, pageSize)` | 列出任务 | Cookie |
| `GetConfig()` | 获取 AI 配置 | Cookie |
| `GetSingers()` | 获取歌手列表 | Cookie |
| `GetTags()` | 获取风格标签 | Cookie |
| `AutoDownload(taskID, outputPath, interval)` | 等待并下载 | Cookie |
| `CreateSongTaskWithAPIKey(prompt, lyrics, style)` | 创建任务 | API Key |
| `GetSongWorkWithAPIKey(taskID)` | 查询作品 | API Key |

## 项目结构

```
miaoxiang-cli/
├── client.go      # API 客户端实现（HTTP 请求、数据结构、认证）
├── main.go        # CLI 入口（Cobra 命令定义、参数解析）
├── go.mod         # Go 模块定义
├── go.sum         # 依赖校验
├── cookies.txt    # Cookie 文件（gitignore 排除）
├── README.md      # 本文档
└── CHANGELOG.md   # 变更日志
```

## 注意事项

- Cookie 具有时效性，过期后需要重新获取
- 每次提交任务会生成两首歌曲变体，同步模式下载第一首
- 队列等待时间取决于当前平台负载，通常在 3-10 分钟
- 下载的音频格式为 MP3（MPEG ADTS, layer III, 256kbps, 44.1kHz, Stereo）
- `--output` 参数的存在决定了同步/异步模式
