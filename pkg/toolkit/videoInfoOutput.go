package toolkit

import (
	"Bilibili_Downloader/pkg/toolkit/data_struct"
	"fmt"
	"time"
)

// 将 Response 结构体转换为控制台输出格式
func responseToConsoleOutput(response *data_struct.VideoInfoResponse) {
	video := response.Data
	fmt.Println("--------------------------------------------------")
	fmt.Printf("视频信息：\n")
	fmt.Printf("标题: %s\n", video.Title)
	//fmt.Printf("封面图片: %s\n", video.Pic)
	fmt.Printf("发布日期: %s\n", time.Unix(video.Pubdate, 0).Format("2006-01-02"))
	fmt.Printf("分类: %s\n", video.Tname)
	fmt.Printf("简介: %s\n", video.Desc)
	fmt.Printf("\n作者: %s\n", video.Owner.Name)
	fmt.Printf("观看数: %d\n", video.Stat.View)
	fmt.Printf("弹幕数: %d\n", video.Stat.Danmaku)
	fmt.Printf("收藏数: %d\n", video.Stat.Favorite)
	fmt.Printf("投币数: %d\n", video.Stat.Coin)
	fmt.Printf("分享数: %d\n", video.Stat.Share)
	fmt.Printf("点赞数: %d\n", video.Stat.Like)
	fmt.Println("--------------------------------------------------")
}

func ConfirmVideoExplanation(response *data_struct.VideoInfoResponse) {
	// 直接在控制台输出数据
	responseToConsoleOutput(response)
}

func PrintDiversityInformationPart1(info *data_struct.VideoInfoResponse) {
	fmt.Println("\n--------------------------------------------------")
	fmt.Println("分    P    列    表")
	fmt.Println("--------------------------------------------------")
	for i := range info.Data.UgcSeason.Sections[0].Episodes {
		fmt.Printf("%d.%s\n", i+1, info.Data.UgcSeason.Sections[0].Episodes[i].Title)
	}
	fmt.Println("--------------------------------------------------")
}

func PrintDiversityInformationPart2(info *data_struct.VideoInfoResponse) {
	fmt.Println("\n--------------------------------------------------")
	fmt.Println("分    P    列    表")
	fmt.Println("--------------------------------------------------")
	for i := range info.Data.Pages {
		fmt.Printf("%d.%s\n", i+1, info.Data.Pages[i].Part)
	}
	fmt.Println("--------------------------------------------------")
}
