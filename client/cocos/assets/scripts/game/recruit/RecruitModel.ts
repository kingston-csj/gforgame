import GameContext from '../../GameContext';
import { ReqHeroRecruit } from '../../net/protocol/ReqHeroRecruit';
import { ResHeroRecruit } from '../../net/protocol/ResHeroRecruit';

export class RecruitModel {
  private static _instance: RecruitModel;

  public static get instance(): RecruitModel {
    if (!RecruitModel._instance) {
      RecruitModel._instance = new RecruitModel();
    }
    return RecruitModel._instance;
  }

  public doRecruit(times: number): Promise<ResHeroRecruit> {
    return new Promise<ResHeroRecruit>((resolve, reject) => {
      GameContext.wsClient.sendMessage(
        ReqHeroRecruit.cmd,
        {
          counter: times,
        },
        (msg: ResHeroRecruit) => {
          resolve(msg);
        }
      );
    });
  }
}
