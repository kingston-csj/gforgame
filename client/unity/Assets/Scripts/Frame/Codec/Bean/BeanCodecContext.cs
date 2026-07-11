using Frame.Commons.Utils;

namespace Frame.Codec.Bean
{
    public class BeanCodecContext
    {
        /// <summary>
        /// 集合字段支持的序列化模式，默认为严格同构模式
        /// </summary>
        static CollectionSerializeMode CollectionSerializeMode = CollectionSerializeMode.STRICT_HOMOGENEOUS;

        /// <summary>
        /// 消息工厂，仅当 collectionSerializeMode为SUB_CLASS_POLYMORPHIC才需要
        /// </summary>
        public static IMessageFactory MessageFactory;
    }
}