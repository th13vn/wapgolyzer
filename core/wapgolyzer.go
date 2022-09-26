package wapgolyzer

import (
	"bytes"
	"encoding/json"
	"log"
)

type Categories struct {
	Name     string `json:"name"`
	Priority int    `json:"priority"`
}

type Result struct {
	URL          string `json:"url"`
	Technologies []Tech `json:"technologies"`
}

type Tech struct {
	Name    string   `json:"name"`
	Version string   `json:"version"`
	Type    []string `json:"type"`
}

func Run(url string, fgpFile string) ([]byte, error) {
	///////////////////////////////////////
	// Check url
	if err := validateURL(url); err != nil {
		return nil, err
	}
	///////////////////////////////////////

	///////////////////////////////////////
	// Prepare for detect technologies
	var technologies []Tech
	result := &Result{
		URL:          url,
		Technologies: []Tech{},
	}
	fgp := &Fingerprints{
		Apps: make(map[string]Signature),
		Cats: make(map[string]Categories),
	}
	if err := fgp.loadFingerprintsFile(fgpFile); err != nil {
		return nil, err
	}
	// Get request to URL
	resp, body := getRequest(url)
	///////////////////////////////////////

	////////////////////////////////////////
	// Proccess
	// header
	log.Println("Header processing")
	headersLower := toLowerHeader(*resp)
	technologies = append(technologies, fgp.checkHeader(headersLower)...)
	log.Println("Header processed")
	// cookies
	log.Println("Cookies processing")
	cookies := getSetCookie(headersLower)
	if len(cookies) > 0 {
		parsedCookies := parseCookie(cookies)
		technologies = append(technologies, fgp.checkCookies(parsedCookies)...)
	}
	log.Println("Cookies processed")
	// body
	log.Println("Body processing")
	bodyLower := bytes.ToLower(body)
	technologies = append(technologies, fgp.checkBody(bodyLower)...)
	log.Println("Body processed")
	// implies
	log.Println("Implies processing")
	technologies = append(technologies, fgp.checkImplies(technologies)...)
	log.Println("Implies processed")
	////////////////////////////////////////

	/////////////////////////////////////////////////////
	// Return result as byte slice
	result.Technologies = appendTechUnique(result.Technologies, technologies)
	byteResult, _ := json.MarshalIndent(result, "", "  ")
	return byteResult, nil
	/////////////////////////////////////////////////////
}
