package main

import (
	"Bilibili_Downloader/internal"
	"Bilibili_Downloader/internal/sso"
	"Bilibili_Downloader/internal/video_processing"
	"Bilibili_Downloader/pkg/httpclient"
	"Bilibili_Downloader/pkg/toolkit"
	"Bilibili_Downloader/pkg/toolkit/data_struct"
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func CatchVideoInfo(VideoInfoUrl string) (*data_struct.VideoInfoResponse, error) {
	//获取视频信息
	data, err := internal.CatchData(VideoInfoUrl)
	if err != nil {
		log.Printf("获取视频信息数据错误: %v\n\n", err)
		fmt.Println("视频信息数据获取异常，请检查网络连接或前往log文件查看详情.")
		return nil, err
	}

	//视频信息反序列化
	Response, err := internal.ProcessResponse(data, 0)
	if err != nil {
		log.Printf("处理视频详情发生错误: %v\n\n", err)
		fmt.Println("视频详情数据处理发生错误，请携带log文件向开发者反馈！")
		return nil, err
	}
	videoInfoResponse, ok := Response.(*data_struct.VideoInfoResponse)
	if !ok {
		log.Println("视频详情数据类型断言失败")
		fmt.Println("程序运行发生异常，请携带log日志文件联系开发者！")
		return nil, fmt.Errorf("视频详情数据类型断言失败")
	}
	return videoInfoResponse, nil
}

func main() {
	// 创建一个新的读取器
	reader := bufio.NewReader(os.Stdin)

	logFile := internal.InitLog()
	defer func() {
		if err := logFile.Close(); err != nil {
			log.Println("log文件close失败:", err)
		}
	}()

	defer func() {
		fmt.Printf("程序执行完毕，请按Enter键退出...")
		_, _ = reader.ReadString('\n')
		log.Println("程序执行完毕，正常退出")
	}()

	//提交清理计划
	defer internal.CacheClean()

	//初始化客户端
	if !httpclient.Init() {
		err := sso.HandleQRCodeLogin()
		if err != nil {
			fmt.Println("处理二维码登录失败:", err)
			return
		}
	}

	internal.HandleUpdate()

	//主逻辑循环
	for {
		//获取用户BV号输入并检查
		BVid := toolkit.CatchAndCheckBVid()
		VideoInfoUrl := fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?bvid=%s", BVid)

		videoInfoResponse, err := CatchVideoInfo(VideoInfoUrl)
		if err != nil {
			break
		}
		//打印视频详细信息并进行确认
		toolkit.ClearScreen()
		toolkit.ConfirmVideoExplanation(videoInfoResponse)

		Default := -2
		var actions map[string]map[string]int64
		if videoInfoResponse.Data.UgcSeason.Sections != nil || len(videoInfoResponse.Data.Pages) > 1 {
			var part int64
			//确认是否使用多分P选择/连续下载功能
			fmt.Printf("检测到视频含有分P，是否使用多分P选择/连续下载功能？(y/n):")
			if toolkit.YesOrNo() {
				Default = 0
				if len(videoInfoResponse.Data.Pages) > 1 {
					toolkit.PrintDiversityInformationPart2(videoInfoResponse)
					part = 2
				} else {
					toolkit.PrintDiversityInformationPart1(videoInfoResponse)
					part = 1
				}
				if actions = toolkit.GetMaps(videoInfoResponse, part); actions == nil {
					return
				}
			} else {
				actions = map[string]map[string]int64{
					videoInfoResponse.Data.Title: {
						videoInfoResponse.Data.Bvid: videoInfoResponse.Data.Cid,
					},
				}
			}
			//TODO:缺少在分P列表中对当前视频的突出显示
		} else {
			actions = map[string]map[string]int64{
				videoInfoResponse.Data.Title: {
					videoInfoResponse.Data.Bvid: videoInfoResponse.Data.Cid,
				},
			}
		}

		for title, video := range actions {
			for bvid, cid := range video {
				//请求视频取流地址
				DownloadURL := fmt.Sprintf("https://api.bilibili.com/x/player/wbi/playurl?bvid=%s&cid=%d&fnval=4048", bvid, cid)
				data, err := internal.CatchData(DownloadURL)
				if err != nil {
					log.Printf("获取下载信息数据发生错误: %v\n\n", err)
					fmt.Println("视频下载信息获取异常，请检查网络连接或前往log文件查看详情.")
					return
				}

				//反序列化视频流信息
				newResponse, err := internal.ProcessResponse(data, 1)
				if err != nil {
					log.Printf("处理下载信息发生错误: %v\n\n", err)
					fmt.Println("视频下载信息处理发生错误，请携带log文件向开发者反馈！")
					return
				}
				downloadInfoResponse, ok := newResponse.(*data_struct.DownloadInfoResponse)
				if !ok {
					log.Println("视频下载数据类型断言失败")
					fmt.Println("程序运行发生异常，请携带log日志文件联系开发者！")
					break
				}

				//处理用户选择
				if len(downloadInfoResponse.Data.AcceptDescription) == 0 || downloadInfoResponse.Data.AcceptDescription[0] == `试看` {
					log.Println("当前账号可能没有观看（下载）该视频的权限，无法获取视频下载地址.")
					fmt.Println("\n当前账号可能没有观看（下载）该视频的权限，无法获取视频下载地址.")
				} else {
					if Default == 0 {
						fmt.Printf("当前为多分P下载模式，是否默认下载可获取的最高分辨率?(y/n)；")
						if toolkit.YesOrNo() {
							Default = 1
						} else {
							Default = -1
						}
					}
					videoIndex, _, resolutionDescription := toolkit.ObtainUserResolutionSelection(int64(Default), title, downloadInfoResponse)

					//请求视频下载
					if err := internal.DownloadFile(downloadInfoResponse.Data.Dash.Video[videoIndex].BackupURL[0], downloadInfoResponse.Data.Dash.Audio[0].BackupURL[0], ""); err != nil {
						log.Printf("请求下载失败：%s\n", err)
						fmt.Println("请求下载失败，请检查网络连接或前往log文件查看详情.")
					}

					//视频音频混流转码
					fmt.Printf("开始视频转码：\n\n")
					video_processing.Transcoding(title, resolutionDescription)
					time.Sleep(500 * time.Millisecond)
				}
			}
		}

		fmt.Printf("是否继续下载其他视频？(y/n):")
		if isContinue := toolkit.YesOrNo(); isContinue {
			toolkit.ClearScreen()
			continue
		} else {
			break
		}
	}
}
