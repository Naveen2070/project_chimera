using System.Collections.Concurrent;
using notification_service.Services.Interfaces;

namespace notification_service.Services
{
    public class FloraNotificationService : IFloraNotificationService
    {
        private readonly IRMQConsumerService _rmqConsumerService;
        private readonly ConcurrentBag<string> _notifications; // Thread-safe collection

        public FloraNotificationService(IRMQConsumerService rmqConsumerService)
        {
            _rmqConsumerService = rmqConsumerService;
            _notifications = new ConcurrentBag<string>();
        }

        public async Task GetSavedNotificationsAsync()
        {
            await _rmqConsumerService.StartListeningAsync("notification_queue", async (message) =>
            {
                // Store messages in a thread-safe collection
                _notifications.Add(message);
                await Task.CompletedTask; // Ensure async compatibility
            });
        }

        public IEnumerable<string> GetNotifications()
        {
            return _notifications;
        }

        public Task SendNotificationAsync(string message)
        {
            throw new NotImplementedException();
        }
    }
}