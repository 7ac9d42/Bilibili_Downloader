package internal

import (
	"Bilibili_Downloader/pkg/httpclient"
	"Bilibili_Downloader/pkg/toolkit"
	"Bilibili_Downloader/pkg/toolkit/data_struct"
	"encoding/json"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func CatchData(Url string) ([]byte, error) {
	client := httpclient.GetClient()
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return nil, err
	}

	toolkit.SetBilibiliHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("Close resp.Body失败:", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取 URL 失败，状态码: %d", resp.StatusCode)
	}
	log.Printf("请求响应成功，状态码：%v\n", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	return body, nil
}

// prepareDownloadPaths 准备下载文件路径
func prepareDownloadPaths(filepath string) (string, string, error) {
	var filepath1, filepath2 string

	if filepath == "" {
		if err := toolkit.CheckAndCreateCacheDir(); err != nil {
			return "", "", fmt.Errorf("检查并创建临时下载目录失败: %w", err)
		}
		filepath1 = "./download_cache/audio_cache"
		filepath2 = "./download_cache/video_cache"
	} else {
		// 检查字符串末尾是否已经有斜杠
		if !strings.HasSuffix(filepath, "/") {
			// 如果没有，则在末尾添加斜杠
			filepath += "/"
		}
		// 如果指定了文件路径，则在文件路径后添加适当的扩展名
		filepath1 = filepath + "audio_cache"
		filepath2 = filepath + "video_cache"
	}

	return filepath1, filepath2, nil
}

// downloadSingleFile 下载单个文件
func downloadSingleFile(client *http.Client, url string, filePath string, bar *pb.ProgressBar) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	toolkit.SetVideoHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("Close resp.Body失败:", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载请求状态错误: %s", resp.Status)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := out.Close(); err != nil {
			log.Println("Close 文件失败:", err)
		}
	}()

	return toolkit.DownloadAndTrackProgress(resp.Body, out, bar)
}

func DownloadFile(urlVideo string, urlAudio string, filepath string) error {
	audioPath, videoPath, err := prepareDownloadPaths(filepath)
	if err != nil {
		log.Println(err)
		fmt.Println(err)
		return err
	}

	client := httpclient.GetClient()

	fmt.Println("发送下载请求")

	// 先获取文件大小以设置进度条
	reqAudio, err := http.NewRequest("HEAD", urlAudio, nil)
	if err != nil {
		return err
	}
	toolkit.SetVideoHeaders(reqAudio)

	reqVideo, err := http.NewRequest("HEAD", urlVideo, nil)
	if err != nil {
		return err
	}
	toolkit.SetVideoHeaders(reqVideo)

	respAudio, err := client.Do(reqAudio)
	if err != nil {
		return err
	}
	defer respAudio.Body.Close()

	respVideo, err := client.Do(reqVideo)
	if err != nil {
		return err
	}
	defer respVideo.Body.Close()

	totalSize := respAudio.ContentLength + respVideo.ContentLength
	bar := pb.StartNew(int(totalSize))
	bar.Set(pb.SIBytesPrefix, true)

	fmt.Println("正在下载，请耐心等待...")
	log.Println("视频下载开始")

	// 下载音频
	if err := downloadSingleFile(client, urlAudio, audioPath, bar); err != nil {
		return fmt.Errorf("音频下载失败: %w", err)
	}

	// 下载视频
	if err := downloadSingleFile(client, urlVideo, videoPath, bar); err != nil {
		return fmt.Errorf("视频下载失败: %w", err)
	}

	bar.Finish()

	toolkit.ClearScreen()
	fmt.Println("下载完毕！")
	log.Println("视频下载成功")
	return nil
}

// ProcessResponse 处理 JSON 响应并返回 Response 结构体
func ProcessResponse(data []byte, flag int) (interface{}, error) {
	var err error
	if flag == 0 {
		var response data_struct.VideoInfoResponse
		err = json.Unmarshal(data, &response)
		if err != nil {
			return nil, fmt.Errorf("视频信息解组失败: %w", err)
		}
		fmt.Println("视频信息数据解组正常！")
		log.Println("视频信息数据解组正常")
		return &response, nil
	} else if flag == 1 {
		var response data_struct.DownloadInfoResponse
		err = json.Unmarshal(data, &response)
		if err != nil {
			return nil, fmt.Errorf("下载信息解组失败: %w", err)
		}
		fmt.Println("下载信息数据解组正常！")
		log.Println("下载信息数据解组正常")
		return &response, nil
	}
	return nil, fmt.Errorf("不支持的 flag 值: %d", flag)
}
