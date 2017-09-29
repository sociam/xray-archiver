#!/usr/bin/env bash

cat apps_NO_icons | while read line 
do
  if [[ -z $(find /export/a/apps/$line -name "*.apk" -type f 2>/dev/null) ]]; then
    echo $line
  fi

done
