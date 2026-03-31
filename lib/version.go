package lib

import "fmt"

var (
	// Logo
	Logo = `
ReadCLI V%s
https://github.com/lvshp/glance
`
)

func DisplayVersion(version string) {
	fmt.Printf(Logo, version)
}
