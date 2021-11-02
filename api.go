package main

type ResData struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Ttl     int                    `json:"ttl"`
	Data    map[string]interface{} `json:"data"`
}

const JuryUrl = "https://api.bilibili.com/x/credit/v2/jury/jury"

type JuryInfo struct {
	Uname       string `json:"uname" mapstructure:"uname"`
	Face        string `json:"face" mapstructure:"face"`
	CaseTotal   int    `json:"case_total" mapstructure:"case_total"`
	TermEnd     int    `json:"term_end" mapstructure:"term_end"`
	Status      int    `json:"status" mapstructure:"status"`
	ApplyStatus int    `json:"apply_status" mapstructure:"apply_status"`
}

const NextUrl = "https://api.bilibili.com/x/credit/v2/jury/case/next"

const InfoUrl = "https://api.bilibili.com/x/credit/v2/jury/case/info"

type CommentType int

const (
	Comment = 1
	Danmu   = 4
)

func (c CommentType) String() string {
	switch c {
	case Comment:
		return "Comment"
	case Danmu:
		return "Danmu"
	default:
		return "Undefined"
	}
}

type CaseInfo struct {
	CaseId    string      `json:"case_id" mapstructure:"case_id"`
	CaseType  CommentType `json:"case_type" mapstructure:"case_type"`
	VoteItems []struct {
		Vote     int    `json:"vote" mapstructure:"vote"`
		VoteText string `json:"vote_text" mapstructure:"vote_text"`
	} `json:"vote_items" mapstructure:"vote_items"`
	DefaultVote int `json:"default_vote" mapstructure:"default_vote"`
	Status      int `json:"status" mapstructure:"status"`
	OriginStart int `json:"origin_start" mapstructure:"origin_start"`
	Avid        int `json:"avid" mapstructure:"avid"`
	Cid         int `json:"cid" mapstructure:"cid"`
	VoteCd      int `json:"vote_cd" mapstructure:"vote_cd"`
	CaseInfo    struct {
		Comment struct {
			Uname   string `json:"uname" mapstructure:"uname"`
			Face    string `json:"face" mapstructure:"face"`
			Content string `json:"content" mapstructure:"content"`
		} `json:"comment" mapstructure:"comment"`
		DanmuImg string `json:"danmu_img" mapstructure:"danmu_img"`
	} `json:"case_info" mapstructure:"case_info"`
}

const OpinionUrl = "https://api.bilibili.com/x/credit/v2/jury/case/opinion"

type Opinion struct {
	Total int `json:"total" mapstructure:"total"`
	List  []struct {
		Opid       int    `json:"opid" mapstructure:"opid"`
		Mid        int    `json:"mid" mapstructure:"mid"`
		Uname      string `json:"uname" mapstructure:"uname"`
		Face       string `json:"face" mapstructure:"face"`
		Vote       int    `json:"vote" mapstructure:"vote"`
		VoteText   string `json:"vote_text" mapstructure:"vote_text"`
		Content    string `json:"content" mapstructure:"content"`
		Anonymous  int    `json:"anonymous" mapstructure:"anonymous"`
		Like       int    `json:"like" mapstructure:"like"`
		Hate       int    `json:"hate" mapstructure:"hate"`
		LikeStatus int    `json:"like_status" mapstructure:"like_status"`
		VoteTime   int    `json:"vote_time" mapstructure:"vote_time"`
	} `json:"list" mapstructure:"list"`
}

const VoteUrl = "https://api.bilibili.com/x/credit/v2/jury/vote"
