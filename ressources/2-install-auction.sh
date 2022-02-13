echo 'Installing auction...'

cd chaincodes/auction/go/
go mod init auctionChaincode.go
go mod tidy
go mod vendor

cd ../../../
export PATH=${PWD}/../fabric-samples/bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}/../fabric-samples/config

peer lifecycle chaincode package auction.tar.gz --path chaincodes/auction/go --lang golang --label auction_1.0

export CORE_PEER_TLS_ENABLED=true

export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Supplier1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/users/Admin@supplier1.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:7051

OUTPUT=$(peer lifecycle chaincode install auction.tar.gz)
PCKGID=$(echo $OUTPUT | sed -n -e 's/^.* auction_1.0://p')
export CC_PACKAGE_ID=auction_1.0:$PCKID

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID public1channel --collections-config chaincodes/auction/go/collections_config_public1.json --name auction --signature-policy "OR('Supplier1MSP.member','Buyer1MSP.member','Buyer2MSP.member','Buyer3MSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Supplier2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/users/Admin@supplier2.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:9051

OUTPUT=$(peer lifecycle chaincode install auction.tar.gz)
PCKGID=$(echo $OUTPUT | sed -n -e 's/^.* auction_1.0://p')
export CC_PACKAGE_ID=auction_1.0:$PCKID

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID public2channel --collections-config chaincodes/auction/go/collections_config_public2.json --name auction --signature-policy "OR('Supplier2MSP.member','Buyer1MSP.member','Buyer2MSP.member','Buyer3MSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Buyer1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/users/Admin@buyer1.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:11051

OUTPUT=$(peer lifecycle chaincode install auction.tar.gz)
PCKGID=$(echo $OUTPUT | sed -n -e 's/^.* auction_1.0://p')
export CC_PACKAGE_ID=auction_1.0:$PCKID

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID public1channel --collections-config chaincodes/auction/go/collections_config_public1.json --name auction --signature-policy "OR('Supplier1MSP.member','Buyer1MSP.member','Buyer2MSP.member','Buyer3MSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID public2channel --collections-config chaincodes/auction/go/collections_config_public2.json --name auction --signature-policy "OR('Supplier2MSP.member','Buyer1MSP.member','Buyer2MSP.member','Buyer3MSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Buyer2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/users/Admin@buyer2.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:13051

OUTPUT=$(peer lifecycle chaincode install auction.tar.gz)
PCKGID=$(echo $OUTPUT | sed -n -e 's/^.* auction_1.0://p')
export CC_PACKAGE_ID=auction_1.0:$PCKID
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID public1channel --collections-config chaincodes/auction/go/collections_config_public1.json --name auction --signature-policy "OR('Supplier1MSP.member','Buyer1MSP.member','Buyer2MSP.member','Buyer3MSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID public2channel --collections-config chaincodes/auction/go/collections_config_public2.json --name auction --signature-policy "OR('Supplier2MSP.member','Buyer1MSP.member','Buyer2MSP.member','Buyer3MSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

export PEER0_COMPANY_CA=${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Buyer3MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_COMPANY_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/users/Admin@buyer3.steelplatform.com/msp
export CORE_PEER_ADDRESS=localhost:15051

OUTPUT=$(peer lifecycle chaincode install auction.tar.gz)
PCKGID=$(echo $OUTPUT | sed -n -e 's/^.* auction_1.0://p')
export CC_PACKAGE_ID=auction_1.0:$PCKID
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID public1channel --collections-config chaincodes/auction/go/collections_config_public1.json --name auction --signature-policy "OR('Supplier1MSP.member','Buyer1MSP.member','Buyer2MSP.member','Buyer3MSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --channelID public2channel --collections-config chaincodes/auction/go/collections_config_public2.json --name auction --signature-policy "OR('Supplier2MSP.member', 'Buyer1MSP.member','Buyer2MSP.member','Buyer3MSP.member')" --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

peer lifecycle chaincode checkcommitreadiness --channelID public1channel --name auction --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem --output json

peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --collections-config chaincodes/auction/go/collections_config_public1.json --channelID public1channel --name auction --signature-policy "OR('Supplier1MSP.member','Buyer1MSP.member','Buyer2MSP.member','Buyer3MSP.member')" --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt --peerAddresses localhost:11051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/ca.crt --peerAddresses localhost:13051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/ca.crt --peerAddresses localhost:15051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/tls/ca.crt

peer lifecycle chaincode querycommitted --channelID public1channel --name auction --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem

peer lifecycle chaincode checkcommitreadiness --channelID public2channel --name auction --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem --output json

peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.steelplatform.com --collections-config chaincodes/auction/go/collections_config_public2.json --channelID public2channel --name auction --signature-policy "OR('Supplier2MSP.member','Buyer1MSP.member','Buyer2MSP.member','Buyer3MSP.member')" --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/tls/ca.crt --peerAddresses localhost:11051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/ca.crt --peerAddresses localhost:13051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/ca.crt --peerAddresses localhost:15051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/tls/ca.crt

peer lifecycle chaincode querycommitted --channelID public2channel --name auction --cafile ${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem


echo 'Done 🥳'