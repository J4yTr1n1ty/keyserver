package internal

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ProtonMail/gopenpgp/v3/crypto"

	"github.com/J4yTr1n1ty/keyserver/pkg/internal/storage"
)

func VerifyKey(publicKeyArmored string) (*crypto.Key, error) {
	key, err := crypto.NewKeyFromArmored(publicKeyArmored)
	if err != nil {
		return nil, errors.New("Invalid PGP public key")
	}

	return key, nil
}

func VerifyMessage(fingerprint, message string) (string, error) {
	key, err := storage.GetKeyByFingerprint(fingerprint)
	if err != nil {
		return "", err
	}
	if key == nil {
		return "", fmt.Errorf("key not found for fingerprint: %s", fingerprint)
	}

	pgp := crypto.PGP()
	verifier, err := pgp.Verify().VerificationKey(key).New()
	if err != nil {
		return "", fmt.Errorf("failed to create verifier: %v", err)
	}
	if verifier == nil {
		return "", fmt.Errorf("verifier is nil")
	}

	verifyResult, err := verifier.VerifyCleartext([]byte(message))
	if err != nil {
		return "", fmt.Errorf("failed to verify message: %v", err)
	}

	if sigErr := verifyResult.SignatureError(); sigErr != nil {
		return "", sigErr
	}

	creationTime := time.Unix(verifyResult.SignatureCreationTime(), 0)

	return creationTime.Format(time.RFC822), nil
}

// GetKeyIdentities returns all entity id strings and their names
func GetKeyIdentities(publicKeyArmored string) ([]string, []string, error) {
	key, err := crypto.NewKeyFromArmored(publicKeyArmored)
	if err != nil {
		return nil, nil, err
	}

	keyRing, err := crypto.NewKeyRing(key)
	if err != nil {
		return nil, nil, err
	}

	var identities []string
	var identityNames []string
	for _, key := range keyRing.GetKeys() {
		identities = append(identities, key.GetFingerprint())
	}

	for _, identity := range keyRing.GetIdentities() {
		identityNames = append(identityNames, identity.Name)
	}

	if len(identities) == 0 {
		return nil, nil, errors.New("no identities found in the provided armored key")
	}

	if len(identityNames) == 0 {
		return nil, nil, errors.New("no identity names found in the provided armored key")
	}

	return identities, identityNames, nil
}
