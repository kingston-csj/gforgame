/**
 * 任务——一键领取所有奖励
 */
export  class ResQuestTakeAllRewards {
     public static cmd: number =  762; 
        
        /** 任务类型 1主线，2日常 */
        public  category : number;
        
        /** 奖励vo */
        public  rewardVos : Array<RewardVo>;
        
        /** 今日活跃度 */
        public  dailyScore : number;
        
        /** 本周活跃度 */
        public  weeklyScore : number;
        
        /** 已领取的任务id列表 */
        public  questIds : Array<number>;
        
        /** 总分数 */
        public  score : number;
        
    
}