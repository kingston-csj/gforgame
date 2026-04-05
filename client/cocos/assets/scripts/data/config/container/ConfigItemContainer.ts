import ConfigContainer from '../ConfigContainer';
import PropData from '../model/PropData';

export default class ConfigItemContainer extends ConfigContainer<PropData> {
  public constructor() {
    super(PropData, PropData.fileName);
  }
}
