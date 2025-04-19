import { ConfigContext } from '../../data/config/container/ConfigContext';
import HeroData from '../../data/config/model/HeroData';
import HeroLevelData from '../../data/config/model/HeroLevelData';
import HerostageData from '../../data/config/model/HerostageData';
import GameContext from '../../GameContext';
import { HeroVo } from '../../net/protocol/items/HeroVo';
import ReqHeroCombine from '../../net/protocol/ReqHeroCombine';
import { ReqHeroUpLevel } from '../../net/protocol/ReqHeroUpLevel';
import { ReqHeroUpStage } from '../../net/protocol/ReqHeroUpStage';
import { ReqPlayerUpLevel } from '../../net/protocol/ReqPlayerUpLevel';
import { ReqPlayerUpStage } from '../../net/protocol/ReqPlayerUpStage';
import { ResHeroUpLevel } from '../../net/protocol/ResHeroUpLevel';
import { ResHeroUpStage } from '../../net/protocol/ResHeroUpStage';
import { ResPlayerUpLevel } from '../../net/protocol/ResPlayerUpLevel';
import { ResPlayerUpStage } from '../../net/protocol/ResPlayerUpStage';
import ResHeroCombine from '../../net/ResHeroCombine';
import { AttributeBox } from '../attribute/attributebox';
import { PurseModel } from '../main/PurseModel';

export class HeroBoxModel {
  private static instance: HeroBoxModel;
  private constructor() {}

  private heros: Map<number, HeroVo> = new Map();

  private quality2Pics: Map<number, string> = new Map();

  // 英雄属性变化回调
  private heroAttrChangedCallbacks: (() => void)[] = [];

  public static getInstance(): HeroBoxModel {
    if (!HeroBoxModel.instance) {
      HeroBoxModel.instance = new HeroBoxModel();
      HeroBoxModel.instance.quality2Pics = new Map();
      HeroBoxModel.instance.quality2Pics.set(0, 'quality_gold');
      HeroBoxModel.instance.quality2Pics.set(1, 'quality_red');
      HeroBoxModel.instance.quality2Pics.set(2, 'quality_purse');
      HeroBoxModel.instance.quality2Pics.set(3, 'quality_pink');
      HeroBoxModel.instance.quality2Pics.set(4, 'quality_blue');
      HeroBoxModel.instance.quality2Pics.set(5, 'quality_green');
    }
    return HeroBoxModel.instance;
  }

  public reset(heros: Map<number, HeroVo>) {
    this.heros = heros;
    for (const hero of this.heros.values()) {
      hero.attrBox = new AttributeBox(hero.attrs);
    }
  }

  public getHero(id: number): HeroVo {
    return this.heros.get(id);
  }

  public addHero(hero: HeroVo) {
    this.heros.set(hero.id, hero);
    hero.attrBox = new AttributeBox(hero.attrs);

    this.notifyHeroAttrChanged();
  }

  public hasHero(id: number): boolean {
    return this.heros.has(id);
  }

  public getHeroes(): Array<HeroVo> {
    return Array.from(this.heros.values());
  }

  public getQualityPicture(quality: number): string {
    return this.quality2Pics.get(quality);
  }

  public onHeroAttrChanged(callback: () => void) {
    this.heroAttrChangedCallbacks.push(callback);
  }

  private notifyHeroAttrChanged() {
    this.heroAttrChangedCallbacks.forEach((callback) => callback());
  }

  public calcUpLevel(hero: HeroVo): number {
    let heroData: HeroData = ConfigContext.configHeroContainer.getRecord(hero.id);
    let currLevel = hero.level;
    let heroLevelData: HeroLevelData = ConfigContext.configHeroLevelContainer.getRecord(currLevel);
    let heroStageData: HerostageData = ConfigContext.configHeroStageContainer.getRecordByStage(
      hero.stage
    );
    let myGold = PurseModel.getInstance().gold;
    let canUpLevel = 0;
    let cost = heroLevelData.cost;

    let master = this.getMaster();
    while (myGold >= cost) {
      canUpLevel++;
      myGold -= cost;
      cost += heroLevelData.cost;
      currLevel++;
      if (canUpLevel >= 10) {
        break;
      }
      if (currLevel >= ConfigContext.configHeroLevelContainer.getMaxLevel()) {
        break;
      }
      if (currLevel >= heroStageData.max_level) {
        break;
      }
      if (heroData.quality !== 0) {
        // 普通英雄等级不能超过主公
        if (currLevel >= master.level) {
          break;
        }
      }
      heroLevelData = ConfigContext.configHeroLevelContainer.getRecord(currLevel);
    }

    if (canUpLevel >= 10) {
      return 10;
    } else if (canUpLevel >= 5) {
      return 5;
    } else if (canUpLevel >= 1) {
      return 1;
    } else {
      return 0;
    }
  }

  public getMaster(): HeroVo {
    for (const hero of this.heros.values()) {
      let heroData: HeroData = ConfigContext.configHeroContainer.getRecord(hero.id);
      if (heroData.quality === 0) {
        return hero;
      }
    }
    return null;
  }

  public requestUpLevel(heroId: number, toLevel: number): Promise<number> {
    let heroData = ConfigContext.configHeroContainer.getRecord(heroId);

    return new Promise<number>((resolve, reject) => {
      if (heroData.quality == 0) {
        GameContext.instance.WebSocketClient.sendMessage(
          ReqPlayerUpLevel.cmd,
          {
            heroId: heroId,
            toLevel: toLevel,
          },
          (msg: ResPlayerUpLevel) => {
            resolve(msg.code);
          }
        );
      } else {
        GameContext.instance.WebSocketClient.sendMessage(
          ReqHeroUpLevel.cmd,
          {
            heroId: heroId,
            toLevel: toLevel,
          },
          (msg: ResHeroUpLevel) => {
            resolve(msg.code);
          }
        );
      }
    });
  }

  public requestUpStage(heroId: number): Promise<number> {
    let heroData = ConfigContext.configHeroContainer.getRecord(heroId);

    return new Promise<number>((resolve, reject) => {
      if (heroData.quality === 0) {
        GameContext.instance.WebSocketClient.sendMessage(
          ReqPlayerUpStage.cmd,
          {},
          (msg: ResPlayerUpStage) => {
            resolve(msg.code);
          }
        );
      } else {
        GameContext.instance.WebSocketClient.sendMessage(
          ReqHeroUpStage.cmd,
          { heroId: heroId },
          (msg: ResHeroUpStage) => {
            resolve(msg.code);
          }
        );
      }
    });
  }

  public requestCombine(heroId: number): Promise<number> {
    return new Promise<number>((resolve, reject) => {
      GameContext.instance.WebSocketClient.sendMessage(
        ReqHeroCombine.cmd,
        { heroId: heroId },
        (msg: ResHeroCombine) => {
          resolve(msg.code);
        }
      );
    });
  }
}
