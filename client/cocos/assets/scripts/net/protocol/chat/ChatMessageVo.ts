/**
 * 聊天消息vo
 */
export  class ChatMessageVo {
    
        
        /** 消息id */
        public  id : string;
        
        /** 发送频道：1个人 2世界 */
        public  channel : number;
        
        /**  */
        public  senderId : string;
        
        /**  */
        public  senderName : string;
        
        /**  */
        public  senderHead : number;
        
        /**  */
        public  receiverId : string;
        
        /**  */
        public  timestamp : number;
        
        /**  */
        public  content : string;
        
    
}