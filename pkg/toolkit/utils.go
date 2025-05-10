package toolkit

import (
	"Bilibili_Downloader/pkg/toolkit/data_struct"
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// 全局变量保存多分P下载过程中是否忽略HDR的选择
var globalIgnoreHDRSet = false
var globalIgnoreHDR = false

// ResetHDRState 重置HDR忽略状态
func ResetHDRState() {
	globalIgnoreHDRSet = false
	globalIgnoreHDR = false
}

func CatchAndCheckBVid() string {
	//正则对BV号进行基本检查
	BECheck := regexp.MustCompile(`^BV[1-9A-HJ-NP-Za-km-z]{10}$`)
	var BVid string
	for {
		fmt.Printf("请输入需要下载视频的BV号：")
		if _, err := fmt.Scanln(&BVid); err != nil {
			ClearScreen()
			fmt.Println("输入读取错误，请重试！")
			log.Println("读取输入错误：", err)
			continue
		}
		if BECheck.MatchString(BVid) {
			break
		} else {
			ClearScreen()
			fmt.Println("BV号格式错误，请检查格式后重试！")
		}
	}
	return BVid
}

// 抽取HDR检测与逻辑处理函数
func handleHDRIfNeeded(Default int64, choose int, effectiveDefinition []int, resolutions map[int]string) int {
	videoCode := effectiveDefinition[choose]
	resolutionDescription := resolutions[videoCode]
	if Default == 1 && resolutionDescription == "HDR" {
		// 如果已经记录过用户选择，直接复用
		if globalIgnoreHDRSet {
			if globalIgnoreHDR {
				for i, code := range effectiveDefinition {
					if resolutions[code] != "HDR" {
						choose = i
						return choose
					}
				}
				fmt.Println("未找到其他非HDR画质选项，继续使用HDR画质")
				return choose
			} else {
				return choose
			}
		}
		// 未记录时询问用户，并记忆其选择
		fmt.Printf("\n--------【请注意：HDR画质在不受支持的设备上播放将会有显著偏色现象】--------\n")
		fmt.Printf("是否忽略HDR画质?(Y/n):")
		if YesOrNo() {
			globalIgnoreHDR = true
			// 查找其他非HDR选项
			for i, code := range effectiveDefinition {
				if resolutions[code] != "HDR" {
					choose = i
					globalIgnoreHDRSet = true
					return choose
				}
			}
			fmt.Println("未找到其他非HDR画质选项，继续使用HDR画质")
		} else {
			globalIgnoreHDR = false
		}
		globalIgnoreHDRSet = true
	}
	return choose
}

// ObtainUserResolutionSelection 获取用户分辨率选择
func ObtainUserResolutionSelection(Default int64, title string, downloadInfoResponse *data_struct.DownloadInfoResponse) (int, int, string) {
	definition := make(map[int]string, 10)
	for i := 0; i < len(downloadInfoResponse.Data.AcceptDescription); i++ {
		definition[downloadInfoResponse.Data.AcceptQuality[i]] = downloadInfoResponse.Data.AcceptDescription[i]
	}

	effectiveDefinition := make([]int, 0, 10)
	effectiveDefinitionMap := make(map[int]bool)
	for i := range downloadInfoResponse.Data.Dash.Video {
		id := downloadInfoResponse.Data.Dash.Video[i].ID
		if !effectiveDefinitionMap[id] {
			effectiveDefinition = append(effectiveDefinition, id)
			effectiveDefinitionMap[id] = true
		}
	}

	var choose int
	if Default != 1 {
		for {
			fmt.Println("\n当前下载的视频为：", title)
			fmt.Println("\n请选择想要下载的分辨率：(ps:此处仅显示当前登录账号有权获取的所有分辨率选项)")
			for i := range effectiveDefinition {
				fmt.Println(i+1, definition[effectiveDefinition[i]])
			}
			fmt.Printf("请输入分辨率前的序号(单个数字)：")
			if _, err := fmt.Scanln(&choose); err != nil {
				ClearScreen()
				log.Println("读取输入发生错误")
				fmt.Println("读取输入发生错误,请检查输入格式后重试，若问题依旧，请携带日志log文件向开发者反馈！")
				continue
			}
			if choose < 1 || choose > len(effectiveDefinition) {
				ClearScreen()
				fmt.Println("输入错误，请检查输入后重试！")
				continue
			}
			choose -= 1
			break
		}
	} else {
		fmt.Println("当前下载的视频为：", title)
		choose = 0
	}

	resolutions := map[int]string{
		6:   "240P",
		16:  "360P",
		32:  "480P",
		64:  "720P",
		74:  "720P60",
		80:  "1080P",
		112: "1080P+",
		116: "1080P60",
		120: "4K",
		125: "HDR",
		126: "杜比视界",
		127: "8K超高清",
	}

	choose = handleHDRIfNeeded(Default, choose, effectiveDefinition, resolutions)

	// 根据最终选择更新 videoIndex、videoCode 及分辨率描述
	var videoIndex int
	for i := range downloadInfoResponse.Data.Dash.Video {
		if downloadInfoResponse.Data.Dash.Video[i].ID == effectiveDefinition[choose] {
			videoIndex = i
			break
		}
	}
	videoCode := effectiveDefinition[choose]
	resolutionDescription := resolutions[videoCode]

	return videoIndex, videoCode, resolutionDescription
}

func GetMaps(info *data_struct.VideoInfoResponse, kind int64) map[string]map[string]int64 {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Tip：若需下载所有分P，可直接输入序号0")
	fmt.Printf("请输入用逗号分隔的整数序号（支持中英逗号）：")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	cleanedInput := strings.ReplaceAll(input, "，", ",")

	parts := strings.Split(cleanedInput, ",")

	if kind == 1 {
		outerMap := make(map[string]map[string]int64, len(info.Data.UgcSeason.Sections[0].Episodes))

		for _, part := range parts {
			// 移除可能的空格
			part = strings.TrimSpace(part)
			if num, err := strconv.Atoi(part); err == nil && num == 0 {
				outerMap = make(map[string]map[string]int64, len(info.Data.UgcSeason.Sections[0].Episodes))
				for i := 0; i < len(info.Data.UgcSeason.Sections[0].Episodes); i++ {
					innerMap := make(map[string]int64, 1)
					innerMap[info.Data.UgcSeason.Sections[0].Episodes[i].Bvid] = info.Data.UgcSeason.Sections[0].Episodes[i].Page.Cid
					outerMap[info.Data.UgcSeason.Sections[0].Episodes[i].Title] = innerMap
				}
			} else if num, err := strconv.Atoi(part); err == nil {
				innerMap := make(map[string]int64, 1)
				innerMap[info.Data.UgcSeason.Sections[0].Episodes[num-1].Bvid] = info.Data.UgcSeason.Sections[0].Episodes[num-1].Page.Cid
				outerMap[info.Data.UgcSeason.Sections[0].Episodes[num-1].Title] = innerMap
			} else {
				fmt.Printf("跳过无效输入： %s\n", part)
			}
		}
		return outerMap

	} else if kind == 2 {
		outerMap := make(map[string]map[string]int64, len(info.Data.Pages))

		for _, part := range parts {
			// 移除可能的空格
			part = strings.TrimSpace(part)
			if num, err := strconv.Atoi(part); err == nil && num == 0 {
				outerMap = make(map[string]map[string]int64, len(info.Data.Pages))
				for i := 0; i < len(info.Data.Pages); i++ {
					innerMap := make(map[string]int64, 1)
					innerMap[info.Data.Bvid] = info.Data.Pages[i].Cid
					outerMap[info.Data.Pages[i].Part] = innerMap
				}
			} else if num, err := strconv.Atoi(part); err == nil {
				innerMap := make(map[string]int64, 1)
				innerMap[info.Data.Bvid] = info.Data.Pages[num-1].Cid
				outerMap[info.Data.Pages[num-1].Part] = innerMap
			} else {
				fmt.Printf("跳过无效输入： %s\n", part)
			}
		}
		return outerMap
	} else {
		fmt.Println("程序运行异常，请携带日志log文件联系开发者！")
		log.Println("打印分P信息异常，不支持的kind值")
	}
	return nil
}
