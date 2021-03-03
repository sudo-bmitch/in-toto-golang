package in_toto

import (
	"context"
	"crypto/x509"
	"log"

	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

//GetSVID grabs the x.509 context.
func GetSVID(ctx context.Context, socketPath string) Key {

	var k Key

	//*x509.Certificate

	client, err := workloadapi.New(ctx, workloadapi.WithAddr(socketPath))
	if err != nil {
		log.Fatalf("Unable to create workload API client: %v", err)
	}
	defer client.Close()

	svidContext, err := client.FetchX509Context(ctx)
	if err != nil {
		log.Fatalf("Error grabbing x.509 context: %v", err)
	}

	log.Printf("using svid %v", svidContext.DefaultSVID().ID.String())

	svid, keyBytes, err := svidContext.DefaultSVID().Marshal()
	if err != nil {
		log.Fatalf("Error marshaling certificate: %v", err)
	}

	k.KeyVal.Private = string(keyBytes)
	k.KeyVal.Public = string(svid)
	k.KeyVal.Certificate = string(svid)
	k.Scheme = ecdsaSha2nistp384
	k.KeyType = ecdsaKeyType //this should be known from the ASN.1 data
	k.KeyIDHashAlgorithms = []string{"sha256", "sha512"}
	err = k.generateKeyID()
	if err != nil {
		log.Fatalf("Error generating keyID: %v", err)
	}

	return k
}

func GetTrustBundle(ctx context.Context, socketPath string) []*x509.Certificate {
	client, err := workloadapi.New(ctx, workloadapi.WithAddr(socketPath))
	if err != nil {
		log.Fatalf("Unable to create workload API client: %v", err)
	}
	defer client.Close()

	bundles, err := client.FetchX509Bundles(ctx)
	if err != nil {
		log.Fatalf("Error fetching x.509 bundles: %v", err)
	}

	certs := []*x509.Certificate{}
	for _, bundle := range bundles.Bundles() {
		certs = append(certs, bundle.X509Authorities()...)
	}

	return certs
}
