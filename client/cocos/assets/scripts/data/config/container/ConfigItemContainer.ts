import ConfigContainer from '../ConfigContainer';
import ItemData from '../model/ItemData';

export default class ConfigItemContainer extends ConfigContainer<ItemData> {
  private static _instance: ConfigItemContainer | null = null;

  private constructor() {
    super(ItemData, ItemData.fileName);
  }

  public static getInstance(): ConfigItemContainer {
    if (!ConfigItemContainer._instance) {
      ConfigItemContainer._instance = new ConfigItemContainer();
    }
    return ConfigItemContainer._instance;
  }
}
