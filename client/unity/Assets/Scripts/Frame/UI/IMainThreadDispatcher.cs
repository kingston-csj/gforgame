using System;
using System.Threading.Tasks;

namespace Nova.Ui
{
    public interface IMainThreadDispatcher
    {
        Task EnqueueAsync(Action action);
        Task<T> EnqueueAsync<T>(Func<T> func);
    }
}