# 妙响 CLI (miaoxiang-cli)

[抖音妙响音乐生成平台](https://music.douyin.com/studio/create) 命令行工具，支持 AI 歌曲生成、任务查询、作品下载，提供 Cookie 和 API Key 两种认证模式。

## 特性

- 两种认证模式：Cookie（模拟浏览器）和 API Key（官方开放接口）
- 同步/异步两种任务提交模式
- 一次生成多首歌曲，全部自动下载，以歌曲标题命名
- 查询可用模型、AI 歌手列表、风格标签
- 管理全部历史任务
- 纯音乐模式（无需歌词）
- 零外部运行时依赖，单二进制文件即可运行
- 跨平台支持：macOS / Linux / Windows

## 下载

从 [GitHub Releases](https://github.com/ahkimkoo/miaoxiang-cli/releases) 下载对应平台的可执行文件。

```bash
# macOS / Linux
chmod +x miaoxiang-cli-darwin-arm64
sudo cp miaoxiang-cli-darwin-arm64 /usr/local/bin/miaoxiang-cli
```

### 从源码编译

需要 Go 1.21+。

```bash
git clone https://github.com/ahkimkoo/miaoxiang-cli.git
cd miaoxiang-cli
go build -o miaoxiang-cli .
```

## 认证

CLI 支持两种认证方式，Cookie 模式为默认推荐方式。

### Cookie 模式（推荐）

**获取 Cookie**：使用 Chrome 打开 https://music.douyin.com/studio/create，登录后 F12 → Application → Cookies，复制为 `name=value; name=value` 格式。

**配置 Cookie**：

```bash
# 放在可执行文件同目录（自动加载）
cp cookies.txt /usr/local/bin/cookies.txt

# 或通过参数指定
miaoxiang-cli --cookies /path/to/cookies.txt gen ...
```

### API Key 模式

```bash
miaoxiang-cli --api-key "YOUR_API_KEY" gen --lyrics "xxx" --style "流行"
```

## 使用

```bash
# 带歌词生成（同步下载）
miaoxiang-cli gen --lyrics "春风吹过田野" --style "民谣" --singer "叶雪如" --output ./songs/

# 纯音乐（无需歌词）
miaoxiang-cli gen --style "电子/舞曲" --instrumental --output ./songs/

# 异步提交（返回任务ID）
miaoxiang-cli gen --lyrics "夜空中最亮的星" --style "流行"

# 查询
miaoxiang-cli models     # 查看可用模型
miaoxiang-cli singers    # 查看AI歌手
miaoxiang-cli tags       # 查看风格标签
miaoxiang-cli jobs       # 列出任务
miaoxiang-cli job --id <taskID>   # 查询单个任务
```

### 下载输出

同步模式（指定 `--output`）会下载全部生成的歌曲，输出 JSON：

```json
[
  {
    "file": "/abs/path/春风.wav",
    "title": "春风",
    "duration": "03:20"
  },
  {
    "file": "/abs/path/春风-2.wav",
    "title": "春风",
    "duration": "03:58"
  }
]
```

- `--output` 传目录（以 `/` 结尾）→ 用歌曲标题命名存入该目录
- `--output` 传文件路径 → 用歌曲标题命名存入同目录
- 标题重复自动加 `-2`、`-3` 后缀
- 进度信息输出到 stderr，JSON 结果输出到 stdout

## 命令参考

### 全局参数

| 参数 | 说明 |
|------|------|
| `--api-key string` | API Key |
| `--cookies string` | Cookie 文件路径（默认自动搜索） |

### `gen` — 提交生成任务

| 参数 | 必填 | 说明 |
|------|------|------|
| `--lyrics string` | 非 `--instrumental` 时 | 歌词内容 |
| `--style string` | 是 | 风格标签（用 `tags` 查看） |
| `--singer string` | 否 | AI 歌手名称（用 `singers` 查看） |
| `--model int` | 否 | 模型类型（用 `models` 查看，默认 7） |
| `--instrumental` | 否 | 纯音乐模式 |
| `--output string` | 否 | 输出路径，设置后同步等待并下载 |
| `--interval int` | 否 | 轮询间隔秒数（默认 15） |

### `models` — 查看可用模型

返回 `ModelType`（用于 `--model`）和 `ModelVersion`。

### `singers` — 查看 AI 歌手

返回 `Name`、`Description`（音域/风格）、`Gender`（0=男/1=女）、`DemoURL`。

### `tags` — 查看风格标签

返回 `TagName`（中英文）。`gen --style` 使用全称。

### `job` — 查询单个任务

| 参数 | 必填 | 说明 |
|------|------|------|
| `--id string` | 是 | 任务 ID |
| `--output string` | 否 | 设置后同步等待下载 |
| `--interval int` | 否 | 轮询间隔秒数 |

### `jobs` — 列出所有任务

列出最近 50 个任务。

## 任务状态

| 状态码 | 说明 |
|--------|------|
| 1 | 排队中 |
| 2 | 处理中 |
| 3 | 已完成 |
| 5 | 失败 |

## 注意事项

- Cookie 具有时效性，过期后需重新获取
- 每次任务可能生成多首歌曲变体，同步模式全部下载
- 队列等待时间通常 3-10 分钟
- `--model` 用 `ModelType` 值（如 7），不是 `ModelID`
- `--style` 用 `TagName` 全称（如 `"电子/舞曲"`），不是 `TagID`
- `--singer` 用歌手 `Name`（如 `"叶雪如"`），不是 `SingerID`
