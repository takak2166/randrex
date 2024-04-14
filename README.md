# randrex

## Summary

A tool that randomly displays strings matching a regular expression.  
正規表現にマッチする文字列をランダムに表示するツール

## Requirement

Requires Go version ≥ 1.21 for build

## Usage

```bash
git clone https://github.com/takak2166/randrex.git
go build
./randrex help
```

## Example

```bash
$ ./randrex generate -n 5 "hoge_[a-z0-9]{8}"
hoge_li832ku6
hoge_ww7j7av0
hoge_xq4nib4l
hoge_f61qj60c
hoge_s474j4o0
```