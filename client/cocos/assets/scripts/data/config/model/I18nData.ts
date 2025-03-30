import BaseConfigItem from '../BaseConfigItem';

export default class I18nData extends BaseConfigItem {
  public static fileName: string = 'i18nData';

  private _content: string;
  public get content(): string {
    return this._content;
  }

  public constructor(data: any) {
    super(data);
    this._content = data['content'];
  }
}
