import { RankInfo } from './items/RankInfo';

export class ResRankQuery {
  public static cmd: number = 7002;

  public type: number;
  public records: RankInfo[];
  public myRecord: RankInfo;
}
