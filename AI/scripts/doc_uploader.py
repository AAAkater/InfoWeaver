from core.config import settings
from db.rabbitmq_client import FileUploadMessage, rabbitmq_client


def main():
    for i in range(5):
        rabbitmq_client.publish(
            settings.RABBITMQ_QUEUE,
            FileUploadMessage(
                event=f"File uploaded {i}",
                file_id=i,
                minio_path=f"files/file_{i}.txt",
            ),
        )


if __name__ == "__main__":
    main()
