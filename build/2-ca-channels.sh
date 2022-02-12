echo 'Resgistering certificates and creating channels...'

mkdir channel-artifacts
docker-compose -f docker/docker-compose-steelplatform-ca.yaml up -d
echo 'Raising dockers...'
sleep 4

export PATH=${PWD}/../fabric-samples/bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}/configtx

. ./organizations/fabric-ca/registerEnroll.sh && createSupplier1
. ./organizations/fabric-ca/registerEnroll.sh && createSupplier2
. ./organizations/fabric-ca/registerEnroll.sh && createBuyer1
. ./organizations/fabric-ca/registerEnroll.sh && createBuyer2
. ./organizations/fabric-ca/registerEnroll.sh && createBuyer3
. ./organizations/fabric-ca/registerEnroll.sh && createLogistics
. ./organizations/fabric-ca/registerEnroll.sh && createOrderer

configtxgen -profile Public1ApplicationGenesis -outputBlock ./channel-artifacts/public1channel.block -channelID public1channel
configtxgen -profile Public2ApplicationGenesis -outputBlock ./channel-artifacts/public2channel.block -channelID public2channel
configtxgen -profile LogisticsApplicationGenesis -outputBlock ./channel-artifacts/logisticschannel.block -channelID logisticschannel

echo 'Done 🥳'