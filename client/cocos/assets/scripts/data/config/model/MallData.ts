
       import BaseConfigItem from '../BaseConfigItem';
            
export class ConsumesDef {
    public type: string;
    public value: string;
}

export class RewardsDef {
    public type: string;
    public value: string;
}

        export default class MallData extends BaseConfigItem {
          public static fileName:string = "mallData";
        
            private _name: string;
            public get name():string {return this._name;}
            
            private _type: number;
            public get type():number {return this._type;}
            
            private _consumes: Array<ConsumesDef>;
            public get consumes():Array<ConsumesDef> {return this._consumes;}
            
            private _rewards: Array<RewardsDef>;
            public get rewards():Array<RewardsDef> {return this._rewards;}
            
        public constructor(data:any) {
            super(data);
                this._name = data['name']
        this._type = data['type']
        this._consumes = data['consumes']
        this._rewards = data['rewards']
    }
}
