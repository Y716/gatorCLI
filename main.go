package main

import (
	"fmt"

	"github.com/Y716/gatorcli/gatorcli/internal/config"
)

func main() {
	cfgFile, err := config.Read()
	if err != nil {
		fmt.Printf("Error found: %v\n", err)
	}

	cfgFile.SetUser("Yasin")

	newCfgFile, err := config.Read()
	if err != nil {
		fmt.Printf("Error found: %v\n", err)
	}
	fmt.Println(newCfgFile.DbURL)
	fmt.Println(newCfgFile.CurrentUserName)
}
