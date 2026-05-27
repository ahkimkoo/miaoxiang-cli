---
name: miaoxiang-cli
description: Use when generating AI music, creating songs with lyrics, producing instrumental tracks, querying AI singer lists, style tags, or available models via the miaoxiang-cli command-line tool
---

# miaoxiang-cli

妙响音乐生成平台 CLI，基于抖音音乐 API。支持 Cookie 和 API Key 两种认证。

## Auth

`--api-key` > `--cookies` > 可执行文件同目录 `cookies.txt` > 当前目录 `cookies.txt`

## Commands

```bash
miaoxiang-cli gen --lyrics "歌词" --style "风格" --singer "歌手" --output ./songs/
miaoxiang-cli gen --style "电子/舞曲" --instrumental --output ./songs/
miaoxiang-cli models       # 查看可用模型（ModelType 用于 --model）
miaoxiang-cli singers      # 查看 AI 歌手
miaoxiang-cli tags         # 查看风格标签
miaoxiang-cli jobs         # 列出任务
miaoxiang-cli job --id <ID> # 查询单个任务，加 --output 等待下载
```

## gen

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--lyrics` | 非 instrumental 时 | — | 歌词 |
| `--style` | 是 | — | 风格标签全称（`tags` 返回的 TagName） |
| `--singer` | 否 | — | AI 歌手名称（`singers` 返回的 Name） |
| `--instrumental` | 否 | false | 纯音乐 |
| `--model` | 否 | 7 | ModelType 值（`models` 返回） |
| `--output` | 否 | — | 设置后同步等待下载全部歌曲 |
| `--interval` | 否 | 15 | 轮询秒数 |

## Download Output

`--output` 设置后下载全部歌曲，stdout 输出 JSON，stderr 输出进度：

```json
[{"file":"/abs/path/歌曲名.wav","title":"歌曲名","duration":"03:20"}]
```

- 目录路径（`./songs/`）→ 以歌曲标题命名存入
- 文件路径（`./song.wav`）→ 以歌曲标题命名存入同目录
- 同名加 `-2`、`-3` 后缀

## Status

| Code | Meaning |
|------|---------|
| 1 | 排队中 |
| 2 | 处理中 |
| 3 | 已完成 |
| 5 | 失败 |

## Common Mistakes

- `--model` 用 `ModelType`（如 7），不是 `ModelID`（固定 0）
- `--style` 用 `TagName` 全称（如 `"电子/舞曲"`），不是 `TagID`
- `--singer` 用 `Name`（如 `"叶雪如"`），不是 `SingerID`
- singers API 返回 `AISingers`（复数），不是 `Singers`
- cookies.txt：`name=value; name=value` 单行分号分隔
