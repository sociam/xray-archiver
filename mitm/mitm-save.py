#!/usr/bin/env python
#
# Simple script showing how to read a mitmproxy dump file
# execute as follows: 
# mitmdump -s mitm-save.py -p 8081 
#
# procedure for collecting data:
# 0. install mitmproxy ca as a root ca on device
# 1. terminate all running apps on device
# 2. set http proxy to <this mitm host>:8081 
# 3. edit mitm-save-config.json for appropriate app types
# 4. run mitmdump -s mitm-save.py -p 8081
# 5. launch target app on device
# 6. do stuff on device for awhile, then terminate proxy then terminate app

from mitmproxy import flow
from mitmproxy.models import HTTPResponse
from netlib.http import Headers
import csv, json, time, os, urllib, random

config = json.loads(''.join(open('./mitm-config.json', 'r').readlines())) # seriously, python? 
runid = ''.join([random.choice('abcdefghijklmnopqrstuvwxyz') for r in xrange(5)])
OUTFILE = '/'.join([config["destdir"],'-'.join([config['app'],config['platform'],config['version'],runid])+'.csv'])
(app, device, platform, version, researcher) = (config['app'], config['device'], config['platform'], config['version'], config['researcher'])

print "logging ", app, platform, version, " to ", OUTFILE

def openfile():
    newFile = not os.path.isfile(OUTFILE) 
    f = open(OUTFILE,'a')
    writer = csv.writer(f)
    if newFile:
        writer.writerow(['app', 'version', 'device', 'platform', 'researcher', 'time', 'runid', 'host', 'url', 'method', 'headers', 'body'])
    return f,writer

def request(context, flow):
    # pretty_host takes the "Host" header of the request into account,
    # which is useful in transparent mode where we usually only have the IP
    # otherwise.

    # Method 1: Answer with a locally generated response
    # print " host ", flow.request.pretty_host, flow.request.url, flow.request.headers

    f,writer = openfile()
    writer.writerow([app, version, device, platform, researcher, int(time.time()*1000), runid, flow.request.pretty_host, urllib.quote(flow.request.url), flow.request.method, urllib.quote(json.dumps(dict(flow.request.headers))), urllib.quote(flow.request.body)]) 
    f.close()
 