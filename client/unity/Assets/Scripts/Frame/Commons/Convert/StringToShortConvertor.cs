using System;
using Game.Util;

namespace Nova.Commons.Convert
{
    public class StringToShortConvertor : GenericConverter
    {
        public Type GetSourceType()
        {
            return typeof(string);
        }

        public Type GetTargetType()
        {
            return typeof(short);
        }

        public object Convert(object source)
        {
            return NumberUtil.ShortValue(source);
        }
    }
}