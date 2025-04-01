using System;
using System.Threading.Tasks;

namespace notification_service.Services.Interfaces
{
    public interface IRMQConsumerService : IAsyncDisposable
    {
        Task StartListeningAsync(string queueName, Func<string, ulong, Task> messageHandler);
        Task AckMessage(ulong deliveryTag);
    }
}
