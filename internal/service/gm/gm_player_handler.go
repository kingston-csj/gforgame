package gm

import (
	"fmt"

	commonerrors "github.com/forfun/gforgame/common/errors"
	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/common/util/conv"
	"github.com/forfun/gforgame/common/util/jsonutil"
	"github.com/forfun/gforgame/internal/constants"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	mysqldb "github.com/forfun/gforgame/internal/infra/persistence"
	playerservice "github.com/forfun/gforgame/internal/service/player"
)

type PlayerGmHandler struct {
	player *playerservice.PlayerService
}

func NewPlayerGmHandler(playerService *playerservice.PlayerService) *PlayerGmHandler {
	return &PlayerGmHandler{
		player: playerService,
	}
}

func (h *PlayerGmHandler) RegisterTo(gm *GmService) {
	gm.Register("reset", "重置玩家数据", "reset", h.handleReset)
	gm.Register("level", "修改等级", "level 100", h.handleLevel)
	gm.Register("clone", "克隆玩家", "clone 1001", h.handleClone)
}

func (h *PlayerGmHandler) handleReset(player *playerdomain.Player, params string) *commonerrors.BusinessError {
	player.Reset()
	return nil
}

func (h *PlayerGmHandler) handleLevel(player *playerdomain.Player, params string) *commonerrors.BusinessError {
	player.Level = conv.Int32Value(params)
	h.player.GetPlayerProfileById(player.Id)
	return nil
}

func (h *PlayerGmHandler) handleClone(player *playerdomain.Player, params string) *commonerrors.BusinessError {
	targetId := params
	if player.Id == targetId {
		return commonerrors.NewBusinessError(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}

	var to playerdomain.Player
	json, err := jsonutil.StructToJSON(player)
	if err != nil {
		return commonerrors.NewBusinessError(constants.I18N_COMMON_INTERNAL_ERROR)
	}
	err = jsonutil.JsonToStruct(json, &to)
	if err != nil {
		return commonerrors.NewBusinessError(constants.I18N_COMMON_INTERNAL_ERROR)
	}
	to.Id = targetId
	to.Name = h.player.RandomName()
	h.player.SavePlayer(&to)

	var scenes []playerdomain.Scene
	err = mysqldb.Db.Where(fmt.Sprintf("id like '%s%%'", player.Id)).Find(&scenes).Error
	if err != nil {
		logger.Error("gm reset scene fail", err)
		return commonerrors.NewBusinessError(constants.I18N_COMMON_INTERNAL_ERROR)
	}
	return nil
}
