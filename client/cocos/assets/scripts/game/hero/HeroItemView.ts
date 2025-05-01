import { _decorator, Label, Node } from 'cc';

import { ConfigContext } from '../../data/config/container/ConfigContext';
import { HeroVo } from '../../net/protocol/items/HeroVo';
import AssetResourceFactory from '../../ui/AssetResourceFactory';
import { BaseUiView } from '../../ui/BaseUiView';
import R from '../../ui/R';
import { UiUtil } from '../../ui/UiUtil';
import { TipsPaneController } from '../common/TipsPaneController';
import { HeroBoxModel } from './HeroBoxModel';
const { ccclass, property } = _decorator;

@ccclass('HeroItem')
export class HeroItem extends BaseUiView {
  @property(Node)
  private kuang: Node;

  @property(Label)
  private heroName: Label;

  @property(Node)
  private icon: Node;

  @property(Node)
  private camp: Node;

  @property(Node)
  private upLevelBtn: Node;

  @property(Node)
  private upStageBtn: Node;

  @property(Label)
  private upLevel: Label;

  @property(Label)
  private level: Label;

  @property(Label)
  private fighting: Label;

  @property(Node)
  private fightGroup: Node;

  private hero: HeroVo;

  protected start(): void {
    this.registerClickEvent(this.upLevelBtn, () => {
      let canUpLevel = HeroBoxModel.getInstance().calcUpLevel(this.hero);
      HeroBoxModel.getInstance()
        .requestUpLevel(this.hero.id, this.hero.level + canUpLevel)
        .then((code) => {
          if (code === 0) {
            this.hero.level += canUpLevel;
            this.fillData(this.hero);
          } else {
            TipsPaneController.showI18nContent(code);
          }
        });
    });

    this.registerClickEvent(this.upStageBtn, () => {
      HeroBoxModel.getInstance()
        .requestUpStage(this.hero.id)
        .then((code) => {
          if (code === 0) {
            this.hero.stage += 1;
            this.fillData(this.hero);
          } else {
            TipsPaneController.showI18nContent(code);
          }
        });
    });
  }

  public fillData(hero: HeroVo) {
    this.hero = hero;
    let heroData = ConfigContext.configHeroContainer.getRecord(hero.id);
    this.heroName.string = heroData.name;
    this.level.string = hero.level.toString();
    this.fighting.string = hero.fight.toString();
    let heroSpriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Hero);
    UiUtil.fillSpriteContent(this.icon, heroSpriteAtlas.getSpriteFrame(heroData.icon));

    let qualityPicture = HeroBoxModel.getInstance().getQualityPicture(heroData.quality);
    let qualitySpriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Quality);

    UiUtil.fillSpriteContent(this.kuang, qualitySpriteAtlas.getSpriteFrame(qualityPicture));

    let campSpriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Camp);
    UiUtil.fillSpriteContent(this.camp, campSpriteAtlas.getSpriteFrame('camp_' + heroData.camp));

    this.refreshButtonStatus();

    if (heroData.quality === 0) {
      this.fightGroup.active = false;
    }
  }

  public refreshButtonStatus() {
    if (!this.hero) {
      return;
    }
    this.upStageBtn.active = false;
    this.upLevelBtn.active = false;

    let heroStageData = ConfigContext.configHeroStageContainer.getRecordByStage(this.hero.stage);
    let nextStageData = ConfigContext.configHeroStageContainer.getRecordByStage(
      this.hero.stage + 1
    );
    if (this.hero.level == heroStageData.max_level && nextStageData) {
      this.upStageBtn.active = true;
    } else {
      if (this.hero.level < ConfigContext.configHeroLevelContainer.getMaxLevel()) {
        this.upLevelBtn.active = true;
        let times = HeroBoxModel.getInstance().calcUpLevel(this.hero);
        if (times > 1) {
          this.upLevel.string = `升${times}级`;
        } else {
          this.upLevel.string = `升级`;
        }
        this.level.string = this.hero.level.toString();
      }
    }
  }
}
