#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${SUP}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        organizations/ccp-template.json
}

function yaml_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${SUP}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        organizations/ccp-template.yaml | sed -e $'s/\\\\n/\\\n          /g'
}

SUP=1
P0PORT=7051
CAPORT=7054
PEERPEM=organizations/peerOrganizations/supplier1.steelplatform.com/tlsca/tlsca.supplier1.steelplatform.com-cert.pem
CAPEM=organizations/peerOrganizations/supplier1.steelplatform.com/ca/ca.supplier1.steelplatform.com-cert.pem

echo "$(json_ccp $SUP $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/supplier1.steelplatform.com/connection-supplier1.json
echo "$(yaml_ccp $SUP $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/supplier1.steelplatform.com/connection-supplier1.yaml

SUP=2
P0PORT=9051
CAPORT=8054
PEERPEM=organizations/peerOrganizations/supplier2.steelplatform.com/tlsca/tlsca.supplier2.steelplatform.com-cert.pem
CAPEM=organizations/peerOrganizations/supplier2.steelplatform.com/ca/ca.supplier2.steelplatform.com-cert.pem

echo "$(json_ccp $SUP $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/supplier2.steelplatform.com/connection-supplier2.json
echo "$(yaml_ccp $SUP $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/supplier2.steelplatform.com/connection-supplier2.yaml
