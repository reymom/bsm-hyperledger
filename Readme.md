# STEEL PLATFORM

This repository contains the necessary files to construct a network with different suppliers, buyers and a logistic company, plus orderer, with all it's CA.

It contains also the chaincode to create and play auctions, and another one to assign and track the delivery of the winner of the corresponding auctions.

## Steps

### 1. Remove previous traces
```
docker rm -f $(docker ps -a -q)
docker volume prune
docker network prune

sudo rm -rf organizations/fabric-ca/ordererOrg/
sudo rm -rf organizations/fabric-ca/supplier1/
sudo rm -rf organizations/fabric-ca/supplier2/
sudo rm -rf organizations/fabric-ca/buyer1/
sudo rm -rf organizations/fabric-ca/buyer2/
sudo rm -rf organizations/fabric-ca/buyer3/
sudo rm -rf organizations/fabric-ca/logistics/
sudo rm -rf organizations/peerOrganizations
sudo rm -rf organizations/ordererOrganizations
sudo rm -rf system-genesis-block/
sudo rm -rf channel-artifacts/

```

### 2. Raise CA's, channels
```
mkdir channel-artifacts

docker-compose -f docker/docker-compose-steelplatform-ca.yaml up -d

```
```
export PATH=${PWD}/../fabric-samples/bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}/configtx

. ./organizations/fabric-ca/registerEnroll.sh && createSupplier1
. ./organizations/fabric-ca/registerEnroll.sh && createSupplier2
. ./organizations/fabric-ca/registerEnroll.sh && createBuyer1
. ./organizations/fabric-ca/registerEnroll.sh && createBuyer2
. ./organizations/fabric-ca/registerEnroll.sh && createBuyer3
. ./organizations/fabric-ca/registerEnroll.sh && createLogistics
. ./organizations/fabric-ca/registerEnroll.sh && createOrderer

# register channels
configtxgen -profile Public1ApplicationGenesis -outputBlock ./channel-artifacts/public1channel.block -channelID public1channel
configtxgen -profile Public2ApplicationGenesis -outputBlock ./channel-artifacts/public2channel.block -channelID public2channel

configtxgen -profile Logistics11ApplicationGenesis -outputBlock ./channel-artifacts/logistics11channel.block -channelID logistics11channel
configtxgen -profile Logistics12ApplicationGenesis -outputBlock ./channel-artifacts/logistics12channel.block -channelID logistics12channel
configtxgen -profile Logistics13ApplicationGenesis -outputBlock ./channel-artifacts/logistics13channel.block -channelID logistics13channel
configtxgen -profile Logistics21ApplicationGenesis -outputBlock ./channel-artifacts/logistics21channel.block -channelID logistics21channel
configtxgen -profile Logistics22ApplicationGenesis -outputBlock ./channel-artifacts/logistics22channel.block -channelID logistics22channel
configtxgen -profile Logistics23ApplicationGenesis -outputBlock ./channel-artifacts/logistics23channel.block -channelID logistics23channel

```

### 3. Enrol the companies into the channels
```
export FABRIC_CFG_PATH=${PWD}/../fabric-samples/config
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem
export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/tls/server.crt
export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/tls/server.key

docker-compose -f docker/docker-compose-steelplatform.yaml -f docker/docker-compose-couch.yaml up -d

```
```
# - orderer
osnadmin channel join --channelID public1channel --config-block ./channel-artifacts/public1channel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
osnadmin channel join --channelID public2channel --config-block ./channel-artifacts/public2channel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"

osnadmin channel join --channelID logistics11channel --config-block ./channel-artifacts/logistics11channel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
osnadmin channel join --channelID logistics12channel --config-block ./channel-artifacts/logistics12channel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
osnadmin channel join --channelID logistics13channel --config-block ./channel-artifacts/logistics13channel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
osnadmin channel join --channelID logistics21channel --config-block ./channel-artifacts/logistics21channel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
osnadmin channel join --channelID logistics22channel --config-block ./channel-artifacts/logistics22channel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
osnadmin channel join --channelID logistics23channel --config-block ./channel-artifacts/logistics23channel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"

osnadmin channel list -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"

export CORE_PEER_TLS_ENABLED=true

# - supplier1

export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Supplier1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/users/Admin@supplier1.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:7051
peer channel join -b ./channel-artifacts/public1channel.block
peer channel join -b ./channel-artifacts/logistics11channel.block
peer channel join -b ./channel-artifacts/logistics12channel.block
peer channel join -b ./channel-artifacts/logistics13channel.block

# - anchor peer
peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com -c public1channel --tls --cafile "$ORDERER_CA"

cd channel-artifacts

configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json

jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json config_copy.json
jq '.channel_group.groups.Application.groups.Supplier1MSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.supplier1.steelplatform.com","port": 7051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb

configtxlator compute_update --channel_id public1channel --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"public1channel", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json

configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

cd ..

peer channel update -f channel-artifacts/config_update_in_envelope.pb -c public1channel -o localhost:7050  --ordererTLSHostnameOverride orderer.steelplatform.com --tls --cafile "${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem"

# - supplier2
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Supplier2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/users/Admin@supplier2.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:9051
peer channel join -b ./channel-artifacts/public2channel.block
peer channel join -b ./channel-artifacts/logistics21channel.block
peer channel join -b ./channel-artifacts/logistics22channel.block
peer channel join -b ./channel-artifacts/logistics23channel.block

# - buyer1
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Buyer1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/users/Admin@buyer1.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:11051
peer channel join -b ./channel-artifacts/public1channel.block
peer channel join -b ./channel-artifacts/public2channel.block
peer channel join -b ./channel-artifacts/logistics11channel.block
peer channel join -b ./channel-artifacts/logistics21channel.block

# - anchor peer
peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com -c public1channel --tls --cafile "$ORDERER_CA"

cd channel-artifacts

configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json

jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json config_copy.json
jq '.channel_group.groups.Application.groups.Buyer1MSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.buyer1.steelplatform.com","port": 11051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb

configtxlator compute_update --channel_id public1channel --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"public1channel", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json

configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

cd ..

peer channel update -f channel-artifacts/config_update_in_envelope.pb -c public1channel -o localhost:7050  --ordererTLSHostnameOverride orderer.steelplatform.com --tls --cafile "${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem"

# - buyer2
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Buyer2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/users/Admin@buyer2.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:13051
peer channel join -b ./channel-artifacts/public1channel.block
peer channel join -b ./channel-artifacts/public2channel.block
peer channel join -b ./channel-artifacts/logistics12channel.block
peer channel join -b ./channel-artifacts/logistics22channel.block

# - buyer3
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Buyer3MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/users/Admin@buyer3.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:15051
peer channel join -b ./channel-artifacts/public1channel.block
peer channel join -b ./channel-artifacts/public2channel.block
peer channel join -b ./channel-artifacts/logistics13channel.block
peer channel join -b ./channel-artifacts/logistics23channel.block

# - logistics
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="LogisticsMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/users/Admin@logistics.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:17051
peer channel join -b ./channel-artifacts/logistics11channel.block
peer channel join -b ./channel-artifacts/logistics12channel.block
peer channel join -b ./channel-artifacts/logistics13channel.block
peer channel join -b ./channel-artifacts/logistics21channel.block
peer channel join -b ./channel-artifacts/logistics22channel.block
peer channel join -b ./channel-artifacts/logistics23channel.block

```

### 4. Install the chaincodes on the organisations channels

#### 4.1. Auction chaincode
```
cd chaincodes/auction/go/
go mod init auctionChaincode.go
go mod tidy
go mod vendor

cd ../../../
export PATH=${PWD}/../fabric-samples/bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}/../fabric-samples/config

peer lifecycle chaincode package auction.tar.gz --path chaincodes/auction/go --lang golang --label auction_1.0

export CORE_PEER_TLS_ENABLED=true

# -- supplier 1
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Supplier1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/users/Admin@supplier1.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:7051

peer lifecycle chaincode install auction.tar.gz
peer lifecycle chaincode queryinstalled
export CC_PACKAGE_ID=auction_1.0:1623100f07069d5c3b1180d764d26b5d5166589226c41153391ed1a1ceab36d7

# --- in public channel 1 with collection
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID public1channel --collections-config chaincodes/auction/go/collections_config_public1.json --name auction --signature-policy "OR('Supplier1MSP.member','Buyer1MSP.member','Buyer2MSP.member','Buyer3MSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

# -- supplier 2
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Supplier2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/users/Admin@supplier2.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:9051

# --- in public channel 2
peer lifecycle chaincode install auction.tar.gz
export CC_PACKAGE_ID=auction_1.0:1623100f07069d5c3b1180d764d26b5d5166589226c41153391ed1a1ceab36d7
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID public2channel --collections-config chaincodes/auction/go/collections_config_public2.json --name auction --signature-policy "OR('Supplier2MSP.member','Buyer1MSP.member','Buyer2MSP.member','Buyer3MSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

# -- buyer 1
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Buyer1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/users/Admin@buyer1.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:11051

peer lifecycle chaincode install auction.tar.gz
export CC_PACKAGE_ID=auction_1.0:1623100f07069d5c3b1180d764d26b5d5166589226c41153391ed1a1ceab36d7
# --- in public channel 1
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID public1channel --collections-config chaincodes/auction/go/collections_config_public1.json --name auction --signature-policy "OR('Supplier1MSP.member','Buyer1MSP.member','Buyer2MSP.member','Buyer3MSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem
# --- in public channel 2
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID public2channel --collections-config chaincodes/auction/go/collections_config_public2.json --name auction --signature-policy "OR('Supplier2MSP.member','Buyer1MSP.member','Buyer2MSP.member','Buyer3MSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

# -- buyer 2
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Buyer2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/users/Admin@buyer2.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:13051

peer lifecycle chaincode install auction.tar.gz
export CC_PACKAGE_ID=auction_1.0:1623100f07069d5c3b1180d764d26b5d5166589226c41153391ed1a1ceab36d7
# --- in public channel 1
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID public1channel --collections-config chaincodes/auction/go/collections_config_public1.json --name auction --signature-policy "OR('Supplier1MSP.member','Buyer1MSP.member','Buyer2MSP.member','Buyer3MSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem
# --- in public channel 2
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID public2channel --collections-config chaincodes/auction/go/collections_config_public2.json --name auction --signature-policy "OR('Supplier2MSP.member','Buyer1MSP.member','Buyer2MSP.member','Buyer3MSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

# -- buyer 3
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Buyer3MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/users/Admin@buyer3.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:15051

peer lifecycle chaincode install auction.tar.gz
export CC_PACKAGE_ID=auction_1.0:1623100f07069d5c3b1180d764d26b5d5166589226c41153391ed1a1ceab36d7
# --- in public channel 1
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID public1channel --collections-config chaincodes/auction/go/collections_config_public1.json --name auction --signature-policy "OR('Supplier1MSP.member','Buyer1MSP.member','Buyer2MSP.member','Buyer3MSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem
# --- in public channel 2
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID public2channel --collections-config chaincodes/auction/go/collections_config_public2.json --name auction --signature-policy "OR('Supplier2MSP.member', 'Buyer1MSP.member','Buyer2MSP.member','Buyer3MSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

# -- commit
peer lifecycle chaincode checkcommitreadiness --channelID public1channel --name auction --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem --output json

# --- in public channel 1
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --collections-config chaincodes/auction/go/collections_config_public1.json --channelID public1channel --name auction --signature-policy "OR('Supplier1MSP.member','Buyer1MSP.member','Buyer2MSP.member','Buyer3MSP.member')" --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt --peerAddresses localhost:11051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/ca.crt --peerAddresses localhost:13051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/ca.crt --peerAddresses localhost:15051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/tls/ca.crt

peer lifecycle chaincode querycommitted --channelID public1channel --name auction --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

# --- in public channel 2
peer lifecycle chaincode checkcommitreadiness --channelID public2channel --name auction --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem --output json

peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --collections-config chaincodes/auction/go/collections_config_public2.json --channelID public2channel --name auction --signature-policy "OR('Supplier2MSP.member','Buyer1MSP.member','Buyer2MSP.member','Buyer3MSP.member')" --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/tls/ca.crt --peerAddresses localhost:11051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/ca.crt --peerAddresses localhost:13051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/ca.crt --peerAddresses localhost:15051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/tls/ca.crt

peer lifecycle chaincode querycommitted --channelID public2channel --name auction --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

```
##### Test auction chaincode
```
# create a public auction
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem -C public1channel -n auction --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt -c '{"function":"CreateAuction","Args":["false", "", "Stainless Steel", "Sheets", "1000", "100", "1"]}'

# ups, just suppliers can do this
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Supplier1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/users/Admin@supplier1.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:7051
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem -C public1channel -n auction --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt -c '{"function":"CreateAuction","Args":["false", "", "Stainless Steel", "Sheets", "1000", "100", "1"]}'

peer chaincode query -C public1channel -n auction -c '{"Args":["GetAllAuctions","",""]}' | jq .

peer chaincode query -C public1channel -n auction -c '{"Args":["QueryAuction","e64d436d-89fa-491d-a60a-41bf559fe88c"]}' | jq .

# place a bid in a public auction
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem -C public1channel -n auction --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt -c '{"function":"Bid","Args":["false", "e64d436d-89fa-491d-a60a-41bf559fe88c", "", "101"]}'
# ups, just buyers can bid
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Buyer1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/users/Admin@buyer1.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:11051
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem -C public1channel -n auction --peerAddresses localhost:11051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/ca.crt -c '{"function":"Bid","Args":["false", "e64d436d-89fa-491d-a60a-41bf559fe88c", "", "101"]}'

# should be rejected because it's below minPrice
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem -C public1channel -n auction --peerAddresses localhost:11051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/ca.crt -c '{"function":"Bid","Args":["false", "ae91120e-3ee2-45be-b2c5-c9184ef753b0", "", "94"]}'

# another bid
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Buyer2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/users/Admin@buyer2.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:13051
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem -C public1channel -n auction --peerAddresses localhost:13051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/ca.crt -c '{"function":"Bid","Args":["false", "ae91120e-3ee2-45be-b2c5-c9184ef753b0", "", "101"]}'

# finish auction
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Supplier1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/users/Admin@supplier1.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:7051
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem -C public1channel -n auction --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt -c '{"function":"FinishAuction","Args":["false", "ae91120e-3ee2-45be-b2c5-c9184ef753b0", ""]}'

# create a private auction
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem -C public1channel -n auction --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt -c '{"function":"CreateAuction","Args":["true", "12", "Green Steel", "Coils", "2500", "120", "0.1"]}'

```

#### 4.2. Logistics chaincode

```
cd chaincodes/logistics/go/
go mod init logisticsChaincode.go
go mod tidy
go mod vendor

cd ../../../
export PATH=${PWD}/../fabric-samples/bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}/../fabric-samples/config

peer lifecycle chaincode package logistics.tar.gz --path chaincodes/logistics/go --lang golang --label logistics_1.0

export CORE_PEER_TLS_ENABLED=true

# -- supplier 1
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Supplier1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/users/Admin@supplier1.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:7051

# --- in logistics channel 11 12 13
peer lifecycle chaincode install logistics.tar.gz
export CC_PACKAGE_ID=logistics_1.0:f230b1c709c2b9006914ba60b18ee7bd2a5c3359cae7b37423afb88a6a375491

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics11channel --name logistics --signature-policy "OR('Supplier1MSP.member','Buyer1MSP.member','LogisticsMSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics12channel --name logistics --signature-policy "OR('Supplier1MSP.member','Buyer2MSP.member','LogisticsMSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics13channel --name logistics --signature-policy "OR('Supplier1MSP.member','Buyer3MSP.member','LogisticsMSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

# -- supplier 2
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Supplier2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/users/Admin@supplier2.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:9051

# --- in logistics channel 21 22 23
peer lifecycle chaincode install logistics.tar.gz
export CC_PACKAGE_ID=logistics_1.0:f230b1c709c2b9006914ba60b18ee7bd2a5c3359cae7b37423afb88a6a375491

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics21channel --name logistics --signature-policy "OR('Supplier2MSP.member','Buyer1MSP.member','LogisticsMSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics22channel --name logistics --signature-policy "OR('Supplier2MSP.member','Buyer2MSP.member','LogisticsMSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics23channel --name logistics --signature-policy "OR('Supplier2MSP.member','Buyer3MSP.member','LogisticsMSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

# -- buyer 1
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Buyer1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/users/Admin@buyer1.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:11051

# --- in logistics channel 11 and 21
peer lifecycle chaincode install logistics.tar.gz
export CC_PACKAGE_ID=logistics_1.0:f230b1c709c2b9006914ba60b18ee7bd2a5c3359cae7b37423afb88a6a375491

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics11channel --name logistics --signature-policy "OR('Supplier1MSP.member','Buyer1MSP.member','LogisticsMSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics21channel --name logistics --signature-policy "OR('Supplier2MSP.member','Buyer1MSP.member','LogisticsMSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

# -- buyer 2
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Buyer2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/users/Admin@buyer2.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:13051

# --- in logistics channel 12 and 22
peer lifecycle chaincode install logistics.tar.gz
export CC_PACKAGE_ID=logistics_1.0:f230b1c709c2b9006914ba60b18ee7bd2a5c3359cae7b37423afb88a6a375491

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics12channel --name logistics --signature-policy "OR('Supplier1MSP.member','Buyer2MSP.member','LogisticsMSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics22channel --name logistics --signature-policy "OR('Supplier2MSP.member','Buyer2MSP.member','LogisticsMSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

# -- buyer 3
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Buyer3MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/users/Admin@buyer3.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:15051

# --- in logistics channel 13 and 23
peer lifecycle chaincode install logistics.tar.gz
export CC_PACKAGE_ID=logistics_1.0:f230b1c709c2b9006914ba60b18ee7bd2a5c3359cae7b37423afb88a6a375491

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics13channel --name logistics --signature-policy "OR('Supplier1MSP.member','Buyer3MSP.member','LogisticsMSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics23channel --name logistics --signature-policy "OR('Supplier2MSP.member','Buyer3MSP.member','LogisticsMSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

# -- logistics
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="LogisticsMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/users/Admin@logistics.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:17051

# --- logistics channel 11, 12, 13, 21, 22, 23
peer lifecycle chaincode install logistics.tar.gz
export CC_PACKAGE_ID=logistics_1.0:f230b1c709c2b9006914ba60b18ee7bd2a5c3359cae7b37423afb88a6a375491

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics11channel --name logistics --signature-policy "OR('Supplier1MSP.member','Buyer1MSP.member','LogisticsMSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics12channel --name logistics --signature-policy "OR('Supplier1MSP.member','Buyer2MSP.member','LogisticsMSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics13channel --name logistics --signature-policy "OR('Supplier1MSP.member','Buyer3MSP.member','LogisticsMSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics21channel --name logistics --signature-policy "OR('Supplier2MSP.member','Buyer1MSP.member','LogisticsMSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics22channel --name logistics --signature-policy "OR('Supplier2MSP.member','Buyer2MSP.member','LogisticsMSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics23channel --name logistics --signature-policy "OR('Supplier2MSP.member','Buyer3MSP.member','LogisticsMSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

# -- commit

# --- in logistics 11
peer lifecycle chaincode checkcommitreadiness --channelID logistics11channel --name logistics --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem --output json

peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics11channel --name logistics --signature-policy "OR('Supplier1MSP.member','Buyer1MSP.member','LogisticsMSP.member')" --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt --peerAddresses localhost:11051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/ca.crt --peerAddresses localhost:17051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/tls/ca.crt

peer lifecycle chaincode querycommitted --channelID logistics11channel --name logistics --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

# --- in logistics 12
peer lifecycle chaincode checkcommitreadiness --channelID logistics12channel --name logistics --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem --output json

peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics12channel --name logistics --signature-policy "OR('Supplier1MSP.member','Buyer2MSP.member','LogisticsMSP.member')" --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt --peerAddresses localhost:13051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/ca.crt --peerAddresses localhost:17051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/tls/ca.crt

peer lifecycle chaincode querycommitted --channelID logistics12channel --name logistics --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

# --- in logistics 13
peer lifecycle chaincode checkcommitreadiness --channelID logistics13channel --name logistics --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem --output json

peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics13channel --name logistics --signature-policy "OR('Supplier1MSP.member','Buyer3MSP.member','LogisticsMSP.member')" --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt --peerAddresses localhost:15051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/tls/ca.crt --peerAddresses localhost:17051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/tls/ca.crt

peer lifecycle chaincode querycommitted --channelID logistics13channel --name logistics --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

# --- in logistics 21
peer lifecycle chaincode checkcommitreadiness --channelID logistics21channel --name logistics --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem --output json

peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics21channel --name logistics --signature-policy "OR('Supplier2MSP.member','Buyer1MSP.member','LogisticsMSP.member')" --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/tls/ca.crt --peerAddresses localhost:11051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/ca.crt --peerAddresses localhost:17051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/tls/ca.crt

peer lifecycle chaincode querycommitted --channelID logistics21channel --name logistics --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

# --- in logistics 22
peer lifecycle chaincode checkcommitreadiness --channelID logistics22channel --name logistics --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem --output json

peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics22channel --name logistics --signature-policy "OR('Supplier2MSP.member','Buyer2MSP.member','LogisticsMSP.member')" --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/tls/ca.crt --peerAddresses localhost:13051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/ca.crt --peerAddresses localhost:17051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/tls/ca.crt

peer lifecycle chaincode querycommitted --channelID logistics22channel --name logistics --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

# --- in logistics 23
peer lifecycle chaincode checkcommitreadiness --channelID logistics23channel --name logistics --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem --output json

peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID logistics23channel --name logistics --signature-policy "OR('Supplier2MSP.member','Buyer3MSP.member','LogisticsMSP.member')" --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/tls/ca.crt --peerAddresses localhost:15051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/tls/ca.crt --peerAddresses localhost:17051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/tls/ca.crt

peer lifecycle chaincode querycommitted --channelID logistics23channel --name logistics --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

```

##### Test logistics chaincode
```
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Supplier1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/users/Admin@supplier1.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:7051

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem -C logistics11channel -n logistics --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt -c '{"function":"CreateDelivery","Args":["f64d436d-89fa-491d-a60a-41bf559fe88c", "LogisticsMSP", "Spain", "Barcelona", "C/Muntaner", "2"]}'

peer chaincode query -C logistics11channel -n logistics -c '{"Args":["GetAllDeliveries","",""]}' | jq .

peer chaincode query -C logistics11channel -n logistics -c '{"Args":["QueryDelivery","f64d436d-89fa-491d-a60a-41bf559fe88c"]}' | jq .

#change state
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem -C logistics11channel -n logistics --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt -c '{"function":"UpdateDeliveryStatus","Args":["f64d436d-89fa-491d-a60a-41bf559fe88c", "1"]}'

#error, just logistic companies can update state of deliveries
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="LogisticsMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/users/Admin@logistics.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:17051

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem -C logistics11channel -n logistics --peerAddresses localhost:17051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/tls/ca.crt -c '{"function":"UpdateDeliveryStatus","Args":["f64d436d-89fa-491d-a60a-41bf559fe88c", "1"]}'

```

### 5. Organization Networks for running the application
```
cd application/go

rm -rf vendor
rm -rf wallet
rm steelPlatform

cd ../../

./organizations/ccp-generate.sh

cd application/go

GO111MODULE=on go mod vendor

go build cmd/app/steelPlatform.go

go run cmd/app/steelPlatform.go

```