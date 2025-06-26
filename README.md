
# YTAIL

### `ytail` reads log lines from files and sends them to Loki.

---

## 📁 ytail Configuration

This document explains the configuration used for the **ytail** log tailer and its client that pushes logs to Loki.

> 🆕 **Note:** When a new log file matching the pattern is added to the log directory, `ytail` will automatically detect it and start reading from it without needing a restart.

---

## 🐾 Tailer Configuration

| Key           | Description                                            |
|---------------|--------------------------------------------------------|
| `scrape_path` | Directory path where log files are located.           |
| `file_regex`  | Regular expression to match log filenames.            |

---

## 🔁 Client Configuration

| Key              | Description                                        |
|------------------|----------------------------------------------------|
| `retry`          | Number of retries on failure.                      |
| `backoff`        | Initial wait time before retrying after a failure. |
| `max_backoff`    | Maximum wait time between retries.                 |
| `push_url`       | URL endpoint for pushing logs to Loki.             |
| `batch_max_size` | Maximum number of log lines per batch.             |
| `batch_max_wait` | Maximum time to wait before sending a batch.       |
| `labels`         | Key-value pairs added as labels to every log line. |

### Example Labels
```yaml
labels:
  service_name: test
```

---

## ✅ Example Matches

These filenames match the `file_regex` pattern:

- `log-2024-10-10.txt`
- `log-2025-01-01.txt`
- `log-1999-12-31.txt`

---

## ❌ Example Non-Matches

These filenames **do not** match the `file_regex` pattern:

- `logfile-2024-10-10.txt`
- `log-2024-10.txt`
- `log-2024-10-10.log`
- `log-20241010.txt`

---

## 🧪 Regex Pattern Breakdown

### Pattern: `^log-\d{4}-\d{2}-\d{2}\.txt$`

| Regex Component | Meaning                                        |
|-----------------|------------------------------------------------|
| `^`             | Anchors the match to the **start** of string   |
| `log-`          | Matches the literal prefix **`log-`**          |
| `\d{4}`        | Matches **4 digits** (the year)                |
| `-`             | Matches a hyphen `-`                           |
| `\d{2}`        | Matches **2 digits** (month)                   |
| `-`             | Matches another hyphen `-`                     |
| `\d{2}`        | Matches **2 digits** (day)                     |
| `\.txt`        | Matches `.txt` (escaped dot)                   |
| `$`             | Anchors the match to the **end** of string     |

---
## To build the tool
* Clone the directory
* Install golang:1.24
* run "go build -o ./bin/ytail ./cmd"
* execute ytail commnad with -help commnad (if you see this output, you are good to go !!!)

  | Option         | Description                                                          |
  | -------------- | -------------------------------------------------------------------- |
  | `-log.lvl`     | Log level. Available log levels:                                     |
  |                | - `-4` = DEBUG                                                       |
  |                | - `0`  = INFO                                                        |
  |                | - `4`  = WARN                                                        |
  |                | - `8`  = ERROR                                                       |
  | `-config.path` | Path to the configuration file (e.g. `ytail-config.yaml`). Required. |


## 🏁 Summary

This config sets up `ytail` to:

- Monitor log files matching `log-YYYY-MM-DD.txt` in the specified directory.
- Automatically start reading newly added matching files.
- Push logs to a Loki server via HTTP.
- Use client-side batching and retries.
- Add helpful labels such as `service_name` to each log entry.
- ✅ Includes an example configuration file: `ytail-config.yaml` for reference and easy setup.




