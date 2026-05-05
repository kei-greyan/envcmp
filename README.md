# envcmp

A CLI tool to diff `.env` files across environments and flag missing or mismatched keys.

---

## Installation

```bash
go install github.com/yourusername/envcmp@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envcmp.git
cd envcmp
go build -o envcmp .
```

---

## Usage

Compare two `.env` files and highlight differences:

```bash
envcmp .env.development .env.production
```

**Example output:**

```
MISSING in .env.production:
  - DEBUG_MODE
  - LOCAL_DB_URL

MISMATCHED keys:
  - API_URL
      development: http://localhost:3000
      production:  https://api.example.com

OK: 12 keys match across both files
```

### Flags

| Flag | Description |
|------|-------------|
| `--keys-only` | Compare key names only, ignore values |
| `--strict` | Exit with non-zero status if any diff is found |
| `--format json` | Output results as JSON |

---

## Contributing

Pull requests and issues are welcome. Please open an issue before submitting large changes.

---

## License

MIT © [yourusername](https://github.com/yourusername)