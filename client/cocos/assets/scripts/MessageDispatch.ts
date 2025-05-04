import PlayerData from './data/user/PlayerData';
import { FightingUpTipsView } from './game/common/FightingUpTipsView';
import { HeroBoxModel } from './game/hero/HeroBoxModel';
import BagpackModel from './game/item/BagpackModel';
import { MailBoxModel } from './game/mail/MailBoxModel';
import { PurseModel } from './game/main/PurseModel';
import { HeroVo } from './net/protocol/items/HeroVo';
import PushHeroAttrChanged from './net/protocol/PushHeroAttrChanged';
import { PushItemChanged } from './net/protocol/PushItemChanged';
import { PushPlayerFightChange } from './net/protocol/PushPlayerFightChange';
import PushPurseInfo from './net/protocol/PushPurseInfo';
import { ResAllHeroInfo } from './net/protocol/ResAllHeroInfo';
import ResBackpackInfo from './net/protocol/ResBackpackInfo';
import { PushMailAll } from './net/protocol/ResMailList';
export class MessageDispatch {
  // 绑定cmd与对应的handler
  private static handlers: Map<number, Function> = new Map();

  public static init(): void {
    MessageDispatch.register(ResBackpackInfo.cmd, (msg: ResBackpackInfo) => {
      if (msg.items) {
        BagpackModel.getInstance().reset(new Map(msg.items.map((item) => [item.id, item])));
      }
    });

    MessageDispatch.register(PushPurseInfo.cmd, (msg: PushPurseInfo) => {
      if (msg.diamond) {
        PurseModel.getInstance().diamond = msg.diamond;
      }
      if (msg.gold) {
        PurseModel.getInstance().gold = msg.gold;
      }
    });

    MessageDispatch.register(ResAllHeroInfo.cmd, (msg: ResAllHeroInfo) => {
      if (msg.heros) {
        HeroBoxModel.getInstance().reset(new Map(msg.heros.map((hero) => [hero.id, hero])));
      }
    });

    MessageDispatch.register(PushItemChanged.cmd, (msg: PushItemChanged) => {
      BagpackModel.getInstance().changeItemByModelId(msg.itemId, msg.count);
    });

    MessageDispatch.register(PushHeroAttrChanged.cmd, (msg: PushHeroAttrChanged) => {
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
    });

    MessageDispatch.register(PushPlayerFightChange.cmd, (msg: PushPlayerFightChange) => {
      let from = PlayerData.instance.fighting;
      let add = msg.fight - from;
      if (add > 0 && from > 0) {
        FightingUpTipsView.display(from, add);
      }
      PlayerData.instance.fighting = msg.fight;
    });

    MessageDispatch.register(PushMailAll.cmd, (msg: PushMailAll) => {
      if (msg.mails) {
        MailBoxModel.getInstance().reset(new Map(msg.mails.map((mail) => [mail.id, mail])));
      }
    });
  }

  /**
   * 注册消息处理器
   */
  private static register(cmd: number, handler: Function): void {
    this.handlers.set(cmd, handler);
  }

  public static dispatch(cmd: number, msg: any): void {
    const handler = this.handlers.get(cmd);
    if (handler) {
      handler(msg);
    }
  }
}
