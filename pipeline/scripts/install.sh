#!/usr/bin/env bash

set -e
set -u
set -f

usage() {
  echo <<EOF
Usage: $0 [OPTIONS]

options:
FLAGS           DESCRIPTION                                       DEFAULT
--help, -h      Show this help message.                           -
--prefix, -p    The prefix for the installation.                  /usr/local
--data-dir, -d  The directory that the apks will be archived in.  <prefix>/var/xray
EOF
}

PREFIX="/usr/local"
while [[ $# > 2 ]]; do
  case "$1" in
    -h|--help)     usage; exit 0;;
    -d|--data-dir) DATA_DIR="$2"; shift ;;
    -p|--prefix)   PREFIX="$2"; shift ;;
    --)            shift; break;;
    *)             usage; exit 64 ;;
  esac
  shift
done

: ${DATA_DIR:="$PREFIX/var/xray"}

if [[ $UID -ne 0 ]]; then
  echo "Please run this script as root!"
  exit 77
fi

if [[ ! -d analyzer ]] || [[ ! -d archiver ]] || [[ ! -d db ]]; then
  echo "analyzer, archiver and db directories not found! Please run this script from the pipeline directory of the repository!"
  exit 64
fi

if ! type npm &> /dev/null; then
  echo "npm is not installed!"
  exit 4
fi

if [[ ! -x analyzer/analyzer ]]; then
  if ! type go &> /dev/null; then
    echo "analyzer binary doesn't exist and go is not installed!"
    exit 5
  fi
  cd analyzer
  go build
  cd -
fi

if [[ ! -x apiserv/apiserv ]]; then
  if ! type go &> /dev/null; then
    echo "apiserv binary doesn't exist and go is not installed!"
    exit 5
  fi
  cd apiserv
  go build
  cd -
fi

if ! id -u xray &> /dev/null; then
  useradd xray -r -d "$DATA_DIR"
fi

install -o xray -Ddm775 "$DATA_DIR/apk_archive"
install -dm755 "$PREFIX/lib/xray/"

install -Dm755 analyzer/analyzer "$PREFIX/bin/analyzer"
install -Dm755 apiserv/apiserv "$PREFIX/bin/apiserv"

# package.json
install -Dm644 package.json package-lock.json -t "$PREFIX/lib/xray/"

# Downloader
install -Dm644 archiver/downloader/downloader.js -t "$PREFIX/lib/xray/archiver/downloader"

# retriever
install -Dm644 archiver/retriever/retriever.js -t "$PREFIX/lib/xray/archiver/retriever"

# explorer
install -Dm644 archiver/explorer/explorer.js -t "$PREFIX/lib/xray/archiver/explorer"

# db
install -Dm644 db/db.js -t "$PREFIX/lib/xray/db"

# util
install -Dm644 util/logger.js -t "$PREFIX/lib/xray/util"

install -Dm644\
        analyzer/xray-analyzer.service\
        apiserv/xray-apiserv.service\
        archiver/downloader/xray-downloader.service\
        archiver/retriever/xray-retriever.service\
        archiver/explorer/xray-explorer.service\
        -t "$PREFIX/lib/systemd/system/"

cd "$PREFIX/lib/systemd/system"
for f in \
  "xray-apiserv.service" "xray-analyzer.service" "xray-explorer.service"\
  "xray-downloader.service" "xray-retriever.service"; do
  sed -i -e "s:<<<PREFIX>>>:$PREFIX:g" "$f"
done

systemctl daemon-reload

cd "$PREFIX/lib/xray/"
npm install

exit 0
