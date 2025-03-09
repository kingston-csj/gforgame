import PlayerData from './data/user/PlayerData';
import { WebSocketClient } from './net/WebSocketClient';

export default class GameContext {
  protected static _instance: GameContext = new GameContext();

  private _playerData: PlayerData = new PlayerData();

  private _wsClient: WebSocketClient = new WebSocketClient();

  public static get instance(): GameContext {
    return this._instance;
  }

  public get WebSocketClient(): WebSocketClient {
    return this._wsClient;
  }
}
