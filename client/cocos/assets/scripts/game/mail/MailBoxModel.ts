import { BaseModel } from '../../frame/mvc/BaseModel';
import { MailVo } from '../../net/protocol/mail/MailVo';

export class MailBoxModel extends BaseModel {
  private static instance: MailBoxModel = new MailBoxModel();

  private _mails: Map<string, MailVo> = new Map();

  public static STATUS_UNREAD = 1;
  public static STATUS_READ = 2;
  public static STATUS_RECEIVED = 3;

  public static getInstance(): MailBoxModel {
    return this.instance;
  }

  public reset(mails: Map<string, MailVo>): void {
    this._mails = mails;
  }

  public getMails(): MailVo[] {
    return Array.from(this._mails.values());
  }

  public getMail(mailId: string): MailVo {
    return this._mails.get(mailId);
  }

  public deleteMails(mailIds: string[]): void {
    mailIds.forEach((mailId) => {
      this._mails.delete(mailId);
    });
  }
}
