/**
 * 成就——加载所有信息
 */
export  class PushAchievementInfo {
     public static cmd: number =  791; 
        
        /** 积分 */
        public  score : number;
        
        /** 所有任务 */
        public  achievementVos : Array<QuestVo>;
        
    
}