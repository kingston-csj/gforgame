/**
 * 邮件vo
 */
export  class MailVo {
    
        
        /**  */
        public  id : string;
        
        /** 邮件标题， 当TemplateId为0时，需要此字段 */
        public  title : string;
        
        /** 邮件内容， 当TemplateId为0时，需要此字段 */
        public  content : string;
        
        /** 邮件奖励 */
        public  rewards : Array<RewardVo>;
        
        /** 邮件模板id */
        public  templateId : number;
        
        /** 邮件状态 */
        public  status : number;
        
        /** 邮件时间 */
        public  time : number;
        
    
}