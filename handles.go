package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, bundleIndex)
}

func handleAPIFiles(w http.ResponseWriter, r *http.Request) {
	f, err := fileWalk()
	if err != nil {
		log.Fatal(err)
	}
	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(b))
}

func handleAPIIP(w http.ResponseWriter, r *http.Request) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Fprintf(w, "%v\n", localAddr.IP)
}

// handleAPIExternalFiles takes a GET parameter "url" and attempts to
// return the /api/files/ from another NoUSB server.
// This is to combat a cross-origin request.
func handleAPIExternalFiles(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "err: %s", err)
		return
	}
	url := r.FormValue("url")
	if url == "" {
		fmt.Fprintf(w, "Could not retrieve url from form")
		return
	}
	b, err := getRequest(httpsStrip(url) + "/api/files/")
	if err != nil {
		fmt.Fprintf(w, "Could not retrieve files from external server")
		return
	}
	fmt.Fprintf(w, string(b))
}
