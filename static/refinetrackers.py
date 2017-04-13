import json
import re
import os
import urllib
import csv
import subprocess
from subprocess import Popen, PIPE, check_call
from collections import defaultdict

mobtrackers = ['Google-Ads', 'Inmobi', 'Mopub', 'Facebook', 'Chartboost', 'Crashlytics', 'Fabric', 'Heyzap', 'Applovin', 'Unitytechnologies', 'Vungle', 'Fyber', 'Millenialmedia', 'Ironsource', 'Google-Analytics', 'Tapjoy', 'Flurry', 'Google', 'Bitstadium', 'Presage', 'Comscore', 'Umeng', 'Tencent', 'Nexage', 'Mixpanel', 'Appboy', 'Yandex', 'Kochava', 'Urbanairship', 'Amazon-Analytics', 'Daum', 'Milkman', 'Playhaven', 'Umeng', 'Mail.Ru', 'Newrelic', 'Admob', 'Moat', 'Smaato', 'Admarvel', 'Appodeal', 'Cauley', 'Revmob', 'Adbuddiz', 'Igaworks', 'Nativex', 'Smartadserver', 'Mobfox', 'Pubnative', 'Hyprmx', 'Pingstart', 'Tune', 'Baidu', 'Swrve', 'Tnkfactory', 'Appnext', 'Tencent']
webtrackers = ['Google-Analytics', 'Google', 'Doubleclick', 'Facebook', 'Google-Ads', 'Cloudfront', 'Scorecardresearch', 'Appnexus', 'Twitter', 'Amazon', 'Criteo', 'Yahoo', 'Oracle', 'Quantcast', 'Iponweb', 'Aol', 'Adobe', 'Openx', 'Ghostery', 'Rubicon', 'Bluekai', 'Cloudflare', 'Liveramp', 'The Nielsen Company', 'Demdex', 'Truste', 'Newrelic', 'Integral Ad Science', 'Adsrvr', 'Moatads', 'Pubmatic', 'Turn', 'Indexchange', 'Doubleverify', 'Mathtag', 'Chartbeat', 'Optimizely', 'Youtube', 'Tapad', 'Exelate', 'Targus', 'Lotame', 'Sizmek', 'Iasds', 'Neustar', 'Jquery', 'Krxd', 'Amazon-Ads', 'Adroll', 'Brightroll', 'Audiencescience', 'Bootstrap', 'Rocketfuel', 'Dataxu', 'Tubemogul', 'Microsoft', 'Yandex', 'Conversant', 'Videology', 'Flashtalking', 'Akamai', 'Facebook-Ads', 'Adap.Tv', 'Magnetic', 'Dstillery', 'Baidu', 'Smartadserver', 'Taboola', 'H&R Block', 'Tealium', 'Convertro', 'Msn', 'Linkedin', 'Mookie1', 'Contextweb', 'Liverail', 'Yadro', 'Bug', 'Chango', 'Conversant Llc', 'Trueffect', 'Signal-Privacy', 'Crazy Egg', 'Drawbridge Inc', 'Sovrn', 'Ovh', 'Exponential', 'Yieldoptimizer', 'Spotxchange', 'Owneriq', 'Insightexpress', 'Adform', 'Automattic', 'Disqus', 'Outbrain', 'Flipps', 'Centromedia', 'Ignitionone', 'Skimlinks', 'Simpli.Fi', 'Wordpress', 'Openx Technologies', 'Ensighten', 'Visual Website Optimizer', 'Rhythmone', 'Mixpanel', 'Ru5', 'Sonobi', 'Hotjar', 'Voicefive', 'Gssprt', 'At Internet', 'Ml315', 'Radiumone', 'Parsely', 'Improve Digital Bv', 'Gumgum', 'Amgdgt', 'Effective Measure', 'Pinterest', 'Vk', 'Com.Br', 'Co.Uk', 'Adition', 'Exoclick', 'Ixi Corporation', 'Nugg.Ad', 'Umeng', 'Alibaba', 'Marketo', 'Gskinner', 'Gigya', 'Tns', 'Eq Works', 'Cxense Asa', 'Steelhouse', 'Yieldlab', '33Across', 'Mail.Ru', 'Pingdom Ab', 'Smart Adserver', 'Brightcove', 'Postrelease', 'Adriver', 'Pagefair', 'Microad', 'Marin Software', 'Collective', 'Trafficjunky', 'Pixalate', 'Demandbase', 'Yieldbot', 'Visualdna', 'Ioam', 'Signal', 'Shopzilla', 'Com.Cn', 'Sharethrough', 'Rambler', 'Clicktale', 'Sociomantic', 'Datonics', 'Run', 'Myspace', 'Taobao', 'Adgear Technologies', 'Sharethis', 'Mybuys', 'Stroeer Media', 'The Mcgraw-Hill Companies', 'Vizury', 'Abc', 'Dyn', 'Longtail Ad Solutions', 'Netseer', 'Smaato', 'The Adex Gmbh', 'Adblade', 'Mxptint', 'Maxymiser', 'Admeta', 'Mcro', 'Sourcepoint Technologies', 'Ne.Jp', 'Yieldmanager', 'Rfihub', 'Bounce Exchange', 'Wtp', 'Convertmedia', 'Rackspace Us', 'Eyeview', 'Qualtrics', 'Intent Iq', 'Stickyads.Tv', 'Soasta', 'Cedexis Inc', 'Weborama', 'Invite Media', 'Histats', 'Statcounter', 'Paypal', 'Bidr', 'Springserve', 'Yahoo', 'Nexstar Broadcasting', 'Triplelift', 'Polar', 'Adelphic, Inc', 'Deep Forest Media', 'Cedexis', 'Richrelevance', 'Resonate', 'Impact-Ad', 'Jivox', 'Mediade', 'Liveperson', 'Zedo', 'Yabidos', 'Com.Tw', 'Adfox', 'Fout', 'Switch', 'Navegg S.A.', 'Adsnative', 'Eyeota Limited', 'Jsdelivr', 'Consumerinfo.Com', 'Brandscreen Pty Ltd', 'Mythings', 'Polar Mobile Group', 'Answercloud', 'Komoona', 'Opt', 'Admeta', 'Meltdsp', 'Appier Inc', 'Netdna', 'Budgetedbauer', 'Gfk', 'Bidtheatre', 'Research Now', 'Cdn', 'Undertone Networks', 'Viglink', 'Adbrain', 'Adstir', 'Sailthru', 'Vimeo', 'Mediametrie', 'Livefyre', 'Mediaforge', 'Webklipper', 'Ntv.Io', 'Tencent', 'Ib-Ibi']
trackers = set(mobtrackers)|set(webtrackers)

import json
obj  = json.load(open('../mitm_out/company_details.json'))

newtrackers = {}

for olditem in obj:
	domains = json.dumps(obj[olditem]["domains"])
	newtrackers[olditem] = domains

print newtrackers

	# domains = json.dumps(obj[olditem]["domains"])
	# print(olditem,)

# for i in xrange(len(obj)):
#     if obj[i]["id"] == "Google":
#         obj.pop(i)
#         break

# # Output the updated file with pretty JSON                                      
# open("updated-file.json", "w").write(
#     json.dumps(obj, sort_keys=True, indent=4, separators=(',', ': '))
# )
