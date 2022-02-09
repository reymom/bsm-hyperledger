#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp_sup {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${SUP}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        organizations/ccp-template-suppliers.json
}

function yaml_ccp_sup {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${SUP}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        organizations/ccp-template-suppliers.yaml | sed -e $'s/\\\\n/\\\n          /g'
}

function json_ccp_buy {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${BUY}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        organizations/ccp-template-buyers.json
}

function yaml_ccp_buy {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${BUY}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        organizations/ccp-template-buyers.yaml | sed -e $'s/\\\\n/\\\n          /g'
}

function json_ccp_logi {
    local PP=$(one_line_pem $3)
    local CP=$(one_line_pem $4)
    sed -e "s/\${P0PORT}/$1/" \
        -e "s/\${CAPORT}/$2/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        organizations/ccp-template-logistics.json
}

function yaml_ccp_logi {
    local PP=$(one_line_pem $3)
    local CP=$(one_line_pem $4)
    sed -e "s/\${P0PORT}/$1/" \
        -e "s/\${CAPORT}/$2/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        organizations/ccp-template-logistics.yaml | sed -e $'s/\\\\n/\\\n          /g'
}

SUP=1
P0PORT=7051
CAPORT=7054
PEERPEM=organizations/peerOrganizations/supplier1.steelplatform.com/tlsca/tlsca.supplier1.steelplatform.com-cert.pem
CAPEM=organizations/peerOrganizations/supplier1.steelplatform.com/ca/ca.supplier1.steelplatform.com-cert.pem

echo "$(json_ccp_sup $SUP $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/supplier1.steelplatform.com/connection-supplier1.json
echo "$(yaml_ccp_sup $SUP $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/supplier1.steelplatform.com/connection-supplier1.yaml

SUP=2
P0PORT=9051
CAPORT=8054
PEERPEM=organizations/peerOrganizations/supplier2.steelplatform.com/tlsca/tlsca.supplier2.steelplatform.com-cert.pem
CAPEM=organizations/peerOrganizations/supplier2.steelplatform.com/ca/ca.supplier2.steelplatform.com-cert.pem

echo "$(json_ccp_sup $SUP $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/supplier2.steelplatform.com/connection-supplier2.json
echo "$(yaml_ccp_sup $SUP $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/supplier2.steelplatform.com/connection-supplier2.yaml

BUY=1
P0PORT=11051
CAPORT=9054
PEERPEM=organizations/peerOrganizations/buyer1.steelplatform.com/tlsca/tlsca.buyer1.steelplatform.com-cert.pem
CAPEM=organizations/peerOrganizations/buyer1.steelplatform.com/ca/ca.buyer1.steelplatform.com-cert.pem

echo "$(json_ccp_buy $BUY $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/buyer1.steelplatform.com/connection-buyer1.json
echo "$(yaml_ccp_buy $BUY $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/buyer1.steelplatform.com/connection-buyer1.yaml

BUY=2
P0PORT=13051
CAPORT=10054
PEERPEM=organizations/peerOrganizations/buyer2.steelplatform.com/tlsca/tlsca.buyer2.steelplatform.com-cert.pem
CAPEM=organizations/peerOrganizations/buyer2.steelplatform.com/ca/ca.buyer2.steelplatform.com-cert.pem

echo "$(json_ccp_buy $BUY $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/buyer2.steelplatform.com/connection-buyer2.json
echo "$(yaml_ccp_buy $BUY $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/buyer2.steelplatform.com/connection-buyer2.yaml

BUY=3
P0PORT=15051
CAPORT=11054
PEERPEM=organizations/peerOrganizations/buyer3.steelplatform.com/tlsca/tlsca.buyer3.steelplatform.com-cert.pem
CAPEM=organizations/peerOrganizations/buyer3.steelplatform.com/ca/ca.buyer3.steelplatform.com-cert.pem

echo "$(json_ccp_buy $BUY $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/buyer3.steelplatform.com/connection-buyer3.json
echo "$(yaml_ccp_buy $BUY $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/buyer3.steelplatform.com/connection-buyer3.yaml

P0PORT=17051
CAPORT=12054
PEERPEM=organizations/peerOrganizations/logistics.steelplatform.com/tlsca/tlsca.logistics.steelplatform.com-cert.pem
CAPEM=organizations/peerOrganizations/logistics.steelplatform.com/ca/ca.logistics.steelplatform.com-cert.pem

echo "$(json_ccp_logi $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/logistics.steelplatform.com/connection-logistics.json
echo "$(yaml_ccp_logi $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/logistics.steelplatform.com/connection-logistics.yaml