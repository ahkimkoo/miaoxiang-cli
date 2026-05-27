package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	baseURL         = "https://music.douyin.com"
	apiCreate       = "/studio_api/create-studio-task"
	apiGetTask      = "/studio_api/get-studio-task"
	apiConfig       = "/studio_api/aigc/config"
	apiSingers      = "/studio_api/aigc/ai-singers"
	apiTagList      = "/studio_api/tag/list"
	apiWork         = "/studio_api/assets/work"
	apiMultiWorks   = "/studio_api/assets/work-list"
	apiVidPlayInfo  = "/studio_api/video/get-vid-play-info"
	apiCreateSong   = "/studio_api/create-song-task-with-apikey"
	apiGetSongWork  = "/studio_api/assets/get-song-work-with-apikey"
)

// BaseResp 通用响应头
type BaseResp struct {
	ErrorCode int    `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
}

// ---- Cookie Auth Types ----

type CreateStudioTaskRequest struct {
	TaskType                     int                           `json:"TaskType"`
	StudioInspiredCreationParams *StudioInspiredCreationParams `json:"StudioInspiredCreationParams,omitempty"`
}

type StudioInspiredCreationParams struct {
	Prompt              string   `json:"Prompt"`
	Lyrics              string   `json:"Lyrics,omitempty"`
	StudioCreationModel int      `json:"StudioCreationModel"`
	IsInstrumental      bool     `json:"IsInstrumental"`
	TagList             []string `json:"TagList"`
}

type CreateStudioTaskResp struct {
	BaseResp   BaseResp                `json:"baseResp"`
	StatusCode int                     `json:"statusCode"`
	Data       *CreateStudioTaskData   `json:"data,omitempty"`
}

type CreateStudioTaskData struct {
	TaskID     string `json:"TaskID"`
	BizCode    string `json:"BizCode"`
	BizMessage string `json:"BizMessage"`
}

type GetStudioTaskResp struct {
	BaseResp   BaseResp          `json:"baseResp"`
	StatusCode int               `json:"statusCode"`
	Data       *StudioTaskStatus `json:"data,omitempty"`
}

type StudioTaskStatus struct {
	Status         int    `json:"Status"`
	QueueTaskCount string `json:"QueueTaskCount"`
	TaskID         string `json:"TaskID"`
	TaskType       int    `json:"TaskType"`
}

type GetWorkResp struct {
	BaseResp BaseResp `json:"baseResp"`
	Data     *struct {
		UserWorks []UserWork `json:"UserWorks"`
	} `json:"data,omitempty"`
}

type UserWork struct {
	TaskID    string `json:"TaskID"`
	WorkID    string `json:"WorkID"`
	ProjectID string `json:"ProjectID"`
	UserID    string `json:"UserID"`
	TaskType  int    `json:"TaskType"`
	Status    int    `json:"Status"`
	Title     string `json:"Title"`
	Cover     string `json:"Cover"`
	VID       string `json:"VID"`
	Tags      []string `json:"Tags"`
	Duration  string `json:"Duration"`
	Lyrics    string `json:"Lyrics"`
	Prompt    string `json:"Prompt"`
	CreateTime string `json:"CreateTime"`
	UpdateTime string `json:"UpdateTime"`
}

type GetVidPlayInfoReq struct {
	Vids []string `json:"Vids"`
}

type GetVidPlayInfoResp struct {
	BaseResp BaseResp `json:"baseResp"`
	Data     *struct {
		OriginPlayUrls map[string]string `json:"OriginPlayUrls"`
		VideoInfos     *struct {
			Mp3PlayUrls map[string]string `json:"Mp3PlayUrls"`
		} `json:"VideoInfos,omitempty"`
	} `json:"data,omitempty"`
}

type MultiGetWorksResp struct {
	BaseResp BaseResp `json:"baseResp"`
	Data     *struct {
		UserWorks []UserWork `json:"UserWorks"`
		HasMore   bool       `json:"HasMore"`
		Total     string     `json:"Total"`
	} `json:"data,omitempty"`
}

// ---- API Key Types ----

type CreateSongTaskReq struct {
	ApiKey string `json:"ApiKey"`
	Prompt string `json:"Prompt"`
	Lyrics string `json:"Lyrics,omitempty"`
	Style  string `json:"Style,omitempty"`
}

type CreateSongTaskResp struct {
	BaseResp BaseResp `json:"baseResp"`
	TaskID   string   `json:"TaskID,omitempty"`
}

type GetSongWorkResp struct {
	BaseResp BaseResp `json:"baseResp"`
	Data     *struct {
		Works []SongWork `json:"Works"`
	} `json:"data,omitempty"`
}

type SongWork struct {
	TaskID  string `json:"TaskID"`
	WorkID  string `json:"WorkID"`
	Title   string `json:"Title"`
	Status  int    `json:"Status"`
	AudioURL string `json:"AudioURL"`
	VID     string `json:"VID"`
	Duration string `json:"Duration"`
	Lyrics  string `json:"Lyrics"`
	Cover   string `json:"Cover"`
}

// ---- Config Types ----

type AIGCConfigResp struct {
	BaseResp BaseResp `json:"baseResp"`
	Data     *struct {
		SupportModels []ModelInfo `json:"SupportModels"`
	} `json:"data,omitempty"`
}

type ModelInfo struct {
	ModelID      int    `json:"ModelID"`
	ModelType    int    `json:"ModelType"`
	ModelVersion string `json:"ModelVersion"`
	DisplayName  string `json:"DisplayName"`
	IsDefault    bool   `json:"IsDefault"`
}

type AISingerReq struct {
	PitchRange  int `json:"PitchRange,omitempty"`
	AISingerType int `json:"AISingerType,omitempty"`
}

type AISingerResp struct {
	BaseResp   BaseResp   `json:"baseResp"`
	StatusCode int        `json:"statusCode"`
	Data       *SingerData `json:"data,omitempty"`
}

type SingerData struct {
	AISingers []AISinger `json:"AISingers"`
}

type AISinger struct {
	SingerID     string `json:"SingerID"`
	Name         string `json:"Name"`
	Description  string `json:"Description"`
	DemoURL      string `json:"DemoURL,omitempty"`
	CoverURL     string `json:"CoverURL,omitempty"`
	SVCSingerID  string `json:"SVCSingerID"`
	Gender       int    `json:"Gender"`
	ArtistID     string `json:"ArtistID,omitempty"`
	MatchPercent int    `json:"MatchPercent"`
}

type TagListResp struct {
	BaseResp   BaseResp   `json:"baseResp"`
	StatusCode int        `json:"statusCode"`
	Data       *TagData   `json:"data,omitempty"`
}

type TagData struct {
	GenreTags []GenreTag `json:"GenreTags"`
}

type GenreTag struct {
	TagID      string `json:"TagID"`
	TagName    string `json:"TagName"`
	TagNameZh  string `json:"TagNameZh"`
	TagNameEn  string `json:"TagNameEn"`
}

// ---- Client ----

type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

func NewClient(apiKey, cookieFile string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	if cookieFile != "" {
		cookies, err := os.ReadFile(cookieFile)
		if err != nil {
			return nil, fmt.Errorf("读取cookie文件失败: %w", err)
		}
		u, _ := url.Parse("https://music.douyin.com")
		var result []*http.Cookie
		for _, line := range strings.Split(strings.TrimSpace(string(cookies)), ";") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				result = append(result, &http.Cookie{
					Name:  strings.TrimSpace(parts[0]),
					Value: strings.TrimSpace(parts[1]),
				})
			}
		}
		jar.SetCookies(u, result)
	}

	return &Client{
		httpClient: &http.Client{
			Jar:     jar,
			Timeout: 30 * time.Second,
		},
		apiKey:  apiKey,
		baseURL: baseURL,
	}, nil
}

func (c *Client) doRequest(method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "https://music.douyin.com")
		req.Header.Set("Referer", "https://music.douyin.com/studio/create")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Log non-2xx responses to stderr for debugging
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Fprintf(os.Stderr, "[DEBUG] %s %s => HTTP %d\n%s\n", method, path, resp.StatusCode, string(respBody))
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, truncateStr(string(respBody), 200))
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			fmt.Fprintf(os.Stderr, "[DEBUG] %s %s => JSON parse error on response: %s\n", method, path, truncateStr(string(respBody), 300))
			return fmt.Errorf("JSON解析失败: %w (response: %s)", err, truncateStr(string(respBody), 200))
		}
	}
	return nil
}

// truncateStr truncates a string to maxLen for debug output.
func truncateStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// CreateTask 创建歌曲生成任务 (cookie模式)
func (c *Client) CreateTask(prompt, lyrics string, modelID int, tags []string, instrumental bool) (*CreateStudioTaskData, error) {
	req := CreateStudioTaskRequest{
		TaskType: 46, // AI歌曲创作
		StudioInspiredCreationParams: &StudioInspiredCreationParams{
			Prompt:              prompt,
			Lyrics:              lyrics,
			StudioCreationModel: modelID,
			IsInstrumental:      instrumental,
			TagList:             tags,
		},
	}

	var resp CreateStudioTaskResp
	if err := c.doRequest("POST", apiCreate, req, &resp); err != nil {
		return nil, err
	}
	if resp.BaseResp.ErrorCode != 0 {
		return nil, fmt.Errorf("创建任务失败: %s", resp.BaseResp.ErrorMsg)
	}
	if resp.Data == nil {
		return nil, fmt.Errorf("创建任务失败: 无返回数据")
	}
	return resp.Data, nil
}

// GetTaskStatus 查询任务状态 (cookie模式)
func (c *Client) GetTaskStatus(taskID string) (*StudioTaskStatus, error) {
	path := fmt.Sprintf("%s?TaskID=%s", apiGetTask, taskID)
	var resp GetStudioTaskResp
	if err := c.doRequest("GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetWork 获取作品详情 (cookie模式)
func (c *Client) GetWork(taskID string) ([]UserWork, error) {
	path := fmt.Sprintf("%s?TaskID=%s&UserWorkType=2", apiWork, taskID)
	var resp GetWorkResp
	if err := c.doRequest("GET", path, nil, &resp); err != nil {
		return nil, err
	}
	if resp.Data == nil {
		return nil, nil
	}
	return resp.Data.UserWorks, nil
}

// GetVidPlayInfo 获取视频播放URL
func (c *Client) GetVidPlayInfo(vids []string) (*GetVidPlayInfoResp, error) {
	var resp GetVidPlayInfoResp
	if err := c.doRequest("POST", apiVidPlayInfo, GetVidPlayInfoReq{Vids: vids}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListJobs 列出所有任务
func (c *Client) ListJobs(page, pageSize int) ([]UserWork, error) {
	path := fmt.Sprintf("%s?CurrentPage=%d&PageSize=%d&UserWorkType=2", apiMultiWorks, page, pageSize)
	var resp MultiGetWorksResp
	if err := c.doRequest("GET", path, nil, &resp); err != nil {
		return nil, err
	}
	if resp.Data == nil {
		return nil, nil
	}
	return resp.Data.UserWorks, nil
}

// GetConfig 获取AI创作配置
func (c *Client) GetConfig() (*AIGCConfigResp, error) {
	var resp AIGCConfigResp
	if err := c.doRequest("GET", apiConfig, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSingers 获取AI歌手列表
func (c *Client) GetSingers() (*AISingerResp, error) {
	var resp AISingerResp
	if err := c.doRequest("POST", apiSingers, AISingerReq{}, &resp); err != nil {
		return nil, err
	}
	// Debug: log response summary to stderr
	fmt.Fprintf(os.Stderr, "[DEBUG] GetSingers: statusCode=%d, errorCode=%d, errorMsg=%q, hasData=%v",
		resp.StatusCode, resp.BaseResp.ErrorCode, resp.BaseResp.ErrorMsg, resp.Data != nil)
	if resp.Data != nil {
		fmt.Fprintf(os.Stderr, ", singerCount=%d", len(resp.Data.AISingers))
	}
	fmt.Fprintln(os.Stderr)
	return &resp, nil
}

// GetTags 获取风格标签
func (c *Client) GetTags() (*TagListResp, error) {
	var resp TagListResp
	if err := c.doRequest("GET", apiTagList, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateSongTaskWithAPIKey API Key模式创建任务
func (c *Client) CreateSongTaskWithAPIKey(prompt, lyrics, style string) (*CreateSongTaskResp, error) {
	req := CreateSongTaskReq{
		ApiKey: c.apiKey,
		Prompt: prompt,
		Lyrics: lyrics,
		Style:  style,
	}
	var resp CreateSongTaskResp
	if err := c.doRequest("POST", apiCreateSong, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSongWorkWithAPIKey API Key模式获取作品
func (c *Client) GetSongWorkWithAPIKey(taskID string) (*SongWork, error) {
	path := fmt.Sprintf("%s?ApiKey=%s&TaskID=%s", apiGetSongWork, c.apiKey, taskID)
	var resp GetSongWorkResp
	if err := c.doRequest("GET", path, nil, &resp); err != nil {
		return nil, err
	}
	if resp.Data != nil && len(resp.Data.Works) > 0 {
		return &resp.Data.Works[0], nil
	}
	return nil, nil
}

// TaskStatusText 返回状态文本
func TaskStatusText(status int) string {
	switch status {
	case 1:
		return "排队中"
	case 2:
		return "处理中"
	case 3:
		return "已完成"
	case 5:
		return "失败"
	default:
		return fmt.Sprintf("未知(%d)", status)
	}
}

// DownloadFile 下载文件
func DownloadFile(fileURL, destPath string) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// DownloadResult 单个文件的下载结果
type DownloadResult struct {
	File     string `json:"file"`
	Title    string `json:"title"`
	Duration string `json:"duration"`
}

// buildOutputPaths 根据 outputPath 和作品列表生成不冲突的文件路径
// outputPath 可以是目录（以/结尾）或文件路径；以作品标题命名，重复时自动加后缀
func buildOutputPaths(outputPath string, works []UserWork) []string {
	ext := filepath.Ext(outputPath)
	if ext == "" {
		ext = ".wav"
	}

	var dir string
	if strings.HasSuffix(outputPath, "/") || outputPath == "." {
		dir = outputPath
	} else {
		dir = filepath.Dir(outputPath)
		if dir == "" {
			dir = "."
		}
		os.MkdirAll(dir, 0755)
	}

	usedNames := map[string]int{}
	var paths []string
	for _, w := range works {
		name := w.Title
		if name == "" {
			name = fmt.Sprintf("song_%d", w.WorkID)
		}
		// 清理文件名中的非法字符
		name = strings.Map(func(r rune) rune {
			if r == '/' || r == '\\' || r == ':' || r == '*' || r == '?' || r == '"' || r == '<' || r == '>' || r == '|' {
				return '_'
			}
			return r
		}, name)

		base := name + ext
		if count := usedNames[base]; count > 0 {
			base = fmt.Sprintf("%s-%d%s", name, count+1, ext)
		}
		usedNames[base]++

		paths = append(paths, filepath.Join(dir, base))
	}
	return paths
}

// AutoDownload 等待任务完成并下载全部作品
func (c *Client) AutoDownload(taskID string, outputPath string, interval time.Duration) error {
	fmt.Fprintf(os.Stderr, "等待任务 %s 完成...\n", taskID)
	for {
		status, err := c.GetTaskStatus(taskID)
		if err != nil {
			return fmt.Errorf("查询状态失败: %w", err)
		}

		fmt.Fprintf(os.Stderr, "  状态: %s, 队列: %s\n", TaskStatusText(status.Status), status.QueueTaskCount)

		if status.Status == 3 {
			works, err := c.GetWork(taskID)
			if err != nil {
				return fmt.Errorf("获取作品失败: %w", err)
			}
			if len(works) == 0 {
				return fmt.Errorf("任务完成但未找到作品")
			}

			var vids []string
			for _, w := range works {
				if w.VID != "" {
					vids = append(vids, w.VID)
				}
			}
			if len(vids) == 0 {
				return fmt.Errorf("作品没有VID")
			}

			playInfo, err := c.GetVidPlayInfo(vids)
			if err != nil {
				return fmt.Errorf("获取播放URL失败: %w", err)
			}

			paths := buildOutputPaths(outputPath, works)
			var results []DownloadResult

			for i, work := range works {
				audioURL := ""
				if playInfo.Data != nil {
					if playInfo.Data.VideoInfos != nil {
						audioURL = playInfo.Data.VideoInfos.Mp3PlayUrls[work.VID]
					}
					if audioURL == "" {
						audioURL = playInfo.Data.OriginPlayUrls[work.VID]
					}
				}
				if audioURL == "" {
					fmt.Fprintf(os.Stderr, "  跳过 %s: 无音频URL (VID=%s)\n", work.Title, work.VID)
					continue
				}

				fmt.Fprintf(os.Stderr, "  下载: %s -> %s\n", work.Title, paths[i])
				if err := DownloadFile(audioURL, paths[i]); err != nil {
					fmt.Fprintf(os.Stderr, "  失败: %s: %v\n", work.Title, err)
					continue
				}

				absPath, _ := filepath.Abs(paths[i])
				results = append(results, DownloadResult{
					File:     absPath,
					Title:    work.Title,
					Duration: work.Duration,
				})
			}

			return printJSON(results)
		}

		if status.Status == 5 {
			return fmt.Errorf("任务失败")
		}

		time.Sleep(interval)
	}
}
