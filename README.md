# portwatch

Lightweight CLI daemon that monitors open ports and alerts on unexpected changes.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start the daemon with a default scan interval of 30 seconds:

```bash
portwatch start
```

Specify a custom interval and define allowed ports:

```bash
portwatch start --interval 60 --allow 22,80,443
```

When an unexpected port is detected, portwatch will alert you in the terminal:

```
[ALERT] New open port detected: 8080 (PID: 3921, Process: node)
[ALERT] Port closed: 443
```

### Commands

| Command | Description |
|---|---|
| `start` | Start the monitoring daemon |
| `snapshot` | Print a snapshot of currently open ports |
| `diff` | Compare current ports against the last snapshot |

### Flags

| Flag | Default | Description |
|---|---|---|
| `--interval` | `30` | Scan interval in seconds |
| `--allow` | none | Comma-separated list of expected ports |
| `--log` | stdout | Path to log file |

## Requirements

- Go 1.21+
- Linux or macOS

## License

MIT © 2024 yourusername