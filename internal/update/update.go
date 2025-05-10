package update

import (
	"Bilibili_Downloader/pkg/config"
	"Bilibili_Downloader/pkg/toolkit"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/cheggaaa/pb/v3"
	"log"
	"net/http"
	"os"
	"strings"
)

// checkForUpdate 检查更新
func checkForUpdate(VersionUpdateApi string, currentVersion string) (string, string, error) {
	resp, err := http.Get(VersionUpdateApi)
	if err != nil {
		log.Printf("failed to get latest release: %s\n", err)
		fmt.Printf("检查更新失败，请检查网络环境\n\n")
		return "?", "", nil
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("close body failed")
		}
	}()

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("failed to get latest release: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	latestVersionStr := result["tag_name"].(string)
	downloadURL := result["assets"].([]interface{})[0].(map[string]interface{})["browser_download_url"].(string)

	// 使用 semver 库比较版本
	currentSemver, err := semver.NewVersion(strings.TrimPrefix(currentVersion, "v"))
	if err != nil {
		return "", "", err
	}
	latestSemver, err := semver.NewVersion(strings.TrimPrefix(latestVersionStr, "v"))
	if err != nil {
		return "", "", err
	}

	if latestSemver.GreaterThan(currentSemver) {
		return downloadURL, latestVersionStr, nil // 有更新
	}
	return "", "", nil // 没有更新
}

// 下载更新
func downloadUpdate(downloadURL string, downloadVersion string) (string, error) {
	resp, err := http.Get(downloadURL)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("close body failed")
		}
	}()

	tempFile, err := os.Create("update_temp")
	if err != nil {
		return "", err
	}

	bar := pb.StartNew(int(resp.ContentLength))
	bar.Set(pb.SIBytesPrefix, true)
	err = toolkit.DownloadAndTrackProgress(resp.Body, tempFile, bar)
	if err != nil {
		return "", err
	}
	bar.Finish()

	if err := tempFile.Close(); err != nil {
		log.Println("close tempFile failed")
		return "", err
	}

	newProgramName := "BiliBili_Downloader_" + downloadVersion + ".exe"
	err = os.Rename("update_temp", newProgramName)
	if err != nil {
		return "", err
	}

	return newProgramName, nil
}

// 处理更新替换
func update(oldVersionPath string) {
	if err := os.Remove(oldVersionPath); err != nil {
		log.Println("删除旧版本失败:", err)
		return
	}
	log.Println("成功删除旧版本，完成程序更新")
	fmt.Println("完成程序更新")
}

// CheckAndUpdate 返回值 -1：err；0; 正常无更新或跳过更新; 1：正常且更新
func CheckAndUpdate() (int, string) {
	if len(os.Args) > 1 && os.Args[1] == "--update" {
		// 被老版本启动用于更新
		if len(os.Args) < 3 {
			log.Println("更新参数缺失")
			fmt.Println("无法正确删除旧版本，可手动删除旧版本")
		}
		update(os.Args[2])
	}

	fmt.Println("读取配置...")
	currentConfig, err := config.ReadConfig()
	if err != nil {
		log.Println("读取配置时发生错误:", err)
		fmt.Println("读取配置时发生错误:", err)
		return -1, ""
	}
	if currentConfig.CurrentVersion != config.CurrentVersion {
		log.Println("内置版本信息与config配置不一致，尝试重建config配置")
		if err := os.Remove(".\\config\\config.json"); err != nil {
			log.Println("删除旧config失败")
		}
		if err := config.CreateConfig(); err != nil {
			log.Println("重建config失败：", err)
			fmt.Println("配置异常，请尝试删除config目录后重启！")
		}
		currentConfig.CurrentVersion = config.CurrentVersion
	}

	fmt.Println("检查更新...")
	downloadURL, latestVersion, err := checkForUpdate(currentConfig.VersionUpdateApi, currentConfig.CurrentVersion)
	if err != nil {
		log.Println("检查更新时发生错误:", err)
		fmt.Println("检查更新时发生错误:", err)
		return -1, ""
	}

	if downloadURL == "" {
		fmt.Printf("当前已是最新版本！\n\n")
		return 0, ""
	} else if downloadURL == "?" {
		return 0, ""
	} else {
		fmt.Printf("发现新版本，当前版本: %s，最新版本: %s\n", currentConfig.CurrentVersion, latestVersion)
		fmt.Print("是否下载更新? (Y/n): ")
		if toolkit.YesOrNo() {
			if newProgramName, err := downloadUpdate(downloadURL, latestVersion); err != nil {
				log.Println("更新失败:", err)
				fmt.Println("更新失败:", err)
			} else {
				if err := os.Remove(".\\config\\config.json"); err != nil {
					log.Println("删除旧config失败")
				}
				fmt.Println("更新成功，请重新启动程序！")
				return 1, newProgramName
			}
		} else {
			fmt.Printf("跳过更新\n\n")
			return 0, ""
		}
	}
	return -1, ""
}
