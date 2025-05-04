import { _decorator, EventKeyboard, Input, input, KeyCode, Label, Node } from 'cc';

import { BaseUiView } from '../../ui/BaseUiView';
import { BagPanelController } from '../item/BagPanelController';

import PlayerData from '../../data/user/PlayerData';
import AssetResourceFactory from '../../ui/AssetResourceFactory';
import R from '../../ui/R';
import { UiUtil } from '../../ui/UiUtil';
import { NumberUtils } from '../../utils/NumberUtils';
import { GmPaneController } from '../gm/GmPaneController';
import { HeroMainPaneController } from '../hero/HeroMainPaneController';
import { MailPaneController } from '../mail/MailPaneController';
import { RecruitPaneController } from '../recruit/RecruitPaneController';
const { ccclass, property } = _decorator;

@ccclass('MainPaneView')
export class MainPaneView extends BaseUiView {
  @property(Node)
  bagPane: Node;

  @property(Node)
  recruitePane: Node;

  @property(Node)
  mailPane: Node;

  @property(Node)
  heroPane: Node;

  @property(Node)
  mainPane: Node;

  @property(Label)
  diamondLabel: Label;

  @property(Label)
  goldLabel: Label;

  @property(Node)
  headIcon: Node;

  @property(Label)
  fightingLabel: Label;

  @property(Label)
  nameLabel: Label;

  private _isAltPressed: boolean = false;
  private _isGPressed: boolean = false;
  private _isMPressed: boolean = false;

  protected start(): void {
    this.registerClickEvent(this.bagPane, this.onBagClick, this);
    this.registerClickEvent(this.recruitePane, this.onRecruiteClick, this);
    this.registerClickEvent(this.heroPane, this.onHeroClick, this);
    this.registerClickEvent(this.mainPane, this.onMainClick, this);
    this.registerClickEvent(this.mailPane, this.onMailClick, this);

    this.nameLabel.string = PlayerData.instance.name;
    this.fightingLabel.string = '战力：' + NumberUtils.formatNumber(PlayerData.instance.fighting);
    let campSpriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Item);
    UiUtil.fillSpriteContent(
      this.headIcon,
      campSpriteAtlas.getSpriteFrame(this.getHeadIconByCamp(PlayerData.instance.camp))
    );

    // 监听键盘 同时按下alt+g+m 打开gm命令
    input.on(Input.EventType.KEY_DOWN, this.onKeyDown, this);
    input.on(Input.EventType.KEY_UP, this.onKeyUp, this);
  }

  private onKeyDown(event: EventKeyboard) {
    if (event.keyCode === KeyCode.ALT_LEFT || event.keyCode === KeyCode.ALT_RIGHT) {
      this._isAltPressed = true;
    }
    if (event.keyCode === KeyCode.KEY_G) {
      this._isGPressed = true;
    }
    if (event.keyCode === KeyCode.KEY_M) {
      this._isMPressed = true;
    }

    // 检查是否所有需要的键都被按下
    if (this._isAltPressed && this._isGPressed && this._isMPressed) {
      GmPaneController.switchDisplay();
    }
  }

  private onKeyUp(event: EventKeyboard) {
    if (event.keyCode === KeyCode.ALT_LEFT || event.keyCode === KeyCode.ALT_RIGHT) {
      this._isAltPressed = false;
    }
    if (event.keyCode === KeyCode.KEY_G) {
      this._isGPressed = false;
    }
    if (event.keyCode === KeyCode.KEY_M) {
      this._isMPressed = false;
    }
  }

  private getHeadIconByCamp(camp: number): string {
    return 1000 + camp + '';
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

  onMailClick() {
    MailPaneController.openUi();
  }
}
