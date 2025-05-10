package toolkit

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// CheckAndCreateCacheDir 检查并创建下载缓存目录
func CheckAndCreateCacheDir() error {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %v", err)
	}

	cacheDirPath := filepath.Join(currentDir, "download_cache")

	// 检查目录是否存在
	if _, err := os.Stat(cacheDirPath); os.IsNotExist(err) {
		// 目录不存在，创建目录
		err = os.MkdirAll(cacheDirPath, 0755)
		if err != nil {
			return fmt.Errorf("临时下载目录创建失败: %v", err)
		}
		fmt.Println("成功创建缓存目录:", cacheDirPath)
		log.Println("建立缓存目录正常:", cacheDirPath)
	} else if err != nil {
		// 如果 os.Stat 返回了错误，但不是 os.IsNotExist
		return fmt.Errorf("检查目录时发生错误: %v", err)
	} else {
		// 目录存在
		fmt.Println("缓存目录已经存在，继续使用:", cacheDirPath)
		log.Println("缓存目录已经存在，继续使用:", cacheDirPath)
	}

	return nil
}

func RemoveCacheDir() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %v", err)
	}

	cacheDirPath := filepath.Join(currentDir, "download_cache")

	if err := os.RemoveAll(cacheDirPath); err != nil {
		return fmt.Errorf("移除临时cache目录失败: %v", err)
	}
	log.Println("成功移除cache目录:", cacheDirPath)
	return nil
}

// CheckAndCreateDir 检查并创建指定目录
func CheckAndCreateDir(dir string) error {
	configDir := dir

	// 检查目录是否存在
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// 目录不存在，创建目录
		err := os.Mkdir(configDir, 0755)
		if err != nil {
			return fmt.Errorf("无法创建目录 %s: %v", configDir, err)
		}
		fmt.Println("目录已创建:", configDir)
	} else if err != nil {
		// 其他错误
		return fmt.Errorf("检查目录 %s 时出错: %v", configDir, err)
	} else {
		// 目录已存在
		fmt.Println("目录已存在:", configDir)
	}
	return nil
}

// CheckAndCleanFileName 检查文件名是否包含不允许的字符，并进行清理
func CheckAndCleanFileName(fileName string) string {
	disallowedChars := []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
	// 检查文件名中的每个字符
	for _, char := range disallowedChars {
		if strings.Contains(fileName, char) {
			// 替换不允许的字符为下划线
			fileName = strings.ReplaceAll(fileName, char, "_")
		}
	}
	return fileName
}
