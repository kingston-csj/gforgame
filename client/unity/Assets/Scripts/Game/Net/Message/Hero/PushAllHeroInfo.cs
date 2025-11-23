using Nova.Net.Socket;

namespace Game.Net.Message.Hero
{
    [MessageMeta(Cmd = 2003)]
    public class PushAllHeroInfo : Nova.Net.Socket.Message
    {
        public HeroInfo[] heros;
        
    }
}