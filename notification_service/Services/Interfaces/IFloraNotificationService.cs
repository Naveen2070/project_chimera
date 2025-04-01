using System.Runtime.CompilerServices;

namespace notification_service.Services.Interfaces
{
    public interface IFloraNotificationService
    {
        Task SendNotificationAsync(string message);

        IAsyncEnumerable<string> GetFloraNotificationsStreamAsync(string clientId, CancellationToken cancellationToken);

        IEnumerable<string> GetNotifications();
    }
}
