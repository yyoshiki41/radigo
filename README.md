# radigo

Record a radiko program.

[![godoc](https://godoc.org/github.com/yyoshiki41/radigo?status.svg)](https://godoc.org/github.com/yyoshiki41/radigo)
[![build](https://travis-ci.org/yyoshiki41/radigo.svg?branch=master)](https://travis-ci.org/yyoshiki41/radigo)
[![codecov](https://codecov.io/gh/yyoshiki41/radigo/branch/master/graph/badge.svg)](https://codecov.io/gh/yyoshiki41/radigo)
[![go report](https://goreportcard.com/badge/github.com/yyoshiki41/radigo)](https://goreportcard.com/report/github.com/yyoshiki41/radigo)

## Installation

・Go 1.7 or newer

```bash
$ go get github.com/yyoshiki41/radigo/cmd/radigo/...
```

## Requirements

- swfextract
- ffmpeg

## Configuration

```bash
$ make init
```

#### - Optional (radiko premium)

If use the [area free](http://radiko.jp/rg/premium/), set the environment variables.

- `RADIKO_MAIL`
- `RADIKO_PASSWORD`

## Usage

```bash
$ radigo help
usage: radigo [--version] [--help] <command> [<args>]

Available commands are:
    area    Get available station ids
    rec     Record a radiko program
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
/tmp/radigo/result.mp3
```

#### - cleanup

```bash
$ make cleanup
```

## License 

The MIT License

## Author

Yoshiki Nakagawa
