/**
 * 
 */
export  class PushSigninInfo {
     public static cmd: number =  3099; 
        
        /** 本月总天数 */
        public  daysInMonth : number;
        
        /** 今天第几天 */
        public  nthDay : number;
        
        /** 已签到天数 */
        public  signinDays : Array<number>;
        
        /** 是否已补签 */
        public  signInMakeUp : Dictionary<int32,int>;
        
    
}