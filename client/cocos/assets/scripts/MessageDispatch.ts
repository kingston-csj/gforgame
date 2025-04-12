import { HeroBoxModel } from './game/hero/HeroBoxModel';
import BagpackModel from './game/item/BagpackModel';
import { PurseModel } from './game/main/PurseModel';
import { HeroVo } from './net/MsgItems/HeroVo';
import PushHeroAttrChanged from './net/PushHeroAttrChanged';
import { PushItemChanged } from './net/PushItemChanged';
import PushPurseInfo from './net/PushPurseInfo';
import { ResAllHeroInfo } from './net/ResAllHeroInfo';
import ResBackpackInfo from './net/ResBackpackInfo';

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
        HeroBoxModel.getInstance().addHero(hero);
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
