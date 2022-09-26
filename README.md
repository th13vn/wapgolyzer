# WAPGOLYZER

Wapgolyzer is tool that able to detect web techonologies by given URL

## How it works

Wapgolyzer use Wappalyzer's set of rules to detect. It sends a GET request to the given URL and uses HTTP Header, HTTP Response content to analyses.

## Installation

```bash
> go install github.com/th13ntc/wapgolyzer/cmd/wapgolyzer
```

## Usage

```bash
Usage: wapgolyzer [options] <url>
  -d    Default json file that hostname is file name
  -o string
        Output file with json format
  -u    Update fingerprints from Wappalyzer Github
```

## Example

```bash
> wapgolyzer https://radzad.com
```

Ouput:

```
{
  "url": "https://radzad.com",
  "technologies": [
    {
      "name": "litespeed",
      "version": "",
      "type": [
        "web servers"
      ]
    },
    {
      "name": "http/3",
      "version": "",
      "type": [
        "miscellaneous"
      ]
    },
    {
      "name": "php",
      "version": "7.4.22",
      "type": [
        "programming languages"
      ]
    },
    {
      "name": "wordpress",
      "version": "",
      "type": [
        "cms",
        "blogs"
      ]
    },
    {
      "name": "bootstrap",
      "version": "",
      "type": [
        "ui frameworks"
      ]
    },
    {
      "name": "elementor",
      "version": "3.4.8",
      "type": [
        "page builders"
      ]
    },
    {
      "name": "wp rocket",
      "version": "",
      "type": [
        "caching",
        "wordpress plugins"
      ]
    },
    {
      "name": "jquery",
      "version": "",
      "type": [
        "javascript libraries"
      ]
    },
    {
      "name": "extendify",
      "version": "4.3.9",
      "type": [
        "wordpress plugins"
      ]
    },
    {
      "name": "contact form 7",
      "version": "5.5.6",
      "type": [
        "wordpress plugins"
      ]
    },
    {
      "name": "woocommerce",
      "version": "",
      "type": [
        "ecommerce",
        "wordpress plugins"
      ]
    },
    {
      "name": "swiper",
      "version": "",
      "type": [
        "javascript libraries"
      ]
    },
    {
      "name": "recaptcha",
      "version": "",
      "type": [
        "security"
      ]
    },
    {
      "name": "wordfence login security",
      "version": "1.0.9",
      "type": [
        "wordpress plugins",
        "security"
      ]
    },
    {
      "name": "mysql",
      "version": "",
      "type": [
        "databases"
      ]
    }
  ]
}
```

## Implement code

```go
package main

import (
	"fmt"

	wapgolyzer "github.com/th13ntc/wapgolyzer/core"
)

func main() {
	input := "https://example.com"
	fgpFile := wapgolyzer.SetupFingerprintsFile() // $HOME/.config/wapgolyzer/fingerprints.json
	wapgolyzer.UpdateFingerprintsFile(fgpFile)
	arrayByte, err := wapgolyzer.Run(input, fgpFile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(arrayByte))
}
```

## Future work

- Read js and css content to analyse

## Changelog
