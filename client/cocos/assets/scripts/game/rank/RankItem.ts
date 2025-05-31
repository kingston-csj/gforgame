import { _decorator, Label } from 'cc';
import { BaseUiView } from '../../frame/mvc/BaseUiView';
import { RankInfo } from '../../net/protocol/items/RankInfo';
const { ccclass, property } = _decorator;

@ccclass('RankItem')
export class RankItem extends BaseUiView {
  @property(Label)
  order: Label;

  @property(Label)
  ownerName: Label;

  @property(Label)
  score: Label;

  protected start(): void {}

  public fillData(rankInfo: RankInfo): void {
    this.order.string = rankInfo.order.toString();
    this.ownerName.string = rankInfo.name;
    this.score.string = rankInfo.value.toString();
  }
}
