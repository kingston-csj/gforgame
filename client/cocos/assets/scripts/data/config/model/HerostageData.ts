
       import BaseConfigItem from '../BaseConfigItem';
            
        export default class HerostageData extends BaseConfigItem {
          public static fileName:string = "herostageData";
        
            private _max_level: number;
            public get max_level():number {return this._max_level;}
            
            private _stage: number;
            public get stage():number {return this._stage;}
            
            private _cost: number;
            public get cost():number {return this._cost;}
            
            private _hp: number;
            public get hp():number {return this._hp;}
            
            private _attack: number;
            public get attack():number {return this._attack;}
            
            private _defense: number;
            public get defense():number {return this._defense;}
            
            private _speed: number;
            public get speed():number {return this._speed;}
            
        public constructor(data:any) {
            super(data);
                this._max_level = data['max_level']
        this._stage = data['stage']
        this._cost = data['cost']
        this._hp = data['hp']
        this._attack = data['attack']
        this._defense = data['defense']
        this._speed = data['speed']
    }
}
