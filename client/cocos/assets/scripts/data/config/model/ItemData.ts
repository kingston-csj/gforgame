
       import BaseConfigItem from '../BaseConfigItem';
            
        export default class ItemData extends BaseConfigItem {
          public static fileName:string = "itemData";
        
            private _type: number;
            public get type():number {return this._type;}
            
            private _name: string;
            public get name():string {return this._name;}
            
            private _quality: number;
            public get quality():number {return this._quality;}
            
            private _tips: string;
            public get tips():string {return this._tips;}
            
            private _icon: string;
            public get icon():string {return this._icon;}
            
        public constructor(data:any) {
            super(data);
                this._type = data['type']
        this._name = data['name']
        this._quality = data['quality']
        this._tips = data['tips']
        this._icon = data['icon']
    }
}
