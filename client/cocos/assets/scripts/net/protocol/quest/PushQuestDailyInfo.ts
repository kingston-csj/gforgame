/**
 * 
 */
export  class PushQuestDailyInfo {
     public static cmd: number =  798; 
        
        /** 已领取的档位索引（0为未领取) */
        public  dailyRewardIndex : number;
        
        /** 今日活跃度 */
        public  dailyScore : number;
        
        /** 所有任务 */
        public  quests : Array<QuestVo>;
        
    
}