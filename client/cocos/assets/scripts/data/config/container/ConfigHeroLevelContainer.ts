import ConfigContainer from '../ConfigContainer';

import HeroLevelData from '../model/HeroLevelData';

export default class ConfigHeroLevelContainer extends ConfigContainer<HeroLevelData> {
  public constructor() {
    super(HeroLevelData, HeroLevelData.fileName);
  }
}
