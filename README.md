# skjald

`skjald` is a tool to help enable users to take notes about their day based on a pomodoro-style cadance. Users simply define their work periods in the config file, and record notes in the `skjald` window periodically, which is recorded at period ends in a note file for reference.

## Config file

The config file lives in `~/.skjald`, is `YAML` formatted, and will set values if not present. 

Current config file expected options:
```
period_work: 25
period_rest: 5
period_lunch: 30
```
