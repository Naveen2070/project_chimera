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

using notification_service.Services.Interfaces;
using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using System.Text;

namespace notification_service.Services
{
    public class RMQConsumerService : IRMQConsumerService, IAsyncDisposable
    {
        private readonly IConfiguration _configuration;
        private readonly IConnection _connection;
        private readonly IChannel _channel;
        private AsyncEventingBasicConsumer _consumer = null!;

        private RMQConsumerService(IConfiguration configuration, IConnection connection, IChannel channel)
        {
            _configuration = configuration;
            _connection = connection;
            _channel = channel;
        }

        public static async Task<RMQConsumerService> CreateAsync(IConfiguration configuration)
        {
            var factory = new ConnectionFactory
            {
                HostName = configuration["RabbitMQ:HostName"]!,
                UserName = configuration["RabbitMQ:UserName"]!,
                Password = configuration["RabbitMQ:Password"]!,
                ConsumerDispatchConcurrency = 1
            };

            var connection = await factory.CreateConnectionAsync();
            var channel = await connection.CreateChannelAsync();

            return new RMQConsumerService(configuration, connection, channel);
        }

        private async Task InitializeQueueAsync(string queueName)
        {
            await _channel.QueueDeclareAsync(queueName, durable: true, exclusive: false, autoDelete: false, arguments: null);
        }

        public async Task StartListeningAsync(string queueName, Func<string, ulong, Task> messageHandler)
        {
            await InitializeQueueAsync(queueName); // Declare queue dynamically

            _consumer = new AsyncEventingBasicConsumer(_channel);
            _consumer.ReceivedAsync += async (model, ea) =>
            {
                var body = ea.Body.ToArray();
                var message = Encoding.UTF8.GetString(body);
                Console.WriteLine($"Message received: {message}");

                try
                {
                    await messageHandler.Invoke(message, ea.DeliveryTag);
                }
                catch (Exception ex)
                {
                    Console.WriteLine($"Error processing message: {ex.Message}");
                    // Decide if you want to reject the message using BasicNack
                }
            };

            await _channel.BasicConsumeAsync(queue: queueName, autoAck: false, consumer: _consumer);
        }

        public async Task AckMessage(ulong deliveryTag)
        {
            Console.WriteLine($"Acking message: {deliveryTag}");
            await _channel.BasicAckAsync(deliveryTag, multiple: false);
        }

        public async ValueTask DisposeAsync()
        {
            if (_channel != null)
            {
                await _channel.CloseAsync();
                await _channel.DisposeAsync();
            }
            if (_connection != null)
            {
                await _connection.CloseAsync();
                await _connection.DisposeAsync();
            }
        }
    }
}
