import { ConfigContext } from '../../data/config/container/ConfigContext';
import HeroLevelData from '../../data/config/model/HeroLevelData';
import { HeroVo } from '../../net/MsgItems/HeroVo';
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
    let heroLevelData: HeroLevelData = ConfigContext.configHeroLevelContainer.getRecord(hero.level);
    let myGold = PurseModel.getInstance().gold;

    let canUpLevel = 0;
    let cost = heroLevelData.cost;
    while (myGold >= cost) {
      canUpLevel++;
      myGold -= cost;
      cost += heroLevelData.cost;
      if (canUpLevel >= 10) {
        break;
      }
    }

    if (canUpLevel >= 10) {
      return 10;
    } else if (canUpLevel >= 5) {
      return 5;
    } else {
      return 1;
    }
  }
}
