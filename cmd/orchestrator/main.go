package main

import (
	"github.com/spf13/cobra"
	"math/rand"
	"simple2gis/intenal/orchestrator"
	"time"
)

//Да, использование кобры тут излишне, но, надо же показать, что умею ей пользоваться (вайпер совесть не позволила запихать)
var rootCmd = &cobra.Command{
	Use:   "s2gis",
	Short: "Root command",
	Long: "This command is the main one, that is, it is entered as an entry point to the CLI application, for example, " +
		"as the main command in Git (git merge..., git pull..., etc)",
	Run: orchestrator.RunOrchestrator,
}

var c = make(chan int)

func init() {
	rootCmd.PersistentFlags().String("http", "127.0.0.1:9999", "HTTP Server addr")
	rootCmd.PersistentFlags().String("db", "guidebook", "DB name")
	rand.Seed(time.Now().UnixNano())
	rootCmd.Execute()
}

func main() {

}