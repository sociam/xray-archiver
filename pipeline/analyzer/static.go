package main

import "C"

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

	//"github.com/sociam/xray-archiver/pipeline/util"
	"../util"
)

var badHosts = map[string]util.Unit{
	"example.com": util.Unit{},
}

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

var hostre = regexp.MustCompile("(?i)https?://((?:[^@!/:>[:cntrl:]\\s]+)\\.(?:(?:a(?:a(?:a|rp)|b(?:arth|b(?:ott|vie)|le|ogado|udhabi|[bc])|c(?:ademy|c(?:enture|ountants?)|o|t(?:ive|or))|d(?:ac|s|ult)|e(?:g|ro|tna)|f(?:amilycompany|l|rica)|g(?:akhan|ency)|i(?:go?|r(?:bus|force|tel))|kdn|l(?:faromeo|i(?:baba|pay)|l(?:finanz|state|y)|s(?:ace|tom))|m(?:e(?:rican(?:express|family)|x)|fam|ica|sterdam)|n(?:alytics|droid|quan|z)|ol|p(?:artments|p(?:le)?)|quarelle|r(?:a(?:b|mco)|chi|my|pa|te?)|s(?:da|ia|sociates)|t(?:hleta|torney)|u(?:ction|di(?:ble|o)?|spost|t(?:hor|os?))|vianca|ws|xa|zure|[c-gilmoq-uwxz])|b(?:a(?:by|idu|n(?:a(?:mex|narepublic)|[dk])|r(?:c(?:elona|lay(?:card|s))|efoot|gains)?|s(?:(?:e|ket)ball)|uhaus|yern)|b(?:va|[ct])|c[gn]|e(?:a(?:ts|uty)|er|ntley|rlin|st(?:buy)?|t)|harti|i(?:ble|ke|ngo?|[doz])|l(?:a(?:ck(?:friday)?|nco)|o(?:ckbuster|(?:omber)?g)|ue)|m[sw]|n(?:l|pparibas)|o(?:ats|ehringer|fa|nd|o(?:k(?:ing)?|ts)|s(?:ch|t(?:ik|on))|utique|[motx])|r(?:adesco|idgestone|o(?:adway|(?:k|th)er)|ussels)|u(?:dapest|gatti|ild(?:ers)?|siness|y|zz)|zh|[abd-jmnorstvwyz])|c(?:a(?:fe|l(?:l|vinklein)|m(?:era|p)|n(?:cerresearch|on)|p(?:etown|ital(?:one)?)|r(?:avan|ds|eers?|tier|[es])|s(?:eih|ino|[aeh])|t(?:ering|holic)|[blmrt])|b(?:re|[ans])|e(?:nter|rn|[bo])|f[ad]|h(?:a(?:n(?:n?el)|se|t)|eap|intai|loe|r(?:istmas|ome|ysler)|urch)|i(?:priani|rcle|sco|t(?:adel|ic|yeats|[iy]))|l(?:aims|eaning|i(?:ck|ni(?:c|que))|o(?:thing|ud)|ub(?:med)?)|o(?:ach|des|ffee|l(?:(?:leg|ogn)e)|m(?:cast|m(?:bank|unity)|p(?:a(?:ny|re)|uter)|sec)?|n(?:dos|s(?:truction|ulting)|t(?:act|ractors))|o(?:king(?:channel)?|[lp])|rsica|u(?:ntry|pons?|rses))|r(?:edit(?:card|union)?|icket|own|s|uises?)|sc|uisinella|y(?:(?:mr|o)u)|[acdf-ik-oru-z])|d(?:a(?:bur|nce|t(?:ing|sun|[ae])|[dy])|clk|ds|e(?:al(?:er|s)?|gree|l(?:ivery|l|oitte|ta)|mocrat|nt(?:al|ist)|si(?:gn)?|v)|hl|i(?:amonds|et|gital|rect(?:ory)?|s(?:co(?:unt|ver)|h)|y)|np|o(?:c(?:s|tor)|dge|ha|mains|wnload|[gt])|rive|tv|u(?:bai|ck|n(?:lop|s)|pont|rban)|v(?:ag|r)|[ejkmoz])|e(?:a(?:rth|t)|co|d(?:eka|u(?:cation)?)|m(?:ail|erck)|n(?:ergy|gineer(?:ing)?|terprises)|p(?:ost|son)|quipment|r(?:icsson|ni)|s(?:q|(?:tat|uranc)e)|tisalat|u(?:rovision|s)|ve(?:nts|rbank)|x(?:change|p(?:ert|osed|ress)|traspace)|[cegr-u])|f(?:a(?:ge|i(?:l|rwinds|th)|mily|ns?|rm(?:ers)?|s(?:hion|t))|e(?:dex|edback|rr(?:ari|ero))|i(?:at|d(?:elity|o)|lm|na(?:l|nc(?:e|ial))|r(?:(?:eston|mdal)?e)|sh(?:ing)?|t(?:ness)?)|l(?:i(?:ckr|ghts|r)|o(?:rist|wers)|y)|o(?:o(?:d(?:network)?|tball)|r(?:d|ex|sale|um)|undation|[ox])|r(?:e(?:e|senius)|l|o(?:gans|nt(?:(?:doo|ie)r)))|tr|u(?:ji(?:tsu|xerox)|nd?|rniture|tbol)|yi|[ijkmor])|g(?:a(?:ll(?:ery|o|up)|mes?|rden|[lp])|biz|dn|e(?:a|nt(?:ing)?|orge)|gee|i(?:fts?|v(?:es|ing))|l(?:a(?:de|ss)|e|ob(?:al|o))|m(?:ail|bh|[ox])|o(?:daddy|l(?:dpoint|[df])|o(?:d(?:hands|year)|g(?:le)?)|[optv])|r(?:a(?:inger|(?:phic|ti)s)|een|ipe|o(?:cery|up))|u(?:ardian|cci|ge|i(?:de|tars)|ru)|[abd-ilmnp-uwy])|h(?:a(?:ir|mburg|ngout|us)|bo|dfc(?:bank)?|e(?:alth(?:care)?|l(?:p|sinki)|r(?:e|mes))|gtv|i(?:phop|samitsu|tachi|v)|kt|o(?:ckey|l(?:dings|iday)|me(?:depot|goods|s(?:ense)?)|n(?:da|eywell)|rse|s(?:pital|t(?:ing)?)|t(?:el(?:e?s)|mail)|use|[tw])|sbc|tc|ughes|y(?:att|undai)|[kmnrtu])|i(?:bm|c(?:bc|[eu])|eee|fm|kano|m(?:amat|db|mo(?:bilien)?)|n(?:dustries|f(?:initi|o)|s(?:(?:titut|ur(?:anc)?)e)|t(?:e(?:(?:rnationa)?l)|uit)|vestments|[gkt])|piranga|rish|s(?:elect|maili|t(?:anbul)?)|t(?:au|v)|veco|wc|[del-oq-t])|j(?:a(?:guar|va)|c[bp]|e(?:ep|tzt|welry)|io|l[cl]|mp|nj|o(?:b(?:s|urg)|[ty])|p(?:morgan|rs)|u(?:egos|niper)|[emop])|k(?:aufen|ddi|erry(?:(?:hotel|logistic|propertie)s)|fh|i(?:nd(?:er|le)|tchen|wi|[am])|o(?:eln|matsu|sher)|p(?:mg|n)|r(?:e?d)|uokgroup|yoto|[eghimnprwyz])|l(?:a(?:caixa|dbrokes|m(?:borghini|er)|n(?:c(?:aster|ia|ome)|d(?:rover)?|xess)|salle|t(?:ino|robe)|wyer|[tw])|ds|e(?:ase|clerc|frak|g(?:al|o)|xus)|gbt|i(?:aison|dl|fe(?:(?:insuranc|styl)e)?|ghting|ke|lly|m(?:ited|o)|n(?:coln|de|k)|psy|v(?:e|ing)|xil)|o(?:ans?|c(?:ker|us)|ft|l|ndon|tt[eo]|ve)|pl(?:financial)?|tda?|u(?:ndbeck|pin|x(?:e|ury))|[abcikr-vy])|m(?:a(?:cys|drid|i(?:f|son)|keup|n(?:agement|go)|r(?:ket(?:ing|s)?|riott|shalls)|serati|ttel|[np])|ba|c(?:d(?:onalds)?|kinsey)|e(?:dia|et|lbourne|m(?:e|orial)|nu|rckmsd|tlife|[dno])|i(?:ami|crosoft|n[it]|tsubishi|[lt])|l[bs]|ma|o(?:bi(?:l[ey])?|da|n(?:ash|ey|ster|tblanc)|par|r(?:mon|tgage)|scow|to(?:rcycles)?|vi(?:e|star)|[eimv])|sd|t[nr]|u(?:seum|tual)|[acdeghk-z])|n(?:a(?:b|dex|goya|me|t(?:ionwide|ura)|vy)|ba|e(?:t(?:bank|flix|work)|ustar|w(?:holland|s)|x(?:t(?:direct)?|us)|[ctw])|fl|go|hk|i(?:co|k(?:e|on)|nja|ssa[ny])|o(?:kia|rt(?:hwesternmutual|on)|w(?:ruz|tv)?)|r[aw]|tt|yc|[acefgilopruz])|o(?:b(?:i|server)|ff(?:ice)?|kinawa|l(?:ayan(?:group)?|dnavy|lo)|m(?:ega)?|n(?:(?:lin|yoursid)e|[egl])|oo|pen|r(?:a(?:(?:cl|ng)e)|g(?:anic)?|igins)|saka|t(?:suka|t)|vh)|p(?:a(?:ge|mperedchef|n(?:asonic|erai)|r(?:is|s|t(?:ners|[sy]))|ssagens|y)|ccw|et|fizer|h(?:armacy|d|ilips|o(?:ne|to(?:graphy|s)?)|ysio)|i(?:aget|c(?:s|t(?:et|ures))|n[gk]|oneer|zza|[dn])|l(?:a(?:ce|y(?:station)?)|u(?:mbing|s))|nc|o(?:hl|ker|litie|rn|st)|r(?:a(?:merica|xi)|ess|ime|o(?:ductions|gressive|mo|pert(?:ies|y)|tection|[df])|udential|[ou])|ub|wc|[ae-hk-nrstwy])|q(?:a|pon|ue(?:bec|st)|vc)|r(?:a(?:cing|dio|id)|e(?:a(?:d|l(?:estate|t(?:or|y)))|cipes|d(?:stone|umbrella)|hab|i(?:sen?|t)|liance|nt(?:als)?|p(?:air|ort|ublican)|st(?:aurant)?|views?|xroth|[dn])|i(?:c(?:h(?:ardli)?|oh)|ghtathome|[lop])|mit|o(?:c(?:her|ks)|deo|gers|om)|svp|u(?:gby|hr|n)|we|yukyu|[eosuw])|s(?:a(?:arland|fe(?:ty)?|kura|l(?:e|on)|ms(?:club|ung)|n(?:dvik(?:coromant)?|ofi)|po|rl|ve|xo|[ps])|b[is]|c(?:h(?:aeffler|midt|o(?:larships|ol)|ule|warz)|ience|johnson|o[rt]|[ab])|e(?:a(?:rch|t)|cur(?:e|ity)|ek|lect|ner|rvices|ven|xy|[swx])|fr|h(?:a(?:ngrila|rp|w)|ell|i(?:(?:ksh)?a)|o(?:es|pping|uji|wtime|[pw])|riram)|i(?:lk|n(?:a|gles)|te)|k(?:in|ype|[iy])|ling|m(?:art|ile)|ncf|o(?:c(?:cer|ial)|ft(?:bank|ware)|hu|l(?:ar|utions)|n[gy]|y)|p(?:ace|iegel|ot|readbetting)|r[lt]|t(?:a(?:da|ples|r(?:hub)?|t(?:e(?:bank|farm)|oil))|c(?:group)?|o(?:ckholm|r(?:(?:ag)?e))|ream|ud(?:io|y)|yle)|u(?:cks|pp(?:l(?:ies|y)|ort)|r(?:f|gery)|zuki)|w(?:atch|i(?:ftcover|ss))|y(?:dney|mantec|stems)|[a-eg-ortuvxyz])|t(?:a(?:ipei|lk|obao|rget|t(?:a(?:motors|r)|too)|xi|[bx])|ci|dk|e(?:am|ch(?:nology)?|l(?:e(?:city|fonica))?|masek|nnis|va)|h(?:d|eat(?:er|re))|i(?:aa|ckets|enda|ffany|ps|r(?:es|ol))|j(?:(?:max)?x)|kmaxx|mall|o(?:day|kyo|ols|p|ray|shiba|tal|urs|wn|y(?:ota|s))|r(?:a(?:d(?:e|ing)|ining|vel(?:channel|ers(?:insurance)?)?)|ust|v)|u(?:be|i|nes|shu)|vs|[cdfghj-ortvwz])|u(?:b(?:ank|s)|connect|n(?:i(?:com|versity)|o)|ol|ps|[agksyz])|v(?:a(?:cations|n(?:a|guard))|e(?:gas|ntures|r(?:isign|sicherung)|t)|i(?:ajes|deo|king|llas|rgin|s(?:a|ion|ta(?:print)?)|v[ao]|[gnp])|laanderen|o(?:dka|l(?:kswagen|vo)|t(?:ing|[eo])|yage)|uelos|[aceginu])|w(?:a(?:l(?:es|mart|ter)|ng(?:gou)?|rman|tch(?:es)?)|e(?:ather(?:channel)?|b(?:cam|er|site)|d(?:ding)?|i(?:bo|r))|hoswho|i(?:en|ki|lliamhill|n(?:dows|e|ners)?)|me|o(?:lterskluwer|odside|r(?:ks?|ld)|w)|t[cf]|[fs])|x(?:box|erox|finity|i(?:(?:hua)?n)|n--(?:1(?:1b4c3d|ck2e1b|qqw23a)|2scrj9c|3(?:0rr7y|bst00m|ds443g|e0b707e|hcrj9c|oq18vl8pn36a|pxu8k)|4(?:2c2d9a|5(?:br(?:5cyl|j9c)|q11c)|gbrim)|5(?:4b7fta0cc|5q(?:w42g|x5d)|(?:su34j936bgs|tzm5)g)|6(?:frz82g|qq986b3xl)|8(?:0a(?:dxhks|o21a|qecdr1a|s(?:ehdb|wg))|y0a063a)|9(?:0a(?:3ac|e|is)|dbq2a|et52u|krt00a)|b(?:4w605ferd|ck1b9a5dre4c)|c(?:1avg|2br7g|ck2b3b|g4bki|lchc0ea0b2g2a9gcd|zr(?:694b|s0t|u2d))|d1a(?:cj3b|lf)|e(?:1a4c|ckvdtc9d|fvy88h|stv75g)|f(?:ct429k|hbei|iq(?:228c5hs|64b|(?:s8|z9)s)|jq720a|lw351e|pcrj9c3d|z(?:c2c9e2c|ys8d69uvgm))|g(?:2xx48c|ckr3f0f|ecrj9c|k3at1e)|h(?:2br(?:eg3eve|j9c(?:8c)?)|xt814e)|i(?:1b6b1a6a2e|mr513n|o0a7i)|j(?:1a(?:ef|mh)|6w193g|lq61u9w7b|vr189m)|k(?:crx77d1x4a|p(?:r(?:(?:w13|y57)d)|u(?:716f|t3i)))|l(?:1acc|gbbat1ad8j)|m(?:gb(?:9awbf|a(?:3a(?:3ejt|4f16a)|7c0bbn0a|a(?:kc7dvf|m7a8h)|b2bd|i9azgqp6j|yh7gpa)|b(?:9fbpob|h1a(?:71e)?)|c(?:0a9azcg|a7dzdo)|erp4a5d4ar|gu82a|i4ecexp|pl2fh|t(?:3dhd|x2b)|x4cd0ab)|ix891f|k1bu44c|xtq1m)|n(?:gb(?:c5azd|e9e0a|rx)|ode|qv7f(?:s00ema)?|yqy26a)|o(?:3cw4h|gbpf8fl)|p(?:1a(?:cf|i)|bt977c|gbs0dh|ssy2u)|q(?:9jyb4c|cka1pmc|xam)|r(?:hqv96g|ovu88b|vc1e0am3e)|s(?:9brj9c|es554g)|t(?:60b56a|ckwe|iq49xqyj)|unup4y|v(?:ermgensberat(?:(?:er-ct|ung-pw)b)|hquv|uq861b)|w(?:4r(?:85el8fhu5dnra|s40l)|gb(?:h1c|l6a))|x(?:hq521b|kc2(?:al3hye2a|dl3a5ee0h))|y(?:9a3aq|fro4i67o|gbi2ammx)|zfr164b)|peria|xx|yz)|y(?:a(?:chts|hoo|maxun|ndex)|o(?:dobashi|ga|kohama|u(?:tube)?)|un|[et])|z(?:a(?:ppos|ra)|ero|ip(?:po)?|one|uerich|[amw]))))")

//export simpleAnalyze
func simpleAnalyze(dir string) ([]string) {
	
	cmd := exec.Command("strings", "-n", "11", path.Join(dir, "classes.dex"))

	out, _ := cmd.Output()
	//if err != nil {
		//return out
	//}
	hostre.Longest()
	matches := hostre.FindAllSubmatch(out, -1)

	uncleanUrls := make([]string, 0, len(matches))
	for _, v := range matches {
		// badHosts is defined in data.go
		s := string(v[1])
		if _, ok := badHosts[s]; !ok {
			uncleanUrls = append(uncleanUrls, s)
		}
	}

	urls := util.Dedup(uncleanUrls)

	return urls
}

//export checkReflect
func checkReflect(dir string) bool {
	cmd := exec.Command("grep", "-Paqh",
		"\\x00\\x00\\x00.Ljava/lang/reflect[/a-zA-Z]*;\\x00\\x00\\x00",
		"--", path.Join(dir, "classes.dex"))

	out, err := cmd.Output()
	if err != nil && strings.TrimSpace(string(out)) != "" {
		fmt.Printf("Error checking for reflection: output below\n%s\n\n", string(out))
		return false
	}

	return (err == nil)
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
				paths[path.Dir(fname)] = util.Unit{}
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

func main() {}
