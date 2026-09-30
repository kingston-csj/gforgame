import { BaseModel } from '../../frame/mvc/BaseModel';
import { FriendApplyVo } from '../../net/protocol/friend/FriendApplyVo';
import { FriendVo } from '../../net/protocol/friend/FriendVo';
import GameContext from '../../GameContext';
import { ReqFriendSearchPlayers } from '../../net/protocol/friend/ReqFriendSearchPlayers';
import { ResFriendSearchPlayers } from '../../net/protocol/friend/ResFriendSearchPlayers';
import { ReqFriendApply } from '../../net/protocol/friend/ReqFriendApply';
import { ResFriendApply } from '../../net/protocol/friend/ResFriendApply';

//定义数据结构
export interface FriendItem {
    id: string;
    name: string;
    head: number;
   
}

export interface FriendApplyItem {
    id : string;

    fromId : string;

    fromName : string;

    fromHead : number;

    targetId : string;

    targetName : string;

    /** 申请结果：1同意 2拒绝 0未处理 */
    status : number;
    time : number;
}

export class FriendModel extends BaseModel {
    public static getInstance(): FriendModel {
        if (!this._instance) {
            this._instance = new FriendModel();
        }
        return this._instance;
    }
    private static _instance: FriendModel;

    //三份核心数据
    public friendList: FriendItem[] = [];
    public applyList: FriendApplyItem[] = [];
    public searchResult: FriendItem[] = [];

    /** 请求好友列表 */
    public refreshFriendList(friends: Array<FriendVo>) {
        this.friendList = friends.map((item) => ({
            id: item.id,
            name: item.name,
            head: item.head,
        }));
        this.notifyChange("friendList", this.friendList);
    }

    /** 请求好友申请列表 */
    public async requestApplyList() {
        //const resp = await net.post("getApplyList")
        //this.applyList = resp.list;
        this.notifyChange("applyList", this.applyList);
    }

    /** 搜索好友 */
    public requestSearchFriend(keyword: string) {
                  GameContext.wsClient.sendMessage(
            ReqFriendSearchPlayers.cmd,
            {
              keyword: keyword,
            },
            (msg: ResFriendSearchPlayers) => {
                this.searchResult = msg.items;
                this.notifyChange("searchResult", msg.items);
            }
          );
    }

     public requestApplyFriend(uid: string) {
                  GameContext.wsClient.sendMessage(
            ReqFriendApply.cmd,
            {
              uid: uid,
            },
            (msg: ResFriendApply) => {
            }
          );
    }


}
