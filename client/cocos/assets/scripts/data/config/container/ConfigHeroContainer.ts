import ConfigContainer from '../ConfigContainer';
import HeroData from '../model/HeroData';
export default class ConfigHeroContainer extends ConfigContainer<HeroData> {
  public constructor() {
    super(HeroData, HeroData.fileName);
  }
}
