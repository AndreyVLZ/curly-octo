package main

import (
	"flag"
	"log"

	"github.com/AndreyVLZ/curly-octo/server"
)

func main() {
	addr := flag.String("a", server.AddresDefault, "адрес подключения")
	tmpDir := flag.String("d", server.TmpDirDefault, "куда сохранять файла")
	key := flag.String("k", "", "ключ")
	exp := flag.Duration("e", server.ExpiresAtDefault, "время жизни токена")
	bufSize := flag.Int("b", server.SendBufSizeDefault, "размер буфера для чтения/записи файла")
	flag.Parse()

	cfg := server.NewConfig(*addr, *tmpDir, *key, *exp, *bufSize)

	server := server.New(cfg)

	if err := server.Start(); err != nil {
		log.Printf("server start: %v\n", err)
	}

	server.Stop()
}
