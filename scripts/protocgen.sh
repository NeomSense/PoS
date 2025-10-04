#!/usr/bin/env bash

set -e

echo "Generating gogo proto code"
cd proto
proto_dirs=$(find . -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    echo "Processing $file"
    protoc \
      -I "." \
      -I "../" \
      --gocosmos_out=../ \
      --grpc-gateway_out=../ \
      --grpc-gateway_opt=logtostderr=true \
      --grpc-gateway_opt=allow_repeated_fields_in_body=true \
      "$file"
  done
done

cd ..
echo "Proto generation complete!"
