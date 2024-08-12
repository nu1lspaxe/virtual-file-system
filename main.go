package main

import (
	"fmt"
	"system/pkg"
)

func main() {
	path, err := pkg.GetManPath()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	pkg.SetManInfo(path, "./vfs.1")

	fmt.Println("Running the main CLI application...")
}
