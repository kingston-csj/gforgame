using System;
using Game.Util;

namespace Nova.Commons.Convert
{
    public class StringToIntConvertor : GenericConverter
    {
        public Type GetSourceType()
        {
            return typeof(string);
        }

        public Type GetTargetType()
        {
            return typeof(int);
        }

        public object Convert(object source)
        {
            return NumberUtil.IntValue(source);
        }
    }
}