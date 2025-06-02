package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"rushkv/server"
)

func main() {
	var (
		nodeID   = flag.String("id", "node1", "Node ID")
		address  = flag.String("addr", "localhost", "Server address")
		port     = flag.Int("port", 8080, "Server port")
		dataPath = flag.String("data", "./data", "Data directory")
	)
	flag.Parse()

	// 创建服务器
	srv, err := server.NewRushKVServer(*nodeID, *address, *port, *dataPath)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// 处理优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down server...")
		srv.Stop()
		os.Exit(0)
	}()

	// 启动服务器
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
