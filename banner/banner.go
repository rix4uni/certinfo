package banner

import (
	"fmt"
)

// prints the version message
const version = "v0.0.5"

func PrintVersion() {
	fmt.Printf("Current certinfo version %s\n", version)
}

// Prints the Colorful banner
func PrintBanner() {
	banner := `
                     __   _         ____     
  _____ ___   _____ / /_ (_)____   / __/____ 
 / ___// _ \ / ___// __// // __ \ / /_ / __ \
/ /__ /  __// /   / /_ / // / / // __// /_/ /
\___/ \___//_/    \__//_//_/ /_//_/   \____/ 
`
	fmt.Printf("%s\n%50s\n\n", banner, "Current certinfo version "+version)
}
