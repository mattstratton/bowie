// A pretty changelog generator
package main

import (
	"fmt"
	"os"

	"github.com/mattstratton/bowie/bowielib"
)

func main() {
	if err := bowielib.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
