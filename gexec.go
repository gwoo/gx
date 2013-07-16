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
	a := AuthHandler(w,r)
	if !a {
		fmt.Fprintf(w, "%s", "Not Authorized.")
		return
	}
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

func AuthHandler(w http.ResponseWriter, r *http.Request) bool {
	url := r.URL
	for k, v := range r.Header {
		fmt.Printf("  %s = %s\n", k, v[0])
	}
	auth, ok := r.Header["Authorization"]
	if !ok {
		w.Header().Add("WWW-Authenticate", fmt.Sprintf("basic realm=\"%s\"", *host))
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("Unauthorized access to %s", url)
		return false
	}
	encoded := strings.Split(auth[0], " ")
	if len(encoded) != 2 || encoded[0] != "Basic" {
		log.Printf("Strange Authorizatoion %q", auth)
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	decoded, err := base64.StdEncoding.DecodeString(encoded[1])
	if err != nil {
		log.Printf("Cannot decode %q: %s", auth, err)
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	parts := strings.Split(string(decoded), ":")
	if len(parts) != 2 {
		log.Printf("Unknown format for credentials %q", decoded)
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if parts[0] == *username && parts[1] == *password {
		return true
	}
	w.Header().Add("WWW-Authenticate", fmt.Sprintf("basic realm=\"%s\"", *host))
	w.WriteHeader(http.StatusUnauthorized)
	log.Printf("Unauthorized access to %s", url)
	return false
}

func main() {
	flag.Parse()
	http.HandleFunc("/", Handler)
	log.Printf("Connected to %s:%d", *host, *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}