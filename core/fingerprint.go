package wapgolyzer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Fingerprints struct {
	// Apps is organized as <name, fingerprint>
	Apps map[string]Signature  `json:"technologies"`
	Cats map[string]Categories `json:"categories"`
}

// Fingerprint is a single piece of information about a tech
type Signature struct {
	Cats    []int                  `json:"cats"`
	HTML    interface{}            `json:"html"`
	CSS     interface{}            `json:"css"`
	Script  interface{}            `json:"scriptSrc"`
	Implies interface{}            `json:"implies"`
	Headers map[string]string      `json:"headers"`
	Cookies map[string]string      `json:"cookies"`
	JS      map[string]string      `json:"js"`
	Meta    map[string]interface{} `json:"meta"`
}

const (
	technologyURL = "https://raw.githubusercontent.com/AliasIO/wappalyzer/master/src/technologies/%s.json"
	categoriesURL = "https://raw.githubusercontent.com/wappalyzer/wappalyzer/master/src/categories.json"
)

func makeTechFileURLs() []string {
	files := []string{"_", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

	var urls []string
	for _, item := range files {
		urls = append(urls, fmt.Sprintf(technologyURL, item))
	}
	return urls
}

func (fgp *Fingerprints) loadFingerprintsFile(FgpFile string) error {
	// Load technologies to fgp variables
	////////////////////////////////////////
	_, err := os.Stat(FgpFile)
	if err == nil {
		log.Printf("Loading " + FgpFile + "\n")
		b, _ := ioutil.ReadFile(FgpFile)
		_ = json.Unmarshal(b, &fgp)
		log.Printf("Loaded " + FgpFile + "\n")
		return nil
	}
	return err
	////////////////////////////////////////
}

func UpdateFingerprintsFile(fgpFile string) {
	fgp := &Fingerprints{
		Apps: make(map[string]Signature),
		Cats: make(map[string]Categories),
	}
	// Start crawl
	//////////////////////////////////////////////////////////////////////////
	log.Printf("Start crawl fingerprints from Wappalyzer\n")
	// technologies
	technologyURLs := makeTechFileURLs()
	for _, technologyURL := range technologyURLs {
		if err := fgp.gatherFingerprintsFromURL(technologyURL); err != nil {
			log.Fatalf("Could not gather technology file %s: %v\n", technologyURL, err)
		}
	}
	// categories
	if err := fgp.gatherFingerprintsFromURL(categoriesURL); err != nil {
		log.Fatalf("Could not gather fingerprints %s: %v\n", categoriesURL, err)
	}
	log.Printf("Read fingerprints from the server\n")
	log.Printf("Starting normalizing of %d fingerprints...\n", len(fgp.Apps))
	//////////////////////////////////////////////////////////////////////////

	// Write file
	////////////////////////////////////////////////////
	txt, _ := json.MarshalIndent(fgp, "", " ")
	_ = ioutil.WriteFile(fgpFile, txt, 0644)
	log.Printf("Done\n")
	////////////////////////////////////////////////////
}

func (fgp *Fingerprints) gatherFingerprintsFromURL(URL string) error {
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if strings.Contains(URL, "categories.json") {
		if err := mapCategories(fgp, data); err != nil {
			return err
		}
	} else {
		if err := mapApps(fgp, data); err != nil {
			return err
		}
	}

	return nil
}
