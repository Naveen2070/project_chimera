import asyncio
from aio_pika import connect, Message, IncomingMessage
from aio_pika.abc import AbstractConnection, AbstractChannel, AbstractQueue


class RpcConsumer:
    connection: AbstractConnection
    channel: AbstractChannel
    queue: AbstractQueue

    def __init__(self, amqp_url: str, queue_name: str) -> None:
        """
        Initialize the RPC consumer.
        :param amqp_url: The RabbitMQ server URL.
        :param queue_name: The name of the queue to listen to for incoming messages.
        """
        self.amqp_url = amqp_url
        self.queue_name = queue_name

    async def connect(self) -> None:
        """
        Establish connection to RabbitMQ and declare the queue.
        """
        self.connection = await connect(self.amqp_url)
        self.channel = await self.connection.channel()
        self.queue = await self.channel.declare_queue(self.queue_name, durable=True)
        print(f"Connected to queue: {self.queue_name}")

    async def on_request(self, message: IncomingMessage) -> None:
        """
        Handle incoming RPC requests, process them, and send a response to the `reply_to` queue.
        :param message: The incoming message from the queue.
        """
        async with message.process():
            print(f"Received request: {message.body.decode()}")

            # Process the incoming message to create a response
            request_data = message.body.decode()
            response_data = await self.process_request(request_data)

            if not message.reply_to:
                print("Error: `reply_to` queue not specified in the message")
                return

            await self.channel.default_exchange.publish(
                Message(
                    body=response_data.encode(),
                    content_type="text/plain",
                    correlation_id=message.correlation_id,
                ),
                routing_key=message.reply_to,
            )
            print(
                f"Sent response to {message.reply_to} with correlation_id: {message.correlation_id}"
            )

    async def process_request(self, request_data: str) -> str:
        """
        Process the incoming request and generate a response.
        This function should contain the business logic.
        :param request_data: The incoming request data.
        :return: The response data.
        """
        return f"Processed: {request_data}"

    async def start(self) -> None:
        """
        Start consuming messages from the queue.
        """
        await self.queue.consume(self.on_request, no_ack=False)
        print(f"Waiting for RPC requests on queue: {self.queue_name}")

    async def stop(self) -> None:
        """
        Stop consuming messages from the queue.
        """
        await self.queue.cancel()
        await self.channel.close()
        await self.connection.close()
        print(f"Stopped consuming from queue: {self.queue_name}")
