/**
 * 任务——领取档位奖
 */
export  class ResQuestTakeProgressReward {
     public static cmd: number =  753; 
        
        /** 任务类型 2每日，5每周，6公会 */
        public  type : number;
        
        /** 已领取的档位索引（0为未领取) */
        public  rewardIndex : number;
        
        /** 奖励vo */
        public  rewardVos : Array<RewardVo>;
        
    
}