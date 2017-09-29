#!/usr/bin/env bash

cat apps_icon | while read line 
do
  find /export/a/apps/$line -name "*.apk" -type f -delete  
  echo $line
done
