package constants

// 任务状态枚举
const (
	// 状态--未完成
	QuestStatusInit = int8(0)
	// 状态--已完成未领奖
	QuestStatusFinished = int8(1)
	// 状态--已领奖
	QuestStatusRewarded = int8(2)
)

// 任务分类枚举
const (

	/**
	 * 任务分类-主线
	 */
	QuestCategoryMain = int32(1)
	/**
	 * 任务分类-日常
	 */
	QuestCategoryDaily = int32(2)
	/**
	 * 任务分类-通行证
	 */
	QuestCategoryPass = int32(3)
	/**
	 * 任务分类-修炼
	 */
	QuestCategoryTrain = int32(4)

	/**
	 * 任务分类-每周
	 */
	QuestCategoryWeekly = int32(5)
	/**
	 * 任务分类-公会
	 */
	QuestCategoryGuild = int32(6)

	/**
	 * 任务分类-成就
	 */
	QuestCategoryAchievement = int32(9)
)

// 任务类型枚举
const (
	/**
	 * 招募
	 */
	QuestTypeRecruit = 1

	/**
	 * 英雄升级
	 */
	QuestTypeHeroUpLevel = 2

	/**
	 * 英雄升级
	 */
	QuestTypeHeroUpStage = 3

	/**
	 * 掌门升级
	 */
	QuestTypeMasterUpLevel = 4

	/**
	 * 装备升级
	 */
	QuestTypeEquipUpLevel = 5

	/**
	 * 闯关(主线)
	 */
	QuestTypePassGuanka = 6

	/**
	 * 通过XX波
	 */
	QuestTypePassGuankaRound = 7

	/**
	 * 集市购买
	 */
	QuestTypeMallBuy = 8

	/**
	 * 消耗金币
	 */
	QuestTypeGoldConsume = 9

	/**
	 * 消耗钻石
	 */
	QuestTypeDiamondConsume = 10

	/**
	 * 登录
	 */
	QuestTypeLogin = 11

	/**
	 * 开宝箱
	 */
	QuestTypeOpenBox = 12

	/**
	 * 挑战副本
	 */
	QuestTypeFuben = 13

	/**
	 * 竞技场
	 */
	QuestTypeArena = 14

	/**
	 * 英雄总量
	 */
	QuestTypeHeroSum = 15

	/**
	 * 挂机结算
	 */
	QuestTypeIdleSettle = 16


	/**
	 * 客户端事件
	 */
	QuestTypeClientEvent = 23
)