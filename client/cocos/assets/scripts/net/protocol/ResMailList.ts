import { MailVo } from './items/MailVo';

export class PushMailAll {
  public static cmd: number = 6011;
  public mails: MailVo[] = [];
}
