package main

import (
	"github.com/sociam/xray-archiver/pipeline/util"
)

var unit = util.Unit{}

var badHosts = map[string]util.Unit{
	"example.com": unit,
}

var mobTrackers = map[string]util.Unit{
	"Inmobi": unit, "Mopub": unit, "Facebook": unit, "Chartboost": unit,
	"Crashlytics": unit, "Heyzap": unit, "Applovin": unit,
	"Unitytechnologies": unit, "Vungle": unit, "Fyber": unit,
	"Millenialmedia": unit, "Ironsource": unit, "Tapjoy": unit,
	"Flurry": unit, "Google": unit, "Bitstadium": unit, "Presage": unit,
	"Comscore": unit, "Umeng": unit, "Tencent": unit, "Nexage": unit,
	"Mixpanel": unit, "Appboy": unit, "Yandex": unit, "Kochava": unit,
	"Urbanairship": unit, "Daum": unit, "Milkman": unit, "Playhaven": unit,
	"Mail.Ru": unit, "Newrelic": unit, "Admob": unit, "Moat": unit,
	"Smaato": unit, "Admarvel": unit, "Appodeal": unit, "Cauley": unit,
	"Revmob": unit, "Adbuddiz": unit, "Igaworks": unit, "Nativex": unit,
	"Smartadserver": unit, "Mobfox": unit, "Pubnative": unit,
	"Hyprmx": unit, "Pingstart": unit, "Tune": unit, "Baidu": unit,
	"Swrve": unit, "Tnkfactory": unit, "Appnext": unit,
}

var webTrackers = map[string]util.Unit{
	"Google": unit, "Doubleclick": unit, "Facebook": unit,
	"Cloudfront": unit, "Scorecardresearch": unit, "Appnexus": unit,
	"Twitter": unit, "Amazon": unit, "Criteo": unit, "Yahoo": unit,
	"Oracle": unit, "Quantcast": unit, "Iponweb": unit, "Adobe": unit,
	"Openx": unit, "Ghostery": unit, "Rubicon": unit, "Bluekai": unit,
	"Cloudflare": unit, "Liveramp": unit, "The Nielsen Company": unit,
	"Demdex": unit, "Truste": unit, "Newrelic": unit,
	"Integral Ad Science": unit, "Adsrvr": unit, "Moatads": unit,
	"Pubmatic": unit, "Turn": unit, "Indexchange": unit,
	"Doubleverify": unit, "Mathtag": unit, "Chartbeat": unit,
	"Optimizely": unit, "Youtube": unit, "Tapad": unit, "Exelate": unit,
	"Targus": unit, "Lotame": unit, "Sizmek": unit, "Iasds": unit,
	"Neustar": unit, "Jquery": unit, "Krxd": unit, "Adroll": unit,
	"Brightroll": unit, "Audiencescience": unit, "Bootstrap": unit,
	"Rocketfuel": unit, "Dataxu": unit, "Tubemogul": unit,
	"Microsoft": unit, "Yandex": unit, "Conversant": unit,
	"Videology": unit, "Flashtalking": unit, "Akamai": unit,
	"Adap.Tv": unit, "Magnetic": unit, "Dstillery": unit, "Baidu": unit,
	"Smartadserver": unit, "Taboola": unit, "H&R Block": unit,
	"Tealium": unit, "Convertro": unit, "MSN": unit, "Linkedin": unit,
	"Mookie1": unit, "Contextweb": unit, "Liverail": unit, "Yadro": unit,
	"Bug": unit, "Chango": unit, "Conversant Llc": unit, "Trueffect": unit,
	"Signal-Privacy": unit, "Crazy Egg": unit, "Drawbridge Inc": unit,
	"Sovrn": unit, "Ovh": unit, "Exponential": unit, "Yieldoptimizer": unit,
	"Spotxchange": unit, "Owneriq": unit, "Insightexpress": unit,
	"Adform": unit, "Automattic": unit, "Disqus": unit, "Outbrain": unit,
	"Flipps": unit, "Centromedia": unit, "Ignitionone": unit,
	"Skimlinks": unit, "Simpli.Fi": unit, "Wordpress": unit,
	"Openx Technologies": unit, "Ensighten": unit,
	"Visual Website Optimizer": unit, "Rhythmone": unit,
	"Mixpanel": unit, "Sonobi": unit, "Hotjar": unit, "Voicefive": unit,
	"Gssprt": unit, "At Internet": unit, "Ml314": unit, "Radiumone": unit,
	"Parsely": unit, "Improve Digital Bv": unit, "Gumgum": unit,
	"Amgdgt": unit, "Effective Measure": unit, "Pinterest": unit,
	"Vk": unit, "Adition": unit, "Exoclick": unit, "Ixi Corporation": unit,
	"Nugg.Ad": unit, "Umeng": unit, "Alibaba": unit, "Marketo": unit,
	"Gskinner": unit, "Gigya": unit, "Tns": unit, "Eq Works": unit,
	"Cxense Asa": unit, "Steelhouse": unit, "Yieldlab": unit,
	"33Across": unit, "Mail.Ru": unit, "Pingdom Ab": unit,
	"Smart Adserver": unit, "Brightcove": unit, "Postrelease": unit,
	"Adriver": unit, "Pagefair": unit, "Microad": unit,
	"Marin Software": unit, "Collective": unit, "Trafficjunky": unit,
	"Pixalate": unit, "Demandbase": unit, "Yieldbot": unit,
	"Visualdna": unit, "Ioam": unit, "Signal": unit, "Shopzilla": unit,
	"Sharethrough": unit, "Rambler": unit, "Clicktale": unit,
	"Sociomantic": unit, "Datonics": unit, "Run": unit, "Myspace": unit,
	"Taobao": unit, "Adgear Technologies": unit, "Sharethis": unit,
	"Mybuys": unit, "Stroeer Media": unit,
	"The Mcgraw-Hill Companies": unit, "Vizury": unit, "Abc": unit,
	"Dyn": unit, "Longtail Ad Solutions": unit, "Netseer": unit,
	"Smaato": unit, "The Adex Gmbh": unit, "Adblade": unit,
	"Mxptint": unit, "Maxymiser": unit, "Admeta": unit, "Mcro": unit,
	"Sourcepoint Technologies": unit, "Yieldmanager": unit,
	"Rfihub": unit, "Bounce Exchange": unit, "Wtp": unit,
	"Convertmedia": unit, "Rackspace Us": unit, "Eyeview": unit,
	"Qualtrics": unit, "Intent Iq": unit, "Stickyads.Tv": unit,
	"Soasta": unit, "Cedexis Inc": unit, "Weborama": unit,
	"Invite Media": unit, "Histats": unit, "Statcounter": unit,
	"Paypal": unit, "Bidr": unit, "Springserve": unit,
	"Nexstar Broadcasting": unit, "Triplelift": unit,
	"Polar": unit, "Adelphic Inc": unit, "Deep Forest Media": unit,
	"Cedexis": unit, "Richrelevance": unit, "Resonate": unit,
	"Impact-Ad": unit, "Jivox": unit, "Mediade": unit,
	"Liveperson": unit, "Zedo": unit, "Yabidos": unit, "Adfox": unit,
	"Fout": unit, "Switch": unit, "Navegg S.A.": unit, "Adsnative": unit,
	"Eyeota Limited": unit, "Jsdelivr": unit,
	"Consumerinfo.Com": unit, "Brandscreen Pty Ltd": unit,
	"Mythings": unit, "Polar Mobile Group": unit,
	"Answercloud": unit, "Komoona": unit, "Meltdsp": unit,
	"Appier Inc": unit, "Netdna": unit, "Budgetedbauer": unit,
	"Gfk": unit, "Bidtheatre": unit, "Research Now": unit,
	"Undertone Networks": unit, "Viglink": unit, "Adbrain": unit,
	"Adstir": unit, "Sailthru": unit, "Vimeo": unit, "Mediametrie": unit,
	"Livefyre": unit, "Mediaforge": unit, "Webklipper": unit,
	"Tencent": unit, "Ib-Ibi": unit,
}

var trackers = util.Combine(mobTrackers, webTrackers)
