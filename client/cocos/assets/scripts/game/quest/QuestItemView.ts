import { _decorator, instantiate, Label, Node, Prefab, Slider } from 'cc';
import { BaseUiView } from '../../frame/mvc/BaseUiView';
import { RedDotComponent } from '../../frame/reddot/RedDotCompoent';
import { RedDotManager } from '../../frame/reddot/RedDotManager';
import { MailVo } from '../../net/protocol/items/MailVo';
import { TimeUtils } from '../../utils/TimeUtils';
import { RewardItem } from '../reward/RewardItem';
import QuestPanelController from './QuestPanelController';
import QuestVo from '../../net/protocol/items/QuestVo';
import { ConfigContext } from '../../data/config/container/ConfigContext';
import { RewardUtil } from '../reward/RewardUtil';
const { ccclass, property } = _decorator;


@ccclass('QuestItemView')
export class QuestItemView extends BaseUiView {
  @property(Label)
  title: Label;

  @property(Label)
  progressTxt: Label;

  @property(Slider)
  progressSlider: Slider;

  @property(Node)
  rewardContainer: Node;

  @property(Prefab)
  rewardItemPrefab: Prefab;

  @property(Node)
  redDot: Node;

  protected start(): void {
  }


  public fillData(quest: QuestVo): void {
    this.rewardContainer.children.forEach((child) => {
      child.destroy();
    });
    let rewardItemSize = { width: 50, height: 50 };
    let questData = ConfigContext.configQuestContainer.getRecord(quest.id);
    let rewards = RewardUtil.parseReward(questData.rewards);
    for (const reward of rewards) {
      const rewardItem = instantiate(this.rewardItemPrefab);
      rewardItem.setParent(this.rewardContainer);
      rewardItem.getComponent(RewardItem).fillData(reward, rewardItemSize);
    }
    this.title.string = questData.desc;
    this.progressTxt.string = `${quest.progress}/${questData.target}`;
    this.progressSlider.progress = quest.progress / Number(questData.target);
  }
}
