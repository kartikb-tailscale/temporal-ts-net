# temporal-start-dev-ext

Extension command for Temporal CLI that adds:

- `temporal start-dev` as a shortcut for `temporal server start-dev`
- automatic local codec server startup
- optional tailnet exposure with `--tailscale`

By default the built-in codec server applies a stateless `zlib + base64` envelope:

- `/encode`: zlib-compress payload bytes, then base64-encode them
- `/decode`: base64-decode payload bytes, then zlib-decompress them

The codec preserves and restores the original payload encoding metadata.
It does not require or validate key IDs and passes any existing
`encryption-key-id` metadata through unchanged.

Default codec server bind address is `127.0.0.1:8081`.
Override with `--codec-port` when needed.

## Install

Build the extension binary:

```bash
go build -o ./bin/temporal-start_dev ./cmd/temporal-start_dev
```

Add `./bin` to your `PATH` and verify discovery:

```bash
temporal help --all
```

You should see `start-dev` listed as an extension command.

## Usage

Start local dev server with codec server:

```bash
temporal start-dev
```

Expose dev server on Tailscale tailnet:

```bash
temporal start-dev \
    --tailscale \
    --tailscale-hostname your-dev-host
```

`--tsnet` and related `--tsnet-*` flags are also accepted aliases.

Pass any `temporal server start-dev` flags through directly:

```bash
temporal start-dev \
    --port 7234 \
    --ui-port 8234 \
    --db-filename /tmp/temporal-dev.db
```

## Extension flags

- `--tailscale`: enable tsnet listener and proxy
- `--tailscale-hostname`: tsnet hostname (default `temporal-dev`)
- `--tailscale-authkey`: auth key for non-interactive auth (or set `TS_AUTHKEY`)
- `--tailscale-state-dir`: local state dir for tsnet node
- `--codec-port`: local codec server port (default `8081`)

All non-extension flags are forwarded to `temporal server start-dev`.
