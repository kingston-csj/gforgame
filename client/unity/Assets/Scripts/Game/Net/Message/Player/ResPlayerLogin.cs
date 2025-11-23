namespace Game.Net.Message
{
    using Nova.Net.Socket;
    /// <summary>
    /// 玩家登录响应
    /// </summary>
    [MessageMeta(Cmd = 154)]
    public class ResPlayerLogin : Response
    {
        public string name;

        public int fighting;

        public int camp;
        
    }
}