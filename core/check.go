package wapgolyzer

import (
	"bytes"
	"log"

	"golang.org/x/net/html"
)

func (fgp *Fingerprints) checkImplies(techArray []Tech) []Tech {
	var technologies []Tech
	for _, tech := range techArray {
		if implies := fgp.Apps[tech.Name].Implies; implies != nil {
			fgp.matchImplies(&technologies, tech, implies)
		}
	}
	return technologies
}

func (fgp *Fingerprints) checkHeader(headersLower map[string]string) []Tech {
	var technologies []Tech
	for app, fingerprint := range fgp.Apps {
		fgp.matchMapString(&technologies, fingerprint.Headers, headersLower, app)
	}
	return technologies
}

func (fgp *Fingerprints) checkCookies(parsedCookies map[string]string) []Tech {
	var technologies []Tech
	for app, fingerprint := range fgp.Apps {
		fgp.matchMapString(&technologies, fingerprint.Cookies, parsedCookies, app)
	}
	return technologies
}

func (fgp *Fingerprints) checkBody(body []byte) []Tech {
	var technologies []Tech
	//////////////////////////////////////////////
	// Check HTML
	log.Println("HTML processing")
	stringBody := string(body) //unsafe convert byte slice to string
	technologies = append(technologies, fgp.checkHTML(stringBody)...)
	log.Println("HTML processed")
	//////////////////////////////////////////////

	//////////////////////////////////////////////
	// Tokenize the HTML body to check
	log.Println("Token HTML processing")
	tokenizer := html.NewTokenizer(bytes.NewReader(body))
	technologies = append(technologies, fgp.checkTags(*tokenizer)...)
	log.Println("Token HTML processed")
	//////////////////////////////////////////////

	return technologies
}

func (fgp *Fingerprints) checkHTML(body string) []Tech {
	var technologies []Tech
	for app, fingerprint := range fgp.Apps {
		if fingerprint.HTML != nil {
			fgp.matchInterfaceString(&technologies, fingerprint.HTML, body, app)
		}
	}

	return technologies
}

func (fgp *Fingerprints) checkTags(tokenizer html.Tokenizer) []Tech {
	var technologies []Tech
	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			return technologies
		case html.StartTagToken:
			token := tokenizer.Token()
			switch token.Data {
			case "script":
				technologies = append(technologies, fgp.checkScriptSrc(token)...)
				// check content of javascript
				if tokenType := tokenizer.Next(); tokenType == html.TextToken {
					if data := tokenizer.Token().Data; data != "" {
						// check JS is future work
						// technologies = append(technologies, fgp.checkJS(data)...)
						continue
					}
				}
			case "meta":
				technologies = append(technologies, fgp.checkMeta(token)...)
			}
		case html.SelfClosingTagToken:
			token := tokenizer.Token()
			switch token.Data {
			case "meta":
				technologies = append(technologies, fgp.checkMeta(token)...)
			case "style":
				// check CSS is future work
				//technologies = append(technologies, fgp.checkCSS(data)...)
				continue
			}

		}
	}
}

func (fgp *Fingerprints) checkScriptSrc(token html.Token) []Tech {
	var technologies []Tech
	source, found := getScriptSource(token)
	if found {
		for app, fingerprint := range fgp.Apps {
			if fingerprint.Script != nil && source != "" {
				fgp.matchInterfaceString(&technologies, fingerprint.Script, source, app)
			}
		}
	}
	return technologies
}

func (fgp *Fingerprints) checkMeta(token html.Token) []Tech {
	var technologies []Tech
	name, content, found := getMetaNameAndContent(token)
	if found {
		for app, fingerprint := range fgp.Apps {
			fgp.matchMeta(&technologies, fingerprint.Meta, name, content, app)
		}
	}
	return technologies
}

/////
// Check Javascript and CSS content is future work

// func (fgp *Fingerprints) checkJS(data string) []Tech {
// 	var technologies []Tech
// 	//fmt.Println(data)
// 	return technologies
// }

// func (fgp *Fingerprints) checkCSS(data string) []Tech {
// 	var technologies []Tech
// 	//fmt.Println(data)
// 	return technologies
// }
