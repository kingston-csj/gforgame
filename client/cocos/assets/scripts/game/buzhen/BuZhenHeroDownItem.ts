import { _decorator, Label, Node } from 'cc';
import { ConfigContext } from '../../data/config/container/ConfigContext';
import { HeroVo } from '../../net/protocol/items/HeroVo';
import AssetResourceFactory from '../../ui/AssetResourceFactory';
import { BaseUiView } from '../../ui/BaseUiView';
import R from '../../ui/R';
import { UiUtil } from '../../ui/UiUtil';
import { ComfirmPaneController } from '../common/ComfirmPaneController';
import { HeroBoxModel } from '../hero/HeroBoxModel';
import { BuZhenPaneController } from './BuZhenPaneController';
const { ccclass, property } = _decorator;

@ccclass('BuZhenHeroDownItem')
export class BuZhenHeroDownItem extends BaseUiView {
  @property(Node)
  private kuang: Node;

  @property(Label)
  private heroName: Label;

  @property(Node)
  private icon: Node;

  @property(Node)
  private camp: Node;

  @property(Node)
  private gou: Node;

  private heroId: number = 0;

  public fillData(hero: HeroVo) {
    this.heroId = hero.id;
    let heroData = ConfigContext.configHeroContainer.getRecord(hero.id);
    this.heroName.string = heroData.name;
    let heroSpriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Hero);
    UiUtil.fillSpriteContent(this.icon, heroSpriteAtlas.getSpriteFrame(heroData.icon));
    let campSpriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Camp);
    UiUtil.fillSpriteContent(this.camp, campSpriteAtlas.getSpriteFrame('camp_' + heroData.camp));

    let qualityPicture = HeroBoxModel.getInstance().getQualityPicture(heroData.quality);
    let qualitySpriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Quality);

    UiUtil.fillSpriteContent(this.kuang, qualitySpriteAtlas.getSpriteFrame(qualityPicture));

    this.registerClickEvent(this.kuang, this.onClick, this);
    if (hero.position > 0) {
      this.gou.active = true;
    }
  }

  private onClick(): void {
    let hero = HeroBoxModel.getInstance().getHero(this.heroId);
    if (hero.position === 0) {
      let emptyPos = HeroBoxModel.getInstance().getEmptyPostion();
      if (emptyPos > 0) {
        // 上阵
        HeroBoxModel.getInstance()
          .requestUpFight(this.heroId, emptyPos)
          .then((code) => {
            if (code === 0) {
              hero.position = emptyPos;
              this.gou.active = true;
              BuZhenPaneController.refreshLineupHeros();
            }
          });
      }
    } else {
      ComfirmPaneController.show('下阵', '确定要下阵吗？', () => {
        HeroBoxModel.getInstance()
          .requestOffFight(this.heroId)
          .then((code) => {
            if (code === 0) {
              this.gou.active = false;
              hero.position = 0;
              ComfirmPaneController.hide();
              BuZhenPaneController.refreshLineupHeros();
            }
          });
      });
    }
  }
}
