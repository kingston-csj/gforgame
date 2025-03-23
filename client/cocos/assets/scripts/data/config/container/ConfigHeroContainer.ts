import ConfigContainer from '../ConfigContainer';
import Config_heroData from '../model/HeroData';

export default class ConfigHeroContainer extends ConfigContainer<Config_heroData> {
  private static _instance: ConfigHeroContainer | null = null;

  private constructor() {
    super(Config_heroData, Config_heroData.fileName);
  }

  public static getInstance(): ConfigHeroContainer {
    if (!ConfigHeroContainer._instance) {
      ConfigHeroContainer._instance = new ConfigHeroContainer();
    }
    return ConfigHeroContainer._instance;
  }
}
