package connection

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/rs/zerolog/log"

	fabricConfig "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
)

func GetGatewayObjects(org Organization) (*gateway.Gateway, NetworkContract, error) {
	wallet, e := gateway.NewFileSystemWallet("wallet")
	if e != nil {
		log.Err(e).Msg("Failed to create wallet")
	}

	if !wallet.Exists(string(org) + "-wallet") {
		e = populateWallet(wallet, org)
		if e != nil {
			log.Err(e).Msg("Failed to populate wallet contents")
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"organizations",
		"peerOrganizations",
		string(org)+".steelplatform.com",
		"connection-"+string(org)+".yaml",
	)

	gw, e := gateway.Connect(
		gateway.WithConfig(fabricConfig.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, string(org)+"-wallet"),
	)
	if e != nil {
		log.Err(e).Msg("Failed to connect to gateway")
	}
	defer gw.Close()

	var (
		gwContractStruct    GatewayContract
		gwNetwork           *gateway.Network
		networkContractsMap = make(NetworkContract)
	)
	orgNetworks := org.getNetworks()
	for _, channel := range orgNetworks {
		gwNetwork, e = gw.GetNetwork(string(channel))
		if e != nil {
			log.Err(e).Msg("Failed to get network")
		} else {
			contract := channel.GetContract()
			gwContract := gwNetwork.GetContract(string(contract))
			gwContractStruct = GatewayContract{
				Name:       contract,
				GwContract: gwNetwork.GetContract(string(contract)),
			}
			networkContractsMap[channel] = gwContractStruct
			log.Info().Msgf("Loaded Contract: %s in network %s", gwContract.Name(), channel)
		}
	}

	return gw, networkContractsMap, nil
}

func populateWallet(wallet *gateway.Wallet, org Organization) error {
	log.Info().Msg("============ Populating wallet ============")
	credPath := filepath.Join(
		"..",
		"..",
		"organizations",
		"peerOrganizations",
		string(org)+".steelplatform.com",
		"users",
		"User1@"+string(org)+".steelplatform.com",
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

	identity := gateway.NewX509Identity(org.getMSP(), string(cert), string(key))

	return wallet.Put(string(org)+"-wallet", identity)
}
