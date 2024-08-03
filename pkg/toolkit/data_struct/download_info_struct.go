package data_struct

type SegmentBase struct {
	Initialization string `json:"Initialization"`
	IndexRange     string `json:"indexRange"`
}

type Video struct {
	ID             int         `json:"id,omitempty"`
	BaseUrl        string      `json:"baseUrl,omitempty"`  //下载地址
	BaseURL        string      `json:"base_url,omitempty"` //备用下载地址
	BackupUrl      []string    `json:"backupUrl,omitempty"`
	BackupURL      []string    `json:"backup_url,omitempty"`
	Bandwidth      int         `json:"bandwidth,omitempty"`
	MimeType       string      `json:"mimeType,omitempty"`
	Mime_Type      string      `json:"mime_type,omitempty"`
	Codecs         string      `json:"codecs,omitempty"`
	Width          int         `json:"width,omitempty"`
	Height         int         `json:"height,omitempty"`
	FrameRate      string      `json:"frameRate,omitempty"`
	Frame_Rate     string      `json:"frame_rate,omitempty"`
	Sar            string      `json:"sar,omitempty"`
	StartWithSap   int         `json:"startWithSap,omitempty"`
	Start_With_Sap int         `json:"start_with_sap,omitempty"`
	SegmentBase    SegmentBase `json:"SegmentBase,omitempty"`
	Segment_Base   SegmentBase `json:"segment_base,omitempty"`
	Codecid        int         `json:"codecid,omitempty"`
}

type Audio struct {
	ID             int         `json:"id,omitempty"`
	BaseUrl        string      `json:"baseUrl,omitempty"`
	BaseURL        string      `json:"base_url,omitempty"`
	BackupUrl      []string    `json:"backupUrl,omitempty"`
	BackupURL      []string    `json:"backup_url,omitempty"`
	Bandwidth      int         `json:"bandwidth,omitempty"`
	MimeType       string      `json:"mimeType,omitempty"`
	Mime_Type      string      `json:"mime_type,omitempty"`
	Codecs         string      `json:"codecs,omitempty"`
	Width          int         `json:"width,omitempty"`
	Height         int         `json:"height,omitempty"`
	FrameRate      string      `json:"frameRate,omitempty"`
	Frame_Rate     string      `json:"frame_rate,omitempty"`
	Sar            string      `json:"sar,omitempty"`
	StartWithSap   int         `json:"startWithSap,omitempty"`
	Start_With_Sap int         `json:"start_with_sap,omitempty"`
	SegmentBase    SegmentBase `json:"SegmentBase,omitempty"`
	Segment_Base   SegmentBase `json:"segment_base,omitempty"`
	Codecid        int         `json:"codecid,omitempty"`
}

type Dash struct {
	Duration        int     `json:"duration,omitempty"`
	MinBufferTime   float64 `json:"minBufferTime,omitempty"`
	Min_Buffer_Time float64 `json:"min_buffer_time,omitempty"`
	Video           []Video `json:"video,omitempty"`
	Audio           []Audio `json:"audio,omitempty"`
}

type SupportFormat struct {
	Quality        int      `json:"quality"`
	Format         string   `json:"format"`
	NewDescription string   `json:"new_description"`
	DisplayDesc    string   `json:"display_desc"`
	Superscript    string   `json:"superscript"`
	Codecs         []string `json:"codecs"`
}

type Durl struct {
	//Ahead     string           `json:"ahead,omitempty"`
	BackupURL []string `json:"backup_url,omitempty"`
	Length    int64    `json:"length,omitempty"` //视频长度（单位毫秒）
	Order     int64    `json:"order,omitempty"`  //分段序号
	Size      int64    `json:"size,omitempty"`   //大小（单位Byte）
	URL       string   `json:"url,omitempty"`
	//Vhead     string           `json:"vhead,omitempty"`
}

type Data struct {
	Quality           int      `json:"quality"`
	Format            string   `json:"format"` //文字描述板格式（格式+分辨率:mp4720）
	TimeLength        int      `json:"timelength"`
	AcceptFormat      string   `json:"accept_format"`
	AcceptDescription []string `json:"accept_description"`
	AcceptQuality     []int    `json:"accept_quality"`
	VideoCodecid      int      `json:"video_codecid"`
	//SeekType          string          `json:"seek_type"`   //???
	Dash           Dash            `json:"dash,omitempty"`
	Durl           []Durl          `json:"durl,omitempty"`
	SupportFormats []SupportFormat `json:"support_formats"`
	//LastPlayTime      int             `json:"last_play_time"`
	//LastPlayCid       int             `json:"last_play_cid"`
	//ViewInfo interface{} `json:"view_info"`
}

type DownloadInfoResponse struct {
	Data Data `json:"data"`
}
