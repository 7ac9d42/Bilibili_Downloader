package video_processing

import (
	"Bilibili_Downloader/pkg/toolkit"
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed ffmpeg/ffmpeg
var embeddedFFmpeg embed.FS

func extractFFmpeg() (string, error) {
	// 读取嵌入的ffmpeg二进制文件
	data, err := embeddedFFmpeg.ReadFile("ffmpeg/ffmpeg")
	if err != nil {
		return "", err
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "ffmpeg")
	if err != nil {
		return "", err
	}

	// 写入ffmpeg二进制文件
	ffmpegPath := filepath.Join(tempDir, "ffmpeg.exe")
	err = os.WriteFile(ffmpegPath, data, 0644)
	if err != nil {
		return "", err
	}

	return ffmpegPath, nil
}

func Transcoding(videoName string, definition string) {
	// 提取ffmpeg
	ffmpegPath, err := extractFFmpeg()
	if err != nil {
		log.Println("ffmpeg释放错误:", err)
		return
	}
	defer func() {
		if err := os.RemoveAll(filepath.Dir(ffmpegPath)); err != nil {
			log.Println("移除临时FFmpeg释放目录失败：", err)
			fmt.Println("移除FFmpeg释放目录失败，详情见log日志，可手动移除，路径：", ffmpegPath)
		}
	}()

	// 定义要处理的文件目录和输出的MP4文件名
	inputDir := "./download_cache"
	if err := toolkit.CheckAndCreateDir("./Download"); err != nil {
		log.Println("视频输出目录检查或创建失败：", err)
	}
	outputFile := fmt.Sprintf("./Download/%s.mp4", toolkit.CheckAndCleanFileName(fmt.Sprintf("%s(%s)", videoName, definition)))

	// 获取所有的cache文件
	caches, err := filepath.Glob(filepath.Join(inputDir, "*"))
	if err != nil {
		log.Println("查找Cache文件出错:", err)
		return
	}

	if len(caches) == 0 {
		log.Println("未找到缓存文件")
		return
	}

	// 创建ffmpeg命令
	cmd := exec.Command(ffmpegPath, "-f", "mp4", "-i", caches[0], "-i", caches[1], "-c:a", "copy", "-c:v", "copy")

	// 将输出重定向到标准输出和标准错误
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//将outputFile作为参数提供而不是直接构建命令，避免可能出现的由于特殊字符导致的错误解释
	cmd.Args = append(cmd.Args, outputFile)

	// 运行ffmpeg命令
	err = cmd.Run()
	toolkit.ClearScreen()
	if err != nil {
		log.Println("ffmpeg运行失败:", err)
		fmt.Println("\n视频文件转码失败，请携带日志文件(log)联系开发者！")
		return
	} else {
		log.Println("视频文件转码成功")
		fmt.Println("视频文件转码成功！")
	}
}
