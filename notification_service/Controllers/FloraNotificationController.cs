//	Copyright 2025 Naveen R
//
//		Licensed under the Apache License, Version 2.0 (the "License");
//		you may not use this file except in compliance with the License.
//		You may obtain a copy of the License at
//
//		http://www.apache.org/licenses/LICENSE-2.0
//
//		Unless required by applicable law or agreed to in writing, software
//		distributed under the License is distributed on an "AS IS" BASIS,
//		WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//		See the License for the specific language governing permissions and
//		limitations under the License.

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
