---
name: hatch
description: "Open a hatch: serve a dev server, static page, report, or prototype on a private tailnet URL with the hatch CLI. Use when asked to open, close, or extend a hatch, or to host, preview, or share anything web-based for the user to view from another device."
agent-config-sync: true
---

**Hatch** is a command-line tool on this machine (`hatch`). "Open a hatch" means: run the thing the user wants to see under `hatch up` and give them the URL it prints. The URL works from any device on the user's tailnet, including their laptop and phone.

It is not the in-app browser preview, `localhost`, a `file://` path, `tailscale serve`, cloudflared, or ngrok. Don't offer any of those as a substitute. If `command -v hatch` finds nothing, say so and stop.

## Open a hatch

### 1. Pick what runs and its name

Hatch assigns a free port (20000-29999) and passes it to the command as `PORT`. The app must listen on **`127.0.0.1:$PORT` only**. Hatch already holds the tailnet address on that port for TLS, so an app that binds every interface (`0.0.0.0`, `::`, or Go's `":"+port`) fails with "address already in use".

| What to show | Command after `hatch up ... --` |
|---|---|
| Static HTML, a report, mocks | `sh -c 'exec python3 -m http.server "$PORT" --bind 127.0.0.1 --directory <dir>'` |
| App that reads `PORT` and binds loopback | `<dev command>` |
| App that takes a port or host flag | `sh -c 'exec <app> --host 127.0.0.1 --port "$PORT"'` |
| Vite | `sh -c 'exec npx vite --host 127.0.0.1 --port "$PORT" --strictPort'` |
| Next.js | `sh -c 'exec npx next dev -H 127.0.0.1 -p "$PORT"'` |

- Hatch execs the command without a shell, so a bare `$PORT` in the argument list is never expanded. Wrap it in `sh -c '...'`, and use `exec` inside so signals reach the app.
- For static files, serve the smallest directory that holds the page and its assets. Not `/tmp`, not `$HOME`. If the page links assets by relative path, serve from the directory those paths start at.
- Check the app's listen code before you start it. If it binds every interface and has no host flag or env var, **stop**. Tell the user and propose the small app change that binds loopback. Do not put a proxy or forwarder in front of it.
- Always pass `--name <short-label>` (lowercase letters, digits, `-`). The default is the git branch, which in agent worktrees is a meaningless id. The label shows up in `hatch ls`, not in the URL.
- When the app builds absolute links or redirects, point it at the public origin: hatch sets `HATCH_URL` in the child's env.

**Done when:** you have one command that starts the app on loopback at `$PORT`, and a name.

### 2. Start it detached

```bash
setsid nohup hatch up --name <label> --minutes 240 -- <command> \
  > /tmp/hatch-<label>.log 2>&1 < /dev/null &
sleep 4; head -2 /tmp/hatch-<label>.log; hatch ls
```

- `hatch up` blocks for as long as the route lives. Start it with `setsid nohup ... &` as above. A background task from the agent's own tool (Claude Code `run_in_background`, for example) is killed when the tool's time limit hits, and that takes the hatch down with it. That cut off a 2-hour hatch before.
- `--minutes` defaults to 30, too short for most reviews. Use 240 unless the user gave a time. `--minutes 0` keeps it up until the process is stopped; use it only when the user asks for that. `hatch up -h` lists the flag. An old note that says otherwise is out of date.
- `--port N` without a command routes to a server that is already running. That works only if the server listens on `127.0.0.1:N` and N is in 20000-29999. Starting the app under hatch is the default.

**Done when:** the log's first line is a URL and `hatch ls` shows your name.

### 3. Check it, then hand it over

```bash
curl -s -o /dev/null -w '%{http_code}\n' '<url><path>'
```

This machine is on the tailnet, so curl the real URL, not `127.0.0.1`.

| Result | Cause |
|---|---|
| `200` (or the app's normal code) | Working |
| `000`, curl exit 35, `packet length too long` | The app answered plain HTTP on the tailnet port: it binds every interface. See step 1. |
| `502` "registered on port N, but nothing answered" | The app hasn't started yet, crashed, or listens on another port. Read the log. |
| Log shows `address already in use`, `command exited with status 1` | The app tried to bind every interface. See step 1. |
| `register: ... (is hatch serve running?)` | The daemon is down. Run `systemctl status hatch` and `journalctl -u hatch -n 30 --no-pager`, report what you find, and stop. Don't start `hatch serve` yourself. |
| `<label> is already registered by another process` | Run `hatch ls`. If the owner is your own earlier hatch, stop it (see below). Otherwise pick another name. Use `--replace` only on a route you own. |

Give the user the full link, including the path to the page (`https://<machine>.<tailnet>.ts.net:<port>/docs/report.html`), how long it stays up, and that the device must be on the tailnet. Every new `hatch up` gets a new port, so the old link stops working.

**Done when:** the URL with its path returns the page through hatch, and the user has the link.

## Close or extend a hatch

Stop the `hatch up` process. It stops the app and removes the route:

```bash
pkill -f '^hatch up --name <label> '
hatch ls
```

- Anchor the pattern with `^hatch up`. An unanchored `pkill -f "<label>"` also matches the agent's own shell command and kills it (exit 144).
- `hatch down <label>` alone does not close a live hatch. The running `hatch up` re-registers within 5 seconds. Use `hatch down` only to clear a route whose process is already gone.
- To extend, stop it and open it again with a longer `--minutes`. Tell the user the new URL.
- "Clean up" or "close out the thread" includes the hatch: stop it and confirm `hatch ls` no longer lists it.

**Done when:** `hatch ls` doesn't list the label.

## Reference

- `hatch ls --json` prints `{"zone": ..., "routes": [{"name", "port", "url", "owner", "expires", ...}]}`. It's an object, not a list.
- Env in the child: `PORT`, `HATCH_PORT`, `HATCH_URL`, `HATCH_NAME`.
- The browser's `Host` is `<machine>.<tailnet>.ts.net:<port>`, passed through unchanged. Dev servers with a host allowlist need that name. For Vite, set `server.allowedHosts: ['.ts.net']`. On "blocked request" or "invalid host header" errors, add it there.
- WebSockets, HMR, and server-sent events pass through.
- Daemon ports are off limits for `--port`: 53, 80, 7080, 7081, 8008, 8080, 8880.
- Source and design notes: the `NUSAlaska/hatch` repository's `README.md`.
