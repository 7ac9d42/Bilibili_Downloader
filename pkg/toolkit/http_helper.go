package toolkit

import (
	"io"
	"net/http"

	"github.com/cheggaaa/pb/v3"
)

// SetBilibiliHeaders 设置B站API请求的通用头信息
func SetBilibiliHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://www.bilibili.com/")
	req.Header.Set("Origin", "https://www.bilibili.com/")
}

// SetVideoHeaders 设置视频请求的通用头信息
func SetVideoHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Referer", "https://www.bilibili.com/video")
	req.Header.Set("Origin", "https://www.bilibili.com")
}

// DownloadAndTrackProgress 从响应体下载内容并跟踪进度
func DownloadAndTrackProgress(reader io.ReadCloser, writer io.Writer, bar *pb.ProgressBar) error {
	// 创建一个代理的reader，用于更新进度条
	proxyReader := bar.NewProxyReader(reader)
	_, err := io.Copy(writer, proxyReader)
	return err
}
