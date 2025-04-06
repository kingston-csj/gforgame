import { _decorator, Label, Node, Sprite } from 'cc';

import { ConfigContext } from '../../data/config/container/ConfigContext';
import GameContext from '../../GameContext';
import { HeroVo } from '../../net/MsgItems/HeroVo';
import { ReqHeroUpLevel } from '../../net/ReqHeroUpLevel';
import { ResHeroUpLevel } from '../../net/ResHeroUpLevel';
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
  public btn: Node;

  @property(Label)
  public upLevel: Label;

  @property(Label)
  public level: Label;
  @property(Label)
  public fighting: Label;

  public hero: HeroVo;

  protected start(): void {
    this.registerClickEvent(this.btn, () => {
      let canUpLevel = HeroBoxModel.getInstance().calcUpLevel(this.hero);
      GameContext.instance.WebSocketClient.sendMessage(
        ReqHeroUpLevel.cmd,
        {
          heroId: this.hero.id,
          toLevel: this.hero.level + canUpLevel,
        },
        (msg: ResHeroUpLevel) => {
          if (msg.code === 0) {
            this.hero.level += canUpLevel;
            this.fillData(this.hero);
          } else {
            TipsPaneController.openUi(msg.code);
          }
        }
      );
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
  }

  public updateUpLevelBtn(times: number) {
    if (!this.hero) {
      return;
    }

    if (times > 1) {
      this.upLevel.string = `升${times}级`;
    } else {
      this.upLevel.string = `升级`;
    }
    this.level.string = this.hero.level.toString();
  }
}
