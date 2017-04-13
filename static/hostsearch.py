import json
import re
import os
import urllib
import csv
import subprocess
from subprocess import Popen, PIPE, check_call

def appendJSON(apk,trackers):

	line = '"%s": %s,' % (apk, json.dumps(trackers))

def findurl(domain,apk):
	domain = str(domain)
	print domain
	apk = str(apk)
	print apk
	p = Popen(['grep', '-F', '-r',  domain, apk], stdin=PIPE, stdout=PIPE, stderr=PIPE,
          bufsize=-1)
	output, error = p.communicate()
	if p.returncode == 0:
		print 'found domain: %s' % domain
		return 1
	elif p.returncode == 1:
		print 'did not find domain: %s' % domain
		return 0
	else:
		assert p.returncode > 1
		print 'error'
		return 0


outfile = open('output.txt', 'a')

# domain = 'facebook'
# apk = 'air.BattleofZombies'
# findurl(domain,apk)

def searchAPK(apk):
	trackers = []
	for tracker in companyinfo:
		# remember to put back to companyinfo[tracker]['domains']
		for domain in companyinfo[tracker]["domains"]:
			print domain
			if domain:
				print "looking for %s in %s" % (domain,apk)
				if findurl(domain,apk):
					print "found"
					trackers.append(tracker)
				else:
					print "didn't find"
			else:
				print "no domains found for this company"
	return trackers


# # open all the different data sources to be combined

companyinfo = open('../mitm_out/company_details.json', 'rb')
companyinfo = json.load(companyinfo)

directory = '../../apks/testbatch/'
for filename in os.listdir(directory):
	print filename
	fullpath = (directory + filename)
	print fullpath
	os.system("java -jar apktool.jar d %s" % fullpath)
	unpackedname = filename.rsplit( ".", 1 )[ 0 ]
	print unpackedname
	trackers = searchAPK(unpackedname)
	trackers = str(trackers)
	outfile.write(trackers)
	os.system("rm -r %s" % unpackedname)



# """
# apptrackers_infile = csv.reader(open('../like-to-like/300_mobile_trackers.csv'))
# apptrackers = {}
# for row in apptrackers_infile:
# 	k, v = row
# 	apptrackers[k] = v

# webtrackers_infile = csv.reader(open('../like-to-like/300_web_trackers.csv'))
# webtrackers = {}
# for row in webtrackers_infile:
# 	k, v = row
# 	k = str(k.title())
# 	webtrackers[k] = v

# def compare(first, second):
# 	firstset = set(first.keys())
# 	secondset = set(second.keys())
# 	print 'The following are web only:'
# 	print firstset - secondset
# 	print 'the following are mobile only:' % second
# 	print secondset - firstset
# 	sharedKeys = set(firstset.intersection(secondset))
# 	print 'the following are both:'
# 	print sharedKeys

# def getparent(tracker):
# 	if tracker in companyinfo:
# 		parent = companyinfo[tracker]['parent']
# 		return parent
# 	else:
# 		return tracker

# for tracker in apptrackers:
# 	parent = getparent(tracker)
# 	if parent != '':
# 		row = '%s,%s,%s\n' % (tracker, parent, apptrackers[tracker])
# 		outfile_app.write(row)
# 	else:
# 		row = '%s,%s,%s\n' % (tracker, tracker, apptrackers[tracker])
# 		outfile_app.write(row)

# for tracker in webtrackers:
# 	parent = getparent(tracker)
# 	if parent != '':
# 		row = '%s,%s,%s\n' % (tracker, parent, webtrackers[tracker])
# 		outfile_web.write(row)
# 	else:
# 		row = '%s,%s,%s\n' % (tracker, tracker, webtrackers[tracker])
# 		outfile_web.write(row)
# 		"""