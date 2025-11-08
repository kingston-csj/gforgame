using System;

/*
 * 通过特性 [MessageHandler] 自动注册消息处理方法
 *
 * 示例：
 * [MessageHandler]
 * public void HandleLoginResponse(ResPlayerLogin res)
 * {
 *     Debug.Log("收到登录响应: " + res.playerId);
 * }
 */
namespace Nova.Net.Socket
{
    /// <summary>
    /// 自动注册消息处理方法的特性
    /// 用于标记需要自动注册的消息处理方法
    /// </summary>
    [AttributeUsage(AttributeTargets.Method)]
    public class MessageHandler : Attribute
    {
    }
    
}