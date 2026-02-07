using System;

namespace Nova.Commons.Convert
{
    public interface GenericConverter
    {
        Type GetSourceType();

        Type GetTargetType();

        object Convert(object source);
    }
}