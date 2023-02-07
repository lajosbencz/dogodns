# DOGODNS

DigitalOcean Dynamic DNS Client

---

## Usage

### Arguments

```bash
Usage of dogodns:
      --domain string   Domain name (default "")
      --dry             Dry run, dont commit any changes
      --interval int    Check interval (default 60)
      --pip string      Public IP fetch URL (default "https://api.ipify.org/?format=raw")
      --token string    DigitalOcean API R+W token (default "")
      --ttl int         Record TTL (default 300)
```

### Environment variables

```
DOGODNS_TOKEN
DOGODNS_DOMAIN
DOGODNS_PIP
DOGODNS_TTL
DOGODNS_INTERVAL
```

### Config file

Located in `/etc/dogodns/dogodns.yaml` or `./dogodns.yaml`

```yaml
# DigitalOcean API R+W token
token: bogus
# Domain name
domain: home.example.me
# Public IP URL
pip: https://api.ipify.org/?format=raw
# Record TTL
ttl: 300
# Check interval
interval: 60
```