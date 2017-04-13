import os, json, csv, glob, sys
reload(sys)
#sys.setdefaultencoding('utf-8')

with open('xray/mitm_out/host_by_app.json') as data_file:
	apps = json.load(data_file)
with open('xray/mitm_out/company_details.json') as data_file2:
	companies = json.load(data_file2)

f = csv.writer(open("tracker_by_app.csv", "a"))

for company in companies:
	print companies[company]['id']

for key in apps.items():
	app = key
	appname = app[0]
#	print appname
	for key in app:
		hosts = key
#		print type(hosts)
		for key in hosts:
			host = key
			if len(host) > 1:
				for company in companies:
					if host == companies[company]['domains'][0]:
						f.writerow([appname,companies[company]['id']])