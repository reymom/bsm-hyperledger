#!/bin/bash

function createSupplier1() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/supplier1.steelplatform.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:7054 --caname ca-supplier1 --tls.certfiles "${PWD}/organizations/fabric-ca/supplier1/tls-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-supplier1.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-supplier1.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-supplier1.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-supplier1.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/msp/config.yaml"

  infoln "Registering peer0"
  set -x
  fabric-ca-client register --caname ca-supplier1 --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/supplier1/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-supplier1 --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/supplier1/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-supplier1 --id.name supplier1admin --id.secret supplier1adminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/supplier1/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca-supplier1 -M "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/msp" --csr.hosts peer0.supplier1.steelplatform.com --tls.certfiles "${PWD}/organizations/fabric-ca/supplier1/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca-supplier1 -M "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls" --enrollment.profile tls --csr.hosts peer0.supplier1.steelplatform.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/supplier1/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/server.key"

  mkdir -p "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/msp/tlscacerts"
  cp "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/tlsca"
  cp "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/tlsca/tlsca.supplier1.steelplatform.com-cert.pem"

  mkdir -p "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/ca"
  cp "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/peers/peer0.supplier1.steelplatform.com/msp/cacerts/"* "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/ca/ca.supplier1.steelplatform.com-cert.pem"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:7054 --caname ca-supplier1 -M "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/users/User1@supplier1.steelplatform.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/supplier1/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/users/User1@supplier1.steelplatform.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://supplier1admin:supplier1adminpw@localhost:7054 --caname ca-supplier1 -M "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/users/Admin@supplier1.steelplatform.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/supplier1/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/supplier1.steelplatform.com/users/Admin@supplier1.steelplatform.com/msp/config.yaml"
}

function createSupplier2() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/supplier2.steelplatform.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:8054 --caname ca-supplier2 --tls.certfiles "${PWD}/organizations/fabric-ca/supplier2/tls-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-supplier2.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-supplier2.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-supplier2.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-supplier2.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/msp/config.yaml"

  infoln "Registering peer0"
  set -x
  fabric-ca-client register --caname ca-supplier2 --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/supplier2/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-supplier2 --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/supplier2/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-supplier2 --id.name supplier2admin --id.secret supplier2adminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/supplier2/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname ca-supplier2 -M "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/msp" --csr.hosts peer0.supplier2.steelplatform.com --tls.certfiles "${PWD}/organizations/fabric-ca/supplier2/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname ca-supplier2 -M "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/tls" --enrollment.profile tls --csr.hosts peer0.supplier2.steelplatform.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/supplier2/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/tls/server.key"

  mkdir -p "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/msp/tlscacerts"
  cp "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/tlsca"
  cp "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/tlsca/tlsca.supplier2.steelplatform.com-cert.pem"

  mkdir -p "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/ca"
  cp "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/peers/peer0.supplier2.steelplatform.com/msp/cacerts/"* "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/ca/ca.supplier2.steelplatform.com-cert.pem"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:8054 --caname ca-supplier2 -M "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/users/User1@supplier2.steelplatform.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/supplier2/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/users/User1@supplier2.steelplatform.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://supplier2admin:supplier2adminpw@localhost:8054 --caname ca-supplier2 -M "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/users/Admin@supplier2.steelplatform.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/supplier2/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/supplier2.steelplatform.com/users/Admin@supplier2.steelplatform.com/msp/config.yaml"
}

function createBuyer1() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/buyer1.steelplatform.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:9054 --caname ca-buyer1 --tls.certfiles "${PWD}/organizations/fabric-ca/buyer1/tls-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-buyer1.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-buyer1.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-buyer1.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-buyer1.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/msp/config.yaml"

  infoln "Registering peer0"
  set -x
  fabric-ca-client register --caname ca-buyer1 --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/buyer1/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-buyer1 --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/buyer1/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-buyer1 --id.name buyer1admin --id.secret buyer1adminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/buyer1/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:9054 --caname ca-buyer1 -M "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/msp" --csr.hosts peer0.buyer1.steelplatform.com --tls.certfiles "${PWD}/organizations/fabric-ca/buyer1/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:9054 --caname ca-buyer1 -M "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls" --enrollment.profile tls --csr.hosts peer0.buyer1.steelplatform.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/buyer1/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/server.key"

  mkdir -p "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/msp/tlscacerts"
  cp "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/tlsca"
  cp "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/tlsca/tlsca.buyer1.steelplatform.com-cert.pem"

  mkdir -p "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/ca"
  cp "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/peers/peer0.buyer1.steelplatform.com/msp/cacerts/"* "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/ca/ca.buyer1.steelplatform.com-cert.pem"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:9054 --caname ca-buyer1 -M "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/users/User1@buyer1.steelplatform.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/buyer1/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/users/User1@buyer1.steelplatform.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://buyer1admin:buyer1adminpw@localhost:9054 --caname ca-buyer1 -M "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/users/Admin@buyer1.steelplatform.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/buyer1/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/buyer1.steelplatform.com/users/Admin@buyer1.steelplatform.com/msp/config.yaml"
}

function createBuyer2() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/buyer2.steelplatform.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:10054 --caname ca-buyer2 --tls.certfiles "${PWD}/organizations/fabric-ca/buyer2/tls-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-10054-ca-buyer2.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-10054-ca-buyer2.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-10054-ca-buyer2.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-10054-ca-buyer2.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/msp/config.yaml"

  infoln "Registering peer0"
  set -x
  fabric-ca-client register --caname ca-buyer2 --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/buyer2/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-buyer2 --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/buyer2/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-buyer2 --id.name buyer2admin --id.secret buyer2adminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/buyer2/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:10054 --caname ca-buyer2 -M "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/msp" --csr.hosts peer0.buyer2.steelplatform.com --tls.certfiles "${PWD}/organizations/fabric-ca/buyer2/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:10054 --caname ca-buyer2 -M "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls" --enrollment.profile tls --csr.hosts peer0.buyer2.steelplatform.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/buyer2/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/server.key"

  mkdir -p "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/msp/tlscacerts"
  cp "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/tlsca"
  cp "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/tlsca/tlsca.buyer2.steelplatform.com-cert.pem"

  mkdir -p "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/ca"
  cp "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/peers/peer0.buyer2.steelplatform.com/msp/cacerts/"* "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/ca/ca.buyer2.steelplatform.com-cert.pem"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:10054 --caname ca-buyer2 -M "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/users/User1@buyer2.steelplatform.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/buyer2/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/users/User1@buyer2.steelplatform.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://buyer2admin:buyer2adminpw@localhost:10054 --caname ca-buyer2 -M "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/users/Admin@buyer2.steelplatform.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/buyer2/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/buyer2.steelplatform.com/users/Admin@buyer2.steelplatform.com/msp/config.yaml"
}

function createBuyer3() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/buyer3.steelplatform.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:11054 --caname ca-buyer3 --tls.certfiles "${PWD}/organizations/fabric-ca/buyer3/tls-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-11054-ca-buyer3.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-11054-ca-buyer3.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-11054-ca-buyer3.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-11054-ca-buyer3.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/msp/config.yaml"

  infoln "Registering peer0"
  set -x
  fabric-ca-client register --caname ca-buyer3 --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/buyer3/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-buyer3 --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/buyer3/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-buyer3 --id.name buyer3admin --id.secret buyer3adminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/buyer3/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:11054 --caname ca-buyer3 -M "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/msp" --csr.hosts peer0.buyer3.steelplatform.com --tls.certfiles "${PWD}/organizations/fabric-ca/buyer3/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:11054 --caname ca-buyer3 -M "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/tls" --enrollment.profile tls --csr.hosts peer0.buyer3.steelplatform.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/buyer3/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/tls/server.key"

  mkdir -p "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/msp/tlscacerts"
  cp "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/tlsca"
  cp "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/tlsca/tlsca.buyer3.steelplatform.com-cert.pem"

  mkdir -p "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/ca"
  cp "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/peers/peer0.buyer3.steelplatform.com/msp/cacerts/"* "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/ca/ca.buyer3.steelplatform.com-cert.pem"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:11054 --caname ca-buyer3 -M "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/users/User1@buyer3.steelplatform.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/buyer3/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/users/User1@buyer3.steelplatform.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://buyer3admin:buyer3adminpw@localhost:11054 --caname ca-buyer3 -M "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/users/Admin@buyer3.steelplatform.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/buyer3/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/buyer3.steelplatform.com/users/Admin@buyer3.steelplatform.com/msp/config.yaml"
}

function createLogistics() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/logistics.steelplatform.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:12054 --caname ca-logistics --tls.certfiles "${PWD}/organizations/fabric-ca/logistics/tls-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-12054-ca-logistics.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-12054-ca-logistics.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-12054-ca-logistics.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-12054-ca-logistics.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/msp/config.yaml"

  infoln "Registering peer0"
  set -x
  fabric-ca-client register --caname ca-logistics --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/logistics/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-logistics --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/logistics/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-logistics --id.name logisticsadmin --id.secret logisticsadminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/logistics/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:12054 --caname ca-logistics -M "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/msp" --csr.hosts peer0.logistics.steelplatform.com --tls.certfiles "${PWD}/organizations/fabric-ca/logistics/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:12054 --caname ca-logistics -M "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/tls" --enrollment.profile tls --csr.hosts peer0.logistics.steelplatform.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/logistics/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/tls/server.key"

  mkdir -p "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/msp/tlscacerts"
  cp "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/tlsca"
  cp "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/tlsca/tlsca.logistics.steelplatform.com-cert.pem"

  mkdir -p "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/ca"
  cp "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/peers/peer0.logistics.steelplatform.com/msp/cacerts/"* "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/ca/ca.logistics.steelplatform.com-cert.pem"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:12054 --caname ca-logistics -M "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/users/User1@logistics.steelplatform.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/logistics/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/users/User1@logistics.steelplatform.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://logisticsadmin:logisticsadminpw@localhost:12054 --caname ca-logistics -M "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/users/Admin@logistics.steelplatform.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/logistics/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/logistics.steelplatform.com/users/Admin@logistics.steelplatform.com/msp/config.yaml"
}

function createOrderer() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/ordererOrganizations/steelplatform.com

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/ordererOrganizations/steelplatform.com

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:13054 --caname ca-orderer --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-13054-ca-orderer.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-13054-ca-orderer.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-13054-ca-orderer.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-13054-ca-orderer.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/ordererOrganizations/steelplatform.com/msp/config.yaml"

  infoln "Registering orderer"
  set -x
  fabric-ca-client register --caname ca-orderer --id.name orderer --id.secret ordererpw --id.type orderer --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the orderer admin"
  set -x
  fabric-ca-client register --caname ca-orderer --id.name ordererAdmin --id.secret ordererAdminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the orderer msp"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:13054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp" --csr.hosts orderer.steelplatform.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/steelplatform.com/msp/config.yaml" "${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/config.yaml"

  infoln "Generating the orderer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:13054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/tls" --enrollment.profile tls --csr.hosts orderer.steelplatform.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/tls/ca.crt"
  cp "${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/tls/signcerts/"* "${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/tls/server.crt"
  cp "${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/tls/keystore/"* "${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/tls/server.key"

  mkdir -p "${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts"
  cp "${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem"

  mkdir -p "${PWD}/organizations/ordererOrganizations/steelplatform.com/msp/tlscacerts"
  cp "${PWD}/organizations/ordererOrganizations/steelplatform.com/orderers/orderer.steelplatform.com/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/steelplatform.com/msp/tlscacerts/tlsca.steelplatform.com-cert.pem"

  infoln "Generating the admin msp"
  set -x
  fabric-ca-client enroll -u https://ordererAdmin:ordererAdminpw@localhost:13054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/steelplatform.com/users/Admin@steelplatform.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/steelplatform.com/msp/config.yaml" "${PWD}/organizations/ordererOrganizations/steelplatform.com/users/Admin@steelplatform.com/msp/config.yaml"
}
