package sso

import (
	"Bilibili_Downloader/pkg/config"
	"Bilibili_Downloader/pkg/httpclient"
	"Bilibili_Downloader/pkg/toolkit"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

// RequestQRCode 请求并获取二维码信息
func RequestQRCode(client *http.Client) (string, string, error) {
	/*resp, err := client.Get("https://passport.bilibili.com/x/passport-login/web/qrcode/generate")
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()*/

	req, err := http.NewRequest("GET", "https://passport.bilibili.com/x/passport-login/web/qrcode/generate", nil)
	if err != nil {
		return "", "", err
	}
	// 设置自定义请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://www.bilibili.com/")
	req.Header.Set("Origin", "https://www.bilibili.com/")
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("未成功关闭响应client:", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var qrResp QRCodeResponse
	err = json.Unmarshal(body, &qrResp)
	if err != nil {
		return "", "", err
	}

	if qrResp.Code != 0 {
		return "", "", fmt.Errorf("申请二维码失败: %s", qrResp.Message)
	}

	return qrResp.Data.QRCodeKey, qrResp.Data.URL, nil
}

// PollQRCodeStatus 轮询二维码状态
func PollQRCodeStatus(client *http.Client, token string) (int, []*http.Cookie, error) {
	qrStatusURL := fmt.Sprintf("https://passport.bilibili.com/x/passport-login/web/qrcode/poll?qrcode_key=%s", url.QueryEscape(token))
	resp, err := client.Get(qrStatusURL)
	if err != nil {
		return 0, nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("未成功关闭请求client:", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	var pollResp PollingResponse
	err = json.Unmarshal(body, &pollResp)
	if err != nil {
		return 0, nil, err
	}

	if pollResp.Code != 0 {
		return pollResp.Data.Code, nil, fmt.Errorf("轮询扫码状态失败: %s", pollResp.Message)
	}

	return pollResp.Data.Code, resp.Cookies(), nil
}

// HandleQRCodeLogin 处理二维码登录流程
func HandleQRCodeLogin() error {
	// 创建 config jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("创建 config jar 失败: %v", err)
	}
	client := &http.Client{Jar: jar}
	httpclient.ChangeClinet(client)

	token, qrURL, err := RequestQRCode(client)
	if err != nil {
		return err
	}

	err = DisplayQRCodeInTerminal(qrURL)
	if err != nil {
		return err
	}

	for {
		status, cookies, err := PollQRCodeStatus(client, token)
		if err != nil {
			return err
		}

		switch status {
		case 86101: // 未扫码
			time.Sleep(2 * time.Second)
			continue
		case 86038: // 二维码超时或失效
			return fmt.Errorf("二维码失效或超时")
		case 86090: // 已扫描未确认
			fmt.Println("二维码已扫描，等待确认")
		case 0: // 登录成功
			toolkit.ClearScreen()
			fmt.Println("登录成功")
			log.Println("登录成功")
			config.StoreCookies(cookies)
			return nil
		default:
			return fmt.Errorf("未知状态码: %d", status)
		}
	}
}
