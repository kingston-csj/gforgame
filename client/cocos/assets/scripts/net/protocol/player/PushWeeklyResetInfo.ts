/**
 * 玩家每周重置信息推送
 */
export  class PushWeeklyResetInfo {
    
        
        /** 充值累计积分 */
        public  weeklyRechargeSum : number;
        
        /** 充值礼包奖励领取状态, 格式为 id1=status,id2=status */
        public  weeklyGiftRewards : string;
        
        /** 充值每周限购商品信息, 格式为 id1=次数1,id2=次数2 */
        public  rechargeBuyTimes : string;
        
    
}