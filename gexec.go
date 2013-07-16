// Copyright 2013 GWoo. All rights reserved.
// The BSD License http://opensource.org/licenses/bsd-license.php.
package main

import (
	"fmt"
	"log"
	"strings"
	"encoding/base64"
	"net/http"
	"os"
	"os/exec"
	"flag"
)

var port = flag.Int("port", 8080, "Port for the server.")
var host = flag.String("host", "localhost", "Host for the server.")
var username = flag.String("username", "demo", "Username for basic auth.")
var password = flag.String("password", "test", "Password for basic auth.")

func Handler(w http.ResponseWriter, r *http.Request) {
	src := r.URL.Path[len("/"):]
	command, err := Save(src)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	output, err := Exec(command)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "%s", output)
}

func Save(src string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		log.Printf("Cannot decode %q: %s", src, err)
		return "", err
	}
	file, err := os.Create("/tmp/" + src)
	if err != nil {
		log.Printf("Cannot save %q: %s", src, err)
		return "", err
	}
	r := strings.NewReplacer("\r", "")
	file.WriteString(r.Replace(string(decoded)))
	file.Chmod(0777)
	file.Close()
	return string(src), nil
}

func Exec(command string) (string, error) {
	c := exec.Command("/tmp/" + command)
	output, err := c.CombinedOutput()
	if err != nil {
		log.Printf("Cannot run %s", err)
		return "", err
	}
	return string(output), nil
}

func EncodeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,
		"<html><title>Encode a Script</title>"+
		"<body style='font-family: Monospace;'>")

	body := r.FormValue("body")
	if body != "" {
		fmt.Fprintf(w,
			"<h2>Done. <a href='%s'>Click Here</a></h2>",
			base64.URLEncoding.EncodeToString([]byte(body)))
	}
	fmt.Fprintf(w,
		"<h2>Encode</h2>"+
		"<form action=\"/encode\" method=\"POST\">"+
		"<p><textarea name=\"body\" style='height:200px;width:400px'>%s</textarea></p>"+
		"<p><input type=\"submit\" value=\"Submit\"></p>"+
		"</form>", body)
	fmt.Fprintf(w, "</body></html>")
}

func AuthHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL
		for k, v := range r.Header {
			fmt.Printf("  %s = %s\n", k, v[0])
		}
		auth, ok := r.Header["Authorization"]
		if !ok {
			log.Printf("Unauthorized access to %s", url)
			w.Header().Add("WWW-Authenticate", fmt.Sprintf("basic realm=\"%s\"", *host))
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Not Authorized.")
			return
		}
		encoded := strings.Split(auth[0], " ")
		if len(encoded) != 2 || encoded[0] != "Basic" {
			log.Printf("Strange Authorization %q", auth)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		decoded, err := base64.StdEncoding.DecodeString(encoded[1])
		if err != nil {
			log.Printf("Cannot decode %q: %s", auth, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		parts := strings.Split(string(decoded), ":")
		if len(parts) != 2 {
			log.Printf("Unknown format for credentials %q", decoded)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if parts[0] == *username && parts[1] == *password {
			fn(w, r)
			return
		}
		log.Printf("Unauthorized access to %s", url)
		w.Header().Add("WWW-Authenticate", fmt.Sprintf("basic realm=\"%s\"", *host))
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Not Authorized.")
		return
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/favicon.ico", http.NotFound)
	http.HandleFunc("/", AuthHandler(Handler))
	http.HandleFunc("/encode", AuthHandler(EncodeHandler))
	log.Printf("Connected to %s:%d", *host, *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}