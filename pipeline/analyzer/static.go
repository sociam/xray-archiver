package main

import (
	// "encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sociam/xray-archiver/pipeline/util"
)

// AndroidManifest is a struct representing the interesting parts of the
// AndroidManifest.xml in APKs
type AndroidManifest struct {
	Package     string            `xml:"package,attr"`
	Perms       []util.Permission `xml:"uses-permission"`
	Sdk23Perms  []util.Permission `xml:"uses-permission-sdk-23"`
	Application manifestApp       `xml:"application"`
}

type manifestApp struct {
	Icon string `xml:"icon,attr"`
}

func parseManifest(app *util.App) (manifest *AndroidManifest, gotIcon bool, err error) {
	manifest = &AndroidManifest{}
	manifestFile, err := os.Open(path.Join(app.OutDir(), "AndroidManifest.xml"))
	if err != nil {
		return nil, false, err
	}
	bytes, err := ioutil.ReadAll(manifestFile)
	if err != nil {
		return nil, false, err
	}
	err = xml.Unmarshal(bytes, manifest)
	if err != nil {
		return nil, false, err
	}

	if manifest.Package != "" {
		app.ID = manifest.Package
	}

	split := strings.SplitN(manifest.Application.Icon, "/", 2)
	if len(split) != 2 {
		return manifest, false, nil
	}
	locn, name := split[0], split[1]
	if len(locn) < 1 {
		return manifest, false, nil
	}
	locn = path.Join(app.OutDir(), "res", locn[1:]) // /tmp/<outdir>/res/{mipmap,drawable}
	name = name + ".png"                            // icon_katana.png

	var matches []string
	if matches, err = filepath.Glob(path.Join(locn+"-*xxxdpi*", name)); err == nil && len(matches) > 0 {
	} else if matches, err = filepath.Glob(path.Join(locn+"-*xxdpi*", name)); err == nil && len(matches) > 0 {
	} else if matches, err = filepath.Glob(path.Join(locn+"-*xdpi*", name)); err == nil && len(matches) > 0 {
	} else if matches, err = filepath.Glob(path.Join(locn+"-*hdpi*", name)); err == nil && len(matches) > 0 {
	} else if matches, err = filepath.Glob(path.Join(locn+"-*tvdpi*", name)); err == nil && len(matches) > 0 {
	} else if matches, err = filepath.Glob(path.Join(locn+"-*mdpi*", name)); err == nil && len(matches) > 0 {
	} else if matches, err = filepath.Glob(path.Join(locn+"*", name)); err == nil && len(matches) > 0 {
	} else {
		return manifest, false, nil
	}

	err = os.Rename(matches[0], path.Join(app.AppDir(), "icon.png"))
	if err != nil {
		fmt.Printf("Failed to rename icon %s to %s", matches[0], path.Join(app.AppDir(), "icon.png"))
		return manifest, false, nil
	}

	return manifest, true, nil
}

func (manifest *AndroidManifest) getPerms() []util.Permission {
	return append(manifest.Perms, manifest.Sdk23Perms...)
}

type company struct {
	ID           string   `json:"id"`
	Name         string   `json:"company"`
	Domains      []string `json:"domains"`
	Founded      string   `json:"founded"`
	Acquired     string   `json:"acquired in"`
	CType        string   `json:"type"`
	TypeTag      string   `json:"typetag"`
	Jurisdiction string   `json:"jurisdiction_code"`
	Parent       string   `json:"parent"`
	Capital      string   `json:"capital"`
	Equity       string   `json:"equity"`
	Size         string   `json:"size"`
	DataSource   string   `json:"data source"`
	Description  string   `json:"description"`
}

var hostregex = regexp.MustCompile("^(?i)(?:http|https)://((?:[^!/:>[:cntrl:][:space:]]+)(?:[\.](?:aaa|aarp|abb|abbott|abbvie|abogado|abudhabi|ac|academy|accenture|accountant|accountants|aco|active|actor|ad|adac|ads|adult|ae|aeg|aero|aetna|af|afl|ag|agakhan|agency|ai|aig|airbus|airforce|airtel|akdn|al|alibaba|alipay|allfinanz|ally|alsace|alstom|am|amica|amsterdam|analytics|android|anquan|ao|apartments|app|apple|aq|aquarelle|ar|aramco|archi|army|arpa|arte|as|asia|associates|at|attorney|au|auction|audi|audible|audio|author|auto|autos|avianca|aw|aws|ax|axa|az|azure|ba|baby|baidu|band|bank|bar|barcelona|barclaycard|barclays|barefoot|bargains|bauhaus|bayern|bb|bbc|bbva|bcg|bcn|bd|be|beats|beer|bentley|berlin|best|bet|bf|bg|bh|bharti|bi|bible|bid|bike|bing|bingo|bio|biz|bj|black|blackfriday|blog|bloomberg|blue|bm|bms|bmw|bn|bnl|bnpparibas|bo|boats|boehringer|bom|bond|boo|book|boots|bosch|bostik|bot|boutique|br|bradesco|bridgestone|broadway|broker|brother|brussels|bs|bt|budapest|bugatti|build|builders|business|buy|buzz|bv|bw|by|bz|bzh|ca|cab|cafe|cal|call|cam|camera|camp|cancerresearch|canon|capetown|capital|car|caravan|cards|care|career|careers|cars|cartier|casa|cash|casino|cat|catering|cba|cbn|cc|cd|ceb|center|ceo|cern|cf|cfa|cfd|cg|ch|chanel|channel|chase|chat|cheap|chintai|chloe|christmas|chrome|church|ci|cipriani|circle|cisco|citic|city|cityeats|ck|cl|claims|cleaning|click|clinic|clinique|clothing|cloud|club|clubmed|cm|cn|co|coach|codes|coffee|college|cologne|com|commbank|community|company|compare|computer|comsec|condos|construction|consulting|contact|contractors|cooking|cool|coop|corsica|country|coupon|coupons|courses|cr|credit|creditcard|creditunion|cricket|crown|crs|cruises|csc|cu|cuisinella|cv|cw|cx|cy|cymru|cyou|cz|dabur|dad|dance|date|dating|datsun|day|dclk|dds|de|deal|dealer|deals|degree|delivery|dell|deloitte|delta|democrat|dental|dentist|desi|design|dev|dhl|diamonds|diet|digital|direct|directory|discount|dj|dk|dm|dnp|do|docs|dog|doha|domains|dot|download|drive|dtv|dubai|dunlop|dupont|durban|dvag|dz|earth|eat|ec|edeka|edu|education|ee|eg|email|emerck|energy|engineer|engineering|enterprises|epost|epson|equipment|er|ericsson|erni|es|esq|estate|et|eu|eurovision|eus|events|everbank|exchange|expert|exposed|express|extraspace|fage|fail|fairwinds|faith|family|fan|fans|farm|fashion|fast|feedback|ferrero|fi|film|final|finance|financial|fire|firestone|firmdale|fish|fishing|fit|fitness|fj|fk|flickr|flights|flir|florist|flowers|flsmidth|fly|fm|fo|foo|football|ford|forex|forsale|forum|foundation|fox|fr|fresenius|frl|frogans|frontier|ftr|fund|furniture|futbol|fyi|ga|gal|gallery|gallo|gallup|game|games|garden|gb|gbiz|gd|gdn|ge|gea|gent|genting|gf|gg|ggee|gh|gi|gift|gifts|gives|giving|gl|glass|gle|global|globo|gm|gmail|gmbh|gmo|gmx|gn|gold|goldpoint|golf|goo|goodyear|goog|google|gop|got|gov|gp|gq|gr|grainger|graphics|gratis|green|gripe|group|gs|gt|gu|guardian|gucci|guge|guide|guitars|guru|gw|gy|hamburg|hangout|haus|hdfcbank|health|healthcare|help|helsinki|here|hermes|hiphop|hisamitsu|hitachi|hiv|hk|hkt|hm|hn|hockey|holdings|holiday|homedepot|homes|honda|horse|host|hosting|hoteles|hotmail|house|how|hr|hsbc|ht|htc|hu|hyundai|ibm|icbc|ice|icu|id|ie|ifm|iinet|il|im|imamat|imdb|immo|immobilien|in|industries|infiniti|info|ing|ink|institute|insurance|insure|int|international|investments|io|ipiranga|iq|ir|irish|is|iselect|ismaili|ist|istanbul|it|itau|iwc|jaguar|java|jcb|jcp|je|jetzt|jewelry|jlc|jll|jm|jmp|jnj|jo|jobs|joburg|jot|joy|jp|jpmorgan|jprs|juegos|kaufen|kddi|ke|kerryhotels|kerrylogistics|kerryproperties|kfh|kg|kh|ki|kia|kim|kinder|kindle|kitchen|kiwi|km|kn|koeln|komatsu|kosher|kp|kpmg|kpn|kr|krd|kred|kuokgroup|kw|ky|kyoto|kz|la|lacaixa|lamborghini|lamer|lancaster|land|landrover|lanxess|lasalle|lat|latrobe|law|lawyer|lb|lc|lds|lease|leclerc|legal|lego|lexus|lgbt|li|liaison|lidl|life|lifeinsurance|lifestyle|lighting|like|limited|limo|lincoln|linde|link|lipsy|live|living|lixil|lk|loan|loans|locker|locus|lol|london|lotte|lotto|love|lr|ls|lt|ltd|ltda|lu|lupin|luxe|luxury|lv|ly|ma|madrid|maif|maison|makeup|man|management|mango|market|marketing|markets|marriott|mattel|mba|mc|md|me|med|media|meet|melbourne|meme|memorial|men|menu|meo|metlife|mg|mh|miami|microsoft|mil|mini|mk|ml|mlb|mls|mm|mma|mn|mo|mobi|mobily|moda|moe|moi|mom|monash|money|montblanc|mormon|mortgage|moscow|motorcycles|mov|movie|movistar|mp|mq|mr|ms|mt|mtn|mtpc|mtr|mu|museum|mutual|mutuelle|mv|mw|mx|my|mz|na|nadex|nagoya|name|natura|navy|nc|ne|nec|net|netbank|netflix|network|neustar|new|news|next|nextdirect|nexus|nf|ng|ngo|nhk|ni|nico|nikon|ninja|nissan|nissay|nl|no|nokia|northwesternmutual|norton|now|nowruz|nowtv|np|nr|nra|nrw|ntt|nu|nyc|nz|obi|office|okinawa|olayan|olayangroup|ollo|om|omega|one|ong|onl|online|ooo|oracle|orange|org|organic|origins|osaka|otsuka|ott|ovh|pa|page|pamperedchef|panerai|paris|pars|partners|parts|party|passagens|pccw|pe|pet|pf|pg|ph|pharmacy|philips|photo|photography|photos|physio|piaget|pics|pictet|pictures|pid|pin|ping|pink|pioneer|pizza|pk|pl|place|play|playstation|plumbing|plus|pm|pn|pohl|poker|porn|post|pr|praxi|press|prime|pro|prod|productions|prof|progressive|promo|properties|property|protection|ps|pt|pub|pw|pwc|py|qa|qpon|quebec|quest|racing|re|read|realestate|realtor|realty|recipes|red|redstone|redumbrella|rehab|reise|reisen|reit|ren|rent|rentals|repair|report|republican|rest|restaurant|review|reviews|rexroth|rich|richardli|ricoh|rio|rip|ro|rocher|rocks|rodeo|room|rs|rsvp|ru|ruhr|run|rw|rwe|ryukyu|sa|saarland|safe|safety|sakura|sale|salon|samsung|sandvik|sandvikcoromant|sanofi|sap|sapo|sarl|sas|save|saxo|sb|sbi|sbs|sc|sca|scb|schaeffler|schmidt|scholarships|school|schule|schwarz|science|scor|scot|sd|se|seat|security|seek|select|sener|services|seven|sew|sex|sexy|sfr|sg|sh|sharp|shaw|shell|shia|shiksha|shoes|shop|shouji|show|shriram|si|silk|sina|singles|site|sj|sk|ski|skin|sky|skype|sl|sm|smile|sn|sncf|so|soccer|social|softbank|software|sohu|solar|solutions|song|sony|soy|space|spiegel|spot|spreadbetting|sr|srl|st|stada|star|starhub|statebank|statefarm|statoil|stc|stcgroup|stockholm|storage|store|stream|studio|study|style|su|sucks|supplies|supply|support|surf|surgery|suzuki|sv|swatch|swiss|sx|sy|sydney|symantec|systems|sz|tab|taipei|talk|taobao|tatamotors|tatar|tattoo|tax|taxi|tc|tci|td|tdk|team|tech|technology|tel|telecity|telefonica|temasek|tennis|teva|tf|tg|th|thd|theater|theatre|tickets|tienda|tiffany|tips|tires|tirol|tj|tk|tl|tm|tmall|tn|to|today|tokyo|tools|top|toray|toshiba|total|tours|town|toyota|toys|tr|trade|trading|training|travel|travelers|travelersinsurance|trust|trv|tt|tube|tui|tunes|tushu|tv|tvs|tw|tz|ua|ubs|ug|uk|unicom|university|uno|uol|ups|us|uy|uz|va|vacations|vana|vc|ve|vegas|ventures|verisign|versicherung|vet|vg|vi|viajes|video|vig|viking|villas|vin|vip|virgin|vision|vista|vistaprint|viva|vlaanderen|vn|vodka|volkswagen|vote|voting|voto|voyage|vu|vuelos|wales|walter|wang|wanggou|warman|watch|watches|weather|weatherchannel|webcam|weber|website|wed|wedding|weibo|weir|wf|whoswho|wien|wiki|williamhill|win|windows|wine|wme|wolterskluwer|work|works|world|ws|wtc|wtf|xbox|xerox|xihuan|xin|xn--11b4c3d|xn--1ck2e1b|xn--1qqw23a|xn--30rr7y|xn--3bst00m|xn--3ds443g|xn--3e0b707e|xn--3pxu8k|xn--42c2d9a|xn--45brj9c|xn--45q11c|xn--4gbrim|xn--55qw42g|xn--55qx5d|xn--5tzm5g|xn--6frz82g|xn--6qq986b3xl|xn--80adxhks|xn--80ao21a|xn--80asehdb|xn--80aswg|xn--8y0a063a|xn--90a3ac|xn--90ais|xn--9dbq2a|xn--9et52u|xn--9krt00a|xn--b4w605ferd|xn--bck1b9a5dre4c|xn--c1avg|xn--c2br7g|xn--cck2b3b|xn--cg4bki|xn--clchc0ea0b2g2a9gcd|xn--czr694b|xn--czrs0t|xn--czru2d|xn--d1acj3b|xn--d1alf|xn--e1a4c|xn--eckvdtc9d|xn--efvy88h|xn--estv75g|xn--fct429k|xn--fhbei|xn--fiq228c5hs|xn--fiq64b|xn--fiqs8s|xn--fiqz9s|xn--fjq720a|xn--flw351e|xn--fpcrj9c3d|xn--fzc2c9e2c|xn--fzys8d69uvgm|xn--g2xx48c|xn--gckr3f0f|xn--gecrj9c|xn--h2brj9c|xn--hxt814e|xn--i1b6b1a6a2e|xn--imr513n|xn--io0a7i|xn--j1aef|xn--j1amh|xn--j6w193g|xn--jlq61u9w7b|xn--jvr189m|xn--kcrx77d1x4a|xn--kprw13d|xn--kpry57d|xn--kpu716f|xn--kput3i|xn--l1acc|xn--lgbbat1ad8j|xn--mgb9awbf|xn--mgba3a3ejt|xn--mgba3a4f16a|xn--mgba7c0bbn0a|xn--mgbaam7a8h|xn--mgbab2bd|xn--mgbayh7gpa|xn--mgbb9fbpob|xn--mgbbh1a71e|xn--mgbc0a9azcg|xn--mgbca7dzdo|xn--mgberp4a5d4ar|xn--mgbpl2fh|xn--mgbt3dhd|xn--mgbtx2b|xn--mgbx4cd0ab|xn--mix891f|xn--mk1bu44c|xn--mxtq1m|xn--ngbc5azd|xn--ngbe9e0a|xn--node|xn--nqv7f|xn--nqv7fs00ema|xn--nyqy26a|xn--o3cw4h|xn--ogbpf8fl|xn--p1acf|xn--p1ai|xn--pbt977c|xn--pgbs0dh|xn--pssy2u|xn--q9jyb4c|xn--qcka1pmc|xn--qxam|xn--rhqv96g|xn--rovu88b|xn--s9brj9c|xn--ses554g|xn--t60b56a|xn--tckwe|xn--unup4y|xn--vermgensberater-ctb|xn--vermgensberatung-pwb|xn--vhquv|xn--vuq861b|xn--w4r85el8fhu5dnra|xn--w4rs40l|xn--wgbh1c|xn--wgbl6a|xn--xhq521b|xn--xkc2al3hye2a|xn--xkc2dl3a5ee0h|xn--y9a3aq|xn--yfro4i67o|xn--ygbi2ammx|xn--zfr164b|xperia|xxx|xyz|yachts|yahoo|yamaxun|yandex|ye|yodobashi|yoga|yokohama|you|youtube|yt|yun|za|zappos|zara|zero|zip|zm|zone|zuerich|zw)))[^\.[:alnum:]-]")

func simpleAnalyze(app *util.App) ([]string, error) {
	//TODO: fix error handling

	// //TODO: replace with DB calls
	// var companies map[string]company
	// companyFile, err := os.Open(path.Join(util.Cfg.DataDir, "company_details.json"))
	// if err != nil {
	// 	return []string{}, err
	// }
	// bytes, err := ioutil.ReadAll(companyFile)
	// if err != nil {
	// 	return []string{}, err
	// }
	// err = json.Unmarshal(bytes, &companies)
	// if err != nil {
	// 	return []string{}, err
	// }

	// for name := range companies {
	// 	if _, ok := trackers[name]; !ok {
	// 		delete(companies, name)
	// 	}
	// }

	// getDomainCo := func(host string) *string {
	// 	for _, company := range companies {
	// 		for _, domain := range company.domains {
	// 			if strings.Contains(host, domain) {
	// 				return &company.id
	// 			}
	// 		}
	// 	}
	// 	return nil
	// }

	cmd := exec.Command("strings", "-n", "11", path.Join(app.OutDir(), "classes.dex"))

	out, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}
	matches := hostregex.FindAllSubmatch(out, -1)
	urls := make([]string, 0, len(matches))
	for _, v := range matches {
		if len(v[1]) > 3 {
			urls = append(urls, string(v[1]))
		}
	}

	urls = util.Dedup(urls)

	// var appTrackers []string

	// irrelevant := []string{"app", "identity", "n/a", "other", "", "library"}
	// for name, company := range companies {
	// 	for _, domain := range company.Domains {
	// 		if strings.Contains(string(urls), string(domain)) {
	// 			toAppend := true
	// 			for _, cat := range irrelevant {
	// 				if companies[name].TypeTag == cat {
	// 					toAppend = false
	// 					break
	// 				}
	// 			}
	// 			if toAppend {
	// 				appTrackers = append(appTrackers, name)
	// 			}
	// 		}
	// 	}
	// }

	return urls, nil
}

func checkReflect(app *util.App) error {
	cmd := exec.Command("grep", "-Paqh",
		"\\x00\\x00\\x00.Ljava/lang/reflect[/a-zA-Z]*;\\x00\\x00\\x00",
		"--", path.Join(app.OutDir(), "classes.dex"))

	out, err := cmd.Output()
	if err != nil && strings.TrimSpace(string(out)) != "" {
		fmt.Printf("Error checking for reflection: output below\n%s\n\n", string(out))
		return err
	}

	app.UsesReflect = err == nil

	return nil
}

func findPackages(app *util.App) ([]string, error) {
	// TODO: fix error handling
	paths := make(map[string]util.Unit)
	err := os.Chdir(path.Join(app.OutDir(), "smali"))
	if err != nil {
		return []string{}, err
	}

	err = filepath.Walk(".",
		func(fname string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if path.Ext(fname) == ".smali" {
				paths[path.Dir(fname)] = unit
			}

			return nil
		})

	pkgs := make([]string, 0, 20)
	for path := range paths {
		//pkg := strings.Replace(path, string(os.PathSeparator), ".", -1)
		pkg := strings.Map(func(ch rune) rune {
			if ch == os.PathSeparator {
				return '.'
			}
			return ch
		}, path)
		pkgs = append(pkgs, pkg)
	}

	return pkgs, err
}
