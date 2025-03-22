import { _decorator, Component, Label, Node, Sprite, SpriteFrame, UITransform } from 'cc';
import ItemData from '../../data/config/model/ItemData';
import AssetLoader from '../../ui/AssertLoader';
import { UIViewController } from '../../ui/UiViewController';
import { BagItemInfo } from './BagItemInfo';
import { Item } from '../../data/user/Bagpack';
import ConfigItemContainer from '../../data/config/container/ConfigItemContainer';
const { ccclass, property } = _decorator;

@ccclass('BagItem')
export class BagItem extends UIViewController {
  @property(Sprite)
  public kuang: Sprite;

  @property(Label)
  public itemName: Label;

  @property(Label)
  public amout: Label;

  @property(Node)
  public icon: Node;

  private itemData: ItemData;

  protected start(): void {
    this.registerClickEvent(this.kuang.node, this.showItemDetail, this);
  }

  showItemDetail() {
    let root = this.node.parent.parent.parent.children[2];
    root.active = true;

    let itemInfo = root.getComponent(BagItemInfo);
    itemInfo.fillData(this.itemData);
  }

  public fillData(item: Item) {
    let itemContianer: ConfigItemContainer = ConfigItemContainer.getInstance();
    this.itemData = itemContianer.getRecord(item.id);
    this.itemName.string = this.itemData.name;
    this.amout.string = 'X' + item.count;

    // 使用AssetLoader加载图片资源
    const iconPath = `picture/item/${this.itemData.icon}`;

    // 先确保icon节点上有Sprite组件
    let iconSprite = this.icon.getComponent(Sprite);
    if (!iconSprite) {
      console.log('Adding Sprite component to icon node');
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
            console.log(`Resized sprite to: ${nodeWidth}x${nodeHeight}`);
          }
        }
      } else {
        console.error('加载图片资源失败:', err);
      }
    });
  }
}
