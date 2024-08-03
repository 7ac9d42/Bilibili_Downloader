package data_struct

type VideoInfoResponse struct {
	Data VideoData `json:"data"`
}

type VideoData struct {
	Bvid       string                 `json:"bvid"`
	Cid        int64                  `json:"cid"`
	Desc       string                 `json:"desc"`        //简介
	Duration   int64                  `json:"duration"`    //时长（所有分P的总时长）
	HonorReply map[string]interface{} `json:"honor_reply"` //获得的荣誉
	Owner      Owner                  `json:"owner"`
	Pic        string                 `json:"pic"`     //封面
	Pubdate    int64                  `json:"pubdate"` //视频发布时间戳
	Stat       Stat                   `json:"stat"`    //视频信息
	Tid        int64                  `json:"tid"`     //分区id
	Title      string                 `json:"title"`   //标题
	Tname      string                 `json:"tname"`   //分类
	Videos     int64                  `json:"videos"`  //分P数
	Pages      []Page                 `json:"pages"`
	UgcSeason  UgcSeason              `json:"ugc_season"`
}

type UgcSeason struct {
	Title    string    `json:"title"`
	Intro    string    `json:"intro"`
	Sections []Section `json:"sections"`
}

type Section struct {
	Episodes []Episode `json:"episodes"`
}

type Episode struct {
	Page  Page   `json:"page"`
	Title string `json:"title"`
	Bvid  string `json:"bvid"`
}

/*
//新版简介，留档以备不时之需
type DescV2 data_struct {
	BizID   *int64  `json:"biz_id,omitempty"`
	RawText *string `json:"raw_text,omitempty"`
	Type    *int64  `json:"type,omitempty"`
}
*/

type Owner struct {
	Face string `json:"face"`
	Mid  int64  `json:"mid"`
	Name string `json:"name"`
}

// Page 分P列表
type Page struct {
	Cid        int64         `json:"cid,omitempty"`       //分P cid
	Dimension  PageDimension `json:"dimension,omitempty"` //当前视频分辨率
	Duration   int64         `json:"duration,omitempty"`  //分P时长
	FirstFrame string        `json:"first_frame,omitempty"`
	From       string        `json:"from,omitempty"`
	Page       int64         `json:"page,omitempty"` //分P序号
	Part       string        `json:"part,omitempty"` //分P标题
	//Vid        string        `json:"vid,omitempty"`
	//Weblink    string        `json:"weblink,omitempty"`
}

type PageDimension struct {
	Height int64 `json:"height"`
	//Rotate int64 `json:"rotate"`
	Width int64 `json:"width"`
}

type Stat struct {
	Aid     int64 `json:"aid"`
	Coin    int64 `json:"coin"`
	Danmaku int64 `json:"danmaku"` //弹幕数
	Dislike int64 `json:"dislike"`
	//Evaluation string `json:"evaluation"`  //评分？？
	Favorite int64 `json:"favorite"`
	//HisRank    int64  `json:"his_rank"`
	Like int64 `json:"like"`
	//NowRank    int64  `json:"now_rank"`
	Reply int64 `json:"reply"`
	Share int64 `json:"share"`
	View  int64 `json:"view"`
	//VT         int64  `json:"vt"`
}
