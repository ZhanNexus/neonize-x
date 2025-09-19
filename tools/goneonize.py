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

def __build():
    args = argparse.ArgumentParser()
    args.add_argument("--os", default=platform.system().lower())
    args.add_argument("--arch", default=platform.machine().lower())
    parse = args.parse_args()
    filename = generated_name(parse.os, parse.arch)
    
    for sh in shell:
        subprocess.call(shlex.split(sh), cwd=cwd)
    
    if (Path(cwd) / "defproto").exists():
        shutil.rmtree(f"{cwd}/defproto")
    os.mkdir(f"{cwd}/defproto")
    os.rename(
        f"{cwd}/github.com/krypton-byte/neonize/defproto/",
        f"{cwd}/defproto")
    shutil.rmtree(f"{cwd}/github.com")
    
    # FIX: Environment management
    build_env = os.environ.copy()
    build_env["CGO_ENABLED"] = "1"
    
    subprocess.call(
        shlex.split(f"go build -buildmode=c-shared -ldflags=-s -o {filename} main.go"),
        cwd=cwd,
        env=build_env,  # Use the corrected environment
    )
    
    if (Path(cwd).parent / filename).exists():
        os.remove(os.path.dirname(cwd) + "/" + filename)
    os.rename(f"{cwd}/{filename}", os.path.dirname(cwd) + "/" + filename)

def build_proto():
    with open(cwd + "/Neonize.proto", "rb") as file:
        with open(cwd + "/defproto/Neonize.proto", "wb") as wf:
            wf.write(file.read())
    for sh in shell:
        subprocess.call(shlex.split(sh), cwd=cwd + "/defproto")

def build_neonize():
    os_name = os.environ.get("GOOS") or platform.system().lower()
    arch_name = os.environ.get("GOARCH") or platform.machine().lower()
    arch_name = arch_normalizer(arch_name)
    
    print(f"Building for: os={os_name}, arch={arch_name}")
    
    # FIX: Proper environment setup
    build_env = os.environ.copy()
    build_env["CGO_ENABLED"] = "1"
    build_env["GOOS"] = os_name
    build_env["GOARCH"] = arch_name
    
    # Set CC compiler for Linux if not set
    if os_name == "linux" and "CC" not in build_env:
        if shutil.which("gcc"):
            build_env["CC"] = "gcc"
        else:
            print("Warning: Could not find 'gcc'. Build might fail.")
    
    filename = generated_name(os_name, arch_name)
    print(f"Output file: {filename}")
    
    # FIX: Use the corrected environment and proper build command
    result = subprocess.call(
        shlex.split(f"go build -buildmode=c-shared -ldflags='-s -w' -o {filename} ."),
        cwd=cwd,
        env=build_env,
    )
    
    if result != 0:
        print("Build failed!")
        exit(1)
    
    # Move the built file
    target_dir = Path(cwd).parent / "neonize"
    target_path = target_dir / filename
    source_path = Path(cwd) / filename
    
    if not target_dir.exists():
        target_dir.mkdir(parents=True)
    
    if target_path.exists():
        target_path.unlink()
    
    shutil.move(str(source_path), str(target_path))
    print(f"Build successful! Library moved to: {target_path}")

def build():
    args = argparse.ArgumentParser()
    sub = args.add_subparsers(dest="build", required=True)
    sub.add_parser("goneonize")
    sub.add_parser("proto")
    sub.add_parser("all")
    parse = args.parse_args()
    
    match parse.build:
        case "goneonize":
            build_neonize()
        case "proto":
            build_proto()
        case "all":
            build_proto()
            build_neonize()

def build_android():
    filename = generated_name("android", "arm64")  # Fixed architecture name
    
    for sh in shell:
        subprocess.call(shlex.split(sh), cwd=cwd)
    
    if (Path(cwd) / "defproto").exists():
        shutil.rmtree(f"{cwd}/defproto")
    os.mkdir(f"{cwd}/defproto")
    os.rename(
        f"{cwd}/github.com/krypton-byte/neonize/defproto/",
        f"{cwd}/defproto")
    shutil.rmtree(f"{cwd}/github.com")
    
    # FIX: Proper environment setup for Android
    build_env = os.environ.copy()
    build_env.update({
        "CGO_ENABLED": "1",
        "GOOS": "android",
        "GOARCH": "arm64",
        "CC": "/home/krypton-byte/Pictures/android-ndk-r26b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android28-clang",
        "CXX": "/home/krypton-byte/Pictures/android-ndk-r26b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android28-clang++",
    })
    
    subprocess.call(
        shlex.split(f"go build -buildmode=c-shared -ldflags=-s -o {filename} main.go"),
        cwd=cwd,
        env=build_env,  # Use the corrected environment
    )
    
    if (Path(cwd).parent / filename).exists():
        os.remove(os.path.dirname(cwd) + "/" + filename)
    os.rename(f"{cwd}/{filename}", os.path.dirname(cwd) + "/" + filename)

if __name__ == "__main__":
    build()