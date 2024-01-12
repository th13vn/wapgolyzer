package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	wapgolyzer "github.com/th13vn/wapgolyzer/core"
)

func main() {
	/////////////////////////
	// Set up fingerprints
	fgpFile := wapgolyzer.SetupFingerprintsFile()
	/////////////////////////

	var input, filename string
	var u, d bool
	var o string
	//var f string
	flag.StringVar(&o, "o", "", "Output file with json format")
	flag.BoolVar(&d, "d", false, "Default json file that hostname is file name")
	flag.BoolVar(&u, "u", false, "Update fingerprints from Wappalyzer Github")
	//flag.StringVar(&f, "f", "", "Specific your fingerprints file")
	flag.Parse()
	/////////////////////////
	lasts := flag.Args()
	args := os.Args
	if len(args) < 2 {
		commandUsage()
		return
	}
	if len(lasts) < 1 {
		commandUsage()
		return
	}
	input = lasts[0]
	// if f != "" {
	// 	fgpFile = f
	// }
	if u {
		wapgolyzer.UpdateFingerprintsFile(fgpFile)
	}
	if o != "" {
		tmp := strings.Replace(o, "/", "", -1)
		filename = strings.Replace(tmp, "..", "", -1)
	}
	if d {
		filename = genFilename(input)
	}
	/////////////////////////
	arrayByte, err := wapgolyzer.Run(input, fgpFile)
	if err != nil {
		fmt.Println(err)
		commandUsage()
		return
	}
	fmt.Println(filename)
	if filename != ".json" && filename != "" {
		ioutil.WriteFile(filename, arrayByte, 0644)
		return
	}
	fmt.Println(string(arrayByte))
}

func commandUsage() {
	fmt.Println("Usage: wapgolyzer [options] <url>")
	flag.PrintDefaults()
}

func genFilename(input string) string {
	u, err := url.Parse(input)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return u.Hostname() + ".json"
}
