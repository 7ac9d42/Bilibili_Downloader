package sso

import (
	"fmt"
	"github.com/skip2/go-qrcode"
)

// DisplayQRCodeInTerminal 在终端显示二维码
func DisplayQRCodeInTerminal(url string) error {
	// 生成二维码
	qr, err := qrcode.New(url, qrcode.Low)
	if err != nil {
		return fmt.Errorf("生成二维码失败: %w", err)
	}

	// 将二维码转换为 ASCII 字符
	ascii := qr.ToSmallString(false)

	// 输出二维码
	fmt.Println(ascii)
	return nil
}
