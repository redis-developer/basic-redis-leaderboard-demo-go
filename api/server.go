package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
)

type Server struct {
	srv      *http.Server
	srvNoTLS *http.Server
}

func (server Server) Start(ctxCancel context.CancelFunc) {

	if server.srv.TLSConfig != nil {
		log.Println(fmt.Sprintf("start server tls listen at [%s]", server.srv.Addr))

		go func() {
			log.Println(fmt.Sprintf("start server tls listen at [%s]", server.srvNoTLS.Addr))
			err := server.srvNoTLS.ListenAndServe()
			if err != nil {
				log.Println(err)
			}
		}()

		if err := server.srv.ListenAndServeTLS("", ""); err != nil {
			if err == http.ErrServerClosed {
				log.Println("server closed under request")
			} else {
				log.Println("server closed unexpect:", err)
			}
		}

	} else {
		log.Println(fmt.Sprintf("start server listen at [%s]", server.srv.Addr))

		if err := server.srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Println("server closed under request")
			} else {
				log.Println("server closed unexpect:", err)
			}
		}
	}

	ctxCancel()

}

func (server Server) Close() error {
	err1 := server.srv.Close()
	if err1 != nil {
		log.Println(err1)
	}
	if server.srvNoTLS != nil {
		err2 := server.srvNoTLS.Close()
		if err2 != nil {
			log.Println(err2)
			return err2
		}
	}

	return err1

}

func NewServer(config Config) *Server {

	log.Println("initializing server")

	server := &Server{
		srv: &http.Server{
			Addr:    config.Addr(),
			Handler: router(config.PublicPath()),
		},
	}

	if config.TLSDisabled() == false {

		certManager := autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache("./certs"),
		}

		server.srv.Addr = ":443"
		server.srv.TLSConfig = &tls.Config{
			GetCertificate: certManager.GetCertificate,
		}

		server.srvNoTLS = &http.Server{
			Addr:    ":80",
			Handler: certManager.HTTPHandler(nil),
		}

	}

	return server

}
