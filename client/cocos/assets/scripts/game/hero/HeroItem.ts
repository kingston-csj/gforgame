import { _decorator, Label, Node, Sprite, UITransform } from 'cc';

import ConfigHeroContainer from '../../data/config/container/ConfigHeroContainer';
import ConfigHeroLevelContainer from '../../data/config/container/ConfigHeroLevelContainer';
import HeroLevelData from '../../data/config/model/HeroLevelData';
import GameContext from '../../GameContext';
import { HeroVo } from '../../net/MsgItems/HeroVo';
import { ReqHeroUpLevel } from '../../net/ReqHeroUpLevel';
import { ResHeroUpLevel } from '../../net/ResHeroUpLevel';
import AssetResourceFactory from '../../ui/AssetResourceFactory';
import { BaseUiView } from '../../ui/BaseUiView';
import R from '../../ui/R';
import { TipsPaneController } from '../common/TipsPaneController';
import { PurseModel } from '../main/PurseModel';
import { HeroBoxModel } from './HeroBoxModel';
const { ccclass, property } = _decorator;

@ccclass('HeroItem')
export class HeroItem extends BaseUiView {
  @property(Sprite)
  public kuang: Sprite;

  @property(Label)
  public heroName: Label;

  @property(Label)
  public heroLevel: Label;

  @property(Node)
  public icon: Node;

  @property(Node)
  public btn: Node;

  @property(Label)
  public upLevel: Label;

  @property(Label)
  public level: Label;

  private hero: HeroVo;

  protected start(): void {
    this.registerClickEvent(this.btn, () => {
      let canUpLevel = this.calcUpLevel(this.hero);
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

  private calcUpLevel(hero: HeroVo) {
    let heroLevelContainer: ConfigHeroLevelContainer = ConfigHeroLevelContainer.getInstance();
    let heroLevelData: HeroLevelData = heroLevelContainer.getRecord(hero.level);
    let myGold = PurseModel.getInstance().gold;
    // 根据当前金币量，计算可以升级的次数，1级，5级，10级。每一级的cost不一样
    let canUpLevel = 0;
    let cost = heroLevelData.cost;
    while (myGold >= cost) {
      canUpLevel++;
      myGold -= cost;
      cost += heroLevelData.cost;
      if (canUpLevel >= 10) {
        break;
      }
    }
    if (canUpLevel >= 10) {
      return 10;
    } else if (canUpLevel >= 5) {
      return 5;
    } else {
      return 1;
    }
  }

  public fillData(hero: HeroVo) {
    this.hero = hero;
    let heroContianer: ConfigHeroContainer = ConfigHeroContainer.getInstance();
    let heroData = heroContianer.getRecord(hero.id);
    this.heroName.string = heroData.name;
    this.upLevel.string = hero.level.toString();

    let canUpLevel = this.calcUpLevel(hero);
    if (canUpLevel > 1) {
      this.upLevel.string = `升${canUpLevel}级`;
    } else {
      this.upLevel.string = `升级`;
    }
    this.level.string = hero.level.toString();

    const iconTransform = this.icon.getComponent(UITransform);
    if (!iconTransform) {
      console.warn('Icon node has no UITransform component');
      return;
    }
    // 保存节点当前的尺寸，用于调整图像
    const originalIconWidth = iconTransform.width;
    const originalIconHeight = iconTransform.height;

    let heroSpriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Hero);
    this.icon.getComponent(Sprite).spriteFrame = heroSpriteAtlas.getSpriteFrame(heroData.icon);
    // 获取当前SpriteFrame
    const sprite = this.icon.getComponent(Sprite);
    if (!sprite || !sprite.spriteFrame) {
      console.warn('Icon has no valid sprite frame');
      return;
    }
    // 设置UITransform的contentSize为原始图片尺寸
    iconTransform.setContentSize(originalIconWidth, originalIconHeight);
    let qualityPicture = HeroBoxModel.getInstance().getQualityPicture(heroData.quality);
    let qualitySpriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Quality);
    this.kuang.getComponent(Sprite).spriteFrame = qualitySpriteAtlas.getSpriteFrame(qualityPicture);
  }
}
