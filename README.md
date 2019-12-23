# pd-oncall

A command-line tool for representing Pager Duty oncall schedule.

## usage
```sh
pd-oncall --api-token=API-TOKEN [<flags>] <command> [<args> ...]
```

### flags:
```
--help                   Show context-sensitive help (also try --help-long and --help-man).
--api-token=API-TOKEN    Auth API token; Might be an environment variable PAGERDUTY_API_TOKEN
--api-url=https://api.pagerduty.com/
                         Pager Duty API URL
--table-style="rounded"  Available table styles: rounded, box, colored
--timeout=10s            Timeout for a single http requests to Pager Duty API
--version                Show application version.
```

### commands:
```
help [<command>...]
  Show help.

config [<flags>]
  simple management of a config file

cache [<flags>]
  simple management of a cache file

now
  list currently oncall for schedules in a config file

schedule [<flags>]
  oncall schedule information

report [<flags>]
  generates a simple oncall report

roster [<flags>]
  roster for all known schedules

user [<flags>]
  oncall schedule for a specific user
```
