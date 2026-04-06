export default class GameConstants {
  static readonly Item = {
    // 招募令道具id
    RECRUIT_ID: 2001,
    // 英雄突破材料
    UpStage: 2003,
  };

  static readonly Quest = {
    // 任务类型
    Category: {
      // 主线任务
      MAIN: 1,

      // 日常任务
      DAILY: 2,
    },
  };

  static readonly I18N = {
    ILLEGAL_PARAMS: 1001,
    INTERNAL_ERROR: 1002,
    NOT_FOUND: 1003,
    CONTENT_NOT_ENOUGH: 1004,
    ITEM_NOT_ENOUGH: 2001,
    GOLD_NOT_ENOUGH: 2002,
    DIAMOND_NOT_ENOUGH: 2003,
  };
}
