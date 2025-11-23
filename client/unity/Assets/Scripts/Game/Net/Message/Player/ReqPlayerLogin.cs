
namespace Game.Net.Message
{
    using Nova.Net.Socket;
    /// <summary>
    /// 玩家登录请求
    /// </summary>
    [MessageMeta(Cmd = 103)]
    public class ReqPlayerLogin : Message
    {
        public string playerId;
        
    }
}