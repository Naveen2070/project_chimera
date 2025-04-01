using Microsoft.AspNetCore.Mvc;
using notification_service.Services.Interfaces;

namespace notification_service.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class FloraNotificationController : ControllerBase
    {
        private readonly IFloraNotificationService _floraNotificationService;

        public FloraNotificationController(IFloraNotificationService floraNotificationService)
        {
            _floraNotificationService = floraNotificationService;
        }

        [HttpGet]
        public IActionResult GetFloraNotifications()
        {
            var floraNotifications = _floraNotificationService.GetNotifications();
            return Ok(floraNotifications);
        }

        [HttpGet("flora-notifications-stream")]
        [Produces("text/event-stream")]
        public async Task GetFloraNotification(CancellationToken cancellationToken)
        {
            Response.ContentType = "text/event-stream";
            Response.Headers.Append("Cache-Control", "no-cache");
            Response.Headers.Append("Connection", "keep-alive");

            string clientId = Guid.NewGuid().ToString();

            await foreach (var message in _floraNotificationService.GetFloraNotificationsStreamAsync(clientId, cancellationToken))
            {
                Console.WriteLine($"Sending to {clientId}: {message}");
                await Response.WriteAsync($"data: {message}\n\n");
                await Response.Body.FlushAsync();
            }
        }
    }
}
