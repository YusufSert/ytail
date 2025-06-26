## ✅ Example Matches

- `log-2024-10-10.txt`
- `log-2025-01-01.txt`
- `log-1999-12-31.txt`

## ❌ Example Non-Matches

- `logfile-2024-10-10.txt`
- `log-2024-10.txt`
- `log-2024-10-10.log`
- `log-20241010.txt`

---

## 🧪 Regex Pattern

```regex
^log-\d{4}-\d{2}-\d{2}\.txt$

| Regex Component | Meaning                                      |
| --------------- | -------------------------------------------- |
| `^`             | Anchors the match to the **start** of string |
| `log-`          | Matches the literal prefix **`log-`**        |
| `\d{4}`         | Matches **4 digits** (the year)              |
| `-`             | Matches a hyphen `-`                         |
| `\d{2}`         | Matches **2 digits** (month)                 |
| `-`             | Matches another hyphen `-`                   |
| `\d{2}`         | Matches **2 digits** (day)                   |
| `\.txt`         | Matches `.txt` (escaped `.`)                 |
| `$`             | Anchors the match to the **end** of string   |
