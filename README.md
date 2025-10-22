### Toolchain Installation

Swahili requires the LLVM and Clang toolchain for compilation. Please follow the instructions below depending on your operating system.

#### Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install -y clang llvm
```

#### Linux (Fedora)
```bash
sudo dnf install clang llvm
```

#### Linux (Arch)
```bash
sudo pacman -S clang llvm
```

#### macOS (Homebrew)
```bash
brew update
brew install llvm
# Clang is provided with llvm via Homebrew. For system clang, install Xcode Command Line Tools:
xcode-select --install
```

After installation, you can check your clang and llvm versions:
```bash
clang --version
llvm-config --version
```

Make sure the installed tools are in your PATH. If you installed via Homebrew on macOS, you may need to add the following to your shell profile:
```bash
export PATH="/usr/local/opt/llvm/bin:$PATH"
```

For more details, refer to the official documentation:
- [LLVM Getting Started](https://llvm.org/docs/GettingStarted.html)
- [Clang Documentation](https://clang.llvm.org/docs/)

---

### Setup
```
cd lang
./dev.sh

Swahili Programming Environment

Usage:
  swa [command]

Available Commands:
  compile     Compile the source code to an executable
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  interpret   Interpret the source code
  server      Start the web service
  tokenize    Tokenize the source code

Flags:
  -h, --help   help for swa

Use "swa [command] --help" for more information about a command.

```
