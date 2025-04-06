import { _decorator } from 'cc';

import { BaseController } from '../../ui/BaseController';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import UiViewFactory from '../../ui/UiViewFactory';
import { HeroBoxModel } from './HeroBoxModel';
import { HeroDetailView } from './HeroDetailView';
const { ccclass, property } = _decorator;

@ccclass('HeroDetailController')
export class HeroDetailController extends BaseController {
  private static instance: HeroDetailController;

  private static creatingPromise: Promise<HeroDetailController> | null = null;

  @property(HeroDetailView)
  heroDetailView: HeroDetailView | null = null;

  private heroBoxModel: HeroBoxModel = HeroBoxModel.getInstance();

  protected start(): void {
    this.initView(this.heroDetailView);
    this.heroBoxModel.onHeroAttrChanged(() => {
      this.heroDetailView.updateCurrentHeroData();
    });
  }

  public static openUi(heroId: number) {
    this.getInstance().then((controller) => {
      if (controller.heroDetailView) {
        controller.heroDetailView.display();
        controller.heroDetailView.fillData(controller.heroBoxModel.getHero(heroId));
      }
    });
  }

  public static closeUi() {
    if (!this.instance) {
      return Promise.resolve();
    }
    return this.getInstance().then((controller) => {
      if (controller.heroDetailView) {
        controller.heroDetailView.hide();
      }
    });
  }

  private static getInstance(): Promise<HeroDetailController> {
    if (this.instance) {
      return Promise.resolve(this.instance);
    }
    if (this.creatingPromise) {
      return this.creatingPromise;
    }
    this.creatingPromise = new Promise((resolve) => {
      UiViewFactory.createUi(
        R.Prefabs.HeroDetailPane,
        LayerIdx.layer5,
        (ui: HeroDetailController) => {
          this.instance = ui;
          this.creatingPromise = null;
          resolve(ui);
        }
      );
    });
    return this.creatingPromise;
  }
}
