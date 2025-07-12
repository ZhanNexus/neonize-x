import re
from phonenumbers import parse, PhoneNumberFormat, format_number
from .log import log
from .message import extract_text, get_message_type
from .ffmpeg import FFmpeg
from .sticker import add_exif
from .thumbnail import save_file_to_temp_directory
from .iofile import get_bytes_from_name_or_url
from .calc import AspectRatioMethod
from .jid import build_jid, Jid2String, JIDToNonAD
from .enum import (
    MediaType,
    MediaTypeToMMS,
    BlocklistAction,
    ChatPresence,
    ChatPresenceMedia,
    ClientName,
    ClientType,
    ParticipantChange,
    ParticipantRequestChange,
    PrivacySetting,
    PrivacySettingType,
    ReceiptType,
)


def gen_vcard(name: str, phone_number: str) -> str:
    """
    Generates a vCard string for a contact.

    :param name: Name of the contact.
    :type name: str
    :param phone_number: Phone number of the contact.
    :type phone_number: str
    :return: vCard string for the contact.
    :rtype: str
    """
    inter_phone_number = format_number(
        parse(f"{'+' if phone_number[0] != '+' else ''}{phone_number}"),
        PhoneNumberFormat.INTERNATIONAL,
    )
    return (
        f"BEGIN:VCARD\nVERSION:3.0\nFN:{name}\nitem1.TEL;waid={phone_number}"
        f":{inter_phone_number}\nitem1.X-ABLabel:Ponsel\nEND:VCARD"
    )


def validate_link(link) -> bool:
    """
    Validates if the provided link is a valid URL.

    :param link: The URL to validate.
    :type link: str
    :return: True if the URL is valid, False otherwise.
    :rtype: bool
    """
    url_pattern = re.compile(
        r"^(https?|ftp)://"
        r"(?:(?:[A-Z0-9](?:[A-Z0-9-]{0,61}[A-Z0-9])?\.)+(?:[A-Z]{2,6}\.?|[A-Z0-9-]{2,}\.?)|"
        r"\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}|"
        r"\[?[A-F0-9]*:[A-F0-9:]+]?)"
        r"(?::\d+)?"
        r"(?:/?|[/?]\S+)$",
        re.IGNORECASE,
    )

    return bool(re.match(url_pattern, link))


__all__ = (
    "log",
    "get_message_type",
    "extract_text",
    "FFmpeg",
    "save_file_to_temp_directory",
    "get_bytes_from_name_or_url",
    "AspectRatioMethod",
    "build_jid",
    "Jid2String",
    "JIDToNonAD",
    "MediaType",
    "MediaTypeToMMS",
    "BlocklistAction",
    "ChatPresence",
    "ChatPresenceMedia",
    "ClientName",
    "ClientType",
    "ParticipantChange",
    "ParticipantRequestChange",
    "PrivacySetting",
    "PrivacySettingType",
    "ReceiptType",
    "add_exif",
    "validate_link",
    "gen_vcard",
)
