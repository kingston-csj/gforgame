using System;
using System.Collections.Generic;

namespace Nova.Commons.Convert
{
    public class GenericConversionService : IConversionService
    {
        private readonly Dictionary<ConverterCacheKey, GenericConverter> _cache =
            new Dictionary<ConverterCacheKey, GenericConverter>();

        public GenericConversionService()
        {
            AddConvertor(new StringToBoolConvertor());
            AddConvertor(new StringToShortConvertor());
            AddConvertor(new StringToIntConvertor());
        }

        private void AddConvertor(GenericConverter converter)
        {
            ConverterCacheKey key = new ConverterCacheKey()
            {
                SourceType = converter.GetSourceType(),
                TargetType = converter.GetTargetType()
            };
            if (_cache.ContainsKey(key))
            {
                throw new Exception($"{converter.GetSourceType()} to {converter.GetTargetType()} already exists!");
            }

            _cache.Add(key, converter);
        }


        public bool CanConvert(string sourceValue, Type targetType)
        {
            throw new NotImplementedException();
        }


        public object Convert(string sourceValue, Type targetType)
        {
            if (targetType == typeof(string))
            {
                return sourceValue;
            }
            GenericConverter converter = GetConverter(sourceValue.GetType(), targetType);
            if (converter != null)
            {
                return converter.Convert(sourceValue);
            }

            throw new NotSupportedException($"{sourceValue.GetType()} to {targetType} not supported!");
        }

        private GenericConverter GetConverter(Type sourceType, Type targetType)
        {
            ConverterCacheKey key = new ConverterCacheKey()
            {
                SourceType = sourceType,
                TargetType = targetType
            };
            if (!_cache.ContainsKey(key))
            {
                throw new Exception($"{sourceType} to {targetType} not exists!");
            }

            return _cache[key];
        }
    }
}