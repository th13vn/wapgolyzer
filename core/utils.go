package wapgolyzer

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

const (
	DIR = ".config/wapgolyzer"
)

func mapCategories(fgp *Fingerprints, data []byte) error {
	fgpTmp := Fingerprints{}
	err := json.NewDecoder(bytes.NewReader(bytes.ToLower(data))).Decode(&fgpTmp.Cats)
	if err != nil {
		return err
	}
	for k, v := range fgpTmp.Cats {
		fgp.Cats[k] = v
	}
	return nil
}

func mapApps(fgp *Fingerprints, data []byte) error {
	fgpTmp := Fingerprints{}
	err := json.NewDecoder(bytes.NewReader(bytes.ToLower(data))).Decode(&fgpTmp.Apps)
	if err != nil {
		return err
	}
	for k, v := range fgpTmp.Apps {
		fgp.Apps[k] = v
	}
	return nil
}

func mapNameCategories(fgp *Fingerprints, app string) []string {
	var categories []string
	for _, i := range fgp.Apps[app].Cats {
		categories = append(categories, fgp.Cats[strconv.Itoa(i)].Name)
	}
	return categories
}

func getRequest(url string) (*http.Response, []byte) {
	// Http request to URL
	////////////////////////////////////////
	log.Printf("Send HTTP Request to URL\n")
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		log.Println(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Recived response\n")
	////////////////////////////////////////
	return resp, body
}

func getScriptSource(token html.Token) (string, bool) {
	if len(token.Attr) < 1 {
		return "", false
	}

	var source string
	for _, attr := range token.Attr {
		switch attr.Key {
		case "src":
			source = attr.Val
		}
	}
	if source == "" {
		return "", false
	}
	return source, true
}

func getMetaNameAndContent(token html.Token) (string, string, bool) {
	if len(token.Attr) < 2 {
		return "", "", false
	}

	var name, content string
	for _, attr := range token.Attr {
		switch attr.Key {
		case "name":
			name = attr.Val
		case "content":
			content = attr.Val
		}
	}
	return name, content, true
}

func getSetCookie(headers map[string]string) []string {
	value, ok := headers["set-cookie"]
	if !ok {
		return nil
	}

	var values []string
	for _, v := range strings.Split(value, " ") {
		if v == "" {
			continue
		}
		if strings.Contains(v, ",") {
			values = append(values, strings.Split(v, ",")...)
		} else if strings.Contains(v, ";") {
			values = append(values, strings.Split(v, ";")...)
		} else {
			values = append(values, v)
		}
	}
	return values
}

func parseCookie(cookies []string) map[string]string {
	parsed := make(map[string]string)
	for _, part := range cookies {
		parts := strings.Split(strings.Trim(part, " "), "=")
		if len(parts) < 2 {
			continue
		}
		parsed[parts[0]] = parts[1]
	}
	return parsed
}

func toLowerHeader(resp http.Response) map[string]string {
	headersArray := resp.Header
	headers := make(map[string]string, len(headersArray))
	headersLower := make(map[string]string, len(headersArray))
	builder := &strings.Builder{}
	for key, value := range headersArray {
		for i, v := range value {
			builder.WriteString(v)
			if i != len(value)-1 {
				builder.WriteString(", ")
			}
		}
		headerValue := builder.String()
		headers[key] = headerValue
		builder.Reset()
	}
	for header, value := range headers {
		headersLower[strings.ToLower(header)] = strings.ToLower(value)
	}
	return headersLower
}

func cleanPattern(pattern *string) {
	// Clean Pattern
	// Because error parsing regexp: invalid or unsupported Perl syntax: `(?!` (SA1000)
	for {
		if strings.Contains(*pattern, "(?!") {
			re, _ := regexp.Compile(`\(\?\!`)
			matches := re.FindIndex([]byte(*pattern))
			if len(matches) > 0 {
				*pattern = strings.Replace((*pattern), "?!", "", 1)
				brackets := 0
				var tmpString string
				for i := matches[0]; i < len(*pattern); i++ {
					tmpString = tmpString + string((*pattern)[i])
					if (*pattern)[i] == '(' {
						brackets++
					}
					if (*pattern)[i] == ')' {
						brackets--
						if brackets == 0 {
							*pattern = strings.Replace((*pattern), tmpString, tmpString+"?", 1)
							break
						}
					}
				}
			}
		} else {
			break
		}
	}
}

func appendTechUnique(first []Tech, technologies []Tech) []Tech {
	for _, tech := range technologies {
		exist := false
		for i, entity := range first {
			if entity.Name == tech.Name && (entity.Version == tech.Version || tech.Version == "") {
				exist = true
				break
			}
			if entity.Name == tech.Name && entity.Version == "" {
				exist = true
				first[i].Version = tech.Version
				break
			}
		}
		if !exist {
			first = append(first, tech)
		}
	}
	return first
}

func validateURL(input string) error {
	_, err := url.ParseRequestURI(input)
	if err != nil {
		return err
	}
	return nil
}

func SetupFingerprintsFile() string {
	absoluteDir := filepath.Join(os.Getenv("HOME"), DIR)
	if _, err := os.Stat(absoluteDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(absoluteDir, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
	return "./fingerprints.json"
}
