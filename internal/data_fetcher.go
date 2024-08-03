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

	//client := &http.Client{}
	client := httpclient.GetClient()
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return nil, err
	}

	// 设置自定义请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://www.bilibili.com/")
	req.Header.Set("Origin", "https://www.bilibili.com/")

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

func DownloadFile(urlVideo string, urlAudio string, filepath string) error {
	var filepath1, filepath2 string

	if filepath == "" {
		if err := toolkit.CheckAndCreateCacheDir(); err != nil {
			log.Println("检查并创建临时下载目录失败", err)
			fmt.Println("检查并创建临时下载目录失败")
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

	//client := &http.Client{}
	client := httpclient.GetClient()
	req1, err := http.NewRequest("GET", urlAudio, nil)
	if err != nil {
		return err
	}
	req2, err := http.NewRequest("GET", urlVideo, nil)
	if err != nil {
		return err
	}

	// 设置自定义请求头
	req1.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0")
	req1.Header.Set("Accept", "*/*")
	req1.Header.Set("Referer", "https://www.bilibili.com/vedio")
	req1.Header.Set("Origin", "https://www.bilibili.com")
	// 设置自定义请求头
	req2.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0")
	req2.Header.Set("Accept", "*/*")
	req2.Header.Set("Referer", "https://www.bilibili.com/vedio")
	req2.Header.Set("Origin", "https://www.bilibili.com")

	resp1, err := client.Do(req1)
	if err != nil {
		return err
	}
	resp2, err := client.Do(req2)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp1.Body.Close(); err != nil {
			log.Println("Close resp1.body失败：", err)
		}
		if err := resp2.Body.Close(); err != nil {
			log.Println("Close resp2.body失败：", err)
		}
	}()

	fmt.Println("发送下载请求")
	// 检查HTTP响应状态码
	if resp1.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status of audio: %s", resp1.Status)
	}
	if resp2.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status of video: %s", resp2.Status)
	}

	fmt.Println("正在下载，请耐心等待...")
	log.Println("视频下载开始")

	// 创建文件
	out1, err := os.Create(filepath1)
	if err != nil {
		return err
	}
	defer func() {
		if err := out1.Close(); err != nil {
			log.Println("Close out1文件失败：", err)
		}
	}()

	// 创建文件
	out2, err := os.Create(filepath2)
	if err != nil {
		return err
	}
	defer func() {
		if err := out2.Close(); err != nil {
			log.Println("Close out2文件失败：", err)
		}
	}()

	totalSize := resp1.ContentLength + resp2.ContentLength
	bar := pb.StartNew(int(totalSize))
	bar.Set(pb.SIBytesPrefix, true)

	err = toolkit.DownloadAndTrackProgress(resp1.Body, out1, bar)
	if err != nil {
		return err
	}
	err = toolkit.DownloadAndTrackProgress(resp2.Body, out2, bar)
	if err != nil {
		return err
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
