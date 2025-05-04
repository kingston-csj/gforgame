import { _decorator, instantiate, Label, Node, Prefab } from 'cc';
import { MailVo } from '../../net/protocol/items/MailVo';
import { BaseUiView } from '../../ui/BaseUiView';
import { TimeUtils } from '../../utils/TimeUtils';
import { RewardItem } from '../reward/RewardItem';
import { MailDetailPaneController } from './MailDetailPaneController';
const { ccclass, property } = _decorator;

@ccclass('MailItemView')
export class MailItemView extends BaseUiView {
  @property(Label)
  title: Label;

  @property(Label)
  expiredTime: Label;

  @property(Node)
  rewardContainer: Node;

  @property(Prefab)
  rewardItemPrefab: Prefab;

  private mailId: number = 0;

  protected start(): void {
    this.registerClickEvent(this.node, this.showMailDetail, this);
  }

  private showMailDetail(): void {
    MailDetailPaneController.openUi(this.mailId);
  }

  public fillData(mail: MailVo): void {
    this.mailId = mail.id;
    this.title.string = mail.title;
    let expiredTime = mail.time + 30 * TimeUtils.ONE_DAY;
    this.expiredTime.string = TimeUtils.getLeftTimeTips(expiredTime);
    this.rewardContainer.children.forEach((child) => {
      child.destroy();
    });
    let rewardItemSize = { width: 50, height: 50 };
    for (const reward of mail.rewards) {
      const rewardItem = instantiate(this.rewardItemPrefab);
      rewardItem.setParent(this.rewardContainer);
      rewardItem.getComponent(RewardItem).fillData(reward, rewardItemSize);
    }
  }
}
