# logslice

Fast log file slicer that extracts time-range segments from large structured logs.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

## Usage

Extract log entries within a specific time range:

```bash
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" --file app.log
```

Read from stdin and write to a file:

```bash
cat app.log | logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" > output.log
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--file` | Path to the input log file | stdin |
| `--from` | Start of the time range (RFC3339) | required |
| `--to` | End of the time range (RFC3339) | required |
| `--format` | Timestamp format in logs | `RFC3339` |
| `--field` | JSON field name containing the timestamp | `time` |

### Example Output

```
{"time":"2024-01-15T08:12:43Z","level":"info","msg":"server started"}
{"time":"2024-01-15T08:45:01Z","level":"warn","msg":"high memory usage"}
```

## Features

- Binary search on large log files for fast range extraction
- Supports structured (JSON) and plain-text log formats
- Reads from files or stdin
- Minimal memory footprint — streams output line by line

## License

MIT © 2024 yourusername