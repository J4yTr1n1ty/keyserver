package storage

import (
	"errors"
	"strings"

	"github.com/ProtonMail/gopenpgp/v3/crypto"

	"github.com/J4yTr1n1ty/keyserver/pkg/boot"
	"github.com/J4yTr1n1ty/keyserver/pkg/models"
)

func VerifyKey(publicKeyArmored string) (*crypto.Key, error) {
	// HACK: Is this the best way to do this?
	key, err := crypto.NewKeyFromArmored(publicKeyArmored)
	if err != nil {
		return nil, errors.New("Invalid PGP public key")
	}

	return key, nil
}

func SubmitKey(publicKeyArmored string) error {
	importedKey, err := VerifyKey(publicKeyArmored)
	if err != nil {
		return err
	}
	keyRing, err := crypto.NewKeyRing(importedKey)
	if err != nil {
		return err
	}

	// TODO: Handle multiple keys
	pubKeyArmored, err := keyRing.GetKey(0)
	if err != nil {
		return errors.New("Failed to serialize public key")
	}

	creationTime := keyRing.GetKeys()[0].GetEntity().PrimaryKey.CreationTime

	keyBytes, err := pubKeyArmored.Serialize()
	if err != nil {
		return err
	}

	key := models.Key{
		Fingerprint: keyRing.GetKeys()[0].GetFingerprint(),
		PublicKey:   keyBytes,
		ValidFrom:   creationTime,
		// ValidUntil:  expirationTime, // TODO: Implement expiration time
	}

	// TODO: Possibly reduce the amount of identities stored in the database because of duplicate keys
	for _, identity := range keyRing.GetIdentities() {
		if identity.Email == "" {
			continue
		}
		key.Identities = append(key.Identities, models.Identity{
			Name:  identity.Name,
			Email: identity.Email,
		})
	}

	if len(key.Identities) == 0 {
		return errors.New("No identities found in the provided armored key")
	}

	result := boot.DB.Create(&key)
	if result.Error != nil {
		// TODO: Handle updating key if it already exists and expiration time has changed
		if strings.Contains(result.Error.Error(), "duplicate key") {
			return errors.New("A key for this fingerprint already exists")
		}
		return errors.New("Failed to store key in database")
	}

	return nil
}

func GetUniqueIdentities() []models.Identity {
	var identities []models.Identity
	boot.DB.Select("DISTINCT ON (name, key_fingerprint) *").Order("name, key_fingerprint, created_at DESC").Find(&identities)
	return identities
}

func ListAllKeys() []models.Key {
	var keys []models.Key
	boot.DB.Preload("Identities").Find(&keys)
	return keys
}

func GetKeyByEmail(email string) (*models.Key, error) {
	var key models.Key
	result := boot.DB.Preload("Identities").Joins("JOIN identities ON identities.key_fingerprint = keys.fingerprint").Where("identities.email = ?", email).First(&key)
	if result.Error != nil {
		return nil, result.Error
	}
	return &key, nil
}

func GetKeyByFingerprint(fingerprint string) (*crypto.Key, error) {
	var key models.Key
	result := boot.DB.Where("fingerprint = ?", fingerprint).First(&key)
	if result.Error != nil {
		return nil, result.Error
	}
	keyring, err := crypto.NewKeyRingFromBinary(key.PublicKey)
	if err != nil {
		return nil, err
	}

	// HACK: Per convention, we only have one key
	return keyring.GetKey(0)
}
