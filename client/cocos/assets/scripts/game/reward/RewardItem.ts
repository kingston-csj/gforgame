import { _decorator, Label, Node, Sprite } from 'cc';
import { ConfigContext } from '../../data/config/container/ConfigContext';
import ConfigItemContainer from '../../data/config/container/ConfigItemContainer';
import HeroData from '../../data/config/model/HeroData';
import ItemData from '../../data/config/model/ItemData';
import RewardInfo from '../../net/MsgItems/RewardInfo';
import AssetResourceFactory from '../../ui/AssetResourceFactory';
import { BaseUiView } from '../../ui/BaseUiView';
import R from '../../ui/R';
import { UiUtil } from '../../ui/UiUtil';

const { ccclass, property } = _decorator;

@ccclass('RewardItem')
export class RewardItem extends BaseUiView {
  @property(Sprite)
  public kuang: Sprite;

  @property(Label)
  public itemName: Label;

  @property(Label)
  public amout: Label;

  @property(Node)
  public icon: Node;

  private itemData: ItemData;

  private heroData: HeroData;

  public fillData(item: RewardInfo) {
    if (item.type == 'item') {
      let itemContianer: ConfigItemContainer = ConfigContext.configItemContainer;
      let [id, count] = item.value.split(',');
      this.itemData = itemContianer.getRecord(parseInt(id));
      this.itemName.string = this.itemData.name;
      this.amout.string = 'X' + parseInt(count);

      let spriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Item);
      UiUtil.fillSpriteContent(this.icon, spriteAtlas.getSpriteFrame(this.itemData.icon));
    } else if (item.type == 'hero') {
      this.heroData = ConfigContext.configHeroContainer.getRecord(parseInt(item.value));
      this.itemName.string = this.heroData.name;
      let spriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Hero);
      UiUtil.fillSpriteContent(this.icon, spriteAtlas.getSpriteFrame(this.heroData.icon));
    }
  }
}
