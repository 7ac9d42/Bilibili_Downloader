package config

import (
	"encoding/json"
	"log"
	"os"
)

// Config 结构体用于存储配置信息
type Config struct {
	VersionUpdateApi string `json:"VersionUpdateApi"`
	CurrentVersion   string `json:"CurrentVersion"`
}

// CurrentVersion 当前版本
const CurrentVersion = `v1.4.1`

// 配置文件名
const configFileName = `.\config\config.json`

// CreateConfig 创建默认配置文件
func CreateConfig() error {
	config := Config{
		VersionUpdateApi: "https://api.github.com/repos/KUAPT/Bilibili_Downloader/releases/latest",
		CurrentVersion:   CurrentVersion,
	}

	file, err := os.Create(configFileName)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("关闭文件失败：", err)
		}
	}()

	return json.NewEncoder(file).Encode(config)
}

// ReadConfig 读取配置文件
func ReadConfig() (Config, error) {
	var config Config
	file, err := os.Open(configFileName)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果配置文件不存在，则创建一个默认的配置文件
			if err := CreateConfig(); err != nil {
				return config, err
			}
			// 再次尝试读取
			file, err = os.Open(configFileName)
			if err != nil {
				return config, err
			}
		} else {
			return config, err
		}
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("关闭文件失败：", err)
		}
	}()

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}
