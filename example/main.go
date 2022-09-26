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
