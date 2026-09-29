/**
 * 聊天消息请求
 */
export  class ReqChat {
     public static cmd: number =  1801; 
        
        /**  */
        public  channel : number;
        
        /**  */
        public  target : string;
        
        /**  */
        public  content : string;
        
    
}