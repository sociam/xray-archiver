package main

import (
	"testing"
)

func TestHostRe(t *testing.T) {
	hostre.Longest()

	goodHosts := map[string]string{
		"google.com":        "google.com",
		"tripadvisor.com":   "tripadvisor.com",
		"wow.isaname.co.uk": "wow.isaname.co.uk",
		"test.foob.ar":      "test.foob.ar",
		"!jkfd.com":         "jkfd.com",
	}

	for str, host := range goodHosts {
		if m := hostre.FindStringSubmatch(str); m == nil || m[1] != host {
			t.Errorf("The regex failed to match host %s", host)
		}
	}

	badHosts := []string{
		"http://www",
		"@mfd.com",
		"jkfd@jkfdkd.com",
	}

	for _, str := range badHosts {
		if locn := hostre.FindStringIndex(str); locn != nil {
			t.Errorf("The regex matched bad host in %s: %s", str, str[locn[0]:locn[1]])
		}
	}
}
