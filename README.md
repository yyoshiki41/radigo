# radigo

Record the [radiko.jp](http://radiko.jp/) program.

[![godoc](https://godoc.org/github.com/yyoshiki41/radigo?status.svg)](https://godoc.org/github.com/yyoshiki41/radigo)
[![go report](https://goreportcard.com/badge/github.com/yyoshiki41/radigo)](https://goreportcard.com/report/github.com/yyoshiki41/radigo)
[![CircleCI](https://circleci.com/gh/yyoshiki41/radigo.svg?style=svg)](https://circleci.com/gh/yyoshiki41/radigo)

[![Docker Stars](https://img.shields.io/docker/stars/yyoshiki41/radigo.svg)](https://hub.docker.com/r/yyoshiki41/radigo/)
[![Docker Build Status](https://img.shields.io/docker/build/yyoshiki41/radigo.svg)](https://hub.docker.com/r/yyoshiki41/radigo/tags/)
[![Docker Automated build](https://img.shields.io/docker/automated/yyoshiki41/radigo.svg)](https://hub.docker.com/r/yyoshiki41/radigo/builds/)

_Please refrain from using beyond the range of personal listening._ </br>
__個人での視聴の目的以外で利用しないでください.__

## Installation

### Docker images

```bash
$ docker pull yyoshiki41/radigo
```

You can launch a radigo container and exec `radigo` command.

```bash
$ docker run --name radigo -itd radigo
$ docker attach radigo
root@158057ab4c2a:/tmp$ radigo rec -id=LFR -s=20180114010000
Now downloading..
/
+------------+----------------------------------+
| STATION ID |              TITLE               |
+------------+----------------------------------+
| LFR        | オードリーのオールナイトニッポン |
+------------+----------------------------------+
| Completed!
/tmp/radigo/output/20180114010000-LFR.aac
```

Copy the output from a radigo container to the host (your local machine).

```bash
$ docker cp radigo:/tmp/radigo/output/20180114010000-LFR.aac ./
$ open ./20180114010000-LFR.aac
```

### Building from source

・Go 1.7 or newer

```bash
$ go get github.com/yyoshiki41/radigo/cmd/radigo/...
# Configuration
$ make init
```

#### Requirements

- ffmpeg
- rtmpdump (only if [recording a live streaming radio](#-rec-live))

#### Cleanup

Remove output files.

```bash
$ make clean
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
+-----------------------+---------------+
|         NAME          |  STATION ID   |
+-----------------------+---------------+
| TBSラジオ		| TBS		|
| 文化放送		| QRR		|
| ニッポン放送		| LFR		|
| ラジオNIKKEI第1	| RN1		|
| ラジオNIKKEI第2	| RN2		|
| InterFM897		| INT		|
| TOKYOFM		| FMT		|
| J-WAVE		| FMJ		|
| ラジオ日本		| JORF		|
| bayfm78		| BAYFM78	|
| NACK5			| NACK5		|
| ＦＭヨコハマ		| YFM		|
| 放送大学		| HOUSOU-DAIGAKU|
| NHKラジオ第1（東京）	| JOAK		|
| NHKラジオ第2		| JOAB		|
| NHK-FM（東京）		| JOAK-FM	|
+-----------------------+---------------+
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

Default: `/tmp/radigo`

If you want to change the working dir, set the environment variables.

- `RADIGO_HOME`

#### - radiko premium

If use the [area free](http://radiko.jp/rg/premium/), set the environment variables.

- `RADIKO_MAIL`
- `RADIKO_PASSWORD`

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

The MIT License

## Author

Yoshiki Nakagawa
