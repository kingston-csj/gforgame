using System;
using Game.Util;

namespace Nova.Commons.Convert
{
    public class StringToBoolConvertor : GenericConverter
    {
        public Type GetSourceType()
        {
            return typeof(string);
        }

        public Type GetTargetType()
        {
            return typeof(bool);
        }

        public object Convert(object source)
        {
            return NumberUtil.BooleanValue(source);
        }
    }
}