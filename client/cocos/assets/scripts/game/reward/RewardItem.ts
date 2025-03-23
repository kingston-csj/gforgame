import { _decorator, Component, Label, Node, Sprite, SpriteFrame, UITransform } from 'cc';
import ItemData from '../../data/config/model/ItemData';
import AssetLoader from '../../ui/AssertLoader';
import { UIViewController } from '../../ui/UiViewController';
import ConfigItemContainer from '../../data/config/container/ConfigItemContainer';
import RewardInfo from '../../net/MsgItems/RewardInfo';
import ConfigHeroContainer from '../../data/config/container/ConfigHeroContainer';
import HeroData from '../../data/config/model/HeroData';

const { ccclass, property } = _decorator;

@ccclass('RewardItem')
export class RewardItem extends UIViewController {
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
    let iconPath = '';
    if (item.type == 'item') {
      let itemContianer: ConfigItemContainer = ConfigItemContainer.getInstance();
      let [id, count] = item.value.split(',');
      this.itemData = itemContianer.getRecord(parseInt(id));
      this.itemName.string = this.itemData.name;
      this.amout.string = 'X' + parseInt(count);
      iconPath = `picture/item/${this.itemData.icon}`;
    } else if (item.type == 'hero') {
      iconPath = `picture/hero/${item.value}`;
      let heroContianer: ConfigHeroContainer = ConfigHeroContainer.getInstance();
      this.heroData = heroContianer.getRecord(parseInt(item.value));
      this.itemName.string = this.heroData.name;
    }
    // 使用AssetLoader加载图片资源

    let iconSprite = this.icon.getComponent(Sprite);
    if (!iconSprite) {
      iconSprite = this.icon.addComponent(Sprite);
    }

    // 获取icon节点的UITransform组件，用于设置大小
    const iconTransform = this.icon.getComponent(UITransform);
    if (!iconTransform) {
      console.warn('Icon node has no UITransform component');
      return;
    }

    // 使用AssetLoader加载ImageAsset
    AssetLoader.loadImageAsset(iconPath, (err, imageAsset) => {
      if (!err && imageAsset) {
        // 从ImageAsset创建SpriteFrame
        const spriteFrame = SpriteFrame.createWithImage(imageAsset);

        // 设置Sprite的缩放模式和类型，使图像自适应节点大小
        iconSprite.type = Sprite.Type.SIMPLE;
        iconSprite.sizeMode = Sprite.SizeMode.CUSTOM;
        // 将SpriteFrame设置到Sprite组件上
        iconSprite.spriteFrame = spriteFrame;
        // 保存节点当前的尺寸，用于调整图像
        const nodeWidth = iconTransform.width;
        const nodeHeight = iconTransform.height;

        // 图像要适应节点大小
        if (spriteFrame) {
          // 通过UITransform设置spriteFrame的大小为节点尺寸
          const uiTrans = iconSprite.node.getComponent(UITransform);
          if (uiTrans) {
            uiTrans.setContentSize(nodeWidth, nodeHeight);
          }
        }
      } else {
        console.error('加载图片资源失败:', err);
      }
    });
  }
}
