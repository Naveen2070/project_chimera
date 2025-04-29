//	Copyright 2025 Naveen R
//
//	Licensed under the Apache License, Version 2.0 (the "License");
//	you may not use this file except in compliance with the License.
//	You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS,
//	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	See the License for the specific language governing permissions and
//	limitations under the License.
using RabbitMQ.Client;
using System.Text;
using user_service.Services.Interfaces;

namespace user_service.Services
{
    public class RMQPublisherService : IRMQPublisherService, IAsyncDisposable
    {
        private readonly IConfiguration _configuration;
        private readonly IConnection _connection;
        private readonly IChannel _channel;

        private RMQPublisherService(IConfiguration configuration, IConnection connection, IChannel channel)
        {
            _configuration = configuration;
            _connection = connection;
            _channel = channel;
        }

        public static async Task<RMQPublisherService> CreateAsync(IConfiguration configuration)
        {
            var factory = new ConnectionFactory
            {
                HostName = configuration["RabbitMQ:HostName"] ?? "localhost",
                UserName = configuration["RabbitMQ:UserName"] ?? "guest",
                Password = configuration["RabbitMQ:Password"] ?? "guest"
            };

            var connection = await factory.CreateConnectionAsync();
            var channel = await connection.CreateChannelAsync();

            return new RMQPublisherService(configuration, connection, channel);
        }

        public async Task SendMessageAsync( string message)
        {
            var queueName = _configuration["RabbitMQ:QueueName"] ?? "error_dump_queue";
            await _channel.QueueDeclareAsync(
                queue: queueName,
                durable: true,
                exclusive: false,
                autoDelete: false,
                arguments: null
            );

            var body = Encoding.UTF8.GetBytes(message);

            await _channel.BasicPublishAsync(
                exchange: "",
                routingKey: queueName,
                body: body
            );

            Console.WriteLine($"[RabbitMQ] Sent message to queue '{queueName}': {message}");
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
