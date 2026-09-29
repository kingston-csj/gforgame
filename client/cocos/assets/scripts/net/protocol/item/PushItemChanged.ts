/**
 * 
 */
import { ItemInfo } from "../common/ItemInfo";
export  class PushItemChanged {
     public static cmd: number =  253; 
        
        /**  */
        public  type : string;
        
        /**  */
        public  changed : Array<ItemInfo>;
        
    
}