/*
Copyright © 2022 Mohamed Hammad Youssef mmhy2003@hotmail.com
*/
package main

import (
	"github.com/mmhy2003/underscoreai/cmd"
	"github.com/mmhy2003/underscoreai/config"
)

func main() {
	config.LoadConfig()
	cmd.Execute()
}
