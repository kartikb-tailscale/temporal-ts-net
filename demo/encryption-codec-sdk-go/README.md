# SDK Go codec demo (prod-style wiring)

This demo validates:

- `temporal start-dev` extension runs and auto-starts its built-in codec server.
- Go SDK worker and starter use a remote codec endpoint through runtime config,
  not business logic changes.
- Workflow/activity code remains standard application code.

## Structure

- `internal/demo/workflow.go`: workflow and activity logic (no codec-specific code)
- `internal/platform/client.go`: Temporal client wiring (runtime codec endpoint)
- `worker/main.go`: worker bootstrap
- `starter/main.go`: workflow starter
- `run_demo.sh`: end-to-end runner

## Run

```bash
./run_demo.sh
```

The script:

1. Starts `temporal start-dev` (extension)
2. Uses extension codec endpoint at `http://127.0.0.1:8081` (or `CODEC_PORT` override)
3. Starts worker with `--codec-endpoint`
4. Starts workflow with `--codec-endpoint`

Expected output includes:

```text
workflow_id=encryption-demo-... run_id=... result="hello codec-demo"
```

## About key IDs

The extension's default codec transform is stateless `zlib + base64` and does
not require or validate key IDs. Any existing `encryption-key-id` metadata is
passed through unchanged.
