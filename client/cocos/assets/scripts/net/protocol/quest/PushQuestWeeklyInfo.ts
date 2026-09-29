/**
 * 每周任务主界面信息
 */
export  class PushQuestWeeklyInfo {
     public static cmd: number =  799; 
        
        /** 已领取的档位索引（0为未领取) */
        public  weeklyRewardIndex : number;
        
        /** 本周活跃度 */
        public  weeklyScore : number;
        
        /** 所有任务 */
        public  quests : Array<QuestVo>;
        
    
}