import os
import platform
import shlex
import argparse
import subprocess
import shutil
from pathlib import Path
from typing import Dict
import glob

cwd = (Path(__file__).parent.parent / "goneonize/").__str__()

# Perintah protoc
shell = [
    "protoc --go_out=. --go_opt=paths=source_relative Neonize.proto",
    "protoc --python_out=../../neonize/proto --mypy_out=../../neonize/proto Neonize.proto",
    *[
        f"protoc --python_out=../../neonize/proto --mypy_out=../../neonize/proto {path}"
        for path in glob.glob("*/*.proto", root_dir=cwd + "/defproto")
    ],
]


def arch_normalizer(arch_: str) -> str:
    arch: Dict[str, str] = {
        "aarch64": "arm64",
        "x86_64": "amd64",
    }
    return arch.get(arch_, arch_)


def generated_name(os_name="", arch_name=""):
    os_name = os_name or platform.system().lower()
    arch_name = arch_normalizer(arch_name or platform.machine().lower())
    if os_name == "windows":
        ext = "dll"
    elif os_name == "linux":
        ext = "so"
    elif os_name == "darwin":
        ext = "dylib"
    else:
        ext = "so"
    return f"neonize-{os_name}-{arch_name}.{ext}"


def build_proto():
    # copy Neonize.proto ke defproto
    with open(cwd + "/Neonize.proto", "rb") as file:
        with open(cwd + "/defproto/Neonize.proto", "wb") as wf:
            wf.write(file.read())

    # jalankan semua perintah protoc
    for sh in shell:
        subprocess.call(shlex.split(sh), cwd=cwd + "/defproto")


def build_neonize():
    os_name = os.environ.get("GOOS") or platform.system().lower()
    arch_name = os.environ.get("GOARCH") or platform.machine().lower()
    print(f"os: {os_name}, arch: {arch_name}")

    filename = generated_name(os_name, arch_name)
    print("output:", filename)

    env = os.environ.copy()
    env.update({"CGO_ENABLED": "1"})

    # compile Go jadi .so
    subprocess.call(
        shlex.split(
            f"go build -buildmode=c-shared -ldflags=-s -o {filename} main.go"
        ),
        cwd=cwd,
        env=env,
    )

    target = Path(cwd).parent / "neonize" / filename
    if target.exists():
        target.unlink()

    os.rename(f"{cwd}/{filename}", target)


def build_android():
    filename = generated_name("android", "arm64")

    # prepare environment Android NDK
    env = os.environ.copy()
    env.update({
        "CGO_ENABLED": "1",
        "CC": "/home/krypton-byte/Pictures/android-ndk-r26b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android28-clang",
        "CXX": "/home/krypton-byte/Pictures/android-ndk-r26b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android28-clang++",
    })

    for sh in shell:
        subprocess.call(shlex.split(sh), cwd=cwd)

    if (Path(cwd) / "defproto").exists():
        shutil.rmtree(f"{cwd}/defproto")
    os.mkdir(f"{cwd}/defproto")
    os.rename(
        f"{cwd}/github.com/krypton-byte/neonize/defproto/",
        f"{cwd}/defproto"
    )
    shutil.rmtree(f"{cwd}/github.com")

    subprocess.call(
        shlex.split(
            f"go build -buildmode=c-shared -ldflags=-s -o {filename} main.go"
        ),
        cwd=cwd,
        env=env,
    )

    target = Path(cwd).parent / filename
    if target.exists():
        target.unlink()

    os.rename(f"{cwd}/{filename}", target)


def build():
    args = argparse.ArgumentParser()
    sub = args.add_subparsers(dest="build", required=True)
    sub.add_parser("goneonize")
    sub.add_parser("proto")
    sub.add_parser("all")
    sub.add_parser("android")
    parse = args.parse_args()

    match parse.build:
        case "goneonize":
            build_neonize()
        case "proto":
            build_proto()
        case "all":
            build_proto()
            build_neonize()
        case "android":
            build_android()


if __name__ == "__main__":
    build()
