/**
 * 
 */
import { FriendApplyVo } from './FriendApplyVo';
import { FriendVo } from './FriendVo';
export  class PushFriendInfo {
     public static cmd: number =  1997; 
        
        /**  */
        public  applyItems : Array<FriendApplyVo>;

        /**  */
        public  friendItems : Array<FriendVo>;

        /**  */
        public  friendSum : number;

    
}