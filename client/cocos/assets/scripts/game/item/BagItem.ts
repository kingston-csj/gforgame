import { _decorator, Label, Node, Sprite } from 'cc';
import { ConfigContext } from '../../data/config/container/ConfigContext';
import ConfigItemContainer from '../../data/config/container/ConfigItemContainer';
import ItemData from '../../data/config/model/ItemData';
import { BaseUiView } from '../../frame/mvc/BaseUiView';
import AssetResourceFactory from '../../ui/AssetResourceFactory';
import R from '../../ui/R';
import { UiUtil } from '../../utils/UiUtil';
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
    let root = this.node.parent.parent.parent.parent.children[2];
    if (root) {
      root.active = true;
      let itemInfo = root.getComponent(BagItemInfo);
      itemInfo.fillData(this.itemData);
    }
  }

  public fillData(item: Item) {
    let itemContianer: ConfigItemContainer = ConfigContext.configItemContainer;
    this.itemData = itemContianer.getRecord(item.id);
    this.itemName.string = this.itemData.name;
    this.amout.string = 'X' + item.count;

    let spriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Item);
    UiUtil.fillSpriteContent(this.icon, spriteAtlas.getSpriteFrame(this.itemData.icon));
  }
}
