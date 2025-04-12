import ConfigContainer from '../ConfigContainer';
import ConfigHeroContainer from './ConfigHeroContainer';
import ConfigHeroLevelContainer from './ConfigHeroLevelContainer';
import ConfigHeroStageContainer from './ConfigHeroStageContainer';
import { ConfigI18nContainer } from './ConfigI18nContainer';
import ConfigItemContainer from './ConfigItemContainer';
import ConfigSkillContainer from './ConfigSkillContainer';

export class ConfigContext {
  public static configItemContainer: ConfigItemContainer;
  public static configHeroContainer: ConfigHeroContainer;
  public static configHeroLevelContainer: ConfigHeroLevelContainer;
  public static configSkillContainer: ConfigSkillContainer;
  public static configI18nContainer: ConfigI18nContainer;
  public static configHeroStageContainer: ConfigHeroStageContainer;

  private static readonly containerTypes = [
    ConfigItemContainer,
    ConfigHeroContainer,
    ConfigHeroLevelContainer,
    ConfigHeroStageContainer,
    ConfigSkillContainer,
    ConfigI18nContainer,
  ];

  public static init() {
    // 自动实例化所有Config容器
    for (const type of ConfigContext.containerTypes) {
      let propertyName = type.name;
      if (propertyName.startsWith('Config')) {
        const instance = new type();
        // propertyName首字母小写
        propertyName = propertyName.charAt(0).toLowerCase() + propertyName.slice(1);
        (ConfigContext as any)[propertyName] = instance;
        (instance as ConfigContainer<any>).afterLoad();
      }
    }
  }
}
