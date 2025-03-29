import { _decorator, Label, Node, Sprite, UITransform } from 'cc';
import ConfigHeroContainer from '../../data/config/container/ConfigHeroContainer';
import ConfigItemContainer from '../../data/config/container/ConfigItemContainer';
import HeroData from '../../data/config/model/HeroData';
import ItemData from '../../data/config/model/ItemData';
import RewardInfo from '../../net/MsgItems/RewardInfo';
import AssetResourceFactory from '../../ui/AssetResourceFactory';
import { BaseUiView } from '../../ui/BaseUiView';
import R from '../../ui/R';

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
    // 获取icon节点的UITransform组件，用于设置大小
    const iconTransform = this.icon.getComponent(UITransform);
    if (!iconTransform) {
      console.warn('Icon node has no UITransform component');
      return;
    }
    // 保存节点当前的尺寸，用于调整图像
    const originalIconWidth = iconTransform.width;
    const originalIconHeight = iconTransform.height;

    if (item.type == 'item') {
      let itemContianer: ConfigItemContainer = ConfigItemContainer.getInstance();
      let [id, count] = item.value.split(',');
      this.itemData = itemContianer.getRecord(parseInt(id));
      this.itemName.string = this.itemData.name;
      this.amout.string = 'X' + parseInt(count);

      let spriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Item);
      this.icon.getComponent(Sprite).spriteFrame = spriteAtlas.getSpriteFrame(this.itemData.icon);
    } else if (item.type == 'hero') {
      let heroContianer: ConfigHeroContainer = ConfigHeroContainer.getInstance();
      this.heroData = heroContianer.getRecord(parseInt(item.value));
      this.itemName.string = this.heroData.name;
      let spriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Hero);
      this.icon.getComponent(Sprite).spriteFrame = spriteAtlas.getSpriteFrame(this.heroData.icon);
    }

    // 获取当前SpriteFrame
    const sprite = this.icon.getComponent(Sprite);
    if (!sprite || !sprite.spriteFrame) {
      console.warn('Icon has no valid sprite frame');
      return;
    }
    // 设置UITransform的contentSize为原始图片尺寸
    iconTransform.setContentSize(originalIconWidth, originalIconHeight);
  }
}
