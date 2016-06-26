#!/usr/bin/env python
#
# Simple script showing how to read a mitmproxy dump file
# execute as follows: 
# mitmdump -s mitm-save.py -p 8081 
#

from mitmproxy import flow
from mitmproxy.models import HTTPResponse
from netlib.http import Headers
import csv, json, time, os, urllib

OUTFILE = '/tmp/out.csv'

def openfile():
    newFile = not os.path.isfile(OUTFILE) 
    f = open(OUTFILE,'a')
    writer = csv.writer(f)
    if newFile:
        writer.writerow(['time', 'host', 'url', 'headers', 'body'])
    return f,writer

def request(context, flow):
    # pretty_host takes the "Host" header of the request into account,
    # which is useful in transparent mode where we usually only have the IP
    # otherwise.

    # Method 1: Answer with a locally generated response
    # print " host ", flow.request.pretty_host, flow.request.url, flow.request.headers

    f,writer = openfile()
    writer.writerow([int(time.time()*1000),flow.request.pretty_host, urllib.quote(flow.request.url), urllib.quote(json.dumps(dict(flow.request.headers))), urllib.quote(flow.request.body)]) # json.dumps(dict(flow.request.headers)), flow.request.body])
    f.close()
 