import { _decorator, instantiate, Label, Node, Prefab, ScrollView, Toggle } from 'cc';
import { ConfigContext } from '../../data/config/container/ConfigContext';
import { BaseUiView } from '../../frame/mvc/BaseUiView';
import { HeroVo } from '../../net/protocol/items/HeroVo';
import { NumberUtils } from '../../utils/NumberUtils';
import { HeroBoxModel } from '../hero/HeroBoxModel';
import { BuZhenHeroDownItem } from './BuZhenHeroDownItem';
import { BuZhenHeroUpItem } from './BuZhenHeroUpItem';
const { ccclass, property } = _decorator;

@ccclass('BuZhenView')
export class BuZhenView extends BaseUiView {
  @property(Node)
  closeBtn: Node;

  @property(Node)
  heroContainer: Node;

  @property(Prefab)
  upHeroPrefab: Prefab;

  @property(Prefab)
  downHeroPrefab: Prefab;

  @property(Node)
  public allBtn: Node | null = null;

  @property(Node)
  public shuBtn: Node | null = null;

  @property(Node)
  public weiBtn: Node | null = null;
  @property(Node)
  public wuBtn: Node | null = null;

  @property(Node)
  public haoBtn: Node | null = null;

  private selectedType: number = 0;

  @property(Node)
  public pos1: Node = null;

  @property(Node)
  public pos2: Node = null;

  @property(Node)
  public pos3: Node = null;

  @property(Node)
  public pos4: Node = null;

  @property(Node)
  public pos5: Node = null;

  @property(Label)
  private fight: Label;

  private posList: Array<Node> = [];

  protected start(): void {
    this.registerClickEvent(this.closeBtn, this.hide, this);

    this.registerToggleButtonEvent(this.allBtn, 0);
    this.registerToggleButtonEvent(this.shuBtn, 1);
    this.registerToggleButtonEvent(this.weiBtn, 2);
    this.registerToggleButtonEvent(this.wuBtn, 3);
    this.registerToggleButtonEvent(this.haoBtn, 4);

    this.registerClickEvent(this.closeBtn, this.hide, this);

    this.posList = [this.pos1, this.pos2, this.pos3, this.pos4, this.pos5];
  }

  private registerToggleButtonEvent(button: Node, type: number): void {
    this.registerClickEvent(button, () => {
      this.selectedType = type;
      this.showMyHeros();
    });
  }

  private showMyHeros(): void {
    this.heroContainer.children.forEach((child) => {
      child.destroy();
    });

    let items: Array<HeroVo> = this.getAllItems();
    // 根据品质排序
    items.sort((a, b) => {
      const config1 = ConfigContext.configHeroContainer.getRecord(a.id);
      const config2 = ConfigContext.configHeroContainer.getRecord(b.id);

      return config1.quality - config2.quality;
    });

    for (let i = 0; i < items.length; i++) {
      let item = instantiate(this.downHeroPrefab);
      item.setParent(this.heroContainer);
      item.getComponent(BuZhenHeroDownItem).fillData(items[i]);
    }

    // 滚动到最左边
    const scrollView = this.heroContainer.parent.parent.getComponent(ScrollView);
    if (scrollView) {
      scrollView.scrollToLeft(0.1);
    }
  }

  private getAllItems(): Array<HeroVo> {
    if (this.selectedType === 0) {
      return HeroBoxModel.getInstance()
        .getHeroes()
        .filter((hero) => {
          let heroData = ConfigContext.configHeroContainer.getRecord(hero.id);
          return heroData.quality > 0;
        });
    }
    let items: Array<HeroVo> = HeroBoxModel.getInstance()
      .getHeroes()
      .filter((hero) => {
        let heroData = ConfigContext.configHeroContainer.getRecord(hero.id);
        return heroData.camp === this.selectedType && heroData.quality > 0;
      });
    return items;
  }

  protected onDisplay(): void {
    this.selectedType = 0;
    this.allBtn.getComponent(Toggle).isChecked = true;
    this.showMyHeros();
    this.showLineupHeros();
  }

  public showLineupHeros(): void {
    let items: Array<HeroVo> = this.getAllItems();

    this.posList.forEach((pos) => {
      pos.getChildByName('pos').active = true;
      pos.getChildByName('node').destroyAllChildren();
    });
    for (let i = 0; i < items.length; i++) {
      let hero = items[i];
      let pos = hero.position;
      if (pos > 0) {
        let heroItem = instantiate(this.upHeroPrefab);
        this.posList[pos - 1].getChildByName('pos').active = false;
        this.posList[pos - 1].getChildByName('node').addChild(heroItem);
        heroItem.getComponent(BuZhenHeroUpItem).fillData(hero);
      }
    }
    this.fight.string = NumberUtils.formatNumber(HeroBoxModel.getInstance().getFightPower());
  }

  public hide(): void {
    this.scheduleOnce(() => {
      super.hide();
    });
  }
}
