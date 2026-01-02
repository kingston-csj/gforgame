package protos

type ReqSceneGetData struct {
	_        struct{} `cmd_ref:"CmdSceneReqGetData" type:"req"`
	PlayerId string   `json:"playerId"`
	SceneId  string   `json:"sceneId"`
}

type ReqSceneSetData struct {
	_       struct{} `cmd_ref:"CmdSceneReqSetData" type:"req"`
	SceneId string   `json:"sceneId"`
	Data    string   `json:"data"`
}

type ResSceneGetData struct {
	_    struct{} `cmd_ref:"CmdSceneResGetData" type:"res"`
	Code int32    `json:"code"`
	Data string   `json:"data"`
}

type ResSceneSetData struct {
	_    struct{} `cmd_ref:"CmdSceneResSetData" type:"res"`
	Code int32    `json:"code"`
	Data string   `json:"data"`
}