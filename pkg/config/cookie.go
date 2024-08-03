package config

import (
	"Bilibili_Downloader/pkg/toolkit"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// StoreCookies 存储Cookie
func StoreCookies(cookies []*http.Cookie) {
	if err := toolkit.CheckAndCreateDir("./config"); err != nil {
		log.Println("视频输出目录检查或创建失败：", err)
	}
	// 创建文件存储 cookies
	file, err := os.Create("./config/cookies.json")
	if err != nil {
		fmt.Println("创建Cookie文件失败:", err)
		log.Println("创建Cookie文件失败:", err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("Close Cookie file失败：", err)
		}
	}()

	// 将 cookies 转换为 JSON 格式
	cookiesJSON, err := json.MarshalIndent(cookies, "", "  ")
	if err != nil {
		fmt.Println("转换 cookies 到 JSON 失败:", err)
		log.Println("转换 cookies 到 JSON 失败:", err)
		return
	}

	// 将 JSON 写入文件
	if err := os.WriteFile(file.Name(), cookiesJSON, 0644); err != nil {
		fmt.Println("写入 cookies 文件失败:", err)
		log.Println("写入 cookies 文件失败:", err)
		return
	}

	fmt.Println("Cookies 已保存到:", file.Name())
	log.Println("Cookies 已保存到:", file.Name())
}

// LoadCookies 加载之前保存的 cookies
func LoadCookies() []*http.Cookie {
	// 读取之前保存的 cookies 文件
	content, err := os.ReadFile("./config/cookies.json")
	if err != nil {
		fmt.Println("未成功加载已保存的配置文件:", err)
		log.Println("未成功加载已保存的配置文件:", err)
		return nil
	}

	var cookies []*http.Cookie
	err = json.Unmarshal(content, &cookies)
	if err != nil {
		fmt.Println("解析 cookies 失败:", err)
		log.Println("解析 cookies 失败:", err)
		return nil
	}
	fmt.Println("Cookie加载成功，前十个字符为：", cookies[0].Value[:10])
	log.Println("Cookie加载成功")

	return cookies
}
