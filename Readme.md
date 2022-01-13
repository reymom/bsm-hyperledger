# STEEL PLATFORM

This repository contains the necessary files to construct a network with different suppliers, buyers and a logistic company, plus orderer, with all it's CA.

It contains also the chaincode to create and play auctions, and another one to assign and track the delivery of the winner of the corresponding auctions.

Steps:
```
# remove previous traces
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
sudo rm -rf channel-artifacts/
mkdir channel-artifacts

# raise CA's and enroll organizations
docker-compose -f docker/docker-compose-steelplatform-ca.yaml up -d

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
configtxgen -profile Private112ApplicationGenesis -outputBlock ./channel-artifacts/private112channel.block -channelID private112channel
configtxgen -profile Private123ApplicationGenesis -outputBlock ./channel-artifacts/private123channel.block -channelID private123channel
configtxgen -profile Private212ApplicationGenesis -outputBlock ./channel-artifacts/private212channel.block -channelID private212channel
configtxgen -profile Private223ApplicationGenesis -outputBlock ./channel-artifacts/private223channel.block -channelID private223channel
configtxgen -profile Logistics11ApplicationGenesis -outputBlock ./channel-artifacts/logistics11channel.block -channelID logistics11channel
configtxgen -profile Logistics12ApplicationGenesis -outputBlock ./channel-artifacts/logistics12channel.block -channelID logistics12channel
configtxgen -profile Logistics13ApplicationGenesis -outputBlock ./channel-artifacts/logistics13channel.block -channelID logistics13channel
configtxgen -profile Logistics21ApplicationGenesis -outputBlock ./channel-artifacts/logistics21channel.block -channelID logistics21channel
configtxgen -profile Logistics22ApplicationGenesis -outputBlock ./channel-artifacts/logistics22channel.block -channelID logistics22channel
configtxgen -profile Logistics23ApplicationGenesis -outputBlock ./channel-artifacts/logistics23channel.block -channelID logistics23channel

# register the companies in the channels
export FABRIC_CFG_PATH=${PWD}/../fabric-samples/config
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem
export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/tls/server.crt
export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/tls/server.key

# register companies in the channels
docker-compose -f docker/docker-compose-steelplatform.yaml -f docker/docker-compose-couch.yaml up -d

# - orderer
osnadmin channel join --channelID public1channel --config-block ./channel-artifacts/public1channel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
osnadmin channel join --channelID public2channel --config-block ./channel-artifacts/public2channel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
osnadmin channel join --channelID private112channel --config-block ./channel-artifacts/private112channel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
osnadmin channel join --channelID private123channel --config-block ./channel-artifacts/private123channel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
osnadmin channel join --channelID private212channel --config-block ./channel-artifacts/private212channel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
osnadmin channel join --channelID private223channel --config-block ./channel-artifacts/private223channel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
osnadmin channel join --channelID logistics11channel --config-block ./channel-artifacts/logistics11channel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
osnadmin channel join --channelID logistics12channel --config-block ./channel-artifacts/logistics12channel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
osnadmin channel join --channelID logistics13channel --config-block ./channel-artifacts/logistics13channel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
osnadmin channel join --channelID logistics21channel --config-block ./channel-artifacts/logistics21channel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
osnadmin channel join --channelID logistics22channel --config-block ./channel-artifacts/logistics22channel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
osnadmin channel join --channelID logistics23channel --config-block ./channel-artifacts/logistics23channel.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"

# osnadmin channel list -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"

# - supplier1
export CORE_PEER_TLS_ENABLED=true
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Supplier1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/users/Admin@supplier1.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:7051
peer channel join -b ./channel-artifacts/public1channel.block
peer channel join -b ./channel-artifacts/private112channel.block
peer channel join -b ./channel-artifacts/private123channel.block
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
peer channel join -b ./channel-artifacts/private212channel.block
peer channel join -b ./channel-artifacts/private223channel.block
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
peer channel join -b ./channel-artifacts/private112channel.block
peer channel join -b ./channel-artifacts/private212channel.block
peer channel join -b ./channel-artifacts/logistics11channel.block
peer channel join -b ./channel-artifacts/logistics21channel.block

# - buyer2
export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Buyer2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/users/Admin@buyer2.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:13051
peer channel join -b ./channel-artifacts/public1channel.block
peer channel join -b ./channel-artifacts/public2channel.block
peer channel join -b ./channel-artifacts/private112channel.block
peer channel join -b ./channel-artifacts/private123channel.block
peer channel join -b ./channel-artifacts/private212channel.block
peer channel join -b ./channel-artifacts/private223channel.block
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
peer channel join -b ./channel-artifacts/private123channel.block
peer channel join -b ./channel-artifacts/private223channel.block
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
