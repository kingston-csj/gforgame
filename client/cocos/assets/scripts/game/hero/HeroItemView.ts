import { _decorator, Label, Node, Sprite } from 'cc';

import { ConfigContext } from '../../data/config/container/ConfigContext';
import { HeroVo } from '../../net/protocol/MsgItems/HeroVo';
import AssetResourceFactory from '../../ui/AssetResourceFactory';
import { BaseUiView } from '../../ui/BaseUiView';
import R from '../../ui/R';
import { UiUtil } from '../../ui/UiUtil';
import { TipsPaneController } from '../common/TipsPaneController';
import { HeroBoxModel } from './HeroBoxModel';
const { ccclass, property } = _decorator;

@ccclass('HeroItem')
export class HeroItem extends BaseUiView {
  @property(Sprite)
  public kuang: Sprite;

  @property(Label)
  public heroName: Label;

  @property(Node)
  public icon: Node;

  @property(Node)
  public upLevelBtn: Node;
  @property(Node)
  public upStageBtn: Node;

  @property(Label)
  public upLevel: Label;

  @property(Label)
  public level: Label;
  @property(Label)
  public fighting: Label;

  @property(Node)
  public fightGroup: Node;

  public hero: HeroVo;

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
            TipsPaneController.openUi(code);
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
            TipsPaneController.openUi(code);
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

    this.kuang.getComponent(Sprite).spriteFrame = qualitySpriteAtlas.getSpriteFrame(qualityPicture);

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
