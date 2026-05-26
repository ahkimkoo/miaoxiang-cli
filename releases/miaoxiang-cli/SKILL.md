---
name: miaoxiang-cli
description: Use when generating AI music, creating songs with lyrics, producing instrumental tracks, querying AI singer lists, style tags, or available models via the miaoxiang-cli command-line tool
---

# miaoxiang-cli

妙响音乐生成平台 CLI 工具，基于抖音音乐 API。

## Auth

优先级：`--api-key` > `--cookies` > 可执行文件同目录 `cookies.txt` > 当前目录 `cookies.txt`

## Commands

```bash
miaoxiang-cli gen --lyrics "歌词" --style "风格" --singer "歌手" --output song.wav
miaoxiang-cli gen --style "电子/舞曲" --instrumental --output bgm.wav
miaoxiang-cli models
miaoxiang-cli singers
miaoxiang-cli tags
miaoxiang-cli jobs
miaoxiang-cli job --id <taskID>
```

## gen — 提交生成任务

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--lyrics` | 非 `--instrumental` 时必填 | — | 歌词文本 |
| `--style` | 必填 | — | 风格标签，用 `tags` 查看 |
| `--singer` | 可选 | — | AI 歌手名称，用 `singers` 查看 |
| `--instrumental` | 可选 | false | 纯音乐模式，无需歌词 |
| `--model` | 可选 | 7 | 模型类型，用 `models` 查看 |
| `--output` | 可选 | — | 设置后同步等待并下载到指定路径 |
| `--interval` | 可选 | 15 | 轮询间隔（秒） |

## models — 查看可用模型

返回 `ModelType`（用于 `--model`）、`ModelVersion`。`ModelID` 字段固定为 0，以 `ModelType` 为准。

## singers — 查看 AI 歌手

返回歌手列表：`Name`、`Description`（音域/风格）、`Gender`（0=男/1=女）、`DemoURL`。

## tags — 查看风格标签

返回 `TagName`（中英文）、`TagID`。`gen --style` 使用 `TagName` 的值。

## job / jobs — 查询任务

- `jobs`：列出最近 50 个任务，显示状态、VID、时长
- `job --id <taskID>`：查询单个任务；加 `--output` 同步等待下载

## Task Status

| Status | 含义 |
|--------|------|
| 1 | 排队中 |
| 2 | 处理中 |
| 3 | 已完成 |
| 5 | 失败 |

## API Endpoints (Cookie Mode)

| Path | Method | Purpose |
|------|--------|---------|
| `/studio_api/create-studio-task` | POST | 创建生成任务 |
| `/studio_api/get-studio-task` | GET | 查询任务状态 |
| `/studio_api/aigc/config` | GET | 获取模型列表 |
| `/studio_api/aigc/ai-singers` | POST | 获取 AI 歌手（注意字段是 `AISingers`） |
| `/studio_api/tag/list` | GET | 获取风格标签 |
| `/studio_api/assets/work` | GET | 获取作品详情 |
| `/studio_api/assets/work-list` | GET | 列出所有任务 |
| `/studio_api/video/get-vid-play-info` | POST | 获取音频播放 URL |

Base URL: `https://music.douyin.com`

## Common Mistakes

- `--model` 用 `ModelType` 值（如 7），不是 `ModelID`（固定 0）
- `--style` 用 `tags` 返回的 `TagName` 全称（如 `"电子/舞曲"`），不是 `TagID`
- `--singer` 用歌手 `Name`（如 `"叶雪如"`），不是 `SingerID`
- `AISingers` 接口返回字段是 `AISingers`（复数），不是 `Singers`
- cookies.txt 格式：`name=value; name=value` 单行，分号分隔
