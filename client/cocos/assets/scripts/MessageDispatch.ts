import PlayerData from "./data/user/PlayerData";
import { FightingUpTipsView } from "./game/common/FightingUpTipsView";
import { HeroBoxModel } from "./game/hero/HeroBoxModel";
import BagpackModel from "./game/item/BagpackModel";
import { MailBoxModel } from "./game/mail/MailBoxModel";
import { PurseModel } from "./game/main/PurseModel";
import { HeroVo } from "./net/protocol/items/HeroVo";
import { MailVo } from "./net/protocol/items/MailVo";
import PushHeroAttrChanged from "./net/protocol/PushHeroAttrChanged";
import { PushItemChanged } from "./net/protocol/PushItemChanged";
import { PushPlayerFightChange } from "./net/protocol/PushPlayerFightChange";
import PushPurseInfo from "./net/protocol/PushPurseInfo";
import { ResHeroPushInfo } from "./net/protocol/ResAllHeroInfo";
import PushBackpackInfo from "./net/protocol/PushBackpackInfo";
import { PushMailAll } from "./net/protocol/ResMailList";
import PushQuestRefreshVo from "./net/protocol/PushQuestRefreshVo";
import { QuestModel } from "./game/quest/QuestModel";
import { ConfigContext } from "./data/config/container/ConfigContext";
import GameConstants from "./game/constants/GameConstants";
import eventBus from "./frame/commons/eventbus/EventBus";
import GameEvent from "./game/constants/GameEvent";
import { PushLoadComplete } from "./net/protocol/PushLoadComplete";
import { MainPaneController } from "./game/main/MainPaneController";

// 存储待注册的处理器
const pendingHandlers: Array<{ cmd: number; handler: Function }> = [];

function MessageHandler(cmd: number) {
  return function (
    target: any,
    propertyKey: string,
    descriptor: PropertyDescriptor,
  ) {
    const originalMethod = descriptor.value;
    // 将处理器存储到待注册列表
    pendingHandlers.push({ cmd, handler: originalMethod });
    return descriptor;
  };
}

export class MessageDispatch {
  // 绑定cmd与对应的handler
  private static handlers: Map<number, Function> = new Map();

  public static init(): void {
    // 注册所有待处理的处理器
    pendingHandlers.forEach(({ cmd, handler }) => {
      this.register(cmd, handler);
    });
  }

  @MessageHandler(PushBackpackInfo.cmd)
  private static handleBackpackInfo(msg: PushBackpackInfo) {
    if (msg.items) {
      BagpackModel.getInstance().reset(
        new Map(msg.items.map((item) => [item.uid, item])),
      );
    }
  }

  @MessageHandler(PushPurseInfo.cmd)
  private static handlePurseInfo(msg: PushPurseInfo) {
    if (msg.diamond) {
      PurseModel.getInstance().diamond = msg.diamond;
    }
    if (msg.gold) {
      PurseModel.getInstance().gold = msg.gold;
    }
  }

  @MessageHandler(ResHeroPushInfo.cmd)
  private static handleAllHeroInfo(msg: ResHeroPushInfo) {
    if (msg.heros) {
      HeroBoxModel.getInstance().reset(
        new Map(msg.heros.map((hero) => [hero.id, hero])),
      );
    }
  }

  @MessageHandler(PushItemChanged.cmd)
  private static handleItemChanged(msg: PushItemChanged) {
    if (msg.type == "item") {
      msg.changed.forEach((item) => {
        BagpackModel.getInstance().changeItemByModelId([item]);
      });
    }
  }

  @MessageHandler(PushHeroAttrChanged.cmd)
  private static handleHeroAttrChanged(msg: PushHeroAttrChanged) {
    const hero = HeroBoxModel.getInstance().getHero(msg.heroId);
    if (hero) {
      // 更新英雄属性
      hero.attrs = msg.attrs;
      hero.fight = msg.fight;
      HeroBoxModel.getInstance().addHero(hero);
    } else {
      // 添加英雄
      const hero = new HeroVo();
      hero.id = msg.heroId;
      hero.level = 1;
      hero.fight = msg.fight;
      hero.stage = 0;
      hero.attrs = msg.attrs;
      HeroBoxModel.getInstance().addHero(hero);
    }
  }

  @MessageHandler(PushPlayerFightChange.cmd)
  private static handlePlayerFightChange(msg: PushPlayerFightChange) {
    let from = PlayerData.instance.fighting;
    let add = msg.fight - from;
    if (add > 0 && from > 0) {
      FightingUpTipsView.display(from, add);
    }
    PlayerData.instance.fighting = msg.fight;
  }

  @MessageHandler(PushMailAll.cmd)
  private static handleMailAll(msg: PushMailAll) {
    if (msg.mails) {
      MailBoxModel.getInstance().reset(
        new Map(
          msg.mails.map((mail) => {
            //  转成MailVo实例
            const mailVo = new MailVo();
            Object.assign(mailVo, mail);
            return [mail.id, mailVo];
          }),
        ),
      );
    }
  }

  @MessageHandler(PushQuestRefreshVo.cmd)
  private static handleQuestRefresh(msg: PushQuestRefreshVo) {
    if (msg.quest) {
      const quest = msg.quest;
      // 更新任务
      QuestModel.instance.refreshQuest(quest);
      let questData = ConfigContext.configQuestContainer.getRecord(quest.id);
      if (questData.category == GameConstants.Quest.Category.MAIN) {
        eventBus.emit(GameEvent.MainQuestRefresh);
      }
    }
  }

  @MessageHandler(PushLoadComplete.cmd)
  private static handleLoadComplete(msg: PushLoadComplete) {
    // 加载完成
    MainPaneController.openUi();
    // eventBus.emit(GameEvent.LoadComplete);
  }

  /**
   * 注册消息处理器
   */
  public static register(cmd: number, handler: Function): void {
    this.handlers.set(cmd, handler);
  }

  /**
   * 分发消息
   */
  public static dispatch(cmd: number, msg: any): void {
    const handler = this.handlers.get(cmd);
    if (handler) {
      handler(msg);
    }
  }
}
