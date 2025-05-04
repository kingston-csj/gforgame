import { MailVo } from '../../net/protocol/items/MailVo';
import { BaseModel } from '../../ui/BaseModel';

export class MailBoxModel extends BaseModel {
  private static instance: MailBoxModel = new MailBoxModel();

  private _mails: Map<number, MailVo> = new Map();

  public static STATUS_UNREAD = 1;
  public static STATUS_READ = 2;
  public static STATUS_RECEIVED = 3;

  public static getInstance(): MailBoxModel {
    return this.instance;
  }

  public reset(mails: Map<number, MailVo>): void {
    this._mails = mails;
  }

  public getMails(): MailVo[] {
    return Array.from(this._mails.values());
  }

  public getMail(mailId: number): MailVo {
    return this._mails.get(mailId);
  }

  public deleteMails(mailIds: number[]): void {
    mailIds.forEach((mailId) => {
      this._mails.delete(mailId);
    });
  }
}
