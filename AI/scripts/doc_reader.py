from core.config import settings
from db.rabbitmq_client import FileUploadMessage, rabbitmq_client


def on_message(msg: FileUploadMessage):
    print(msg)


def main():
    rabbitmq_client.consume(settings.RABBITMQ_QUEUE, on_message)


if __name__ == "__main__":
    main()
