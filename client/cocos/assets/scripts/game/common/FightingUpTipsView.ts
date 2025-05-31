import { _decorator, Label, Tween, tween } from 'cc';
import { BaseController } from '../../frame/mvc/BaseController';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import UiViewFactory from '../../ui/UiViewFactory';
const { ccclass, property } = _decorator;

@ccclass('FightingUpTipsView')
export class FightingUpTipsView extends BaseController {
  private static instance: FightingUpTipsView;
  private static creatingPromise: Promise<FightingUpTipsView> | null = null;
  @property(Label)
  public text: Label = null;

  @property(Label)
  public addNum: Label = null;

  protected _tween: Tween<any> = null;

  public static display(fromNum: number, addNum: number) {
    this.getInstance().then((view: FightingUpTipsView) => {
      view.node.active = true;
      view.addNum.string = addNum.toString();

      view.text.string = '战力：' + fromNum;
      view.addNum.string = '+' + addNum;
      if (view._tween != null) {
        view._tween.stop();
        view._tween = null;
      }
      let params = {
        score: fromNum,
        add: addNum,
      };
      let sum = fromNum + addNum;

      view._tween = tween(params)
        .delay(0.2)
        .to(
          1.1,
          {
            score: sum,
            add: 0,
          },
          {
            easing: 'circOut',
            onUpdate: () => {
              view.text.string = '战力：' + Math.floor(params.score);
              view.addNum.string = '+' + Math.floor(params.add);
            },
          }
        )
        .call(() => {
          view._tween = null;
          view.node.active = false;
        })
        .start();
    });
  }

  private static getInstance(): Promise<FightingUpTipsView> {
    if (this.instance) {
      return Promise.resolve(this.instance);
    }
    if (this.creatingPromise) {
      return this.creatingPromise;
    }
    this.creatingPromise = new Promise((resolve) => {
      UiViewFactory.createUi(
        R.Prefabs.FightingUpTipsPane,
        LayerIdx.layer5,
        (ui: FightingUpTipsView) => {
          this.instance = ui;
          this.creatingPromise = null;
          resolve(ui);
        }
      );
    });
    return this.creatingPromise;
  }
}
