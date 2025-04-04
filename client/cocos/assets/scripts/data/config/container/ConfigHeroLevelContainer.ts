import ConfigContainer from '../ConfigContainer';

import HeroLevelData from '../model/HeroLevelData';

export default class ConfigHeroLevelContainer extends ConfigContainer<HeroLevelData> {
  private static _instance: ConfigHeroLevelContainer | null = null;

  private constructor() {
    super(HeroLevelData, HeroLevelData.fileName);
  }

  public static getInstance(): ConfigHeroLevelContainer {
    if (!ConfigHeroLevelContainer._instance) {
      ConfigHeroLevelContainer._instance = new ConfigHeroLevelContainer();
    }
    return ConfigHeroLevelContainer._instance;
  }
}
