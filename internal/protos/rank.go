package protos

type RankInfo struct {
	Id          string `json:"id"`
	Order       int32  `json:"order"`
	Value       int64  `json:"value"`
	Name        string `json:"name"`
	SecondValue int64  `json:"secondValue"`
	ExtraInfo   string `json:"extraInfo"`
}

type ReqRankQuery struct {
	_        struct{} `cmd_ref:"CmdRankReqQuery" type:"req"`
	Type     int32    `json:"type"`
	Start    int32    `json:"start"`
	PageSize int32    `json:"pageSize"`
}

type ResRankQuery struct {
	_        struct{}   `cmd_ref:"CmdRankResQuery" type:"res"`
	Type     int32      `json:"type"`
	Records  []RankInfo `json:"records"`
	MyRecord RankInfo   `json:"myRecord"`
}
