# DOGODNS

This app will use the public IP of it's network
to update the DNS registry of DigitalOcean

---

## Usage

### Run Updater

```bash
dogodns service # --config dogodns.yaml
```

### Available Commands

```bash
Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Creates default config file
  ip          Prints public IP that will be used for A record updates
  service     Run the service
  status      Shows a brief status page

Flags:
  -c, --config string        Path to config file without extension (default "dogodns.yaml")
  -d, --domain stringArray   List of domain names
  -h, --help                 help for dogodns
  -i, --interval int         Interval between checks (default 60)
  -p, --pip string           HTTP URL to fetch public IP from (default "https://ipecho.net/plain")
  -t, --token string         DigitalOcean R+W token
  -l, --ttl int              Domain record TTL (default 300)
```

### Config file

Located at `./dogodns.yaml` by default

```yaml
# DigitalOcean API R+W token
token: <secret>

# Domain names (use either key)
domain: home.example.me
domains:
  - dev.example.me

# Public IP URL
pip: https://ipecho.net/plain

# Record TTL
ttl: 300

# Check interval
interval: 60

```
