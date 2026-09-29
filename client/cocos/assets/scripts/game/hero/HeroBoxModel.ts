import { ConfigContext } from "../../data/config/container/ConfigContext";
import HeroData from "../../data/config/model/HeroData";
import HeroLevelData from "../../data/config/model/HeroLevelData";
import HerostageData from "../../data/config/model/HerostageData";
import GameContext from "../../GameContext";
import { HeroInfo as HeroVo } from "../../net/protocol/hero/HeroInfo";
import { ReqHeroChangePosition } from "../../net/protocol/hero/ReqHeroChangePosition";
import { ReqHeroCombine } from "../../net/protocol/hero/ReqHeroCombine";

import { ReqHeroOffFight } from "../../net/protocol/hero/ReqHeroOffFight";
import { ReqHeroUpFight } from "../../net/protocol/hero/ReqHeroUpFight";
import { ReqHeroLevelUp } from "../../net/protocol/hero/ReqHeroLevelUp";
import { ReqHeroUpStage } from "../../net/protocol/hero/ReqHeroUpStage";
import { ReqPlayerUpLevel } from "../../net/protocol/player/ReqPlayerUpLevel";
import { ReqPlayerUpStage } from "../../net/protocol/player/ReqPlayerUpStage";
import { ResHeroChangePosition } from "../../net/protocol/hero/ResHeroChangePosition";
import { ResHeroCombine } from "../../net/protocol/hero/ResHeroCombine";
import { ResHeroOffFight } from "../../net/protocol/hero/ResHeroOffFight";
import { ResHeroUpFight } from "../../net/protocol/hero/ResHeroUpFight";
import { ResHeroLevelUp } from "../../net/protocol/hero/ResHeroLevelUp";
import { ResHeroUpStage } from "../../net/protocol/hero/ResHeroUpStage";
import { ResPlayerUpLevel } from "../../net/protocol/player/ResPlayerUpLevel";
import { ResPlayerUpStage } from "../../net/protocol/player/ResPlayerUpStage";

import { AttributeBox } from "../attribute/attributebox";
import GameConstants from "../constants/GameConstants";
import BagpackModel from "../item/BagpackModel";
import { PurseModel } from "../main/PurseModel";

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
      HeroBoxModel.instance.quality2Pics.set(0, "quality_gold");
      HeroBoxModel.instance.quality2Pics.set(1, "quality_red");
      HeroBoxModel.instance.quality2Pics.set(2, "quality_purse");
      HeroBoxModel.instance.quality2Pics.set(3, "quality_pink");
      HeroBoxModel.instance.quality2Pics.set(4, "quality_blue");
      HeroBoxModel.instance.quality2Pics.set(5, "quality_green");
    }
    return HeroBoxModel.instance;
  }

  public reset(heros: Map<number, HeroVo>) {
    this.heros = heros;
    for (const hero of this.heros.values()) {
      // hero.attrBox = new AttributeBox(hero.attrs);
    }
  }

  public getHero(id: number): HeroVo {
    return this.heros.get(id);
  }

  public addHero(hero: HeroVo) {
    this.heros.set(hero.id, hero);
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

  public checkCanUpStage(hero: HeroVo): boolean {
    let heroStageData = ConfigContext.configHeroStageContainer.getRecordByStage(
      hero.stage
    );
    let nextStageData = ConfigContext.configHeroStageContainer.getRecordByStage(
      hero.stage + 1
    );
    return hero.level == heroStageData.max_level && nextStageData != null;
  }

  public checkUpStageItem(hero: HeroVo): boolean {
    let heroStageData = ConfigContext.configHeroStageContainer.getRecordByStage(
      hero.stage
    );
    let nextStageData = ConfigContext.configHeroStageContainer.getRecordByStage(
      hero.stage + 1
    );
    let costItemId = GameConstants.Item.UpStage;
    let ownItem = BagpackModel.getInstance().getItemCount(costItemId);
    return ownItem >= nextStageData.cost;
  }

  public calcUpLevel(hero: HeroVo): number {
    let heroData: HeroData = ConfigContext.configHeroContainer.getRecord(
      hero.id
    );
    let currLevel = hero.level;
    let heroLevelData: HeroLevelData =
      ConfigContext.configHeroLevelContainer.getRecord(currLevel);
    let heroStageData: HerostageData =
      ConfigContext.configHeroStageContainer.getRecordByStage(hero.stage);
    let myGold = PurseModel.getInstance().gold;
    let canUpLevel = 0;
    let cost = heroLevelData.cost;

    let master = this.getMaster();
    while (myGold >= cost) {
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
      canUpLevel++;
      myGold -= cost;
      cost += heroLevelData.cost;
      currLevel++;
      if (canUpLevel >= 10) {
        break;
      }
      heroLevelData =
        ConfigContext.configHeroLevelContainer.getRecord(currLevel);
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
      let heroData: HeroData = ConfigContext.configHeroContainer.getRecord(
        hero.id
      );
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
        GameContext.wsClient.sendMessage(
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
        GameContext.wsClient.sendMessage(
          ReqHeroLevelUp.cmd,
          {
            heroId: heroId,
            toLevel: toLevel,
          },
          (msg: ResHeroLevelUp) => {
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
        GameContext.wsClient.sendMessage(
          ReqPlayerUpStage.cmd,
          {},
          (msg: ResPlayerUpStage) => {
            resolve(msg.code);
          }
        );
      } else {
        GameContext.wsClient.sendMessage(
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
      GameContext.wsClient.sendMessage(
        ReqHeroCombine.cmd,
        { heroId: heroId },
        (msg: ResHeroCombine) => {
          resolve(msg.code);
        }
      );
    });
  }

  public requestUpFight(heroId: number, position: number): Promise<number> {
    return new Promise<number>((resolve, reject) => {
      GameContext.wsClient.sendMessage(
        ReqHeroUpFight.cmd,
        { heroId: heroId, position: position },
        (msg: ResHeroUpFight) => {
          resolve(msg.code);
        }
      );
    });
  }

  public requestOffFight(heroId: number): Promise<number> {
    return new Promise<number>((resolve, reject) => {
      GameContext.wsClient.sendMessage(
        ReqHeroOffFight.cmd,
        { heroId: heroId },
        (msg: ResHeroOffFight) => {
          resolve(msg.code);
        }
      );
    });
  }

  public requestChangePostion(
    heroId: number,
    position: number
  ): Promise<ResHeroChangePosition> {
    return new Promise<ResHeroChangePosition>((resolve, reject) => {
      GameContext.wsClient.sendMessage(
        ReqHeroChangePosition.cmd,
        { heroId: heroId, position: position },
        (msg: ResHeroChangePosition) => {
          resolve(msg);
        }
      );
    });
  }

  public getEmptyPostion(): number {
    let used = new Set<number>();
    this.heros.forEach((e) => {
      used.add(e.position);
    });
    for (let i = 1; i <= 5; i++) {
      if (!used.has(i)) {
        return i;
      }
    }
    return 0;
  }

  public getFightPower(): number {
    let power = 0;
    this.heros.forEach((e) => {
      power += e.fight;
    });
    return power;
  }
}
