package internal_test

import (
	"testing"

	"github.com/J4yTr1n1ty/keyserver/pkg/internal"
)

func TestExtractPGPMessageAndSignature(t *testing.T) {
	testString := `-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA256

Test Message
-----BEGIN PGP SIGNATURE-----

iQIzBAEBCAAdFiEE3kR/d1SuxTT3CvPps1HOpXi+0vwFAmbhoSoACgkQs1HOpXi+
0vyxkQ//cGZ/te4JVOZYJfphg3K2t2JvzZ/8/LYVav/zMM6bjzWVwpRfi6zBEs5a
IwnxQgAEWRuK+nPg1grYafcsBBiH2Yy3a1nkHurv2tscmd2lBDx8tJiwSdLKNQTs
uTVlQMP8V4qXP/GpfYrLvrOB2eRsLzpKaWLAqAqL2X8qxfKy8Tm2zvE69TE8jcSx
VAzgGFDap1xnvwtaPuNipa88jGMmj6gTVArA7BFh1EuReh1L7FMbH3Yn1nUn5k8h
47P2dEu5WjGnyz5KPZBWhyntW3CNKtkCpkfrLPPZTyMYytvrYS4pMMa9ZRvG2hUo
JDeqo1nohx6QgWNDhmdtpCwl2nEdv5M9Ckxt0gbPRSvBKUEZxi02Lpk5rTzkrSe5
7VRnrbTIWLDuw2nLn1jbjh3XkuOR8zb6dGPadP4inE4LSnru5NZpgsS3bffUCYhz
OdEubGkkcCVlIdyv9Lu48FKthJH62veJ1Me+7yM0cuXewF/482ycLaeusYIO9QSP
Qbx5T0Gh8Si9k3yNPXBzM7lj8Yf2aek14136fQpIlBSvy/5oGCpqCRJ44wjJWnUV
0eNv6SuQnFXLv5n2VlmDprruNf1r63eBTDh6Xo3K9GHzYZENqTFcBwPXTrxyQ+9r
68y1JbMV7uCjvn37FKRC8vGq4yES9Pl8SaFKoyBwMyG/w9TQAqU=
=LUpM
-----END PGP SIGNATURE-----`

	message, signature, err := internal.ExtractPGPMessageAndSignature(testString)
	if err != nil {
		t.Errorf("Failed to extract PGP message and signature: %v", err)
	}

	expectedMessage := `-----BEGIN PGP MESSAGE-----
  Hash: SHA256

  Test Message
  -----END PGP MESSAGE-----`

	expectedSignature := `-----BEGIN PGP SIGNATURE-----

iQIzBAEBCAAdFiEE3kR/d1SuxTT3CvPps1HOpXi+0vwFAmbhoSoACgkQs1HOpXi+
0vyxkQ//cGZ/te4JVOZYJfphg3K2t2JvzZ/8/LYVav/zMM6bjzWVwpRfi6zBEs5a
IwnxQgAEWRuK+nPg1grYafcsBBiH2Yy3a1nkHurv2tscmd2lBDx8tJiwSdLKNQTs
uTVlQMP8V4qXP/GpfYrLvrOB2eRsLzpKaWLAqAqL2X8qxfKy8Tm2zvE69TE8jcSx
VAzgGFDap1xnvwtaPuNipa88jGMmj6gTVArA7BFh1EuReh1L7FMbH3Yn1nUn5k8h
47P2dEu5WjGnyz5KPZBWhyntW3CNKtkCpkfrLPPZTyMYytvrYS4pMMa9ZRvG2hUo
JDeqo1nohx6QgWNDhmdtpCwl2nEdv5M9Ckxt0gbPRSvBKUEZxi02Lpk5rTzkrSe5
7VRnrbTIWLDuw2nLn1jbjh3XkuOR8zb6dGPadP4inE4LSnru5NZpgsS3bffUCYhz
OdEubGkkcCVlIdyv9Lu48FKthJH62veJ1Me+7yM0cuXewF/482ycLaeusYIO9QSP
Qbx5T0Gh8Si9k3yNPXBzM7lj8Yf2aek14136fQpIlBSvy/5oGCpqCRJ44wjJWnUV
0eNv6SuQnFXLv5n2VlmDprruNf1r63eBTDh6Xo3K9GHzYZENqTFcBwPXTrxyQ+9r
68y1JbMV7uCjvn37FKRC8vGq4yES9Pl8SaFKoyBwMyG/w9TQAqU=
=LUpM
-----END PGP SIGNATURE-----`

	if message != expectedMessage {
		t.Errorf("Incorrect PGP message. \nActual:%s. \nExpected: %s", message, expectedMessage)
	}

	if signature != expectedSignature {
		t.Errorf("Incorrect PGP signature. \nActual:%s. \nExpected: %s", signature, expectedSignature)
	}
}
