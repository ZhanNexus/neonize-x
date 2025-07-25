import ctypes
import logging

from ..proto.Neonize_pb2 import LogEntry

try:
    from colorlog import ColoredFormatter
except Exception:
    ColoredFormatter = None

log = logging.getLogger(__name__)

if ColoredFormatter:
    formatter = ColoredFormatter(
        "%(asctime)s.%(msecs)03d %(log_color)s%[%(name)s %(levelname)s] - %(message)s%(reset)s",
        datefmt="%H:%M:%S",
        log_colors={
            "INFO": "cyan",
            "WARNING": "yellow",
            "ERROR": "red",
            "CRITICAL": "bold_red",
        },
    )
    stream_handler = logging.StreamHandler()
    stream_handler.setFormatter(formatter)
else:
    stream_handler = logging.StreamHandler()

logging.basicConfig(
    format="%(asctime)s.%(msecs)03d [%(name)s %(levelname)s] - %(message)s",
    datefmt="%H:%M:%S",
    level=logging.INFO,
    handlers=[stream_handler],
)

clientlogger = logging.getLogger("whatsmeow.Client")
dblogger = logging.getLogger("Whatsmeow.Database")


def log_whatsmeow(binary: int, size: int):
    log_msg = LogEntry.FromString(ctypes.string_at(binary, size))
    if log_msg.Name == "Client":
        log = clientlogger
    elif log_msg.Name == "Database":
        log = dblogger
    else:
        log = logging.getLogger(f"whatsmeow.{log_msg.Name}")
    getattr(log, log_msg.Level.lower())(log_msg.Message)
