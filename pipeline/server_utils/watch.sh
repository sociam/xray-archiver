#!/usr/bin/env bash

while true; do
  echo -n 'Playstore Apps: '
  psql -h 192.168.0.9 -U apiserv -d xraydb -c  'SELECT count(*) FROM playstore_apps' | grep -Eo '[0-9][0-9]+'

  echo -n ' - The Empire Apps: '
  psql -h 192.168.0.9 -U apiserv -d xraydb -c "SELECT count(*) FROM app_versions WHERE region='uk'" | grep -Eo '[0-9][0-9]+'
  echo -n ' - Colonies Apps: '
  psql -h 192.168.0.9 -U apiserv -d xraydb -c "SELECT count(*) FROM app_versions WHERE region='us'" | grep -Eo '[0-9][0-9]+'

  echo -n 'Free Apps: '
  psql -h 192.168.0.9 -U apiserv -d xraydb -c 'SELECT count(*) FROM playstore_apps WHERE free=True' | grep -Eo '[0-9][0-9]+'




  echo -n 'Downloaded: '
  psql -h 192.168.0.9 -U apiserv -d xraydb -c  'SELECT count(*) FROM app_versions WHERE downloaded=True' | grep -Eo '[0-9][0-9]+'

  echo -n 'Analy{s,z}ed: '
  psql -h 192.168.0.9 -U apiserv -d xraydb -c  'SELECT count(*) FROM app_versions WHERE analyzed=True;' | grep -Eo '[0-9][0-9]+'

  echo -n 'Failed: '
  psql -h 192.168.0.9 -U apiserv -d xraydb -c 'SELECT count(*) FROM app_versions WHERE last_analyze_attempt IS NOT NULL AND analyzed = False' | grep -Eo '[0-9][0-9]+'

  echo -n 'Alternatives: '
  psql -h 192.168.0.9 -U apiserv -d xraydb -c  'SELECT count(*) FROM alt_apps' | grep -Eo '[0-9][0-9]+'


  sleep 30
  echo
done
