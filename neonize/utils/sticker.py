import json
import os
import tempfile
import uuid

from .ffmpeg import AFFmpeg
from .iofile import TemporaryFile
from .platform import is_executable_installed


def add_exif(name: str = "", packname: str = "") -> bytes:
    """
    Adds EXIF metadata to a sticker pack.

    :param name: Name of the sticker pack, defaults to an empty string.
    :type name: str, optional
    :param packname: Publisher of the sticker pack, defaults to an empty string.
    :type packname: str, optional
    :return: Byte array containing the EXIF metadata.
    :rtype: bytes
    """
    json_data = {
        "sticker-pack-id": "com.snowcorp.stickerly.android.stickercontentprovider b5e7275f-f1de-4137-961f-57becfad34f2",
        "sticker-pack-name": name,
        "sticker-pack-publisher": packname,
        "android-app-store-link": "https://play.google.com/store/apps/details?id=com.marsvard.stickermakerforwhatsapp",
        "ios-app-store-link": "https://itunes.apple.com/app/sticker-maker-studio/id1443326857",
    }

    exif_attr = bytes.fromhex("49 49 2A 00 08 00 00 00 01 00 41 57 07 00 00 00 00 00 16 00 00 00")
    json_buffer = json.dumps(json_data).encode("utf-8")
    exif = exif_attr + json_buffer
    exif_length = len(json_buffer)
    exif = exif[:14] + exif_length.to_bytes(4, "little") + exif[18:]
    return exif


def webpmux_is_installed():
    return is_executable_installed("webpmux")


max_sticker_size = 512000
webpmux_is_available = False
if webpmux_is_installed():
    max_sticker_size = 712000
    webpmux_is_available = True


async def convert_to_sticker(
    file: bytes, name="", packname="", enforce_not_broken=False, animated_gif=False, is_webm=False
):
    async with AFFmpeg(file) as ffmpeg:
        sticker = await ffmpeg.cv_to_webp(
            enforce_not_broken=enforce_not_broken,
            animated_gif=animated_gif,
            max_sticker_size=max_sticker_size,
            is_webm=is_webm,
        )
    if not webpmux_is_available:
        return sticker, False

    exif_filename = TemporaryFile(prefix=None, touch=False).__enter__()
    with open(exif_filename.path, "wb") as file:
        file.write(add_exif(name=name, packname=packname))
    temp = tempfile.gettempdir() + "/" + uuid.uuid4().__str__() + ".webp"
    async with AFFmpeg(sticker) as ffmpeg:
        cmd = [
            "webpmux",
            "-set",
            "exif",
            exif_filename.path.__str__(),
            ffmpeg.filepath,
            "-o",
            temp,
        ]
        await ffmpeg.call(cmd)
    exif_filename.__exit__(None, None, None)
    with open(temp, "rb") as file:
        buf = file.read()
    os.remove(temp)
    return buf, True
