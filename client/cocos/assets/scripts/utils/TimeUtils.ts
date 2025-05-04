export class TimeUtils {
  public static ONE_SECOND = 1000;
  public static ONE_MINUTE = 60 * TimeUtils.ONE_SECOND;
  public static ONE_HOUR = 60 * TimeUtils.ONE_MINUTE;
  public static ONE_DAY = 24 * TimeUtils.ONE_HOUR;

  // 获取两个时间之间的天数差
  // 如果不足1天，则返回0
  public static getDiffDays(t1: number, t2: number): number {
    return Math.floor((t2 - t1) / TimeUtils.ONE_DAY);
  }

  public static getDiffHours(t1: number, t2: number): number {
    return Math.floor((t2 - t1) / TimeUtils.ONE_HOUR);
  }

  public static getDiffMinutes(t1: number, t2: number): number {
    return Math.floor((t2 - t1) / TimeUtils.ONE_MINUTE);
  }

  // 获取剩余时间提示
  // 如果不足1天，则返回小时
  // 如果不足1小时，则返回分钟
  public static getLeftTimeTips(time: number): string {
    const now = Date.now();
    const days = TimeUtils.getDiffDays(now, time);
    const hours = TimeUtils.getDiffHours(now, time);
    const minutes = TimeUtils.getDiffMinutes(now, time);
    if (days < 0) {
      return '已过期';
    }
    if (days > 0) {
      return `${days}天后过期 `;
    }
    if (hours > 0) {
      return `${hours}小时后过期`;
    }
    return `${minutes}分钟后过期`;
  }
}
