import { _decorator, instantiate, Label, Node, Prefab, ScrollView } from 'cc';

import { HeroVo } from '../../net/MsgItems/HeroVo';
import { BaseUiView } from '../../ui/BaseUiView';
import { PurseModel } from '../main/PurseModel';
import { HeroBoxModel } from './HeroBoxModel';
import { HeroItem } from './HeroItem';

const { ccclass, property } = _decorator;

@ccclass('HeroMainView')
export class HeroMainView extends BaseUiView {
  @property(Label)
  goldLabel: Label;
  @property(Node)
  heroContainer: Node;
  @property(Prefab)
  heroPrefab: Prefab;

  protected onDisplay(): void {
    this.heroContainer.children.forEach((child) => {
      child.destroy();
    });

    let items: Array<HeroVo> = HeroBoxModel.getInstance().getHeroes();
    for (let i = 0; i < items.length; i++) {
      let item = instantiate(this.heroPrefab);
      item.setParent(this.heroContainer);
      item.getComponent(HeroItem).fillData(items[i]);
    }
    // 自动滑动到最顶部的item
    this.scrollToTop();
    this.goldLabel.string = PurseModel.getInstance().gold.toString();

    // 绑定金币数据
    PurseModel.getInstance().onGoldChange((value) => {
      if (this.goldLabel) {
        this.goldLabel.string = value.toString();
      }
    });
  }

  private scrollToTop() {
    const scrollView = this.heroContainer.parent.parent.getComponent(ScrollView);
    if (scrollView) {
      scrollView.scrollToTop(0.1);
    }
  }
}
