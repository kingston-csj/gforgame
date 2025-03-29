import { _decorator, Label, Node, Sprite, UITransform } from 'cc';
import ConfigItemContainer from '../../data/config/container/ConfigItemContainer';
import ItemData from '../../data/config/model/ItemData';
import AssetResourceFactory from '../../ui/AssetResourceFactory';
import { BaseUiView } from '../../ui/BaseUiView';
import R from '../../ui/R';
import { BagItemInfo } from './BagItemInfo';
import { Item } from './BagpackModel';
const { ccclass, property } = _decorator;

@ccclass('BagItem')
export class BagItem extends BaseUiView {
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
    if (root) {
      root.active = true;
      let itemInfo = root.getComponent(BagItemInfo);
      itemInfo.fillData(this.itemData);
    }
  }

  public fillData(item: Item) {
    let itemContianer: ConfigItemContainer = ConfigItemContainer.getInstance();
    this.itemData = itemContianer.getRecord(item.id);
    this.itemName.string = this.itemData.name;
    this.amout.string = 'X' + item.count;

    const iconTransform = this.icon.getComponent(UITransform);
    if (!iconTransform) {
      console.warn('Icon node has no UITransform component');
      return;
    }
    // 保存节点当前的尺寸，用于调整图像
    const originalIconWidth = iconTransform.width;
    const originalIconHeight = iconTransform.height;

    let spriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Item);
    this.icon.getComponent(Sprite).spriteFrame = spriteAtlas.getSpriteFrame(this.itemData.icon);
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
