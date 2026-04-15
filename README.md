# grpcannon

A lightweight gRPC load-testing CLI with configurable concurrency and latency histograms.

---

## Installation

```bash
go install github.com/yourusername/grpcannon@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/grpcannon.git && cd grpcannon && go build -o grpcannon .
```

---

## Usage

```bash
grpcannon [options] <target>
```

### Example

```bash
grpcannon \
  --proto ./api/service.proto \
  --call helloworld.Greeter/SayHello \
  --data '{"name": "world"}' \
  --concurrency 50 \
  --requests 1000 \
  localhost:50051
```

### Options

| Flag | Default | Description |
|------|---------|-------------|
| `--proto` | | Path to the `.proto` file |
| `--call` | | Fully qualified gRPC method name |
| `--data` | | JSON request payload |
| `--concurrency` | `10` | Number of concurrent workers |
| `--requests` | `200` | Total number of requests to send |
| `--timeout` | `30s` | Per-request timeout |

### Sample Output

```
Summary:
  Total requests:   1000
  Concurrency:      50
  Total time:       4.321s
  Requests/sec:     231.4

Latency histogram:
  p50:   18.2ms
  p90:   45.7ms
  p95:   62.1ms
  p99:  104.3ms
  max:  210.8ms
```

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)