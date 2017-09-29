#!/usr/bin/env bash

psql -h 192.168.0.9 -U apiserv -d xraydb -c  "SELECT count(*) FROM app_versions WHERE analyzed=True AND icon IS NOT NULL;" | grep -Eo '[0-9][0-9]+'

