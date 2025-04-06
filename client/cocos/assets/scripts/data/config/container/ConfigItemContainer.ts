import ConfigContainer from '../ConfigContainer';
import ItemData from '../model/ItemData';

export default class ConfigItemContainer extends ConfigContainer<ItemData> {
  public constructor() {
    super(ItemData, ItemData.fileName);
  }
}
