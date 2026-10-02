import { BaseModel } from '../../frame/mvc/BaseModel';
import { FriendApplyVo } from '../../net/protocol/friend/FriendApplyVo';
import { FriendVo } from '../../net/protocol/friend/FriendVo';
import GameContext from '../../GameContext';
import { ReqFriendSearchPlayers } from '../../net/protocol/friend/ReqFriendSearchPlayers';
import { ResFriendSearchPlayers } from '../../net/protocol/friend/ResFriendSearchPlayers';
import { ReqFriendApply } from '../../net/protocol/friend/ReqFriendApply';
import { ResFriendApply } from '../../net/protocol/friend/ResFriendApply';
import { ReqFriendDealApplyRecord } from '../../net/protocol/friend/ReqFriendDealApplyRecord';
import { ResFriendDealApplyRecord } from '../../net/protocol/friend/ResFriendDealApplyRecord';

export class FriendModel extends BaseModel {
    public static getInstance(): FriendModel {
        if (!this._instance) {
            this._instance = new FriendModel();
        }
        return this._instance;
    }
    private static _instance: FriendModel;

    //三份核心数据
    public friendList: FriendVo[] = [];
    public applyList: FriendApplyVo[] = [];
    public searchResult: FriendVo[] = [];

    /** 请求好友列表 */
    public refreshFriendListCache(friends: Array<FriendVo>) {
        this.friendList = friends;
        this.notifyChange("friendList", this.friendList);
    }

        /** 请求好友申请列表 */
    public refreshApplyListCache(applies: Array<FriendApplyVo>) {
        this.applyList = applies;
        this.notifyChange("applyList", this.applyList);
    }

    /** 请求好友申请列表 */
    public async requestApplyList() {
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


    public requestApplyFriend(uid: string): Promise<ResFriendApply> {
        return new Promise<ResFriendApply>((resolve, reject) => {
        GameContext.wsClient.sendMessage(
            ReqFriendApply.cmd,
            {
            friendId: uid,
            },
            (msg: ResFriendApply) => {
            resolve(msg);
            }
        );
        });
    }

    /**
     * 处理好友申请
     * @param uid 好友id
     * @param status 申请结果：1同意 2拒绝
     * @returns 
     */
    public requestDealApplyFriend(uid: string, status: number): Promise<ResFriendDealApplyRecord> {
        return new Promise<ResFriendDealApplyRecord>((resolve, reject) => {
        GameContext.wsClient.sendMessage(
            ReqFriendDealApplyRecord.cmd,
            {
            applyId: uid,
            status: status,
            },
            (msg: ResFriendDealApplyRecord) => {
            resolve(msg);
            }
        );
        });
    }

}
