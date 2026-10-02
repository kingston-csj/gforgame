import { _decorator, EditBox, Node, Prefab, ScrollView } from 'cc';
const { ccclass, property } = _decorator;
import { BaseUiView } from '../../frame/mvc/BaseUiView';
import { FriendApplyVo } from '../../net/protocol/friend/FriendApplyVo';
import { FriendApplyItemView } from './FriendApplyItemView';
import { instantiate } from 'cc';

@ccclass('FriendApplySubView')
export class FriendApplySubView extends BaseUiView {
    @property(ScrollView)
    scrollView: ScrollView = null!;
    @property(Prefab)
    applyItemPrefab: Prefab = null!;
    @property(Node)
    content: Node = null!;
    
    // ---------------- 向外抛出事件，由Controller赋值监听 ----------------
    /** 点击添加好友 */
    public onApplyFriendClick: ((uid: string, status: number) => void) | null = null;

    private _itemNodes: Node[] = [];

    protected start(): void {
        this.recycleAllItem();
    }

    /**
     * 【对外接口】由SubController调用，刷新好友列表
     */
    public setData(dataList: FriendApplyVo[]) {
        this.recycleAllItem();
        for (const item of dataList) {
            const node = this.createOneItem(item);
            this.content.addChild(node);
            this._itemNodes.push(node);
        }
    }

    private createOneItem(data: FriendApplyVo): Node {
        const node = instantiate(this.applyItemPrefab);
        // 假设item节点上组件叫 FriendItemComp，负责填充name/avatar
        const comp = node.getComponent(  FriendApplyItemView)!;
        // item内部按钮抛出事件
        this.registerClickEvent(comp.addBtn, () => {
            this.onApplyFriendClick?.(data.id, 1);
        }, this);
        this.registerClickEvent(comp.rejectBtn, () => {
            this.onApplyFriendClick?.(data.id, 2);
        }, this);
        comp.fillData(data);
        return node;
    }

    private recycleAllItem() {
        this.content.removeAllChildren();
        this._itemNodes.length = 0;
    }

    protected onHide(): void {
        // 子页签隐藏，清空列表
        this.recycleAllItem();
    }
}
