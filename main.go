package main

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/crewjam/saml/samlsp"
	"github.com/kelseyhightower/envconfig"
)

type Specification struct {
	ENTITYID       string `default:"test"`
	IDPMETADATAURL string `required:"true"`
	SIGNREQUEST    bool   `default:"true"`
}

func hello(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Success!")
	fmt.Fprintf(w, "Welcome, %s!", samlsp.AttributeFromContext(r.Context(), "email"))
}

func main() {
	var s Specification
	err := envconfig.Process("APP", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	keyPair, err := tls.LoadX509KeyPair("myservice.cert", "myservice.key")
	if err != nil {
		panic(err)
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		panic(err)
	}

	idpMetadataURL, err := url.Parse(s.IDPMETADATAURL)
	if err != nil {
		panic(err)
	}
	idpMetadata, err := samlsp.FetchMetadata(context.Background(), http.DefaultClient,
		*idpMetadataURL)
	if err != nil {
		panic(err)
	}

	rootURL, err := url.Parse("http://localhost:8000")
	if err != nil {
		panic(err)
	}

	samlSP, _ := samlsp.New(samlsp.Options{
		URL:         *rootURL,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		IDPMetadata: idpMetadata,
		EntityID:    s.ENTITYID,
		SignRequest: s.SIGNREQUEST,
	})
	app := http.HandlerFunc(hello)
	http.Handle("/hello", samlSP.RequireAccount(app))
	http.Handle("/saml/", samlSP)
	fmt.Println("App is starting, link: http://localhost:8000/hello")
	http.ListenAndServe(":8000", nil)
}
