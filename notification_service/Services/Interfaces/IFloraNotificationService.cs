using System.Runtime.CompilerServices;

namespace notification_service.Services.Interfaces
{
    public interface IFloraNotificationService
    {
        Task SendNotificationAsync(string message);

        IAsyncEnumerable<string> GetFloraNotificationsStreamAsync(CancellationToken cancellationToken);

        IEnumerable<string> GetNotifications();
    }
}
