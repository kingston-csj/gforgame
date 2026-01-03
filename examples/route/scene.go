package route

import (
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/service/scene"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)



type SceneRoute struct {
	network.Base
}

func NewSceneRoute() *SceneRoute {
	return &SceneRoute{}
}

func (ps *SceneRoute) ReqQuery(s *network.Session, index int32, msg *protos.ReqSceneGetData) *protos.ResSceneGetData {
	playerId := msg.PlayerId
	sceneId := msg.SceneId
	scene := scene.GetSceneService().GetOrCreateScene(playerId, sceneId)
	return &protos.ResSceneGetData{
		Code: 0,
		Data: scene.Data}
}

func (ps *SceneRoute) ReqUpdate(s *network.Session, index int32, msg *protos.ReqSceneSetData) *protos.ResSceneSetData {
	player := network.GetPlayerBySession(s).(*playerdomain.Player)
	sceneId := msg.SceneId
	  scene.GetSceneService().UpdateScene(player.Id, sceneId, msg.Data)
	return &protos.ResSceneSetData{
		Code: 0}
}