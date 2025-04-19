import { _decorator, Button, instantiate, Node, Prefab, ScrollView, Toggle } from 'cc';

import { ConfigContext } from '../../data/config/container/ConfigContext';
import HeroData from '../../data/config/model/HeroData';
import { BaseUiView } from '../../ui/BaseUiView';
import { HeroBoxModel } from './HeroBoxModel';
import { HeroDetailController } from './HeroDetailController';
import { HeroItem } from './HeroItemView';
import { HeroLibItemView } from './HeroLibItemView';
const { ccclass, property } = _decorator;

@ccclass('HeroLibView')
export class HeroLibView extends BaseUiView {
  @property(Node)
  heroContainer: Node;

  @property(Prefab)
  heroPrefab: Prefab;

  @property(Node)
  closeBtn: Node;

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

  // 0: 全部,
  private selectedType: number = 0;

  private children: Map<number, HeroItem> = new Map();

  protected start(): void {
    this.registerClickEvent(this.closeBtn, this.hide, this);

    this.registerToggleButtonEvent(this.allBtn, 0);
    this.registerToggleButtonEvent(this.shuBtn, 1);
    this.registerToggleButtonEvent(this.weiBtn, 2);
    this.registerToggleButtonEvent(this.wuBtn, 3);
    this.registerToggleButtonEvent(this.haoBtn, 4);
  }

  private registerToggleButtonEvent(button: Node, type: number): void {
    this.registerClickEvent(button, () => {
      this.selectedType = type;
      this.showItems();
    });
  }

  private showItems(): void {
    this.heroContainer.children.forEach((child) => {
      child.destroy();
    });
    this.children.clear();
    let items: Array<HeroData> = this.getAllItems();
    // 根据品质排序
    items.sort((a, b) => {
      const config1 = ConfigContext.configHeroContainer.getRecord(a.id);
      const config2 = ConfigContext.configHeroContainer.getRecord(b.id);

      const has1 = HeroBoxModel.getInstance().hasHero(a.id);
      const has2 = HeroBoxModel.getInstance().hasHero(b.id);

      if (has1 && !has2) {
        return -1;
      }
      if (!has1 && has2) {
        return 1;
      }

      return config1.quality - config2.quality;
    });
    for (let i = 0; i < items.length; i++) {
      let item = instantiate(this.heroPrefab);
      item.setParent(this.heroContainer);
      item.getComponent(HeroLibItemView).fillData(items[i].id);

      if (HeroBoxModel.getInstance().hasHero(items[i].id)) {
        item
          .getChildByName('ui')
          .getComponent(Button)
          .node.on(Button.EventType.CLICK, () => {
            HeroDetailController.openUi(items[i].id);
          });
      }

      this.children.set(items[i].id, item.getComponent(HeroItem));
    }
    // 自动滑动到最顶部的item
    this.scrollToTop();
  }

  private getAllItems(): Array<HeroData> {
    if (this.selectedType === 0) {
      return ConfigContext.configHeroContainer.getAllRecords().filter((hero) => {
        return hero.quality !== 0;
      });
    }
    let items: Array<HeroData> = ConfigContext.configHeroContainer
      .getAllRecords()
      .filter((hero) => {
        return hero.camp === this.selectedType && hero.quality !== 0;
      });
    return items;
  }

  protected onDisplay(): void {
    this.selectedType = 0;
    this.allBtn.getComponent(Toggle).isChecked = true;
    this.showItems();
  }

  private scrollToTop() {
    const scrollView = this.heroContainer.parent.parent.getComponent(ScrollView);
    if (scrollView) {
      scrollView.scrollToTop(0.1);
    }
  }
}
