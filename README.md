# org2gcal

org2gcal is command line tool to generate `event.json` for [gcal](https://github.com/toniov/gcal-cli).

## Install

```
go get -u github.com/Ladicle/org2gcal
```

## Usage

```
❯❯❯ org2gcal --help
NAME:
   org2gcal - Convert time-log to json format for gcal

Usage:
   org2gcal [date]

DATE:
   format       2006-1-2
   note         This argument is optional. if you do not specify this, date is used today.
```

## Supported time-log format

The leading and trailing white space is ignored.

```
- HH:MM summary
```

## Quickstart

```
❯❯❯ org2gcal 2018-02-21
input time logs:
- 11:30 go to office
- 12:00 lunch
- 13:00 coding
- 23:00 go to bed

❯❯❯ gcal bulk -e events.json
```
