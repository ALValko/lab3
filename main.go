package main

import (
	"avp-3/app"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	logger.Fatalf("Сервер завершил работу с ошибкой: %s\n", app.Run(logger))
}
