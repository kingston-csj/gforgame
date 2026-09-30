import { _decorator, Node, Prefab, ScrollView } from 'cc';
const { ccclass, property } = _decorator;
import { BaseUiView } from '../../frame/mvc/BaseUiView';
import { FriendItem } from './FriendModel';
import { FriendListItemView } from './FriendListItemView';
import { instantiate } from 'cc';

@ccclass('FriendListSubView')
export class FriendListSubView extends BaseUiView {
    @property(ScrollView)
    scrollView: ScrollView = null!;
    @property(Prefab)
    friendItemPrefab: Prefab = null!;
    @property(Node)
    content: Node = null!;

    // ---------------- 向外抛出事件，由Controller赋值监听 ----------------
    /** 点击私聊 */
    public onChatClick: ((uid: string) => void) | null = null;
    /** 点击删除好友 */
    public onDeleteClick: ((uid: string) => void) | null = null;

    private _itemNodes: Node[] = [];

    /**
     * 【对外接口】由SubController调用，刷新好友列表
     */
    public setData(dataList: FriendItem[]) {
        this.recycleAllItem();
        for (const item of dataList) {
            const node = this.createOneItem(item);
            this.content.addChild(node);
            this._itemNodes.push(node);
        }
    }

    private createOneItem(data: FriendItem): Node {
        const node = instantiate(this.friendItemPrefab);
        // 假设item节点上组件叫 FriendItemComp，负责填充name/avatar
        const comp = node.getComponent( FriendListItemView)!;
        // comp.setData(data);

        // item内部按钮抛出事件
        this.registerClickEvent(comp.chatBtnNode, () => {
            this.onChatClick?.(data.id);
        }, this);
        this.registerClickEvent(comp.deleteBtnNode, () => {
            this.onDeleteClick?.(data.id);
        }, this);
        return node;
    }

    private recycleAllItem() {
        for (const n of this._itemNodes) {
            n.destroy();
        }
        this._itemNodes.length = 0;
    }

    protected onHide(): void {
        // 子页签隐藏，清空列表
        this.recycleAllItem();
    }
}
