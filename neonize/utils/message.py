from ..proto import Neonize_pb2 as neonize_proto
from ..proto.waE2E.WAWebProtobufsE2E_pb2 import Message, PollUpdateMessage
from ..types import MessageWithContextInfo


def get_message_type(message: Message) -> str:
    """
    Determines the type of message.

    :param message: The message object.
    :type message: Message
    :raises IndexError: If the message type cannot be determined.
    :return: The type of the message.
    :rtype: str
    """
    msg_fields = message.ListFields()
    field_name = msg_fields[0][0].name
    return field_name


def extract_text(message: Message) -> str:
    """
    Extracts text content from a message.

    :param message: The message object.
    :type message: Message
    :return: The extracted text content.
    :rtype: str
    """
    msg_fields = message.ListFields()

    _, field_value = msg_fields[0]
    if isinstance(field_value, str):
        return field_value
    text_attrs = ["text", "caption", "name", "conversation"]
    for attr in text_attrs:
        if hasattr(field_value, attr):
            val = getattr(field_value, attr)
            if isinstance(val, str) and val.strip():
                return val
    return ""


def get_poll_update_message(message: neonize_proto.Message) -> PollUpdateMessage | None:
    """
    Extracts pollUpdateMessage from event Message
    :param message: The message object.
    :type message: neonize_proto.Message
    :return: The extracted poll update message.
    :rtype: PollUpdateMessage
    """
    msg = message.Message
    if msg.pollUpdateMessage.ListFields():
        pollUpdateMessage: PollUpdateMessage = msg.pollUpdateMessage
        return pollUpdateMessage


def message_has_contextinfo(message: Message) -> bool:
    for field_name, msg in message.ListFields():
        if field_name.name.endswith("Message"):
            break
    else:
        return False
    return type(msg) in MessageWithContextInfo.__constraints__
