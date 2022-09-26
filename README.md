# WAPGOLYZER

Wapgolyzer is tool that able to detect web techonologies by given URL

## How it works

Wapgolyzer use Wappalyzer's set of rules to detect. It sends a GET request to the given URL and uses HTTP Header, HTTP Response content to analyses.

## Installation

```bash
> go install github.com/th13ntc/wapgolyzer/cmd/wapgolyzer@latest
```

## Usage

```bash
Usage: wapgolyzer [options] <url>
  -d    Default json file that hostname is file name
  -o string
        Output file with json format
  -u    Update fingerprints from Wappalyzer Github
```

> Note: You have to use `-u` option for the first time using to crawl fingerprints.json detect

## Example

```bash
> wapgolyzer https://example.com
```

Ouput:

```
{
  "url": "http://example.com",
  "technologies": [
    {
      "name": "azure cdn",
      "version": "",
      "type": [
        "cdn"
      ]
    },
    {
      "name": "amazon ecs",
      "version": "",
      "type": [
        "iaas"
      ]
    },
    {
      "name": "azure",
      "version": "",
      "type": [
        "paas"
      ]
    },
    {
      "name": "amazon web services",
      "version": "",
      "type": [
        "paas"
      ]
    },
    {
      "name": "docker",
      "version": "",
      "type": [
        "containers"
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
