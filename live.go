package main

import (
	"flag"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	dirs := flag.Args()
	if len(dirs) == 0 {
		dirs = []string{"."}
	}
	var log *zap.Logger
	if *debug {
		log, _ = zap.NewDevelopment()
	} else {
		log, _ = zap.NewProduction()
	}

	wsserver := NewWSServer()
	wsserver.Log = log
	http.HandleFunc("/ws", wsserver.Handler())
	http.Handle("/", http.FileServer(http.Dir(".")))
	go func() {
		watchDirs(dirs,
			func(path string) {
				wsserver.Notify(path)
			})
	}()

	addr := "localhost:8080"
	log.Debug("Watching", zap.Strings("dirs", dirs))
	log.Debug("Serving on http://" + addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Error("Server start", zap.Error(err))
	}
}
