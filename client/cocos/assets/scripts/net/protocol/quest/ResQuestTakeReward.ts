/**
 * 任务——领取达标奖
 */
export  class ResQuestTakeReward {
     public static cmd: number =  754; 
        
        /**  */
        public  code : number;
        
        /** 今日活跃度 */
        public  dailyScore : number;
        
        /** 本周活跃度 */
        public  weeklyScore : number;
        
        /** 奖励vo */
        public  rewardVos : Array<RewardVo>;
        
    
}