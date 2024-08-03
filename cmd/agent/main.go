package main

import (
	"context"
	"flag"
	"log"

	"github.com/AndreyVLZ/curly-octo/agent"
)

func main() {
	addr := flag.String("a", agent.AddresDefault, "адрес подключения")
	tmpDir := flag.String("d", agent.TmpDirDefault, "куда сохранять файла")
	secretKey := flag.String("k", "", "ключ")
	chunkSize := flag.Int("c", agent.ChunkBufSizeDefault, "размер буфера для шифрования")
	bufSize := flag.Int("b", agent.ChunkBufSizeDefault, "размер буфера для чтения/записи файла")
	flag.Parse()

	ctx := context.Background()

	cfg := agent.NewConfig(
		*addr, *tmpDir, *secretKey,
		*chunkSize, *bufSize,
	)

	agent, err := agent.New(cfg)
	if err != nil {
		log.Printf("new agent: %v\n", err)

		return
	}

	if err := agent.Start(ctx); err != nil {
		log.Printf("agent start: %v\n", err)

		return
	}

	if err := agent.Stop(ctx); err != nil {
		log.Printf("agent stop: %v\n", err)

		return
	}
}
