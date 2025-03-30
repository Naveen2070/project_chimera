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
    }
}