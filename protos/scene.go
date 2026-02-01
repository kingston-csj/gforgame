package protos

type ReqSceneGetData struct {
	_        struct{} `cmd_ref:"CmdSceneReqGetData"`
	PlayerId string   `json:"playerId"`
	SceneId  string   `json:"sceneId"`
}

type ReqSceneSetData struct {
	_       struct{} `cmd_ref:"CmdSceneReqSetData"`
	SceneId string   `json:"sceneId"`
	Data    string   `json:"data"`
}

type ResSceneGetData struct {
	_    struct{} `cmd_ref:"CmdSceneResGetData"`
	Code int32    `json:"code"`
	Data string   `json:"data"`
}

type ResSceneSetData struct {
	_    struct{} `cmd_ref:"CmdSceneResSetData"`
	Code int32    `json:"code"`
	Data string   `json:"data"`
}