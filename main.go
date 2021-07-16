package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	f "github.com/fauna/faunadb-go/v4/faunadb"
)

const (
	secretPort = "1000"
	proxyPort  = "8085"
)

func main() {
	time.Sleep(15 * time.Second)

	client := f.NewFaunaClient("secret", f.Endpoint("http://localhost:8443"))

	var secret string

	if err := createDatabase(client, &secret); err != nil {
		log.Fatal(err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		provideSecret(&secret, client)
		wg.Done()
	}()

	go func() {
		proxy(&secret)
		wg.Done()
	}()

	wg.Wait()
}

func provideSecret(secret *string, client *f.FaunaClient) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			if err := createDatabase(client, secret); err != nil {
				log.Fatal(err)
			}
		}
		fmt.Fprint(rw, *secret)
	})

	log.Fatal(http.ListenAndServe(":"+secretPort, mux))
}

func proxy(secret *string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		client := &http.Client{}

		req, err := http.NewRequest(r.Method, "http://localhost:8084"+r.URL.Path, r.Body)
		if err != nil {
			log.Fatal(err)
		}
		for k, v := range r.Header {
			req.Header.Add(k, v[0])
		}
		req.Header.Set("Authorization", "Bearer "+*secret)

		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		for k, v := range res.Header {
			rw.Header().Set(k, v[0])
		}

		rw.WriteHeader(res.StatusCode)
		rw.Write(body)
	})

	log.Fatal(http.ListenAndServe(":"+proxyPort, mux))
}
