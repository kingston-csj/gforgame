package protos

type RankInfo struct {
	Id          string `json:"id"`
	Order       int    `json:"order"`
	Value       int64  `json:"value"`
	Name        string `json:"name"`
	SecondValue int64  `json:"secondValue"`
	ExtraInfo   string `json:"extraInfo"`
}

type ReqRankQuery struct {
	_        struct{} `cmd_ref:"CmdRankReqQuery" type:"req"`
	Type     int      `json:"type"`
	Start    int      `json:"start"`
	PageSize int      `json:"pageSize"`
}

type ResRankQuery struct {
	_        struct{}   `cmd_ref:"CmdRankResQuery" type:"res"`
	Type     int        `json:"type"`
	Records  []RankInfo `json:"records"`
	MyRecord RankInfo   `json:"myRecord"`
}
