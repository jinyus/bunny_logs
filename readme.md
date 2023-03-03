CLI to process logs from [bunnycdn](https://bunny.net)

```zsh
Usage of bunny_logs:
  -d value
        domain to exclude from referrers (can be repeated)
  -f value
        log file to process (can be repeated)
```

## How to use:

- download your log from the bunnycdn panel from https://panel.bunny.net

- Run the command:

      ./bunny_logs -f mylog1.log -f mylog2.log -d example.com -d www.example.com
