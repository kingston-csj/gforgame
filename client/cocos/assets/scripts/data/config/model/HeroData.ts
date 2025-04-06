
       import BaseConfigItem from '../BaseConfigItem';
            
        export default class HeroData extends BaseConfigItem {
          public static fileName:string = "heroData";
        
            private _name: string;
            public get name():string {return this._name;}
            
            private _quality: number;
            public get quality():number {return this._quality;}
            
            private _tips: string;
            public get tips():string {return this._tips;}
            
            private _icon: string;
            public get icon():string {return this._icon;}
            
            private _prob: number;
            public get prob():number {return this._prob;}
            
            private _shard: number;
            public get shard():number {return this._shard;}
            
            private _item: number;
            public get item():number {return this._item;}
            
            private _skills: string;
            public get skills():string {return this._skills;}
            
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
                this._name = data['name']
        this._quality = data['quality']
        this._tips = data['tips']
        this._icon = data['icon']
        this._prob = data['prob']
        this._shard = data['shard']
        this._item = data['item']
        this._skills = data['skills']
        this._hp = data['hp']
        this._attack = data['attack']
        this._defense = data['defense']
        this._speed = data['speed']
    }
}
