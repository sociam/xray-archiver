#!/usr/bin/env bash

set -e
set -u
set -f
#set -o pipefail

if [[ -z $1 ]]; then
	DATA_DIR=$1
else
	DATA_DIR=/usr/local/var/xray
fi

if [[ $UID -ne 0 ]]; then
	echo "Please run this script as root!"
	exit 1
fi

if ! type npm; then
	echo "npm is not installed!"
	exit 2
fi

install -dm775 "$DATA_DIR/{apk_archive,apk_unpacked}"
install -Dm644 xray-scraper.service /etc/systemd/systemd/xray-scraper.service
install -Dm644 scrape.js db.js package.json /usr/local/lib/xray/scraper/

cd /usr/local/lib/xray/scraper
npm install

chown xray:xray

exit 0

