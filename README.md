# radigo

Record the [radiko.jp](http://radiko.jp/) program.

[![godoc](https://godoc.org/github.com/yyoshiki41/radigo?status.svg)](https://godoc.org/github.com/yyoshiki41/radigo)
[![go report](https://goreportcard.com/badge/github.com/yyoshiki41/radigo)](https://goreportcard.com/report/github.com/yyoshiki41/radigo)
[![CircleCI](https://circleci.com/gh/yyoshiki41/radigo.svg?style=svg)](https://circleci.com/gh/yyoshiki41/radigo)

[![Docker Stars](https://img.shields.io/docker/stars/yyoshiki41/radigo.svg)](https://hub.docker.com/r/yyoshiki41/radigo/)
[![Docker Build Status](https://img.shields.io/docker/build/yyoshiki41/radigo.svg)](https://hub.docker.com/r/yyoshiki41/radigo/tags/)
[![Docker Automated build](https://img.shields.io/docker/automated/yyoshiki41/radigo.svg)](https://hub.docker.com/r/yyoshiki41/radigo/builds/)

**Please do not use this project for commercial use. Only for your personal, non-commercial use.**</br>
**個人での視聴の目的以外で利用しないでください.**

## Installation

### Docker images

```bash
$ docker pull yyoshiki41/radigo
```

You can launch a radigo container and exec `radigo` command.

```bash
# Mount the volume `"$PWD"/output`(default output path) into `/output` in the container
$ docker run -it \
    -v "$(pwd)"/output:/output \
    yyoshiki41/radigo rec -id=LFR -s=20180401010000
Now downloading..
/
+------------+---------------------------------+
| STATION ID |              TITLE              |
+------------+---------------------------------+
| LFR        |　　  オードリーのオールナイトニッポン   　|
+------------+---------------------------------+
| Completed!
/output/20180401010000-LFR.aac
```

Open the output file created by the container on your local machine.

```bash
$ open $PWD/output/20180401010000-LFR.aac
```

### Build the binary from source

・Go 1.11 or higher

```bash
$ make installdeps
$ make build
$ radigo help
```

Or release binaries are available on [the releases page](https://github.com/yyoshiki41/radigo/releases).

#### Requirements

- ffmpeg
- rtmpdump (only if [recording a live streaming radio](#-rec-live))

### Build docker image from source

```bash
$ make docker-build
$ docker run -it yyoshiki41/radigo help
```

## Usage

```bash
$ radigo help
usage: radigo [--version] [--help] <command> [<args>]

Available commands are:
    area           Get available station ids
    browse         Browse radiko.jp
    browse-live    Browse radiko.jp live
    rec            Record a radiko program
    rec-live       Record a live program
```

### ■ area

```bash
$ radigo area
Area ID: JP13
+--------------------+----------------+
|        NAME        |   STATION ID   |
+--------------------+----------------+
| TBSラジオ           | TBS            |
| ニッポン放送          | LFR            |
| InterFM897         | INT            |
| TOKYO FM           | FMT            |
| J-WAVE             | FMJ            |
| bayfm78            | BAYFM78        |
| NACK5              | NACK5          |
| ＦＭヨコハマ           | YFM            |
| 文化放送            | QRR            |
| ラジオNIKKEI第1      | RN1            |
| ラジオNIKKEI第2      | RN2            |
| NHKラジオ第2         | JOAB           |
| NHK-FM（東京）        | JOAK-FM        |
| NHKラジオ第１(東京)   | JOAK           |
| 放送大学            | HOUSOU-DAIGAKU |
+--------------------+----------------+
```

#### Note

`Area ID` is ISO 3166-2 code that is defined for 47 prefectures.

_c.f._ _[ISO 3166-2:JP - Wikipedia](https://ja.wikipedia.org/wiki/ISO_3166-2:JP)_

### ■ rec

Record the program using the [timefree](http://radiko.jp/#!/fun/timeshift).

```bash
$ radigo rec -id=LFR -s=20161126010000
Now downloading..
+------------+---------------------------------+
| STATION ID |              TITLE              |
+------------+---------------------------------+
| LFR        |　　  オードリーのオールナイトニッポン   　|
+------------+---------------------------------+
Completed!
/tmp/output/20161126010000-LFR.aac
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
/tmp/output/20161205083547-LFR.aac
```

### ■ browse / browse-live

Browse [radiko.jp](http://radiko.jp/).

```bash
$ radigo browse -id=LFR -s=20161126010000
```

```bash
$ radigo browse-live -id=LFR
```

### Optional

#### - working dir

Default: `${PWD}/output`

If you want to change the working dir, set the environment variables.

- `RADIGO_HOME`

#### - radiko premium

If use the [area free](http://radiko.jp/rg/premium/), set the environment variables.

- `RADIKO_MAIL`
- `RADIKO_PASSWORD`

##### e.g.

```bash
# export RADIKO_MAIL="radigo@example.com" && export RADIKO_PASSWORD="password"
$ radigo rec -a=JP13 -id=LFR -s=20161126010000
```

## Known Issues

### `ERROR: Failed to concat aac files`

(caused by the limitation of file descriptors maybe.)

Increase the number of file descriptors.

```bash
$ ulimit -n 16384
```

## Resources

- [Japanese](http://qiita.com/yyoshiki41/items/f81442d7dc2d0ddcf15b)
- [Listen on itunes](http://esola.co/posts/2017/aac-profile/)

## License

radigo is licensed under the GPLv3 license for all open source applications.

Please do not use this project for commercial use, it is not intended to be used for commercial use.

## Author

Yoshiki Nakagawa
