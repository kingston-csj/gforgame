import { _decorator, Label, Node } from 'cc';

import { BaseUiView } from '../../ui/BaseUiView';
import { BagPanelController } from '../item/BagPanelController';

import { HeroMainPaneController } from '../hero/HeroMainPaneController';
import { RecruitPaneController } from '../recruit/RecruitPaneController';
import { PurseModel } from './PurseModel';

const { ccclass, property } = _decorator;

@ccclass('MainPaneView')
export class MainPaneView extends BaseUiView {
  @property(Node)
  bagPane: Node;

  @property(Node)
  recruitePane: Node;

  @property(Node)
  heroPane: Node;

  @property(Node)
  mainPane: Node;

  @property(Label)
  diamondLabel: Label;

  @property(Label)
  goldLabel: Label;

  protected start(): void {
    this.registerClickEvent(this.bagPane, this.onBagClick, this);
    this.registerClickEvent(this.recruitePane, this.onRecruiteClick, this);
    this.registerClickEvent(this.heroPane, this.onHeroClick, this);
    this.registerClickEvent(this.mainPane, this.onMainClick, this);
    PurseModel.getInstance().onGoldChange((value) => {
      this.goldLabel.string = value.toString();
    });
    PurseModel.getInstance().onDiamondChange((value) => {
      this.diamondLabel.string = value.toString();
    });
  }

  onBagClick() {
    BagPanelController.openUi();
    RecruitPaneController.closeUi();
    HeroMainPaneController.closeUi();
  }

  onRecruiteClick() {
    RecruitPaneController.openUi();
    BagPanelController.closeUi();
    HeroMainPaneController.closeUi();
  }

  onHeroClick() {
    HeroMainPaneController.openUi();
    RecruitPaneController.closeUi();
    BagPanelController.closeUi();
  }

  onMainClick() {
    RecruitPaneController.closeUi();
    HeroMainPaneController.closeUi();
    BagPanelController.closeUi();
  }
}
