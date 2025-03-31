using notification_service.Services.Interfaces;
using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using System.Text;
using System.Threading.Tasks;

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

        public async Task StartListeningAsync(string queueName, Func<string, Task> messageHandler)
        {
            await InitializeQueueAsync(queueName); // Declare queue dynamically

            _consumer = new AsyncEventingBasicConsumer(_channel);
            _consumer.ReceivedAsync += async (model, ea) =>
            {
                Console.WriteLine($"Message received: {Encoding.UTF8.GetString(ea.Body.ToArray())}");
                var body = ea.Body.ToArray();
                var message = Encoding.UTF8.GetString(body);
                await messageHandler.Invoke(message);
            };

            await _channel.BasicConsumeAsync(queue: queueName, autoAck: true, consumer: _consumer);
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
