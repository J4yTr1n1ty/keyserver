package storage

import (
	"errors"
	"strings"

	"golang.org/x/crypto/openpgp"

	"github.com/J4yTr1n1ty/keyserver/pkg/boot"
	"github.com/J4yTr1n1ty/keyserver/pkg/models"
)

// Fucntion that returns all entity id strings and their names
func GetKeyIdentities(publicKeyArmored string) ([]string, []string, error) {
	enityList, err := VerifyKey(publicKeyArmored)
	if err != nil {
		return nil, nil, err
	}

	var identities []string
	var identityNames []string
	for _, entity := range enityList {
		identities = append(identities, entity.PrimaryKey.KeyIdString())
		for _, identity := range entity.Identities {
			identityNames = append(identityNames, identity.Name)
		}
	}

	if len(identities) == 0 {
		return nil, nil, errors.New("no identities found in the provided armored key")
	}

	if len(identityNames) == 0 {
		return nil, nil, errors.New("no identity names found in the provided armored key")
	}

	return identities, identityNames, nil
}

func GetUniqueIdentities() []models.Identity {
	var identities []models.Identity
	boot.DB.Distinct("name", "key_fingerprint").Find(&identities).Preload("Key")
	return identities
}

func VerifyKey(publicKeyArmored string) (openpgp.EntityList, error) {
	entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(publicKeyArmored))
	if err != nil {
		return nil, errors.New("invalid PGP public key")
	}

	if len(entityList) == 0 {
		return nil, errors.New("no keys found in the provided armored key")
	}

	for _, entity := range entityList {
		if entity.PrivateKey != nil {
			return nil, errors.New("private key detected, only public keys are accepted")
		}
	}

	return entityList, nil
}

func SubmitKey(publicKeyArmored string) error {
	entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(publicKeyArmored))
	if err != nil {
		return errors.New("Invalid PGP public key")
	}

	if len(entityList) == 0 {
		return errors.New("No keys found in the provided armored key")
	}

	entity := entityList[0]
	if entity.PrivateKey != nil {
		return errors.New("Private key detected, only public keys are accepted")
	}

	key := models.Key{
		Fingerprint: entity.PrimaryKey.KeyIdString(),
		PublicKey:   publicKeyArmored, // TODO: rethink this for storage efficiency
	}

	for _, identity := range entity.Identities {
		key.Identities = append(key.Identities, models.Identity{
			Name:    identity.UserId.Name,
			Comment: identity.UserId.Comment,
			Email:   identity.UserId.Email,
		})
	}

	result := boot.DB.Create(&key)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			return errors.New("A key for this fingerprint already exists")
		}
		return errors.New("Failed to store key in database")
	}

	return nil
}

func ListAllKeys() []models.Key {
	var keys []models.Key
	boot.DB.Preload("Identities").Find(&keys)
	return keys
}

func GetKey(email string) (models.Key, error) {
	var key models.Key
	result := boot.DB.Preload("Identities").Joins("JOIN identities ON identities.key_fingerprint = keys.fingerprint").Where("identities.email = ?", email).First(&key)

	if result.Error == nil {
		return key, nil
	}

	return models.Key{}, errors.New("Key not found")
}
