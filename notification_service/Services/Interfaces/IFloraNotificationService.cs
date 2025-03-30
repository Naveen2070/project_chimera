namespace notification_service.Services.Interfaces
{
    public interface IFloraNotificationService
    {
        Task SendNotificationAsync(string message);

        Task GetSavedNotificationsAsync();

        IEnumerable<string> GetNotifications();
    }
}
