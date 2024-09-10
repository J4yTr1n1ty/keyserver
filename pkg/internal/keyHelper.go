package internal

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

// ExtractPublicKey extracts and serializes only the key part of the PGP public key.
func ExtractPublicKey(pgpKey string) (string, error) {
	// Decode the armored key
	block, err := armor.Decode(bytes.NewReader([]byte(pgpKey)))
	if err != nil {
		return "", fmt.Errorf("failed to decode armored key: %v", err)
	}

	// Parse the key into entities
	entityList, err := openpgp.ReadKeyRing(block.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read key ring: %v", err)
	}

	// Buffer to store the serialized key
	var buf bytes.Buffer

	// Re-armoring the key
	w, err := armor.Encode(&buf, "PGP PUBLIC KEY BLOCK", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create armor encoder: %v", err)
	}

	// Serialize only the key
	for _, entity := range entityList {
		err := entity.Serialize(w)
		if err != nil {
			return "", fmt.Errorf("failed to serialize entity: %v", err)
		}
	}

	// Close the armor encoder
	err = w.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close armor encoder: %v", err)
	}

	// Return the clean key without any comments or extra metadata
	return buf.String(), nil
}
