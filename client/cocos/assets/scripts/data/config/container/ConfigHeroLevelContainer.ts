import ConfigContainer from '../ConfigContainer';

import HeroLevelData from '../model/HeroLevelData';

export default class ConfigHeroLevelContainer extends ConfigContainer<HeroLevelData> {
  private maxLevel: number = 0;

  public constructor() {
    super(HeroLevelData, HeroLevelData.fileName);
  }

  public afterLoad() {
    for (const item of this._datas.values()) {
      if (item.level > this.maxLevel) {
        this.maxLevel = item.level;
      }
    }
  }

  public getMaxLevel(): number {
    return this.maxLevel;
  }
}
