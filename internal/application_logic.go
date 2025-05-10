package internal

import (
	"Bilibili_Downloader/internal/update"
	"Bilibili_Downloader/pkg/toolkit"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func InitLog() *os.File {
	//初始化日志文件
	logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
	// 将日志输出设置到文件
	log.SetOutput(logFile)
	// 设置日志前缀和格式
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("程序运行，开始日志记录")
	return logFile
}

func CacheClean() {
	if err := toolkit.RemoveCacheDir(); err != nil {
		toolkit.ClearScreen()
		log.Println("缓存目录清理失败:", err)
		fmt.Println("缓存目录清理失败，确认需清理时可手动清理或重新运行程序")
	} else {
		fmt.Println("缓存目录清理完毕")
	}
}

func HandleUpdate() {
	//检查更新
	if check, newProgramName := update.CheckAndUpdate(); check == -1 {
		fmt.Println("程序运行异常，请携带log日志联系开发者反馈")
		return
	} else if check == 1 {
		newProgram := fmt.Sprintf(".\\%s", newProgramName)
		cmd := exec.Command("cmd.exe", "/K", "start", "", newProgram, "--update", os.Args[0])
		if err := cmd.Start(); err != nil {
			log.Println("启动新程序失败:", err)
			fmt.Println("启动更新程序失败")
			return
		}
		os.Exit(0)
	} else {
		log.Println("完成检查更新过程")
	}
}
