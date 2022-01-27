package routes

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func populateWallet(wallet *gateway.Wallet) error {
	log.Info().Msg("============ Populating wallet ============")
	credPath := filepath.Join(
		"..",
		"..",
		"organizations",
		"peerOrganizations",
		"supplier1.steelplatform.com",
		"users",
		"User1@supplier1.steelplatform.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Supplier1MSP", string(cert), string(key))
	fmt.Println("identity.MspID = ", identity.MspID)

	return wallet.Put("appUser", identity)
}
