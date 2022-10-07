package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// миддлвейр логгер, пишет данные поступившие из запроса
func LoggerMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer next.ServeHTTP(w, r)
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("Request %s %s from %q\n", r.Method, r.RequestURI, r.RemoteAddr))
		sb.WriteString("Header:\n")
		for k, v := range r.Header {
			sb.WriteString(fmt.Sprintf("\t%s: %q\n", k, v))
		}
		sb.WriteString("FormValues:\n")
		if err := r.ParseForm(); err != nil {
			sb.WriteString(fmt.Sprintf("could not parse form: %s", err))
			log.Println("\t", sb.String())
			return
		}
		for k, v := range r.Form {
			sb.WriteString(fmt.Sprintf("%s: %q\n", k, v))
		}
		log.Println(sb.String())
		return
	})
}
