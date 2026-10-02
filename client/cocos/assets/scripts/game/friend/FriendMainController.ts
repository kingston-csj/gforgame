
import { _decorator } from 'cc';
import { BaseController } from '../../frame/mvc/BaseController';
import { FriendMainView } from './FriendMainView';
import { FriendTabType } from './FriendTabType';
import { FriendModel } from './FriendModel';
import { LayerIdx } from '../../ui/LayerIds';
import R from '../../ui/R';
import UiViewFactory from '../../ui/UiViewFactory';
import { FriendApplyVo } from '../../net/protocol/friend/FriendApplyVo';
import { FriendVo } from '../../net/protocol/friend/FriendVo';



import { TipsPaneController } from '../common/TipsPaneController';

const { ccclass, property } = _decorator;

@ccclass('FriendMainController')
export class FriendMainController extends BaseController  {

    private static instance: FriendMainController;

    private _friendModel: FriendModel = FriendModel.getInstance();
    private _unsubList: Array<() => void> = [];


    private static creatingPromise: Promise<FriendMainController> | null = null;


    @property(FriendMainView)
    private mainView: FriendMainView = null;

    protected bindViewEvents(): void {
        if (!this.view) return;

        //监听view抛出的页签切换事件
        this.mainView.onTabChange = (tab) => {
            this.handlerTabSwitch(tab);
        }

        //订阅model数据变更
        const unsubFriendList = this._friendModel.onChange("friendList", (data) => {
            // 刷新【好友列表子页面】UI
        });
        this._unsubList.push(unsubFriendList);

        const unsubApplyList = this._friendModel.onChange("applyList", (data) => {
            //刷新【申请列表子页面】UI
             this.mainView.friendApplySubView.setData(data as FriendApplyVo[] ) 
        });
        this._unsubList.push(unsubApplyList);

        const unsubSearchResult = this._friendModel.onChange("searchResult", (data) => {
            //刷新【搜索页子页面】UI
            this.mainView.friendSearchSubView.setData(data as FriendVo[] ) 
        });
        this._unsubList.push(unsubSearchResult);

        this.mainView.friendSearchSubView.onAddFriendClick = (uid) => {
            this._friendModel.requestApplyFriend(uid).then((msg) => {
              if (msg.code == 0) {
                  TipsPaneController.showStringContent("申请成功");
                  return;
              }
              TipsPaneController.showI18nContent(msg.code);
            });
        }

        this.mainView.friendApplySubView.onApplyFriendClick = (applyId, status) => {
            this._friendModel.requestDealApplyFriend(applyId ,status).then((msg) => {
              if (msg.code == 0) {
                  TipsPaneController.showStringContent("操作成功");
                  return;
              }
              TipsPaneController.showI18nContent(msg.code);
            });
        }
    }

    /**
     * 处理页签切换业务
     */
    private handlerTabSwitch(tab: FriendTabType) {
        //1.通知view切换UI显示
        this.view!.setActiveTab(tab);

        //2.根据不同页签，触发数据加载
        switch (tab) {
            case FriendTabType.FriendList:
                this.mainView.friendListSubView.setData(this._friendModel.friendList);
                break;
            case FriendTabType.FriendApply:
                this._friendModel.requestApplyList();
                break;
            case FriendTabType.FriendSearch:
                //搜索页，不需要自动请求，等待输入关键词再搜索
                break;
        }
    }

    onDestroy() {
        //组件销毁，批量取消model监听，防止内存泄漏
        this._unsubList.forEach(unsub => unsub());
        this._unsubList.length = 0;
    }

    public static openUi() {
        this.getInstance().then((controller) => {
          if (controller.mainView) {
            controller.mainView.display();
          }
        });
      }
    
      public static closeUi() {
        if (!this.instance) {
          return Promise.resolve();
        }
        return this.getInstance().then((controller) => {
          if (controller.mainView) {
            controller.mainView.hide();
          }
        });
      }
    
      protected start(): void {
        this.initView(this.mainView);
      }
    
      private static getInstance(): Promise<FriendMainController> {
        if (this.instance) {
          return Promise.resolve(this.instance);
        }
        if (this.creatingPromise) {
          return this.creatingPromise;
        }
        this.creatingPromise = new Promise((resolve) => {
          UiViewFactory.createUi(
            R.Prefabs.FriendMain,
            LayerIdx.layer4,
            (ui: FriendMainController) => {
              this.instance = ui;
              this.creatingPromise = null;
              resolve(ui);
            }
          );
        });
        return this.creatingPromise;
      }
}
