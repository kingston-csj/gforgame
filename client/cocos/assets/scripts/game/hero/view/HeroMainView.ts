import { _decorator, Button, instantiate, Label, Node, Prefab, ScrollView } from 'cc';

import { ConfigContext } from '../../../data/config/container/ConfigContext';
import { BaseUiView } from '../../../frame/mvc/BaseUiView';
import { HeroVo } from '../../../net/protocol/items/HeroVo';
import { BuZhenPaneController } from '../../buzhen/BuZhenPaneController';
import BagpackModel from '../../item/BagpackModel';
import { PurseModel } from '../../main/PurseModel';
import { HeroDetailController } from '../controller/HeroDetailController';
import { HeroLibPaneController } from '../controller/HeroLibPaneController';

import GameConstants from '../../constants/GameConstants';
import { HeroBoxModel } from '../HeroBoxModel';
import { HeroItem } from './HeroItemView';
const { ccclass, property } = _decorator;

@ccclass('HeroMainView')
export class HeroMainView extends BaseUiView {
  @property(Label)
  goldLabel: Label;
  @property(Label)
  itemLabel: Label;
  @property(Node)
  heroContainer: Node;
  @property(Prefab)
  heroPrefab: Prefab;
  @property(Node)
  heroLibBtn: Node;
  @property(Node)
  buZhenBtn: Node;

  private children: Map<number, HeroItem> = new Map();

  protected start(): void {
    this.registerClickEvent(
      this.heroLibBtn,
      () => {
        HeroLibPaneController.openUi();
      },
      this
    );
    this.registerClickEvent(
      this.buZhenBtn,
      () => {
        BuZhenPaneController.openUi();
      },
      this
    );
  }

  protected onDisplay(): void {
    this.heroContainer.children.forEach((child) => {
      child.destroy();
    });
    this.children.clear();
    let items: Array<HeroVo> = HeroBoxModel.getInstance().getHeroes();
    // 根据品质排序
    items.sort((a, b) => {
      const config1 = ConfigContext.configHeroContainer.getRecord(a.id);
      const config2 = ConfigContext.configHeroContainer.getRecord(b.id);

      return config1.quality - config2.quality;
    });
    for (let i = 0; i < items.length; i++) {
      let item = instantiate(this.heroPrefab);
      item.setParent(this.heroContainer);
      item.getComponent(HeroItem).fillData(items[i]);
      item
        .getChildByName('ui')
        .getComponent(Button)
        .node.on(Button.EventType.CLICK, () => {
          HeroDetailController.openUi(items[i].id);
        });
      this.children.set(items[i].id, item.getComponent(HeroItem));
    }
    // 自动滑动到最顶部的item
    this.scrollToTop();
    this.goldLabel.string = PurseModel.getInstance().gold.toString();
    this.itemLabel.string = BagpackModel.getInstance()
      .getItemCount(GameConstants.Item.UpStage)
      .toString();
  }

  private scrollToTop() {
    const scrollView = this.heroContainer.parent.parent.getComponent(ScrollView);
    if (scrollView) {
      scrollView.scrollToTop(0.1);
    }
  }

  public updataAllHeroItems() {
    for (let item of this.children.values()) {
      item.refreshButtonStatus();
    }
  }
}
