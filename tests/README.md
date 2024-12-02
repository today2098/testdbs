# Test Guide

Prepare `.env` file here.

```console
$ touch tests/.env
```

Set environment variables about DSN (**MySQL only**).

```bash
# DSN example (minimum)
DSN_TEST="user:password@tcp(localhost:3306)/"
```

Finally, do test.

```console
$ go test -shuffle=on ./...
```
