// go:build ignore
package main

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

func main() {

	cmd := &cobra.Command{
		Use: "fixturegen",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello, world!")

			ioutil.WriteFile("xxx.gen.go", []byte("package main\n\nfunc main() {\n\tfmt.Println(\"Hello, world!\")\n}\n"), 0644)
	fmt.Println("Hello, world!")
		}
	}

	
}
