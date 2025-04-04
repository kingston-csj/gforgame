import { HeroVo } from './MsgItems/HeroVo';

export class ResAllHeroInfo {
  public static cmd = 5003;
  public heros: HeroVo[] = [];
}
