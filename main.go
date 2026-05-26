package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var (
	flagAPIKey     string
	flagCookies    string
	flagLyrics     string
	flagStyle      string
	flagSinger     string
	flagOutput     string
	flagInstrumental bool
	flagModel      int
	flagInterval   int
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "miaoxiang-cli",
		Short: "妙响音乐生成CLI",
		Long:  "妙响音乐生成平台命令行工具，支持歌曲生成、任务查询、作品下载",
	}

	rootCmd.PersistentFlags().StringVar(&flagAPIKey, "api-key", "", "API Key (可选，默认使用cookie)")
	rootCmd.PersistentFlags().StringVar(&flagCookies, "cookies", "", "Cookie文件路径 (默认: 同目录cookies.txt)")

	genCmd := &cobra.Command{
		Use:   "gen",
		Short: "提交歌曲生成任务",
		Long: `提交歌曲生成任务。

示例:
  # 带歌词的歌曲生成（同步下载）
  miaoxiang-cli gen --lyrics "春风吹过田野，花开满山坡" --style "民谣" --singer "叶雪如" --output song.wav

  # 纯音乐生成（不需要歌词）
  miaoxiang-cli gen --style "电子/舞曲" --instrumental --output instrumental.wav

  # 异步提交（不等待，返回任务ID）
  miaoxiang-cli gen --lyrics "夜空中最亮的星" --style "流行"

  # 查看可用模型/歌手/风格
  miaoxiang-cli models
  miaoxiang-cli singers
  miaoxiang-cli tags`,
		RunE: cmdGen,
	}
	genCmd.Flags().StringVar(&flagLyrics, "lyrics", "", "歌词 (非纯音乐时必填)")
	genCmd.Flags().StringVar(&flagStyle, "style", "", "风格标签 (必填，用 tags 命令查看)")
	genCmd.Flags().StringVar(&flagSinger, "singer", "", "AI歌手名称 (可选，用 singers 命令查看)")
	genCmd.Flags().StringVar(&flagOutput, "output", "", "输出文件路径 (设置此参数表示同步等待并下载)")
	genCmd.Flags().BoolVar(&flagInstrumental, "instrumental", false, "纯音乐模式 (无需歌词)")
	genCmd.Flags().IntVar(&flagModel, "model", 7, "模型类型 (用 models 命令查看可选模型, 默认: 7=V5.5)")
	genCmd.Flags().IntVar(&flagInterval, "interval", 15, "轮询间隔(秒)")
	genCmd.MarkFlagRequired("style")

	jobCmd := &cobra.Command{
		Use:   "job",
		Short: "查询单个任务",
		Long: `查询单个任务状态
示例:
  miaoxiang-cli job --id xxx  (返回状态JSON)
  miaoxiang-cli job --id xxx --output song.wav (同步等待并下载)`,
		RunE: cmdJob,
	}
	jobCmd.Flags().StringVar(&flagOutput, "output", "", "输出文件路径 (设置此参数表示同步等待并下载)")
	jobCmd.Flags().IntVar(&flagInterval, "interval", 15, "轮询间隔(秒)")
	jobCmd.Flags().String("id", "", "任务ID")
	jobCmd.MarkFlagRequired("id")

	jobsCmd := &cobra.Command{
		Use:   "jobs",
		Short: "列出所有任务",
		RunE:  cmdJobs,
	}

	singersCmd := &cobra.Command{
		Use:   "singers",
		Short: "获取AI歌手列表",
		RunE:  cmdSingers,
	}

	tagsCmd := &cobra.Command{
		Use:   "tags",
		Short: "获取风格标签",
		RunE:  cmdTags,
	}

	modelsCmd := &cobra.Command{
		Use:   "models",
		Short: "获取可用模型列表",
		RunE:  cmdModels,
	}

	rootCmd.AddCommand(genCmd, jobCmd, jobsCmd, singersCmd, tagsCmd, modelsCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func getClient() (*Client, error) {
	cookieFile := flagCookies
	if cookieFile == "" && flagAPIKey == "" {
		// 尝试同目录下的 cookies.txt
		execPath, err := os.Executable()
		if err == nil {
			cookieFile = filepath.Join(filepath.Dir(execPath), "cookies.txt")
		}
		if _, err := os.Stat(cookieFile); err != nil {
			// 尝试当前工作目录
			cookieFile = "cookies.txt"
			if _, err := os.Stat(cookieFile); err != nil {
				cookieFile = ""
			}
		}
	}
	return NewClient(flagAPIKey, cookieFile)
}

func cmdGen(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	if flagAPIKey != "" {
		// API Key模式
		resp, err := client.CreateSongTaskWithAPIKey(flagLyrics, flagLyrics, flagStyle)
		if err != nil {
			return err
		}
		fmt.Printf("任务已创建: %s\n", resp.TaskID)

		if flagOutput != "" {
			// 同步等待并下载
			return autoDownloadAPIKey(client, resp.TaskID, flagOutput, time.Duration(flagInterval)*time.Second)
		}
		// 异步，返回JSON
		return printJSON(map[string]interface{}{
			"TaskID": resp.TaskID,
			"Status": "created",
		})
	}

	// Cookie模式
	resp, err := client.CreateTask(flagLyrics, "", flagModel, nil, flagInstrumental)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "任务已创建: %s\n", resp.TaskID)

	if flagOutput != "" {
		// 同步等待并下载
		return client.AutoDownload(resp.TaskID, flagOutput, time.Duration(flagInterval)*time.Second)
	}

	// 异步，返回JSON
	return printJSON(map[string]interface{}{
		"TaskID":   resp.TaskID,
		"Status":   "created",
		"TaskType": 46,
	})
}

func cmdJob(cmd *cobra.Command, args []string) error {
	taskID, _ := cmd.Flags().GetString("id")
	client, err := getClient()
	if err != nil {
		return err
	}

	if flagOutput != "" {
		// 同步等待并下载
		return client.AutoDownload(taskID, flagOutput, time.Duration(flagInterval)*time.Second)
	}

	// 查询状态
	status, err := client.GetTaskStatus(taskID)
	if err != nil {
		return err
	}

	// 获取作品详情
	works, _ := client.GetWork(taskID)

	result := map[string]interface{}{
		"TaskID":         taskID,
		"Status":         status.Status,
		"StatusText":     TaskStatusText(status.Status),
		"QueueTaskCount": status.QueueTaskCount,
	}

	if len(works) > 0 {
		var workList []map[string]interface{}
		for _, w := range works {
			workList = append(workList, map[string]interface{}{
				"WorkID":   w.WorkID,
				"Title":    w.Title,
				"Status":   w.Status,
				"VID":      w.VID,
				"Duration": w.Duration,
			})
		}
		result["Works"] = workList
	}

	return printJSON(result)
}

func cmdJobs(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	jobs, err := client.ListJobs(1, 50)
	if err != nil {
		return err
	}

	if len(jobs) == 0 {
		fmt.Println("暂无任务")
		return nil
	}

	for _, j := range jobs {
		fmt.Printf("[%s] %s | %s | VID: %s | %s\n",
			TaskStatusText(j.Status), j.TaskID, j.Title, j.VID, j.Duration)
	}
	return nil
}

func cmdSingers(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	resp, err := client.GetSingers()
	if err != nil {
		return err
	}
	if resp.Data == nil || len(resp.Data.AISingers) == 0 {
		fmt.Println("暂无歌手")
		return nil
	}

	return printJSON(resp.Data.AISingers)
}

func cmdTags(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	resp, err := client.GetTags()
	if err != nil {
		return err
	}
	if resp.Data == nil || len(resp.Data.GenreTags) == 0 {
		fmt.Println("暂无标签")
		return nil
	}

	return printJSON(resp.Data.GenreTags)
}

func cmdModels(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	resp, err := client.GetConfig()
	if err != nil {
		return err
	}
	if resp.Data == nil || len(resp.Data.SupportModels) == 0 {
		fmt.Println("暂无模型配置")
		return nil
	}

	return printJSON(resp.Data.SupportModels)
}

func printJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func autoDownloadAPIKey(client *Client, taskID, outputPath string, interval time.Duration) error {
	fmt.Printf("等待任务 %s 完成...\n", taskID)
	for {
		work, err := client.GetSongWorkWithAPIKey(taskID)
		if err != nil {
			return fmt.Errorf("查询状态失败: %w", err)
		}
		if work != nil && work.Status == 3 {
			if work.AudioURL == "" {
				return fmt.Errorf("任务完成但无音频URL")
			}
			fmt.Printf("下载音频到: %s\n", outputPath)
			if err := DownloadFile(work.AudioURL, outputPath); err != nil {
				return fmt.Errorf("下载失败: %w", err)
			}
			fmt.Printf("下载完成: %s\n", outputPath)
			return nil
		}
		if work != nil && work.Status == 5 {
			return fmt.Errorf("任务失败")
		}
		fmt.Printf("  状态: 处理中\n")
		time.Sleep(interval)
	}
}
