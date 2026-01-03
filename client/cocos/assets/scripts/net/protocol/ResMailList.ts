import { MailVo } from "./items/MailVo";

export class PushMailAll {
  public static cmd: number = 599;
  public mails: MailVo[] = [];
}
