/**
 * 
 */
export  class ChallengeRecord {
    
        
        /** 对手id */
        public  opponentId : string;
        
        /** 对手名称 */
        public  opponentName : string;
        
        /** 对手头像 */
        public  opponentHead : number;
        
        /** 对手战力 */
        public  opponentFighting : number;
        
        /** 挑战时间 */
        public  challengeTime : number;
        
        /** 获胜方id */
        public  winner : string;
        
        /** 得分，获胜为正数， 失败为负数 */
        public  score : number;
        
    
}