#!/bin/bash
FILE=$1

echo "Processing $FILE file..."
java -jar apktool.jar d $FILE
x=$FILE
y=${x%.apk}
grep -r "facebook.com" ${y##*/} > findings.txt
rm -r ${y##*/}