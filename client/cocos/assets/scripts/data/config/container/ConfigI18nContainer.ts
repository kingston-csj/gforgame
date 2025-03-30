import ConfigContainer from '../ConfigContainer';
import I18nData from '../model/I18nData';

export class ConfigI18nContainer extends ConfigContainer<I18nData> {
  private static _instance: ConfigI18nContainer | null = null;

  private constructor() {
    super(I18nData, I18nData.fileName);
  }

  public static getInstance(): ConfigI18nContainer {
    if (!ConfigI18nContainer._instance) {
      ConfigI18nContainer._instance = new ConfigI18nContainer();
    }
    return ConfigI18nContainer._instance;
  }
}
