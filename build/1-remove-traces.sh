echo 'Removing traces...'

docker rm -f $(docker ps -a -q)
docker volume prune
docker network prune

rm -rf */*/*/vendor
rm auction.tar.gz logistics.tar.gz
rm -rf application/go/vendor
rm -rf application/go/wallet
rm application/go/steelPlatform

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

echo 'Done 🥳'