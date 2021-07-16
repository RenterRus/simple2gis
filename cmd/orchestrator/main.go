package main

import (
	"github.com/spf13/cobra"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"simple2gis/intenal/server"
	"simple2gis/intenal/sqlite"
	"sync"
	"syscall"
	"time"
)

//Да, использование кобры тут излишне, но, надо же показать, что умею ей пользоваться (вайпер совесть не позволила запихать)
var rootCmd = &cobra.Command{
	Use:   "s2gis",
	Short: "Root command",
	Long: "This command is the main one, that is, it is entered as an entry point to the CLI application, for example, " +
		"as the main command in Git (git merge..., git pull..., etc)",
	Run: RunOrchestrator,
}

func init() {
	rootCmd.PersistentFlags().String("http", "127.0.0.1:9999", "HTTP Server addr")
	rootCmd.PersistentFlags().String("db", "guidebook", "DB name")
	rand.Seed(time.Now().UnixNano())
	rootCmd.Execute()
}

func main() {

}

// RunOrchestrator Инициализация HTTP сервера, установка соединения с БД и красивое завершение работы
func RunOrchestrator(cmd *cobra.Command, args []string) {
	var wg sync.WaitGroup
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	dbname, _ := cmd.Flags().GetString("db")

	var httpServer *server.HTTPServer

	sqlite.DBClient = sqlite.Initial(dbname)
	wg.Add(3)
	go func() {
		defer wg.Done()
		log.Println("HTTP Server starting")
		addr, _ := cmd.Flags().GetString("http")
		httpServer = server.NewServer(addr)
		httpServer.Start()
	}()
	go func() {
		defer wg.Done()
		go sqlite.DBProc()
	}()
	go func() {
		defer wg.Done()
		<-sign
		httpServer.GraceShutdown()
		log.Println("HTTP Shutdown")
		sqlite.DBClient.DisableConnect()
		log.Println("DB Connection is closed")
	}()

	wg.Wait()
}
