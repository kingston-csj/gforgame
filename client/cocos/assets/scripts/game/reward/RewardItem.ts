import { _decorator, Label, Node, Sprite, UITransform } from "cc";
import { ConfigContext } from "../../data/config/container/ConfigContext";
import ConfigItemContainer from "../../data/config/container/ConfigItemContainer";
import HeroData from "../../data/config/model/HeroData";
import ItemData from "../../data/config/model/ItemData";
import { BaseUiView } from "../../frame/mvc/BaseUiView";
import RewardVo from "../../net/protocol/items/RewardVo";
import AssetResourceFactory from "../../ui/AssetResourceFactory";
import R from "../../ui/R";
import { UiUtil } from "../../utils/UiUtil";

const { ccclass, property } = _decorator;

@ccclass("RewardItem")
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

  public fillData(item: RewardVo, size?: { width: number; height: number }) {
    const uiTrans = this.node.getComponent(UITransform);
    if (uiTrans && size) {
      // 计算缩放比例
      const scaleX = size.width / uiTrans.width;
      const scaleY = size.height / uiTrans.height;
      // 设置父节点尺寸
      uiTrans.setContentSize(size.width, size.height);
      // 设置所有子节点等比缩放
      for (const child of this.node.children) {
        child.setScale(scaleX, scaleY, 1);
      }
    }
    if (item.type == "item") {
      let itemContianer: ConfigItemContainer =
        ConfigContext.configItemContainer;
      let [id, count] = item.value.split("=");
      this.itemData = itemContianer.getRecord(parseInt(id));
      this.itemName.string = this.itemData.name;
      this.amout.string = "X" + parseInt(count);
      let spriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(
        R.Sprites.Item
      );
      UiUtil.fillSpriteContent(
        this.icon,
        spriteAtlas.getSpriteFrame(this.itemData.icon)
      );
    } else if (item.type == "hero") {
      this.heroData = ConfigContext.configHeroContainer.getRecord(
        parseInt(item.value)
      );
      this.itemName.string = this.heroData.name;
      let spriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(
        R.Sprites.Hero
      );
      UiUtil.fillSpriteContent(
        this.icon,
        spriteAtlas.getSpriteFrame(this.heroData.icon)
      );
    } else if (item.type == "currency") {
      // 在道具表，配置特殊货币道具，纯展示
      let [kind, count] = item.value.split("=");
      this.itemName.string = kind;
      this.amout.string = "X" + parseInt(count);
      if (kind == "gold") {
        let spriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(
          R.Sprites.Item
        );
        UiUtil.fillSpriteContent(this.icon, spriteAtlas.getSpriteFrame("9998"));
      } else if (kind == "diamond") {
        let spriteAtlas = AssetResourceFactory.instance.getSpriteAtlas(
          R.Sprites.Item
        );
        UiUtil.fillSpriteContent(this.icon, spriteAtlas.getSpriteFrame("9999"));
      }
    }
  }
}
