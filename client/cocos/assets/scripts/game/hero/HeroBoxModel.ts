import { HeroVo } from '../../net/MsgItems/HeroVo';

export class HeroBoxModel {
  private static instance: HeroBoxModel;
  private constructor() {}

  private heros: Map<number, HeroVo> = new Map();

  private quality2Pics: Map<number, string> = new Map();

  public static getInstance(): HeroBoxModel {
    if (!HeroBoxModel.instance) {
      HeroBoxModel.instance = new HeroBoxModel();
      HeroBoxModel.instance.quality2Pics = new Map();
      HeroBoxModel.instance.quality2Pics.set(0, 'quality_gold');
      HeroBoxModel.instance.quality2Pics.set(1, 'quality_red');
      HeroBoxModel.instance.quality2Pics.set(2, 'quality_purse');
      HeroBoxModel.instance.quality2Pics.set(3, 'quality_pink');
      HeroBoxModel.instance.quality2Pics.set(4, 'quality_blue');
      HeroBoxModel.instance.quality2Pics.set(5, 'quality_green');
    }
    return HeroBoxModel.instance;
  }

  public reset(heros: Map<number, HeroVo>) {
    this.heros = heros;
  }

  public getHero(id: number): HeroVo {
    return this.heros.get(id);
  }

  public addHero(hero: HeroVo) {
    this.heros.set(hero.id, hero);
  }

  public getHeroes(): Array<HeroVo> {
    return Array.from(this.heros.values());
  }

  public getQualityPicture(quality: number): string {
    return this.quality2Pics.get(quality);
  }
}
