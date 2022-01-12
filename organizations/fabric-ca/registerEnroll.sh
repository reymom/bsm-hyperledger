#!/bin/bash

function createSupplier1() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/supplier1.example.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/supplier1.example.com/

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
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/supplier1.example.com/msp/config.yaml"

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
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca-supplier1 -M "${PWD}/organizations/peerOrganizations/supplier1.example.com/peers/peer0.supplier1.example.com/msp" --csr.hosts peer0.supplier1.example.com --tls.certfiles "${PWD}/organizations/fabric-ca/supplier1/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/supplier1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/supplier1.example.com/peers/peer0.supplier1.example.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca-supplier1 -M "${PWD}/organizations/peerOrganizations/supplier1.example.com/peers/peer0.supplier1.example.com/tls" --enrollment.profile tls --csr.hosts peer0.supplier1.example.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/supplier1/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/supplier1.example.com/peers/peer0.supplier1.example.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/supplier1.example.com/peers/peer0.supplier1.example.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/supplier1.example.com/peers/peer0.supplier1.example.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/supplier1.example.com/peers/peer0.supplier1.example.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/supplier1.example.com/peers/peer0.supplier1.example.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/supplier1.example.com/peers/peer0.supplier1.example.com/tls/server.key"

  mkdir -p "${PWD}/organizations/peerOrganizations/supplier1.example.com/msp/tlscacerts"
  cp "${PWD}/organizations/peerOrganizations/supplier1.example.com/peers/peer0.supplier1.example.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/supplier1.example.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/peerOrganizations/supplier1.example.com/tlsca"
  cp "${PWD}/organizations/peerOrganizations/supplier1.example.com/peers/peer0.supplier1.example.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/supplier1.example.com/tlsca/tlsca.supplier1.example.com-cert.pem"

  mkdir -p "${PWD}/organizations/peerOrganizations/supplier1.example.com/ca"
  cp "${PWD}/organizations/peerOrganizations/supplier1.example.com/peers/peer0.supplier1.example.com/msp/cacerts/"* "${PWD}/organizations/peerOrganizations/supplier1.example.com/ca/ca.supplier1.example.com-cert.pem"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:7054 --caname ca-supplier1 -M "${PWD}/organizations/peerOrganizations/supplier1.example.com/users/User1@supplier1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/supplier1/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/supplier1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/supplier1.example.com/users/User1@supplier1.example.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://supplier1admin:supplier1adminpw@localhost:7054 --caname ca-supplier1 -M "${PWD}/organizations/peerOrganizations/supplier1.example.com/users/Admin@supplier1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/supplier1/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/supplier1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/supplier1.example.com/users/Admin@supplier1.example.com/msp/config.yaml"
}

function createSupplier2() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/supplier2.example.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/supplier2.example.com/

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
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/supplier2.example.com/msp/config.yaml"

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
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname ca-supplier2 -M "${PWD}/organizations/peerOrganizations/supplier2.example.com/peers/peer0.supplier2.example.com/msp" --csr.hosts peer0.supplier2.example.com --tls.certfiles "${PWD}/organizations/fabric-ca/supplier2/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/supplier2.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/supplier2.example.com/peers/peer0.supplier2.example.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname ca-supplier2 -M "${PWD}/organizations/peerOrganizations/supplier2.example.com/peers/peer0.supplier2.example.com/tls" --enrollment.profile tls --csr.hosts peer0.supplier2.example.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/supplier2/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/supplier2.example.com/peers/peer0.supplier2.example.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/supplier2.example.com/peers/peer0.supplier2.example.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/supplier2.example.com/peers/peer0.supplier2.example.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/supplier2.example.com/peers/peer0.supplier2.example.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/supplier2.example.com/peers/peer0.supplier2.example.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/supplier2.example.com/peers/peer0.supplier2.example.com/tls/server.key"

  mkdir -p "${PWD}/organizations/peerOrganizations/supplier2.example.com/msp/tlscacerts"
  cp "${PWD}/organizations/peerOrganizations/supplier2.example.com/peers/peer0.supplier2.example.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/supplier2.example.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/peerOrganizations/supplier2.example.com/tlsca"
  cp "${PWD}/organizations/peerOrganizations/supplier2.example.com/peers/peer0.supplier2.example.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/supplier2.example.com/tlsca/tlsca.supplier2.example.com-cert.pem"

  mkdir -p "${PWD}/organizations/peerOrganizations/supplier2.example.com/ca"
  cp "${PWD}/organizations/peerOrganizations/supplier2.example.com/peers/peer0.supplier2.example.com/msp/cacerts/"* "${PWD}/organizations/peerOrganizations/supplier2.example.com/ca/ca.supplier2.example.com-cert.pem"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:8054 --caname ca-supplier2 -M "${PWD}/organizations/peerOrganizations/supplier2.example.com/users/User1@supplier2.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/supplier2/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/supplier2.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/supplier2.example.com/users/User1@supplier2.example.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://supplier2admin:supplier2adminpw@localhost:8054 --caname ca-supplier2 -M "${PWD}/organizations/peerOrganizations/supplier2.example.com/users/Admin@supplier2.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/supplier2/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/supplier2.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/supplier2.example.com/users/Admin@supplier2.example.com/msp/config.yaml"
}

function createBuyer1() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/buyer1.example.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/buyer1.example.com/

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
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/buyer1.example.com/msp/config.yaml"

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
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:9054 --caname ca-buyer1 -M "${PWD}/organizations/peerOrganizations/buyer1.example.com/peers/peer0.buyer1.example.com/msp" --csr.hosts peer0.buyer1.example.com --tls.certfiles "${PWD}/organizations/fabric-ca/buyer1/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/buyer1.example.com/peers/peer0.buyer1.example.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:9054 --caname ca-buyer1 -M "${PWD}/organizations/peerOrganizations/buyer1.example.com/peers/peer0.buyer1.example.com/tls" --enrollment.profile tls --csr.hosts peer0.buyer1.example.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/buyer1/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer1.example.com/peers/peer0.buyer1.example.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/buyer1.example.com/peers/peer0.buyer1.example.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/buyer1.example.com/peers/peer0.buyer1.example.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/buyer1.example.com/peers/peer0.buyer1.example.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/buyer1.example.com/peers/peer0.buyer1.example.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/buyer1.example.com/peers/peer0.buyer1.example.com/tls/server.key"

  mkdir -p "${PWD}/organizations/peerOrganizations/buyer1.example.com/msp/tlscacerts"
  cp "${PWD}/organizations/peerOrganizations/buyer1.example.com/peers/peer0.buyer1.example.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/buyer1.example.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/peerOrganizations/buyer1.example.com/tlsca"
  cp "${PWD}/organizations/peerOrganizations/buyer1.example.com/peers/peer0.buyer1.example.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/buyer1.example.com/tlsca/tlsca.buyer1.example.com-cert.pem"

  mkdir -p "${PWD}/organizations/peerOrganizations/buyer1.example.com/ca"
  cp "${PWD}/organizations/peerOrganizations/buyer1.example.com/peers/peer0.buyer1.example.com/msp/cacerts/"* "${PWD}/organizations/peerOrganizations/buyer1.example.com/ca/ca.buyer1.example.com-cert.pem"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:9054 --caname ca-buyer1 -M "${PWD}/organizations/peerOrganizations/buyer1.example.com/users/User1@buyer1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/buyer1/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/buyer1.example.com/users/User1@buyer1.example.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://buyer1admin:buyer1adminpw@localhost:9054 --caname ca-buyer1 -M "${PWD}/organizations/peerOrganizations/buyer1.example.com/users/Admin@buyer1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/buyer1/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/buyer1.example.com/users/Admin@buyer1.example.com/msp/config.yaml"
}

function createBuyer2() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/buyer2.example.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/buyer2.example.com/

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
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/buyer2.example.com/msp/config.yaml"

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
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:10054 --caname ca-buyer2 -M "${PWD}/organizations/peerOrganizations/buyer2.example.com/peers/peer0.buyer2.example.com/msp" --csr.hosts peer0.buyer2.example.com --tls.certfiles "${PWD}/organizations/fabric-ca/buyer2/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer2.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/buyer2.example.com/peers/peer0.buyer2.example.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:10054 --caname ca-buyer2 -M "${PWD}/organizations/peerOrganizations/buyer2.example.com/peers/peer0.buyer2.example.com/tls" --enrollment.profile tls --csr.hosts peer0.buyer2.example.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/buyer2/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer2.example.com/peers/peer0.buyer2.example.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/buyer2.example.com/peers/peer0.buyer2.example.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/buyer2.example.com/peers/peer0.buyer2.example.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/buyer2.example.com/peers/peer0.buyer2.example.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/buyer2.example.com/peers/peer0.buyer2.example.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/buyer2.example.com/peers/peer0.buyer2.example.com/tls/server.key"

  mkdir -p "${PWD}/organizations/peerOrganizations/buyer2.example.com/msp/tlscacerts"
  cp "${PWD}/organizations/peerOrganizations/buyer2.example.com/peers/peer0.buyer2.example.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/buyer2.example.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/peerOrganizations/buyer2.example.com/tlsca"
  cp "${PWD}/organizations/peerOrganizations/buyer2.example.com/peers/peer0.buyer2.example.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/buyer2.example.com/tlsca/tlsca.buyer2.example.com-cert.pem"

  mkdir -p "${PWD}/organizations/peerOrganizations/buyer2.example.com/ca"
  cp "${PWD}/organizations/peerOrganizations/buyer2.example.com/peers/peer0.buyer2.example.com/msp/cacerts/"* "${PWD}/organizations/peerOrganizations/buyer2.example.com/ca/ca.buyer2.example.com-cert.pem"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:10054 --caname ca-buyer2 -M "${PWD}/organizations/peerOrganizations/buyer2.example.com/users/User1@buyer2.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/buyer2/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer2.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/buyer2.example.com/users/User1@buyer2.example.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://buyer2admin:buyer2adminpw@localhost:10054 --caname ca-buyer2 -M "${PWD}/organizations/peerOrganizations/buyer2.example.com/users/Admin@buyer2.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/buyer2/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer2.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/buyer2.example.com/users/Admin@buyer2.example.com/msp/config.yaml"
}

function createBuyer3() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/buyer3.example.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/buyer3.example.com/

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
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/buyer3.example.com/msp/config.yaml"

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
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:11054 --caname ca-buyer3 -M "${PWD}/organizations/peerOrganizations/buyer3.example.com/peers/peer0.buyer3.example.com/msp" --csr.hosts peer0.buyer3.example.com --tls.certfiles "${PWD}/organizations/fabric-ca/buyer3/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer3.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/buyer3.example.com/peers/peer0.buyer3.example.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:11054 --caname ca-buyer3 -M "${PWD}/organizations/peerOrganizations/buyer3.example.com/peers/peer0.buyer3.example.com/tls" --enrollment.profile tls --csr.hosts peer0.buyer3.example.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/buyer3/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer3.example.com/peers/peer0.buyer3.example.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/buyer3.example.com/peers/peer0.buyer3.example.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/buyer3.example.com/peers/peer0.buyer3.example.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/buyer3.example.com/peers/peer0.buyer3.example.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/buyer3.example.com/peers/peer0.buyer3.example.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/buyer3.example.com/peers/peer0.buyer3.example.com/tls/server.key"

  mkdir -p "${PWD}/organizations/peerOrganizations/buyer3.example.com/msp/tlscacerts"
  cp "${PWD}/organizations/peerOrganizations/buyer3.example.com/peers/peer0.buyer3.example.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/buyer3.example.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/peerOrganizations/buyer3.example.com/tlsca"
  cp "${PWD}/organizations/peerOrganizations/buyer3.example.com/peers/peer0.buyer3.example.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/buyer3.example.com/tlsca/tlsca.buyer3.example.com-cert.pem"

  mkdir -p "${PWD}/organizations/peerOrganizations/buyer3.example.com/ca"
  cp "${PWD}/organizations/peerOrganizations/buyer3.example.com/peers/peer0.buyer3.example.com/msp/cacerts/"* "${PWD}/organizations/peerOrganizations/buyer3.example.com/ca/ca.buyer3.example.com-cert.pem"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:11054 --caname ca-buyer3 -M "${PWD}/organizations/peerOrganizations/buyer3.example.com/users/User1@buyer3.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/buyer3/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer3.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/buyer3.example.com/users/User1@buyer3.example.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://buyer3admin:buyer3adminpw@localhost:11054 --caname ca-buyer3 -M "${PWD}/organizations/peerOrganizations/buyer3.example.com/users/Admin@buyer3.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/buyer3/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/buyer3.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/buyer3.example.com/users/Admin@buyer3.example.com/msp/config.yaml"
}

function createLogistics() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/logistics.example.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/logistics.example.com/

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
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/logistics.example.com/msp/config.yaml"

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
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:12054 --caname ca-logistics -M "${PWD}/organizations/peerOrganizations/logistics.example.com/peers/peer0.logistics.example.com/msp" --csr.hosts peer0.logistics.example.com --tls.certfiles "${PWD}/organizations/fabric-ca/logistics/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/logistics.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/logistics.example.com/peers/peer0.logistics.example.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:12054 --caname ca-logistics -M "${PWD}/organizations/peerOrganizations/logistics.example.com/peers/peer0.logistics.example.com/tls" --enrollment.profile tls --csr.hosts peer0.logistics.example.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/logistics/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/logistics.example.com/peers/peer0.logistics.example.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/logistics.example.com/peers/peer0.logistics.example.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/logistics.example.com/peers/peer0.logistics.example.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/logistics.example.com/peers/peer0.logistics.example.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/logistics.example.com/peers/peer0.logistics.example.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/logistics.example.com/peers/peer0.logistics.example.com/tls/server.key"

  mkdir -p "${PWD}/organizations/peerOrganizations/logistics.example.com/msp/tlscacerts"
  cp "${PWD}/organizations/peerOrganizations/logistics.example.com/peers/peer0.logistics.example.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/logistics.example.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/peerOrganizations/logistics.example.com/tlsca"
  cp "${PWD}/organizations/peerOrganizations/logistics.example.com/peers/peer0.logistics.example.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/logistics.example.com/tlsca/tlsca.logistics.example.com-cert.pem"

  mkdir -p "${PWD}/organizations/peerOrganizations/logistics.example.com/ca"
  cp "${PWD}/organizations/peerOrganizations/logistics.example.com/peers/peer0.logistics.example.com/msp/cacerts/"* "${PWD}/organizations/peerOrganizations/logistics.example.com/ca/ca.logistics.example.com-cert.pem"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:12054 --caname ca-logistics -M "${PWD}/organizations/peerOrganizations/logistics.example.com/users/User1@logistics.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/logistics/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/logistics.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/logistics.example.com/users/User1@logistics.example.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://logisticsadmin:logisticsadminpw@localhost:12054 --caname ca-logistics -M "${PWD}/organizations/peerOrganizations/logistics.example.com/users/Admin@logistics.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/logistics/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/logistics.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/logistics.example.com/users/Admin@logistics.example.com/msp/config.yaml"
}

function createOrderer() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/ordererOrganizations/example.com

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/ordererOrganizations/example.com

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:9054 --caname ca-orderer --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/ordererOrganizations/example.com/msp/config.yaml"

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
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp" --csr.hosts orderer.example.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/example.com/msp/config.yaml" "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/config.yaml"

  infoln "Generating the orderer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls" --enrollment.profile tls --csr.hosts orderer.example.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt"
  cp "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/signcerts/"* "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt"
  cp "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/keystore/"* "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.key"

  mkdir -p "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts"
  cp "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

  mkdir -p "${PWD}/organizations/ordererOrganizations/example.com/msp/tlscacerts"
  cp "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

  infoln "Generating the admin msp"
  set -x
  fabric-ca-client enroll -u https://ordererAdmin:ordererAdminpw@localhost:9054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/example.com/users/Admin@example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/example.com/msp/config.yaml" "${PWD}/organizations/ordererOrganizations/example.com/users/Admin@example.com/msp/config.yaml"
}
