package main

import (
	"fmt"
	"log"
	"net/http"
)

type MyResponseWriter struct {
	http.ResponseWriter
	code int
}

func (mw *MyResponseWriter) WriteHeader(code int) {
	mw.code = code
	mw.ResponseWriter.WriteHeader(code)
}

func RunSomeCode(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got a %s request for: %v", r.Method, r.URL)
		myrw := &MyResponseWriter{ResponseWriter: w, code: -1}
		handler.ServeHTTP(myrw, r)
		log.Println("Response status: ", myrw.code)
	})
}

func BoardHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {

	case http.MethodGet:

		if err := req.ParseForm(); err != nil {
			fmt.Println("Something went bad")
			fmt.Fprintln(w, "Something went bad")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Work in progress ...")
		w.WriteHeader(http.StatusOK)
	}
}

func BoatsHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {

	case http.MethodGet:

		if err := req.ParseForm(); err != nil {
			fmt.Println("Something went bad")
			fmt.Fprintln(w, "Something went bad")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Work in progress ...")
		w.WriteHeader(http.StatusOK)
	}
}

func HitHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {

	case http.MethodPost:

		if err := req.ParseForm(); err != nil {
			fmt.Println("Something went bad")
			fmt.Fprintln(w, "Something went bad")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Work in progress ...")
		w.WriteHeader(http.StatusOK)
	}
}
