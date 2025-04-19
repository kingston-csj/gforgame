import { _decorator, Button, Color, Label, Node, Sprite } from 'cc';

import { ConfigContext } from '../../data/config/container/ConfigContext';
import { HeroVo } from '../../net/protocol/MsgItems/HeroVo';
import AssetResourceFactory from '../../ui/AssetResourceFactory';
import { BaseUiView } from '../../ui/BaseUiView';
import R from '../../ui/R';
import { UiUtil } from '../../ui/UiUtil';
import { TipsPaneController } from '../common/TipsPaneController';
import { GameConstants } from '../GameConstants';
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

  @property(Node)
  public fightArea: Node;

  @property(Node)
  public attrGroup: Node;

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

  @property(Node)
  public upLevelGroup: Node;

  @property(Node)
  public upStageGroup: Node;

  @property(Node)
  public upStageBtn: Node;

  @property(Label)
  public upLevelInfo: Label;

  @property(Label)
  public goldNum: Label;

  @property(Label)
  public upstageItemNum: Label;

  @property(Label)
  public fightNum: Label;

  @property(Label)
  public levelNum: Label;

  @property(Label)
  public description: Label;

  @property(Node)
  camp: Node;

  private hero: HeroVo;

  protected start(): void {
    this.touchArea.on(Button.EventType.CLICK, () => {
      HeroDetailController.closeUi();
    });
    this.skillNameGroup = [this.skill1Name, this.skill2Name, this.skill3Name, this.skill4Name];
    this.skillBtnGroup = [this.skill1Btn, this.skill2Btn, this.skill3Btn, this.skill4Btn];

    PurseModel.getInstance().onChange('gold', (value) => {
      this.goldNum.string = value.toString();
      this.updateLevelOrStageBtn();
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

    let campSpriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Camp);
    UiUtil.fillSpriteContent(this.camp, campSpriteAtlas.getSpriteFrame('camp_' + heroData.camp));

    let skills = heroData.skills.split(';');
    // 技能组
    this.skillNameGroup.forEach((label, index) => {
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

      let skillData = ConfigContext.configSkillContainer.getSkill(Number(skills[index]));
      if (hero.stage < skillData.stage) {
        // 设置按钮为灰色
        btn.getComponent(Sprite).color = Color.GRAY;
      }
    });

    // 主公不显示战斗力和属性面板
    if (heroData.quality === 0) {
      this.fightArea.active = false;
      this.attrGroup.active = false;
      this.description.string = heroData.tips;
      this.description.node.active = true;
    } else {
      this.description.node.active = false;
      this.fightArea.active = true;
      this.attrGroup.active = true;
      this.attr1Label.string = hero.attrBox.getHp().toString();
      this.attr2Label.string = hero.attrBox.getAttack().toString();
      this.attr3Label.string = hero.attrBox.getDefense().toString();
      this.attr4Label.string = hero.attrBox.getSpeed().toString();
    }

    this.fightNum.string = hero.fight.toString();
    this.levelNum.string = hero.level.toString();

    this.upLevelBtn.off(Button.EventType.CLICK);
    this.upLevelBtn.on(Button.EventType.CLICK, () => {
      let canUpLevel = HeroBoxModel.getInstance().calcUpLevel(this.hero);
      if (canUpLevel <= 0) {
        TipsPaneController.showI18nContent(GameConstants.I18N.TIPS_2001);
        return;
      }

      HeroBoxModel.getInstance()
        .requestUpLevel(this.hero.id, this.hero.level + canUpLevel)
        .then((code) => {
          if (code === 0) {
            // 数据模型会触发界面更新
            this.levelNum.string = (this.hero.level + canUpLevel).toString();
            this.updateLevelOrStageBtn();
          } else {
            TipsPaneController.showI18nContent(code);
          }
        });
    });

    this.goldNum.string = PurseModel.getInstance().gold.toString();

    this.upStageBtn.off(Button.EventType.CLICK);
    this.upStageBtn.on(Button.EventType.CLICK, () => {
      HeroBoxModel.getInstance()
        .requestUpStage(this.hero.id)
        .then((code) => {
          if (code === 0) {
            this.updateLevelOrStageBtn();
          } else {
            TipsPaneController.showI18nContent(code);
          }
        });
    });

    this.updateLevelOrStageBtn();
  }

  private updateLevelOrStageBtn() {
    this.upLevelGroup.active = false;
    this.upStageGroup.active = false;
    let heroStageData = ConfigContext.configHeroStageContainer.getRecordByStage(this.hero.stage);
    let nextStageData = ConfigContext.configHeroStageContainer.getRecordByStage(
      this.hero.stage + 1
    );
    if (this.hero.level == heroStageData.max_level && nextStageData) {
      this.upStageGroup.active = true;
    } else {
      if (this.hero.level < ConfigContext.configHeroLevelContainer.getMaxLevel()) {
        this.upLevelGroup.active = true;
        let times = HeroBoxModel.getInstance().calcUpLevel(this.hero);
        if (times > 1) {
          this.upLevelInfo.string = `升${times}级`;
        } else {
          this.upLevelInfo.string = `升级`;
        }
        this.levelNum.string = this.hero.level.toString();
      }
    }
  }
  public onHide(): void {
    this.skillDescPanel.active = false;
  }
}
