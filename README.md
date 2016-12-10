# radigo

Record the [radiko.jp](http://radiko.jp/) programs.

[![godoc](https://godoc.org/github.com/yyoshiki41/radigo?status.svg)](https://godoc.org/github.com/yyoshiki41/radigo)
[![CircleCI](https://circleci.com/gh/yyoshiki41/radigo.svg?style=svg)](https://circleci.com/gh/yyoshiki41/radigo)
[![go report](https://goreportcard.com/badge/github.com/yyoshiki41/radigo)](https://goreportcard.com/report/github.com/yyoshiki41/radigo)

## Installation

・Go 1.7 or newer

```bash
$ go get github.com/yyoshiki41/radigo/cmd/radigo/...
```

## Requirements

- swfextract
- ffmpeg
- rtmpdump (only if [recording a live streaming radio](#-rec-live))

## Configuration

```bash
$ make init
```

### Optional

#### - working dir

Default: `/tmp/radigo`

If you want to change the working dir, set the environment variables.

- `RADIGO_HOME`

#### - radiko premium

If use the [area free](http://radiko.jp/rg/premium/), set the environment variables.

- `RADIKO_MAIL`
- `RADIKO_PASSWORD`

## Usage

```bash
$ radigo help
usage: radigo [--version] [--help] <command> [<args>]

Available commands are:
    area        Get available station ids
    rec         Record a radiko program
    rec-live    Record a live program
```

### ■ area

```bash
$ radigo area
Area ID: JP13
+------------------+----------------+
|       NAME       |   STATION ID   |
+------------------+----------------+
| TBSラジオ         | TBS            |
| ニッポン放送       | LFR            |
| InterFM897       | INT            |
| TOKYO FM         | FMT            |
| J-WAVE           | FMJ            |
| bayfm78          | BAYFM78        |
| NACK5            | NACK5          |
| ＦＭヨコハマ       | YFM            |
+------------------+----------------+
```

### ■ rec

Record the program using the [timefree](http://radiko.jp/#!/fun/timeshift).

```bash
$ radigo rec -id=LFR -s=20161126010000
Now downloading..
+------------+---------------------------------+
| STATION ID |              TITLE              |
+------------+---------------------------------+
| LFR        |　　  オードリーのオールナイトニッポン |
+------------+---------------------------------+
Completed!
/tmp/radigo/output/20161126010000-LFR.mp3
```

### ■ rec-live

Record the live streaming program.

```bash
$ radigo rec-live -id=LFR -t=3600
Now downloading..
+------------+---------------+
| STATION ID | DURATION(SEC) |
+------------+---------------+
| LFR        |          3600 |
+------------+---------------+
Completed!
/tmp/radigo/output/20161205083547-LFR.mp3
```

### Cleanup

Remove cache and force refresh.

```bash
$ make clean
```

## License 

The MIT License

## Author

Yoshiki Nakagawa
