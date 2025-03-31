using System.Collections.Concurrent;
using System.Runtime.CompilerServices;
using notification_service.Services.Interfaces;

namespace notification_service.Services
{
    public class FloraNotificationService : IFloraNotificationService
    {
        private readonly IRMQConsumerService _rmqConsumerService;
        private readonly ConcurrentQueue<string> _notificationQueue = new();

        public FloraNotificationService(IRMQConsumerService rmqConsumerService)
        {
            _rmqConsumerService = rmqConsumerService;
            StartListening();
        }

        private void StartListening()
        {
            _rmqConsumerService.StartListeningAsync("notification_queue", async (message) =>
            {
                Console.WriteLine($"Received message: {message}");
                _notificationQueue.Enqueue(message.ToString());
            });
        }

        public async IAsyncEnumerable<string> GetFloraNotificationsStreamAsync([EnumeratorCancellation] CancellationToken cancellationToken)
        {
            while (!cancellationToken.IsCancellationRequested)
            {
                while (_notificationQueue.TryDequeue(out var message))
                {
                    yield return message;
                }

                await Task.Delay(500, cancellationToken);
            }
        }

        public IEnumerable<string> GetNotifications()
        {
            throw new NotImplementedException();
        }

        public Task SendNotificationAsync(string message)
        {
            throw new NotImplementedException();
        }
    }
}
