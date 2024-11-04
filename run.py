import os
import shutil
import subprocess
import argparse
import sys

def run_native(address=None, port=None):
    command = ["go", "run", "cmd/native/main.go"]

    if address:
        command += ["--address", address]
    if port:
        command += ["--port", str(port)]

    result = subprocess.run(command, text=True)
    if result.returncode != 0:
        sys.exit(result.returncode)

def run_server():
    # build `game.wasm` file
    env = os.environ.copy()
    env.update({
        "GOOS": "js",
        "GOARCH": "wasm"
    })
    result = subprocess.run(["go", "build", "-o", "server/static/game.wasm", "cmd/web/main.go"], 
                            text=True, env=env)
    if result.returncode != 0:
        print(result.stderr)
        sys.exit(result.returncode)

    # copy `watch_exec.js` to `server/static/`
    goroot = subprocess.check_output(["go", "env", "GOROOT"]).decode().strip()
    wasm_exec_path = os.path.join(goroot, "misc", "wasm", "wasm_exec.js")
    shutil.copy(wasm_exec_path, "server/static/")

    subprocess.run(["go", "run", "cmd/server/main.go"], check=True)

def main() -> None:
    parser = argparse.ArgumentParser(description='Run Go applications.')
    
    # Define subparsers for the different commands
    subparsers = parser.add_subparsers(dest='command', required=True)

    # Create a parser for the 'native' command
    native_parser = subparsers.add_parser('native', help='Run the native Go application')
    native_parser.add_argument('--address', type=str, help='The address to bind to')
    native_parser.add_argument('--port', type=str, help='The port to bind to')

    # Create a parser for the 'server' command
    server_parser = subparsers.add_parser('server', help='Run the server Go application')

    # Parse the arguments
    args = parser.parse_args()

    if args.command == 'native':
        run_native(args.address, args.port)
    elif args.command == 'server':
        run_server()

if __name__ == "__main__":
    main()
