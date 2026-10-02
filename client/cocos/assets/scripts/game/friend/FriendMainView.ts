import { _decorator, Node } from 'cc';
const { ccclass, property } = _decorator;
import { BaseUiView } from '../../frame/mvc/BaseUiView';
import { FriendTabType } from './FriendTabType';
import { FriendListSubView } from './FriendListSubView';
import { FriendSearchSubView } from './FriendSearchSubView';
import { FriendApplySubView } from './FriendApplySubView';

@ccclass('FriendMainView')
export class FriendMainView extends BaseUiView {

    @property(FriendListSubView)
    friendListSubView: FriendListSubView = null;
    @property(FriendSearchSubView)
    friendSearchSubView: FriendSearchSubView = null;
    @property(FriendApplySubView)
    friendApplySubView: FriendApplySubView = null;

    @property(Node)
    listPanelTag: Node;      //好友列表页签容器
    @property(Node)
    searchPanelTag: Node;   //搜索页签容器
    @property(Node)
    applyPanelTag: Node;     //申请列表页签容器

    @property(Node)
    closeBtn: Node;

    //页签切换回调，由Controller注册
    public onTabChange: ((tab: FriendTabType) => void) | null = null;

    protected start(): void {
        this.registerClickEvent(this.listPanelTag, () => {
            this.onTabChange?.(FriendTabType.FriendList);
        }, this);

        this.registerClickEvent(this.searchPanelTag, () => {
            this.onTabChange?.(FriendTabType.FriendSearch);
        }, this);

        this.registerClickEvent(this.applyPanelTag, () => {
            this.onTabChange?.(FriendTabType.FriendApply);
        }, this);

        this.registerClickEvent(this.closeBtn, this.hide, this);
    }

    /**
     * 【对外接口】Controller调用，设置当前激活哪个页签，只操作UI
     */
    public setActiveTab(tab: FriendTabType) {
        this.friendListSubView.node.active = tab === FriendTabType.FriendList;
        this.friendSearchSubView.node.active = tab === FriendTabType.FriendSearch;
        this.friendApplySubView.node.active = tab === FriendTabType.FriendApply;
    }

    protected onDisplay(): void {
        //打开面板，默认切好友列表，UI渲染交给view接口
        this.setActiveTab(FriendTabType.FriendList);
    }

    protected onHide(): void {
        //面板关闭，可以清空部分UI状态
    }
}
