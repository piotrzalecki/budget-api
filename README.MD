# budget-api

## Authentication

Every request must include header `X-API-Key: <value of $BUDGET_API_KEY>`.
The server panics on startup if `BUDGET_API_KEY` is unset.

**Note:** The `/health` endpoint does not require authentication.

## Dependencies

# Development Tools

| Tool Name     | Version |
| ------------- | ------- |
| sqlc          | v1.29.0 |
| goose         | v3.24.3 |
| golangci-lint | 2.1.6   |

# Raspberry Pi – Journald Log‑Size Limits for Budget App

> **Why set limits?**\
> Continuous JSON access logs can fill an SD‑card in weeks.\
> `systemd‑journald` supports disk quotas—configure them once and forget.

---

## Recommended Settings

| Key (`journald.conf`) | Suggested value | Reason                                                                             |
| --------------------- | --------------- | ---------------------------------------------------------------------------------- |
| `SystemMaxUse=`       | **200M**        | Caps total persistent journal size on disk. Old logs rotate away.                  |
| `SystemMaxFileSize=`  | **20M**         | Upper bound per individual file. Creates more, smaller chunks—safer on power loss. |
| `RuntimeMaxUse=`      | default         | Applies only to `/run` (tmpfs) logs; default is fine.                              |
| `MaxRetentionSec=`    | _(optional)_    | Use if you prefer time‑based retention (e.g., `2592000` for 30 days).              |

> **200 MB ≈ 7‑10 days** of single‑user traffic at \~30 KB/request.\
> Tune up or down depending on card size and log verbosity.

---

## Step‑by‑Step

1. **Create an override file**

   ```bash
   sudo mkdir -p /etc/systemd/journald.conf.d
   sudo nano /etc/systemd/journald.conf.d/budget-size.conf
   ```

   Paste:

   ```ini
   [Journal]
   SystemMaxUse=200M
   SystemMaxFileSize=20M
   ```

2. **Reload journald**

   ```bash
   sudo systemctl restart systemd-journald
   ```

3. **Verify**

   ```bash
   journalctl --disk-usage
   # → e.g. "Archived and active journals 80.0M, max 200.0M, 120.0M free."
   ```

---

## Ongoing Checks

| Command                               | What it tells you                   |
| ------------------------------------- | ----------------------------------- |
| `journalctl --disk-usage`             | Current on‑disk footprint vs. limit |
| `journalctl -u budgetd -f`            | Tail live API logs                  |
| `ls -lh /var/log/journal/*/*.journal` | See individual rotated files        |

---

## Adjusting Later

- **Need more history?** Increase `SystemMaxUse=` (e.g., `500M`) and restart journald.
- **Card almost full?** Lower the value; journald deletes oldest chunks to comply.
- **Switch to time‑based retention?** Add `MaxRetentionSec=` and remove `SystemMaxUse=`.

---

### TL;DR (one‑liner)

```bash
sudo tee /etc/systemd/journald.conf.d/budget-size.conf <<'EOF'
[Journal]
SystemMaxUse=200M
SystemMaxFileSize=20M
EOF
sudo systemctl restart systemd-journald
journalctl --disk-usage
```
