import { _decorator, Button, Label, Node } from 'cc';

import { ConfigContext } from '../../data/config/container/ConfigContext';
import GameContext from '../../GameContext';
import { HeroVo } from '../../net/MsgItems/HeroVo';
import { ReqHeroUpLevel } from '../../net/ReqHeroUpLevel';
import { ResHeroUpLevel } from '../../net/ResHeroUpLevel';
import AssetResourceFactory from '../../ui/AssetResourceFactory';
import { BaseUiView } from '../../ui/BaseUiView';
import R from '../../ui/R';
import { UiUtil } from '../../ui/UiUtil';
import { TipsPaneController } from '../common/TipsPaneController';
import { PurseModel } from '../main/PurseModel';
import { HeroBoxModel } from './HeroBoxModel';
import { HeroDetailController } from './HeroDetailController';
const { ccclass, property } = _decorator;

@ccclass('HeroDetailView')
export class HeroDetailView extends BaseUiView {
  @property(Node)
  public touchArea: Node;

  @property(Label)
  public heroName: Label;

  @property(Node)
  public heroIcon: Node;

  @property(Label)
  public skill1Name: Label;

  @property(Label)
  public skill2Name: Label;

  @property(Label)
  public skill3Name: Label;

  @property(Label)
  public skill4Name: Label;

  @property(Node)
  public skill1Btn: Node;

  @property(Node)
  public skill2Btn: Node;

  @property(Node)
  public skill3Btn: Node;

  @property(Node)
  public skill4Btn: Node;

  private skillNameGroup: Label[];

  private skillBtnGroup: Node[];

  @property(Label)
  public attr1Label: Label;

  @property(Label)
  public attr2Label: Label;

  @property(Label)
  public attr3Label: Label;

  @property(Label)
  public attr4Label: Label;

  @property(Label)
  public skillDesc: Label;

  @property(Node)
  public skillDescPanel: Node;

  @property(Node)
  public upLevelBtn: Node;

  @property(Label)
  public upLevelInfo: Label;

  @property(Label)
  public goldNum: Label;

  @property(Label)
  public fightNum: Label;

  @property(Label)
  public levelNum: Label;

  private hero: HeroVo;

  protected start(): void {
    this.touchArea.on(Button.EventType.CLICK, () => {
      HeroDetailController.closeUi();
    });
    this.skillNameGroup = [this.skill1Name, this.skill2Name, this.skill3Name, this.skill4Name];
    this.skillBtnGroup = [this.skill1Btn, this.skill2Btn, this.skill3Btn, this.skill4Btn];

    PurseModel.getInstance().onGoldChange((value) => {
      this.goldNum.string = value.toString();
      this.updateUpLevelBtn(HeroBoxModel.getInstance().calcUpLevel(this.hero));
    });
  }

  public updateCurrentHeroData() {
    if (!this.hero) {
      return;
    }
    this.fillData(this.hero);
  }

  public fillData(hero: HeroVo) {
    this.hero = hero;
    let heroData = ConfigContext.configHeroContainer.getRecord(hero.id);
    this.heroName.string = heroData.name;

    let heroSpriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Hero);
    // 设置UITransform的contentSize为原始图片尺寸
    UiUtil.fillSpriteContent(this.heroIcon, heroSpriteAtlas.getSpriteFrame(heroData.icon));

    // 技能组
    this.skillNameGroup.forEach((label, index) => {
      let skills = heroData.skills.split(';');
      let skillData = ConfigContext.configSkillContainer.getSkill(Number(skills[index]));
      label.string = skillData.name;
    });

    this.skillBtnGroup.forEach((btn, index) => {
      // 先移除已存在的事件监听器
      btn.off(Button.EventType.CLICK);
      // 再添加新的事件监听器
      btn.on(Button.EventType.CLICK, () => {
        this.skillDescPanel.active = true;
        let skills = heroData.skills.split(';');
        let skillData = ConfigContext.configSkillContainer.getSkill(Number(skills[index]));
        this.skillDesc.string = skillData.tips;
      });
    });

    this.attr1Label.string = hero.attrBox.getHp().toString();
    this.attr2Label.string = hero.attrBox.getAttack().toString();
    this.attr3Label.string = hero.attrBox.getDefense().toString();
    this.attr4Label.string = hero.attrBox.getSpeed().toString();

    this.fightNum.string = hero.fight.toString();
    this.levelNum.string = hero.level.toString();

    this.upLevelBtn.off(Button.EventType.CLICK);
    this.upLevelBtn.on(Button.EventType.CLICK, () => {
      let canUpLevel = HeroBoxModel.getInstance().calcUpLevel(this.hero);
      GameContext.instance.WebSocketClient.sendMessage(
        ReqHeroUpLevel.cmd,
        {
          heroId: this.hero.id,
          toLevel: this.hero.level + canUpLevel,
        },
        (msg: ResHeroUpLevel) => {
          if (msg.code > 0) {
            TipsPaneController.openUi(msg.code);
          } else {
            // 数据模型会触发界面更新
            this.levelNum.string = (this.hero.level + canUpLevel).toString();
          }
        }
      );
    });
    this.updateUpLevelBtn(HeroBoxModel.getInstance().calcUpLevel(this.hero));
    this.goldNum.string = PurseModel.getInstance().gold.toString();
  }

  public onHide(): void {
    this.skillDescPanel.active = false;
  }

  private updateUpLevelBtn(times: number) {
    if (!this.hero) {
      return;
    }

    if (times > 1) {
      this.upLevelInfo.string = `升${times}级`;
    } else {
      this.upLevelInfo.string = `升级`;
    }
  }
}
