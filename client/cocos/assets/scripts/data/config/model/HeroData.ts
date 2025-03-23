
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
            
        public constructor(data:any) {
            super(data);
                this._name = data['name']
        this._quality = data['quality']
        this._tips = data['tips']
        this._icon = data['icon']
        this._prob = data['prob']
        this._shard = data['shard']
        this._item = data['item']
    }
}
