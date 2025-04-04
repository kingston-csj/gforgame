export class PurseModel {
  private static instance: PurseModel;

  private _diamond: number = 0;
  private _gold: number = 0;

  public static getInstance(): PurseModel {
    if (!PurseModel.instance) {
      PurseModel.instance = new PurseModel();
    }
    return PurseModel.instance;
  }

  // 使用 getter/setter 来触发数据变化通知
  get diamond(): number {
    return this._diamond;
  }

  set diamond(value: number) {
    if (this._diamond !== value) {
      this._diamond = value;
      this.notifyDiamondChange();
    }
  }

  get gold(): number {
    return this._gold;
  }

  set gold(value: number) {
    if (this._gold !== value) {
      this._gold = value;
      this.notifyGoldChange();
    }
  }

  // 数据变化时的回调函数
  private diamondChangeCallbacks: ((value: number) => void)[] = [];
  private goldChangeCallbacks: ((value: number) => void)[] = [];

  // 注册数据变化监听
  onDiamondChange(callback: (value: number) => void) {
    this.diamondChangeCallbacks.push(callback);
    return () => {
      const index = this.diamondChangeCallbacks.indexOf(callback);
      if (index > -1) {
        this.diamondChangeCallbacks.splice(index, 1);
      }
    };
  }

  onGoldChange(callback: (value: number) => void) {
    this.goldChangeCallbacks.push(callback);
    return () => {
      const index = this.goldChangeCallbacks.indexOf(callback);
      if (index > -1) {
        this.goldChangeCallbacks.splice(index, 1);
      }
    };
  }

  // 通知数据变化
  private notifyDiamondChange() {
    this.diamondChangeCallbacks.forEach((callback) => callback(this._diamond));
  }

  private notifyGoldChange() {
    this.goldChangeCallbacks.forEach((callback) => callback(this._gold));
  }
}
