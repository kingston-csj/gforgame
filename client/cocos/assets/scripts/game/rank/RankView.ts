import { _decorator, instantiate, Label, Node, Prefab, Toggle } from 'cc';
import { BaseUiView } from '../../frame/mvc/BaseUiView';
import { RankInfo } from '../../net/protocol/items/RankInfo';
import { RankItem } from './RankItem';
import { RankModel } from './RankModel';
const { ccclass, property } = _decorator;

@ccclass('RankView')
export class RankView extends BaseUiView {
  @property(Prefab)
  private rankItemPrefab: Prefab;

  @property(Node)
  private rankContainer: Node;

  @property(Node)
  private top3Container: Node;

  @property(Node)
  private levelBtn: Node;

  @property(Node)
  private fightBtn: Node;

  private selectedType: number = 0;

  private top3: Node[] = [];

  @property(Node)
  closeBtn: Node;

  protected start(): void {
    for (let i = 1; i <= 3; i++) {
      this.top3.push(this.top3Container.getChildByName(`item${i}`));
    }
    this.registerToggleButtonEvent(this.levelBtn, RankModel.RANK_TYPE_LEVEL);
    this.registerToggleButtonEvent(this.fightBtn, RankModel.RANK_TYPE_FIGHTING);
    this.registerClickEvent(this.closeBtn, this.hide, this);
  }

  private registerToggleButtonEvent(button: Node, type: number): void {
    this.registerClickEvent(button, async () => {
      this.selectedType = type;
      const resRankQuery = await RankModel.getInstance().queryRank(type);
      this.showRankItems(resRankQuery.records);
    });
  }

  protected onDisplay(): void {
    this.selectedType = RankModel.RANK_TYPE_LEVEL;
    this.levelBtn.getComponent(Toggle).isChecked = true;
    // 如何让LevenBtn触发一下点击事件
    this.levelBtn.emit('click', this.levelBtn);
  }

  private showRankItems(records: RankInfo[]): void {
    this.rankContainer.children.forEach((child) => {
      child.destroy();
    });
    for (let i = 0; i < 3; i++) {
      if (records[i]) {
        this.top3[i].getChildByName('name').getComponent(Label).string = records[i].name;
        this.top3[i].getChildByName('score').getComponent(Label).string =
          records[i].value.toString();
      }
    }
    for (let i = 3; i < records.length; i++) {
      const item = instantiate(this.rankItemPrefab);
      item.setParent(this.rankContainer);
      item.getComponent(RankItem).fillData(records[i]);
    }
  }
}
