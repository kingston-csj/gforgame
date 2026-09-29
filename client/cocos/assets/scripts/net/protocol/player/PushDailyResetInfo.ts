/**
 * 玩家每日重置信息推送
 */
export  class PushDailyResetInfo {
     public static cmd: number =  156; 
        
        /** 普通招募次数 */
        public  normalRecruitTimes : number;
        
        /** 高级招募次数 */
        public  highRecruitTimes : number;
        
        /** 商城每日购买次数 */
        public  mallDailyBuyTimes : number;
        
        /** 每日充值金额 */
        public  dailyRechargeSum : number;
        
    
}