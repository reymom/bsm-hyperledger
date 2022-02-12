cd application/go

rm -rf vendor
rm -rf wallet
rm steelPlatform

cd ../../

./organizations/ccp-generate.sh

cd application/go

GO111MODULE=on go mod vendor

ENV_DAL=`echo $DISCOVERY_AS_LOCALHOST`

echo "ENV_DAL:"$DISCOVERY_AS_LOCALHOST

if [ "$ENV_DAL" != "true" ]
then
	export DISCOVERY_AS_LOCALHOST=true
fi

echo "DISCOVERY_AS_LOCALHOST="$DISCOVERY_AS_LOCALHOST

go build cmd/app/steelPlatform.go

echo "run steelPlatform..."

go run cmd/app/steelPlatform.go
