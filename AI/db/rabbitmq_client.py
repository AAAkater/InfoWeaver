from datetime import datetime
from typing import Callable

from pika import BasicProperties, BlockingConnection, URLParameters
from pika.adapters.blocking_connection import BlockingChannel
from pika.exceptions import AMQPChannelError, AMQPConnectionError
from pydantic import BaseModel, Field

from core.config import settings
from utils.logger import logger


class FileUploadMessage(BaseModel):
    event: str = Field(..., description="Event type/name")
    file_id: int = Field(..., description="Unique file identifier")
    minio_path: str = Field(..., description="Path to file in MinIO storage")
    timestamp: datetime = Field(default_factory=datetime.now, description="Event timestamp")


class RabbitmqClient:
    """
    RabbitMQ client for managing connections, channels, and message operations.

    Provides methods for:
    - Publishing messages to queues/exchanges
    - Consuming messages from queues
    - Queue management (declare, delete, purge)
    """

    def __init__(self, url: str):
        self.params = URLParameters(url)
        self.exchange = ""
        try:
            self.conn = BlockingConnection(self.params)
            logger.info("Connected to RabbitMQ")
        except AMQPConnectionError as e:
            logger.error(f"Failed to connect to RabbitMQ: {e}")
            exit(0)

    def channel(self) -> BlockingChannel:
        """
        Get a new channel from the connection.
        """
        return self.conn.channel()

    def publish(self, queue_name: str, message: FileUploadMessage):
        """
        Publish a message to a queue.

        Args:
            queue_name: Name of the queue to publish to
            message: Message to publish (Pydantic model, dict, or string)
        Returns:
            True if message was published successfully, False otherwise
        """

        channel = self.channel()
        try:
            # Declare queue
            channel.queue_declare(queue=queue_name, durable=True)
            # Publish message
            channel.basic_publish(
                exchange=self.exchange,
                routing_key=queue_name,
                body=message.model_dump_json(),
                properties=BasicProperties(content_type="application/json", delivery_mode=2),
            )
            logger.info(f"Message published to queue '{queue_name}'")
        except AMQPChannelError as e:
            logger.error(f"Failed to publish message to queue '{queue_name}': {e}")
            return
        finally:
            if channel.is_open:
                channel.close()

    def consume(self, queue_name: str, callback: Callable[[FileUploadMessage], None]) -> None:
        """
        Consume messages from a queue.

        Args:
            queue_name: Name of the queue to consume from
            callback: Function to call for each message (receives message body as string)
        """

        channel = self.channel()

        def on_message(ch: BlockingChannel, method, properties, body):
            """Internal message handler"""
            try:
                # Decode message
                message = body.decode("utf-8")
                logger.info(f"Received message from queue '{queue_name}': {message}")
                # Call user callback with parsed FileUploadMessage
                callback(FileUploadMessage.model_validate_json(message))
                # Acknowledge message
                ch.basic_ack(delivery_tag=method.delivery_tag)
                logger.info(f"Message acknowledged from queue '{queue_name}'")
            except Exception as e:
                logger.error(f"Error processing message from queue '{queue_name}': {e}")
                ch.basic_nack(delivery_tag=method.delivery_tag, requeue=True)

        try:
            # Declare queue
            channel.queue_declare(queue=queue_name, durable=True)

            # Set QoS
            channel.basic_qos(prefetch_count=1)

            # Set up consumer
            channel.basic_consume(
                queue=queue_name,
                on_message_callback=on_message,
            )

            logger.info(f"Started consuming from queue '{queue_name}'. Waiting for messages...")
            channel.start_consuming()
        except AMQPChannelError as e:
            logger.error(f"Failed to consume from queue '{queue_name}': {e}")
            raise
        finally:
            if channel.is_open:
                channel.close()

    def delete_queue(self, queue_name: str, if_unused: bool = False, if_empty: bool = False) -> None:
        """
        Delete a queue.

        Args:
            queue_name: Name of the queue to delete
            if_unused: Only delete if the queue is unused
            if_empty: Only delete if the queue is empty
        """
        channel = self.channel()
        try:
            channel.queue_delete(queue=queue_name, if_unused=if_unused, if_empty=if_empty)
            logger.info(f"Queue '{queue_name}' deleted")
        except AMQPChannelError as e:
            logger.error(f"Failed to delete queue '{queue_name}': {e}")
            raise
        finally:
            if channel.is_open:
                channel.close()

    def purge_queue(self, queue_name: str) -> int:
        """
        Purge all messages from a queue.

        Args:
            queue_name: Name of the queue to purge

        Returns:
            Number of messages purged
        """
        channel = self.channel()
        try:
            result = channel.queue_purge(queue=queue_name)
            message_count = result.method.message_count
            logger.info(f"Purged {message_count} messages from queue '{queue_name}'")
            return message_count
        except AMQPChannelError as e:
            logger.error(f"Failed to purge queue '{queue_name}': {e}")
            return 0
        finally:
            if channel.is_open:
                channel.close()

    def close(self) -> None:
        """
        Close the connection to RabbitMQ.
        """
        if self.conn and not self.conn.is_closed:
            self.conn.close()
            logger.info("Disconnected from RabbitMQ")


rabbitmq_client = RabbitmqClient(settings.RABBITMQ_URL)
