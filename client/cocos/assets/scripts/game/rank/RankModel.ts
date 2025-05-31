import { BaseModel } from '../../frame/mvc/BaseModel';
import GameContext from '../../GameContext';
import { RankInfo } from '../../net/protocol/items/RankInfo';
import { ReqRankQuery } from '../../net/protocol/ReqRankQuery';
import { ResRankQuery } from '../../net/protocol/ResRankQuery';

export class RankModel extends BaseModel {
  public static readonly RANK_TYPE_LEVEL = 1;
  public static readonly RANK_TYPE_FIGHTING = 2;

  private static instance: RankModel;

  public static getInstance(): RankModel {
    if (!RankModel.instance) {
      RankModel.instance = new RankModel();
    }
    return RankModel.instance;
  }

  private _selectedRankType: number = 0;

  private _records: RankInfo[] = [];

  get selectedRankType(): number {
    return this._selectedRankType;
  }

  set selectedRankType(value: number) {
    this._selectedRankType = value;
  }

  get records(): RankInfo[] {
    return this._records;
  }

  set records(value: RankInfo[]) {
    this._records = value;
    this.notifyChange('records', value);
  }

  public queryRank(type: number): Promise<ResRankQuery> {
    return new Promise<ResRankQuery>((resolve, reject) => {
      GameContext.instance.WebSocketClient.sendMessage(
        ReqRankQuery.cmd,
        { type: type, start: 1, pageSize: 100 },
        (msg: ResRankQuery) => {
          resolve(msg);
        }
      );
    });
  }
}
