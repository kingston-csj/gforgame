export class HeroVo {
  public id: number;
  public level: number;
  public position: number;
  public stage: number;

  constructor(heroId: number) {
    this.id = heroId;
    this.level = 1;
    this.position = 0;
    this.stage = 0;
  }
}
