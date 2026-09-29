/**
 * 推送新聊天消息
 */
export  class PushChatNewMessage {
     public static cmd: number =  1899; 
        
        /**  */
        public  code : number;
        
        /**  */
        public  messages : Array<ChatMessageVo>;
        
    
}