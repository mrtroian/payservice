package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	*http.Server
}

func NewServer() Server {
	return Server{new(http.Server)}
}

func (srv *Server) SetAddr(host string, port int) {
	srv.Addr = fmt.Sprintf("%s:%d", host, port)
}

func (srv *Server) SetRouter(r http.Handler) {
	srv.Handler = r
}

func (srv *Server) SetTLS(cert, key string) error {
	certificate, err := tls.LoadX509KeyPair(cert, key)

	if err != nil {
		return err
	}

	srv.TLSConfig = &tls.Config{
		Certificates:       []tls.Certificate{certificate},
		InsecureSkipVerify: true,
	}
	return nil
}

func (srv *Server) Start() {
	go func() {
		log.Println("server: starting on", srv.Addr)
		if err := srv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Println("server:", err)
			os.Exit(1)
		}
	}()
}

func (srv *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("server: forced shutdown:", err)
	} else {
		log.Println("server: shut down correctly")
	}
}
