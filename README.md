## How to install

### Docker 
It's simpler to use the docker image if you do not want to install llvm manually.

```bash
alias swahili='docker run -w "$(pwd)" -v "$(pwd)":"$(pwd)" ghcr.io/papepathe/swa-lang:master'

swahili compile -s my-code.swa
swahili tokenize -s my-code.swa
swahili parse -s my-code.swa
```

### Manual 

Swahili requires the LLVM and Clang toolchain for compilation.
Please follow the instructions below depending on your operating system.

#### Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install -y clang llvm-19
```

#### Linux (Fedora)
```bash
sudo dnf install clang llvm-19
```

#### Linux (Arch)
```bash
sudo pacman -S clang llvm-19
```

#### macOS (Homebrew)
```bash
brew update
brew install llvm@19
xcode-select --install
```

After installation, you can check your clang and llvm versions:
```bash
clang --version
llvm-config --version
```

Make sure the installed tools are in your PATH. If you installed via 
Homebrew on macOS, you may need to add the following to your shell profile:
```bash
export PATH="/usr/local/opt/llvm/bin:$PATH"
```

For more details, refer to the official documentation:
- [LLVM Getting Started](https://llvm.org/docs/GettingStarted.html)
- [Clang Documentation](https://clang.llvm.org/docs/)
