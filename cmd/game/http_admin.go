package main

import (
	"strings"

	"github.com/forfun/gforgame/internal/constants"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	playerrepo "github.com/forfun/gforgame/internal/infra/repository/player"
	"github.com/gin-gonic/gin"
)

type AdminReqQueryPlayers struct {
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize"`
	Sign     string `form:"sign"` // 模糊搜索关键词：角色ID/昵称
}

// PlayerSimpleVo 玩家简易返回VO
type PlayerSimpleVo struct {
	Account string `json:"account"`
	Name    string `json:"name"`
	Level   int32  `json:"level"`
	Diamond int64  `json:"diamond"`
	Gold    int64  `json:"gold"`
}

type GuildSimpleVo struct {
	Id string `json:"id"`
	Name    string `json:"name"`
	Level   int32  `json:"level"`
	// 成员数量
	MemberCount int32 `json:"memberCount"`
	// 贡献值
	Contribution int32 `json:"contribution"`
	// 会长
	Leader string `json:"leader"`
}

// PageResult 分页统一返回结构
type PageResult struct {
	PageNum    int               `json:"pageNum"`
	PageSize   int               `json:"pageSize"`
	TotalSize  int64             `json:"totalSize"`
	TotalPages int               `json:"totalPages"`
}

type PagePlayerResult struct {
	PageResult
	Content    []*PlayerSimpleVo `json:"content"`
}

type PageGuildResult struct {
	PageResult
	Content    []*GuildSimpleVo `json:"content"`
}

func queryPlayers(c *gin.Context, playerRepo *playerrepo.PlayerRepository, profileService *playerrepo.PlayerProfileService) {
	// 1. 绑定GET查询参数
	var req AdminReqQueryPlayers
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(200, AdminHttpResponse{
			Code: 500,
			Msg:  "params parse failed: " + err.Error(),
		})
		return
	}

	// 2. 分页参数边界矫正 等价 coerceAtLeast / coerceAtMost
	pageNum := req.PageNum
	if pageNum < 1 {
		pageNum = 1
	}
	pageSize := req.PageSize
	if pageSize > 50 {
		pageSize = 50
	}
	if pageSize < 1 {
		pageSize = 10 // 兜底最小分页，防止0报错
	}

	// 3. 读取全部缓存档案
	allProfiles := profileService.GetAllPlayerProfiles()

	// 4. 关键词模糊过滤
	var filterList []*playerdomain.PlayerProfile
	if strings.TrimSpace(req.Sign) != "" {
		keyword := strings.ToLower(req.Sign)
		for _, p := range allProfiles {
			playerIdLow := strings.ToLower(p.Id)
			nameLow := strings.ToLower(p.Name)
			if strings.Contains(playerIdLow, keyword) || strings.Contains(nameLow, keyword) {
				filterList = append(filterList, p)
			}
		}
	} else {
		filterList = allProfiles
	}

	// 5. 分页切片截取
	total := int64(len(filterList))
	skip := (pageNum - 1) * pageSize
	// drop skip + take pageSize
	var pageSlice []*playerdomain.PlayerProfile
	if skip < len(filterList) {
		end := skip + pageSize
		if end > len(filterList) {
			end = len(filterList)
		}
		pageSlice = filterList[skip:end]
	}

	// 6. 转换为返回VO
	pageData := make([]*PlayerSimpleVo, 0, len(pageSlice))
	for _, profile := range pageSlice {
		player := playerRepo.GetPlayer(profile.Id)
		backpack := player.Backpack
		vo := &PlayerSimpleVo{
			Account: profile.Id,
			Name:    profile.Name,
			Level:   profile.Level,
			Diamond: int64(backpack.GetItemCount(constants.ITEM_DIAMOND_ID)),
			Gold:    int64(backpack.GetItemCount(constants.ITEM_GOLD_ID)),
		}
		pageData = append(pageData, vo)
	}

	// 7. 计算总页数
	var totalPages int
	if total%int64(pageSize) == 0 {
		totalPages = int(total / int64(pageSize))
	} else {
		totalPages = int(total/int64(pageSize)) + 1
	}

	// 组装分页结果
	pageResult := &PagePlayerResult{
		PageResult: PageResult{
			PageNum:    pageNum,
			PageSize:   pageSize,
			TotalSize:  total,
			TotalPages: totalPages,
		},
		Content:    pageData,
	}

	// 8. gin返回json
	c.JSON(200, Ok(pageResult))
}
