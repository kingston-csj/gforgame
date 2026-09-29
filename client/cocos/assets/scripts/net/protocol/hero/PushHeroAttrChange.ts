/**
 * 
 */
import { AttrInfo } from "../common/AttrInfo";
export  class PushHeroAttrChange {
     public static cmd: number =  5007; 
        
        /**  */
        public  heroId : number;
        
        /**  */
        public  attrs : Array<AttrInfo>;
        
        /**  */
        public  fight : number;
        
    
}