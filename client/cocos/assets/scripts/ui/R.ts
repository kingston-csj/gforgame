// 资源枚举
export default class R {
  // 预制体枚举
  static readonly Prefabs = {
    AssetResourceFactory: {
      name: 'assetResourceFactory',
      path: 'prefabs/panel/common/AssetResourceFactory',
    },

    LoginPane: {
      name: 'loginPane',
      path: 'prefabs/panel/login/LoginPane',
    },
    MainPane: {
      name: 'mainPane',
      path: 'prefabs/panel/main/MainPane',
    },
    GmPane: {
      name: 'gmPane',
      path: 'prefabs/panel/gm/GmPane',
    },
    BagPane: {
      name: 'bagPane',
      path: 'prefabs/panel/bag/BagPane',
    },
    RecruitPane: {
      name: 'recruitPane',
      path: 'prefabs/panel/recruit/RecruitPane',
    },
    RecruitSettlePane: {
      name: 'recruitSettlePane',
      path: 'prefabs/panel/recruit/RecruitSettlePane',
    },
    TipsPane: {
      name: 'tipsPane',
      path: 'prefabs/panel/common/TipsPane',
    },
    HeroMainPane: {
      name: 'heroMainPane',
      path: 'prefabs/panel/hero/HeroMainPane',
    },
    HeroDetailPane: {
      name: 'heroDetailPane',
      path: 'prefabs/panel/hero/HeroDetailPane',
    },
    FightingUpTipsPane: {
      name: 'fightingUpTipsPane',
      path: 'prefabs/panel/common/FightingUpTipsPane',
    },
    HeroLibPane: {
      name: 'heroLibPane',
      path: 'prefabs/panel/hero/HeroLibPane',
    },
    BuZhenPane: {
      name: 'buZhenPane',
      path: 'prefabs/panel/hero/BuZhenPane',
    },
    ComfirmPane: {
      name: 'comfirmPane',
      path: 'prefabs/panel/common/ComfirmPane',
    },
    MailPane: {
      name: 'mailPane',
      path: 'prefabs/panel/mail/MailPane',
    },
    MailDetailPane: {
      name: 'mailDetailPane',
      path: 'prefabs/panel/mail/MailDetailPane',
    },
    RankPane: {
      name: 'rankPane',
      path: 'prefabs/panel/rank/RankPane',
    },
  };

  // 图片纹理图集枚举
  static readonly Sprites = {
    Hero: 'hero_all',
    Item: 'item_all',
    Quality: 'quality_all',
    Camp: 'camp_all',
  };
}
