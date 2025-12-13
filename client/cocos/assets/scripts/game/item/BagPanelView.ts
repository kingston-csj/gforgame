import { _decorator, instantiate, Node, Prefab, Toggle } from "cc";
import { ConfigContext } from "../../data/config/container/ConfigContext";
import ConfigItemContainer from "../../data/config/container/ConfigItemContainer";
import ItemData from "../../data/config/model/ItemData";
import { BaseUiView } from "../../frame/mvc/BaseUiView";
import { BagItem } from "./BagItem";
import BagpackModel, { Item } from "./BagpackModel";
const { ccclass, property } = _decorator;

@ccclass("BagPanelView")
export class BagPanelView extends BaseUiView {
  @property(Prefab)
  public itemPrefab: Prefab;

  @property(Node)
  public itemContainer: Node;

  @property(Node)
  public allBtn: Node | null = null;

  @property(Node)
  public suiPianBtn: Node | null = null;

  @property(Node)
  public giftBtn: Node | null = null;

  @property(Node)
  public materialBtn: Node | null = null;

  // 道具详情面板
  @property(Node)
  public detailContainer: Node | null = null;

  // 0: 全部, 1: 礼包, 2: 材料, 3: 碎片
  private selectedType: number = 0;

  public start(): void {
    this.registerToggleButtonEvent(this.allBtn, 0);
    this.registerToggleButtonEvent(this.giftBtn, 1);
    this.registerToggleButtonEvent(this.materialBtn, 2);
    this.registerToggleButtonEvent(this.suiPianBtn, 3);
  }

  private registerToggleButtonEvent(button: Node, type: number): void {
    this.registerClickEvent(
      button,
      () => {
        this.detailContainer.active = false;
        this.selectedType = type;
        this.showItems();
      },
      this
    );
  }

  protected onDisplay() {
    this.selectedType = 0;
    this.allBtn.getComponent(Toggle).isChecked = true;
    this.showItems();
  }

  private showItems() {
    // 先清空itemContainer
    this.itemContainer.children.forEach((child) => {
      child.destroy();
    });

    let items: Array<Item> = this.getAllItems();
    for (let i = 0; i < items.length; i++) {
      let item = instantiate(this.itemPrefab);
      item.setParent(this.itemContainer);
      item.getComponent(BagItem).fillData(items[i]);
    }
  }

  private getAllItems(): Array<Item> {
    let itemContianer: ConfigItemContainer = ConfigContext.configItemContainer;
    let items: Array<Item> = BagpackModel.getInstance().getItems();
    let filterItems: Array<Item> = [];
    if (this.selectedType === 0) {
      filterItems = items;
    } else {
      for (let i = 0; i < items.length; i++) {
        let itemData: ItemData = itemContianer.getRecord(items[i].cf_id);
        if (itemData.type === this.selectedType) {
          filterItems.push(items[i]);
        }
      }
    }
    return filterItems.filter((item) => item.count > 0);
  }
}
