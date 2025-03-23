import { _decorator, instantiate, Node, Prefab, Sprite, Toggle } from 'cc';
import { UIViewController } from '../../ui/UiViewController';
import UiView from '../../ui/UiView';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import ConfigItemContainer from '../../data/config/container/ConfigItemContainer';
import ItemData from '../../data/config/model/ItemData';
import { BagItem } from './BagItem';
import { Item } from '../../data/user/Bagpack';
import GameContext from '../../GameContext';
const { ccclass, property } = _decorator;

@ccclass('BagPanelController')
export class BagPanelController extends UIViewController {
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

  private static instance: BagPanelController;

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

  public static openUi() {
    if (BagPanelController.instance) {
      BagPanelController.instance.display();
    } else {
      BagPanelController.instance = new BagPanelController();

      UiView.createUi(R.bagPane, LayerIdx.layer5, (ui: UIViewController) => {
        BagPanelController.instance = ui as BagPanelController;
        ui.display();
      });
    }
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
    let itemContianer: ConfigItemContainer = ConfigItemContainer.getInstance();
    let items: Array<Item> = GameContext.instance.playerData.Bagpack.getItems();
    let filterItems: Array<Item> = [];
    if (this.selectedType === 0) {
      filterItems = items;
    } else {
      for (let i = 0; i < items.length; i++) {
        let itemData: ItemData = itemContianer.getRecord(items[i].id);
        if (itemData.type === this.selectedType) {
          filterItems.push(items[i]);
        }
      }
    }
    return filterItems;
  }
}
