import { _decorator, Button, Color, Label, Node, Sprite } from 'cc';

import { ConfigContext } from '../../data/config/container/ConfigContext';
import AssetResourceFactory from '../../ui/AssetResourceFactory';
import { BaseUiView } from '../../ui/BaseUiView';
import R from '../../ui/R';
import { UiUtil } from '../../ui/UiUtil';
import { TipsPaneController } from '../common/TipsPaneController';
import BagpackModel from '../item/BagpackModel';
import { HeroBoxModel } from './HeroBoxModel';
const { ccclass, property } = _decorator;

@ccclass('HeroLibItemView')
export class HeroLibItemView extends BaseUiView {
  @property(Sprite)
  private kuang: Sprite;

  @property(Label)
  private heroName: Label;

  @property(Node)
  combineBtn: Node;

  @property(Label)
  progressLabel: Label;

  @property(Node)
  private icon: Node;

  @property(Node)
  private camp: Node;

  private heroId: number = 0;

  public start(): void {}

  public fillData(heroId: number) {
    this.heroId = heroId;
    let heroData = ConfigContext.configHeroContainer.getRecord(heroId);
    let heroSpriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Hero);
    UiUtil.fillSpriteContent(this.icon, heroSpriteAtlas.getSpriteFrame(heroData.icon));

    let campSpriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Camp);
    UiUtil.fillSpriteContent(this.camp, campSpriteAtlas.getSpriteFrame('camp_' + heroData.camp));

    let qualityPicture = HeroBoxModel.getInstance().getQualityPicture(heroData.quality);
    let qualitySpriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Quality);
    this.heroName.string = heroData.name;
    this.kuang.getComponent(Sprite).spriteFrame = qualitySpriteAtlas.getSpriteFrame(qualityPicture);

    if (HeroBoxModel.getInstance().hasHero(heroId)) {
      this.node.getChildByName('ui').getChildByName('icon').getComponent(Sprite).color =
        Color.WHITE;
      this.node.getChildByName('ui').getChildByName('kuang').getComponent(Sprite).grayscale = false;
      this.progressLabel.node.active = false;
      this.combineBtn.active = false;
    } else {
      let suipianSum = BagpackModel.getInstance().getItemCount(heroData.item);
      this.progressLabel.string = `${suipianSum}/${heroData.shard}`;
      if (suipianSum >= heroData.shard) {
        this.combineBtn.active = true;
        this.progressLabel.node.active = false;
        this.combineBtn.getComponent(Button).node.on(Button.EventType.CLICK, () => {
          HeroBoxModel.getInstance()
            .requestCombine(this.heroId)
            .then((code) => {
              if (code === 0) {
                TipsPaneController.showStringContent('合成成功');
                this.fillData(this.heroId);
              }
            });
        });
      }
    }
  }
}
