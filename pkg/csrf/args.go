package csrf

import (
	"encoding/base64"
	"log"
	"os"
)

const keySizeBytes = 256

var (
	decodedCSRFKey = ""
)

func Ensure(argCSRFKey string) {

	decoded, err := base64.StdEncoding.DecodeString(argCSRFKey)
	if err != nil {
		log.Fatal(err, "Could not decode CSRF key")
		os.Exit(255)
	}

	if len(decoded) != keySizeBytes {
		log.Fatalf("Could not validate CSRF key. Expected size %d, got %d", keySizeBytes, len(decoded))
	}

	decodedCSRFKey = string(decoded)
}

func Key() string {
	if len(decodedCSRFKey) == 0 {
		log.Fatal("CSRF key was not properly initialized. Run 'csrf.Ensure()' first.")
	}

	return decodedCSRFKey
}
