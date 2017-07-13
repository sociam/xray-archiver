package main

import (
	"errors"
	"fmt"
	"os/exec"
	"path"
)

type Unit struct{}

var unit Unit

func apkPath(apk, ver string) string {
	return path.Join(cfg.ApkDir, fmt.Sprintf("%s-%s.apk", apk, ver))
}

func outDir(apk, ver string) string {
	return path.Join(cfg.UnpackDir, path.Join(apk, ver))
}

func unpack(apk, ver string) error {
	cmd := exec.Command("apktool", "d", apkPath(apk, ver), "-o", outDir(apk, ver))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("Error '%s' unpacking apk; output below:\n%s",
			err.Error(), string(out)))
	}
	return nil
}

func combine(a, b map[string]Unit) map[string]Unit {
	ret := a
	for e, _ := range b {
		ret[e] = unit
	}
	return ret
}

func strmap(args ...string) map[string]Unit {
	ret := make(map[string]Unit)
	for _, e := range args {
		ret[e] = unit
	}

	return ret
}
