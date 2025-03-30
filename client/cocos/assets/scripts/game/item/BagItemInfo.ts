import { _decorator, Component, Label, Sprite, UITransform } from 'cc';
import ItemData from '../../data/config/model/ItemData';
import AssetResourceFactory from '../../ui/AssetResourceFactory';
import R from '../../ui/R';
const { ccclass, property } = _decorator;

@ccclass('BagItemInfo')
export class BagItemInfo extends Component {
  @property(Sprite)
  public icon: Sprite;

  @property(Label)
  public itemName: Label;

  @property(Label)
  public itemDesc: Label;

  public fillData(item: ItemData) {
    // 获取icon节点的UITransform组件，用于设置大小
    const iconTransform = this.icon.getComponent(UITransform);
    if (!iconTransform) {
      console.warn('Icon node has no UITransform component');
      return;
    }
    // 保存节点当前的尺寸，用于调整图像
    const originalIconWidth = iconTransform.width;
    const originalIconHeight = iconTransform.height;

    this.itemName.string = item.name;
    this.itemDesc.string = item.tips;

    let spriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Item);
    this.icon.getComponent(Sprite).spriteFrame = spriteAtlas.getSpriteFrame(item.icon);
    // 获取当前SpriteFrame
    const sprite = this.icon.getComponent(Sprite);
    if (!sprite || !sprite.spriteFrame) {
      console.warn('Icon has no valid sprite frame');
      return;
    }

    // 设置UITransform的contentSize为原始图片尺寸
    iconTransform.setContentSize(originalIconWidth, originalIconHeight);

    // // 使用AssetLoader加载图片资源
    // const iconPath = `picture/item/${item.icon}`;

    // // 先确保icon节点上有Sprite组件
    // let iconSprite = this.icon.getComponent(Sprite);
    // if (!iconSprite) {
    //   console.log('Adding Sprite component to icon node');
    //   iconSprite = this.icon.addComponent(Sprite);
    // }

    // // 获取icon节点的UITransform组件，用于设置大小
    // const iconTransform = this.icon.getComponent(UITransform);
    // if (!iconTransform) {
    //   console.warn('Icon node has no UITransform component');
    //   return;
    // }

    // // 使用AssetLoader加载ImageAsset
    // AssetLoader.loadImageAsset(iconPath, (err, imageAsset) => {
    //   if (!err && imageAsset) {
    //     // 从ImageAsset创建SpriteFrame
    //     const spriteFrame = SpriteFrame.createWithImage(imageAsset);

    //     // 设置Sprite的缩放模式和类型，使图像自适应节点大小
    //     iconSprite.type = Sprite.Type.SIMPLE;
    //     iconSprite.sizeMode = Sprite.SizeMode.CUSTOM;

    //     // 将SpriteFrame设置到Sprite组件上
    //     iconSprite.spriteFrame = spriteFrame;

    //     // 保存节点当前的尺寸，用于调整图像
    //     const nodeWidth = iconTransform.width;
    //     const nodeHeight = iconTransform.height;

    //     // 图像要适应节点大小
    //     if (spriteFrame) {
    //       // 通过UITransform设置spriteFrame的大小为节点尺寸
    //       const uiTrans = iconSprite.node.getComponent(UITransform);
    //       if (uiTrans) {
    //         uiTrans.setContentSize(nodeWidth, nodeHeight);
    //         // console.log(`Resized sprite to: ${nodeWidth}x${nodeHeight}`);
    //       }
    //     }
    //   } else {
    //     console.error('加载图片资源失败:', err);
    //   }
    // });
  }
}
