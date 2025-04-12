
       import BaseConfigItem from '../BaseConfigItem';
            
        export default class SkillData extends BaseConfigItem {
          public static fileName:string = "skillData";
        
            private _skillId: number;
            public get skillId():number {return this._skillId;}
            
            private _stage: number;
            public get stage():number {return this._stage;}
            
            private _name: string;
            public get name():string {return this._name;}
            
            private _hero: string;
            public get hero():string {return this._hero;}
            
            private _tips: string;
            public get tips():string {return this._tips;}
            
        public constructor(data:any) {
            super(data);
                this._skillId = data['skillId']
        this._stage = data['stage']
        this._name = data['name']
        this._hero = data['hero']
        this._tips = data['tips']
    }
}
