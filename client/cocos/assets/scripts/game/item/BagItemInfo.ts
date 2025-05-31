import { _decorator, Component, Label, Node } from 'cc';
import ItemData from '../../data/config/model/ItemData';
import AssetResourceFactory from '../../ui/AssetResourceFactory';
import R from '../../ui/R';
import { UiUtil } from '../../utils/UiUtil';
const { ccclass, property } = _decorator;

@ccclass('BagItemInfo')
export class BagItemInfo extends Component {
  @property(Node)
  public icon: Node;

  @property(Label)
  public itemName: Label;

  @property(Label)
  public itemDesc: Label;

  public fillData(item: ItemData) {
    this.itemName.string = item.name;
    this.itemDesc.string = item.tips;

    let spriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(R.Sprites.Item);
    UiUtil.fillSpriteContent(this.icon, spriteAtlas.getSpriteFrame(item.icon));
  }
}
