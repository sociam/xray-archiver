#!/usr/bin/env python
#
# Simple script showing how to read a mitmproxy dump file
# execute as follows: 
# mitmdump -s mitm-save.py -p 8081 
#

from mitmproxy import flow
from mitmproxy.models import HTTPResponse
from netlib.http import Headers
import csv, json, time, os, urllib, random

config = json.loads(''.join(open('./mitm-save-config.json', 'r').readlines())) # seriously, python? 
rand = ''.join([random.choice('abcdefghijklmnopqrstuvwxyz') for r in xrange(5)])
OUTFILE = '/'.join([config["destdir"],'-'.join([config['app'],config['platform'],config['version'],rand])+'.csv'])
(app, device, platform, version, researcher) = (config['app'], config['device'], config['platform'], config['version'], config['researcher'])

print "logging ", app, platform, version, " to ", OUTFILE

def openfile():
    newFile = not os.path.isfile(OUTFILE) 
    f = open(OUTFILE,'a')
    writer = csv.writer(f)
    if newFile:
        writer.writerow(['app', 'version', 'device', 'platform', 'researcher', 'time', 'host', 'url', 'method', 'headers', 'body'])
    return f,writer

def request(context, flow):
    # pretty_host takes the "Host" header of the request into account,
    # which is useful in transparent mode where we usually only have the IP
    # otherwise.

    # Method 1: Answer with a locally generated response
    # print " host ", flow.request.pretty_host, flow.request.url, flow.request.headers

    f,writer = openfile()
    writer.writerow([app, version, device, platform, researcher, int(time.time()*1000), flow.request.pretty_host, urllib.quote(flow.request.url), flow.request.method, urllib.quote(json.dumps(dict(flow.request.headers))), urllib.quote(flow.request.body)]) 
    f.close()
 