#!/bin/bash
CURRENT_DIR=$1

for x in $(find ${CURRENT_DIR}/TravelTales_GrpcProto/* -type d); do
  protoc -I=${x} -I=${CURRENT_DIR}/grpc-proto/ -I /usr/local/go --go_out=${CURRENT_DIR} \
   --go-grpc_out=${CURRENT_DIR} ${x}/*.proto
done