using System;

namespace Frame.Commons.Utils
{
    public interface IMessageFactory
    {
        void Register(int cmd, Type type);

        int GetMessageCmd(Type type);

        Type GetMessageType(int cmd);

        bool Contains(int cmd);
    }
}