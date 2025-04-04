import { _decorator } from 'cc';

import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import UiViewFactory from '../../ui/UiViewFactory';

import { BaseController } from '../../ui/BaseController';

import { MainPaneView } from './MainPaneView';
import { PurseModel } from './PurseModel';
const { ccclass, property } = _decorator;

@ccclass('MainPaneController')
export class MainPaneController extends BaseController {
  purseModel: PurseModel = PurseModel.getInstance();
  private static instance: MainPaneController;

  @property(MainPaneView)
  private mainView: MainPaneView;

  private static creatingPromise: Promise<MainPaneController> | null = null;

  public static openUi() {
    this.getInstance().then((controller) => {
      if (controller.mainView) {
        controller.mainView.display();
      }
    });
  }

  protected start(): void {
    this.initView(this.mainView);
    this.bindViewToModel();
  }

  private bindViewToModel() {
    if (!this.view) return;

    // 绑定钻石数据
    this.purseModel.onDiamondChange((value) => {
      if (this.view?.diamondLabel) {
        this.view.diamondLabel.string = value.toString();
      }
    });

    // 绑定金币数据
    this.purseModel.onGoldChange((value) => {
      if (this.view?.goldLabel) {
        this.view.goldLabel.string = value.toString();
      }
    });

    // 初始化显示
    this.view.diamondLabel.string = this.purseModel.diamond.toString();
    this.view.goldLabel.string = this.purseModel.gold.toString();
  }

  // 更新数据的方法
  public updateDiamond(value: number) {
    this.purseModel.diamond = value;
  }

  public updateGold(value: number) {
    this.purseModel.gold = value;
  }

  // 获取数据的方法
  public getDiamond(): number {
    return this.purseModel.diamond;
  }

  public getGold(): number {
    return this.purseModel.gold;
  }

  private static getInstance(): Promise<MainPaneController> {
    if (this.instance) {
      return Promise.resolve(this.instance);
    }
    if (this.creatingPromise) {
      return this.creatingPromise;
    }
    this.creatingPromise = new Promise((resolve) => {
      UiViewFactory.createUi(R.Prefabs.MainPane, LayerIdx.layer1, (ui: MainPaneController) => {
        this.instance = ui;
        this.creatingPromise = null;
        resolve(ui);
      });
    });
    return this.creatingPromise;
  }
}
