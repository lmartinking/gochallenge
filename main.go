package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hey!")
}

func init() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)
}

func run() error {
	var (
		port          = flag.Int("port", 8000, "Port to listen on")
		pub_key_path  = flag.String("pubkey", "etc/key.pub", "Public key")
		priv_key_path = flag.String("privkey", "etc/key.priv", "Private key")
	)
	flag.Parse()

	pub_key, err := LoadPubKey(*pub_key_path)
	if err != nil {
		return err
	}

	priv_key, err := LoadPrivKey(*priv_key_path)
	if err != nil {
		return err
	}

	Cfg.PrivKey = priv_key
	Cfg.PubKey = pub_key

	AddAccount("user@example.com", "badpassword")

	log.Info("Listening on port:", *port)

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/login", HandleLogin)
	http.HandleFunc("/refresh", HandleRefresh)
	http.HandleFunc("/protected", ProtectedWrapper(HandleProtected))

	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Fatal("Error running:", err)
		os.Exit(1)
	}
}
