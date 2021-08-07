import (
	"github.com/spf13/cobra"
	"log"
	"simple2gis/intenal/server"
	"simple2gis/intenal/sqlite"
	"sync"
	"syscall"
	"time"
)

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
