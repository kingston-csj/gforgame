import { RedDotManager } from '../../frame/reddot/RedDotManager';
import { HeroBoxModel } from './HeroBoxModel';

export class HeroManager {
  private static instance: HeroManager = new HeroManager();

  private constructor() {}

  public static getInstance(): HeroManager {
    return HeroManager.instance;
  }

  public refreshRedDots() {
    let heros = HeroBoxModel.getInstance().getHeroes();
    for (let hero of heros) {
      let upStageRedDot = false;
      if (
        HeroBoxModel.getInstance().checkCanUpStage(hero) &&
        HeroBoxModel.getInstance().checkUpStageItem(hero)
      ) {
        upStageRedDot = true;
      }
      let upLevelRedDot = HeroBoxModel.getInstance().calcUpLevel(hero) > 0;
      RedDotManager.instance.updateScore(`hero/items/${hero.id}/stage`, upStageRedDot ? 1 : 0);
      RedDotManager.instance.updateScore(`hero/items/${hero.id}/level`, upLevelRedDot ? 1 : 0);
    }
  }
}
