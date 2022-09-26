package wapgolyzer

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func isMatch(value string, pattern string) (bool, string) {
	var re *regexp.Regexp
	var matches [][]string
	var splitted []string
	index := -1
	cleanPattern(&pattern)
	if strings.Contains(pattern, "\\;version:") {
		splitted = strings.Split(pattern, "\\;")
		re, _ = regexp.Compile(splitted[0])
		for _, part := range splitted {
			if strings.HasPrefix(part, "version:") {
				i := strings.TrimPrefix(part, "version:")
				index, _ = strconv.Atoi(string(i[len(i)-1]))
			}
		}
	} else {
		re, _ = regexp.Compile(pattern)
	}
	matches = re.FindAllStringSubmatch(value, -1)
	if len(matches) > 0 {
		if index > -1 {
			return true, matches[0][index]
		}
		return true, ""
	}
	return false, ""
}

func (fgp *Fingerprints) checkMatch(technologies *[]Tech, app string, str string, pattern string) {
	match, version := isMatch(str, pattern)
	if match {
		fgp.handleMatch(technologies, app, version)
	}
}

func (fgp *Fingerprints) handleMatch(technologies *[]Tech, app string, version string) {
	cats := mapNameCategories(fgp, app)
	tech := Tech{
		Name:    app,
		Version: version,
		Type:    cats,
	}
	*technologies = append(*technologies, tech)
}

func (fgp *Fingerprints) matchInterfaceString(technologies *[]Tech, i interface{}, str string, app string) {
	// Handle interface{} with reflect
	//////////////////////////////////////////////////////
	value := reflect.ValueOf(i)
	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice: //case of slice
		for i := 0; i < value.Len(); i++ {
			pattern := value.Index(i).Interface().(string)
			fgp.checkMatch(technologies, app, str, pattern)
		}
	case reflect.String: // case of string
		pattern := value.Interface().(string)
		fgp.checkMatch(technologies, app, str, pattern)
	}
	//////////////////////////////////////////////////////
}

func (fgp *Fingerprints) matchMapString(technologies *[]Tech, m map[string]string, s map[string]string, app string) {
	for key, pattern := range m {
		value, ok := s[key]
		if ok {
			if pattern == "" {
				fgp.handleMatch(technologies, app, "")
			}
			fgp.checkMatch(technologies, app, value, pattern)
		}
	}
}

func (fgp *Fingerprints) matchMeta(technologies *[]Tech, m map[string]interface{}, name string, content string, app string) {
	for k, v := range m {
		if name == k {
			value := reflect.ValueOf(v)
			kind := reflect.TypeOf(v).Kind()
			switch kind {
			case reflect.Slice: //case of slice
				for i := 0; i < value.Len(); i++ {
					pattern := value.Index(i).Interface().(string)
					if pattern == "" {
						fgp.handleMatch(technologies, app, "")
					}
					fgp.checkMatch(technologies, app, content, pattern)
				}
			case reflect.String: // case of string
				pattern := value.Interface().(string)
				if pattern == "" {
					fgp.handleMatch(technologies, app, "")
				}
				fgp.checkMatch(technologies, app, content, pattern)
			}
		}
	}
}

func (fgp *Fingerprints) matchImplies(technologies *[]Tech, tech Tech, implies interface{}) {
	value := reflect.ValueOf(implies)
	kind := reflect.TypeOf(implies).Kind()
	switch kind {
	case reflect.Slice: //case of slice
		for i := 0; i < value.Len(); i++ {
			app := value.Index(i).Interface().(string)
			fgp.handleMatch(technologies, app, "")
		}
	case reflect.String: // case of string
		app := value.Interface().(string)
		fgp.handleMatch(technologies, app, "")
	}
}
