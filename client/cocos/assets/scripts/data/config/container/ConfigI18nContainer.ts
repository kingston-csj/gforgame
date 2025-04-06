import ConfigContainer from '../ConfigContainer';
import I18nData from '../model/I18nData';

export class ConfigI18nContainer extends ConfigContainer<I18nData> {
  public constructor() {
    super(I18nData, I18nData.fileName);
  }
}
