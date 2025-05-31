import { _decorator, Label } from 'cc';
import { BaseUiView } from '../../frame/mvc/BaseUiView';

const { ccclass, property } = _decorator;

@ccclass('TipsView')
export class TipsView extends BaseUiView {
  @property(Label)
  tipsLabel: Label = null!;

  public setTips(tips: string) {
    this.tipsLabel.string = tips;
  }
}
