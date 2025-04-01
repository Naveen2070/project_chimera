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

using System.Collections.Concurrent;
using System.Runtime.CompilerServices;
using notification_service.Services.Interfaces;

namespace notification_service.Services
{
    public class FloraNotificationService : IFloraNotificationService
    {
        private readonly IRMQConsumerService _rmqConsumerService;
        private readonly ConcurrentQueue<string> _notificationQueue = new();
        private readonly ConcurrentDictionary<string, bool> _connectedClients = new();

        public FloraNotificationService(IRMQConsumerService rmqConsumerService)
        {
            _rmqConsumerService = rmqConsumerService;
            StartListening();
        }

        private void StartListening()
        {
            _rmqConsumerService.StartListeningAsync("notification_queue", async (message, deliveryTag) =>
            {
                Console.WriteLine($"Received message: {message}");

                if (_connectedClients.Count > 0)
                {
                    // If clients are connected, broadcast the message immediately
                    foreach (var key in _connectedClients.Keys)
                    {
                        _notificationQueue.Enqueue(message);
                    }
                }
                else
                {
                    // If no clients are connected, store the message for later
                    _notificationQueue.Enqueue(message);
                }

                try
                {
                    await _rmqConsumerService.AckMessage(deliveryTag);
                }
                catch (Exception ex)
                {
                    Console.WriteLine(ex.ToString());
                    throw new Exception(ex.Message);
                }
            });
        }

        public void AddClient(string clientId)
        {
            _connectedClients.TryAdd(clientId, true);
        }

        public void RemoveClient(string clientId)
        {
            _connectedClients.TryRemove(clientId, out _);
        }

        public async IAsyncEnumerable<string> GetFloraNotificationsStreamAsync(
            string clientId,
            [EnumeratorCancellation] CancellationToken cancellationToken)
        {
            AddClient(clientId);

            try
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
            finally
            {
                RemoveClient(clientId);
            }
        }

        public IEnumerable<string> GetNotifications()
        {
            return _notificationQueue.ToArray();
        }

        public Task SendNotificationAsync(string message)
        {
            _notificationQueue.Enqueue(message);
            return Task.CompletedTask;
        }
    }
}
