import { _decorator, EditBox, Node, Prefab, ScrollView } from 'cc';
const { ccclass, property } = _decorator;
import { BaseUiView } from '../../frame/mvc/BaseUiView';
import {  FriendModel } from './FriendModel';
import { FriendVo } from '../../net/protocol/friend/FriendVo';
import { FriendSearchItemView } from './FriendSearchItemView';
import { instantiate } from 'cc';

@ccclass('FriendSearchSubView')
export class FriendSearchSubView extends BaseUiView {
    @property(ScrollView)
    scrollView: ScrollView = null!;
    @property(Prefab)
    searchItemPrefab: Prefab = null!;
    @property(Node)
    content: Node = null!;

    
    @property(EditBox)
    searchBox: EditBox;

    @property(Node)
    searchBtn: Node;

    // ---------------- 向外抛出事件，由Controller赋值监听 ----------------
    /** 点击添加好友 */
    public onAddFriendClick: ((uid: string) => void) | null = null;

    private _itemNodes: Node[] = [];

    protected start(): void {
        this.recycleAllItem();
        this.registerClickEvent(this.searchBtn, () => {
           FriendModel.getInstance().requestSearchFriend(this.searchBox.string)
        }, this);
    }

    /**
     * 【对外接口】由SubController调用，刷新好友列表
     */
    public setData(dataList: FriendVo[]) {
        this.recycleAllItem();
        for (const item of dataList) {
            const node = this.createOneItem(item);
            this.content.addChild(node);
            this._itemNodes.push(node);
        }
    }

    private createOneItem(data: FriendVo): Node {
        const node = instantiate(this.searchItemPrefab);
        // 假设item节点上组件叫 FriendItemComp，负责填充name/avatar
        const comp = node.getComponent(  FriendSearchItemView)!;
        // item内部按钮抛出事件
        this.registerClickEvent(comp.addBtn, () => {
            this.onAddFriendClick?.(data.id);
        }, this);
        comp.fillData(data);
        return node;
    }

    private recycleAllItem() {
        // for (const n of this._itemNodes) {
        //     n.destroy();
        // }
        this.content.removeAllChildren();
        this._itemNodes.length = 0;
    }

    protected onHide(): void {
        // 子页签隐藏，清空列表
        this.recycleAllItem();
    }
}
