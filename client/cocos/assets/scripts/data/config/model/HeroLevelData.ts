
       import BaseConfigItem from '../BaseConfigItem';
            
        export default class HerolevelData extends BaseConfigItem {
          public static fileName:string = "herolevelData";
        
            private _level: number;
            public get level():number {return this._level;}
            
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
                this._level = data['level']
        this._cost = data['cost']
        this._hp = data['hp']
        this._attack = data['attack']
        this._defense = data['defense']
        this._speed = data['speed']
    }
}
