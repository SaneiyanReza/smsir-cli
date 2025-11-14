# 📱 SMS.ir CLI

> A simple message can connect worlds with a single command

A professional command-line tool for interacting with SMS.ir APIs, featuring both a command-line interface (CLI) and an interactive terminal user interface (TUI).

## ✨ Features

- 🎨 **Interactive UI**: Modern terminal interface with smooth animations and intuitive navigation
- 📤 **Send SMS Messages**: Send bulk SMS messages via command-line or interactive UI
- 📊 **Dashboard & Statistics**: View your credit balance and available lines in real-time
- 🌐 **Full Persian/Farsi Support**: Proper UTF-8 handling for Persian text input and display
- 🔧 **Easy Configuration**: Simple setup and management of API credentials
- 🎯 **Dual Interface**: Choose between CLI commands or interactive TUI mode

## 🚀 Quick Start

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/SaneiyanReza/smsir-cli.git
   cd smsir-cli
   ```

2. **Build the project:**
   ```bash
   go build -o smsir cmd/smsir
   ```

3. **Install globally (optional):**
   ```bash
   # On Linux/macOS
   sudo mv smsir /usr/local/bin/
   
   # On Windows, add to PATH
   ```

### Configuration

First, set up your API credentials:

```bash
smsir config set --api-key YOUR_SMS.ir_API_KEY --line YOUR_LINE_NUMBER
```

Or use the interactive configuration:

```bash
smsir menu
# Then select "🔧 Configure API Key & Line Number"
```

## 📖 Usage

### Command Line Mode

#### Send SMS

```bash
# Send to a single number
smsir send -m "Hello World" -t "09120000000"

# Send to multiple numbers
smsir send -m "Hello" -t "09120000000,09121111111"

# With custom line number
smsir send -m "Test" -t "09120000000" -l 90001234
```

#### Check Credit

```bash
smsir credit
```

#### View Available Lines

```bash
smsir lines
```

#### View Configuration

```bash
smsir config show
```

#### Validate Configuration

```bash
smsir config validate
```

### Interactive UI Mode

Launch the interactive menu:

```bash
smsir menu
```

The interactive menu provides:

- 🔧 **Configure API Key & Line Number**: Step-by-step configuration wizard
- 📤 **Send SMS**: Interactive SMS sending with Persian text support
- 📊 **Dashboard**: Real-time view of credit and available lines
- 💻 **Command Line Mode**: Quick access to command help

## 📋 Available Commands

| Command | Description | Flags |
|---------|-------------|-------|
| `config` | Configuration management | `set`, `show`, `validate` |
| `send` | Send SMS message | `-m, --message`, `-t, --to`, `-l, --line` |
| `credit` | Show current credit balance | - |
| `lines` | Show available lines | - |
| `menu` | Launch interactive menu | - |

### Command Details

#### `smsir config`

Manage your API configuration.

```bash
# Set API credentials
smsir config set --api-key YOUR_SMS.ir_KEY --line YOUR_LINE

# Show current configuration (masked)
smsir config show

# Validate configuration
smsir config validate
```

#### `smsir send`

Send SMS messages to one or more recipients.

**Required flags:**
- `-m, --message`: Message text to send
- `-t, --to`: Comma-separated list of mobile numbers

**Optional flags:**
- `-l, --line`: Line number (uses configured line if not provided)

**Examples:**
```bash
# Basic usage
smsir send -m "Hello from SMS.ir CLI" -t "09120000000"

# Multiple recipients
smsir send -m "Bulk message" -t "09120000000,09111111111"

# With Persian text (use quotes)
smsir send -m "سلام دنیا" -t "09120000000"
```

#### `smsir credit`

Display your current SMS credit balance.

```bash
smsir credit
# Output: 💰 Current Credit: 1000 SMS
```

#### `smsir lines`

List all available line numbers.

```bash
smsir lines
# Output: 📞 Available Lines:
#          1. 90001234
#          2. 9981234
```

## 🎨 UI Features

### Interactive Dashboard

The dashboard provides a real-time view of:
- 💰 Current credit balance
- 📞 Available line numbers
- 🔄 Refresh capability (press `r`)

### Send SMS UI

The interactive SMS sending interface features:
- Step-by-step wizard
- Persian/Farsi text input support
- Clipboard paste support (Ctrl+V)
- Real-time validation
- Success/error feedback with detailed results

### Configuration Wizard

Easy setup with:
- Secure input handling
- Validation feedback
- Progress indicators

## 📥 Installation

### 1. Download the latest release

Go to the **[Releases page](https://github.com/SaneiyanReza/smsir-cli/releases/latest)** and download the binary for your OS:

| OS      | File                     |
|---------|--------------------------|
| Windows | `smsir-cli-windows.exe` |
| Linux   | `smsir-cli-linux`       |
| macOS   | `smsir-cli-macos`       |

## 🛠️ Development

### Prerequisites

- Go 1.21 or higher
- SMS.ir API credentials

### Building from Source

```bash
# Clone the repository
git clone https://github.com/SaneiyanReza/smsir-cli.git
cd smsir-cli

# Install dependencies
go mod download

# Build
go build -o smsir cmd/smsir

# Run
./smsir --help
```
## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. Here's how:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [SMS.ir](https://sms.ir) for providing the SMS API
- [Cobra](https://github.com/spf13/cobra) for CLI framework
- [Bubbletea](https://github.com/charmbracelet/bubbletea) for TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) for styling

## 📞 Support

If you have any questions or encounter issues:

- Open an issue on [GitHub](https://github.com/SaneiyanReza/smsir-cli/issues)
- Check the [SMS.ir API Documentation](https://sms.ir/web-service)

---

**A simple message can connect worlds with a single command** 🌍

Made with ❤️ using Go