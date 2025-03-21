#!/bin/bash
cd ./proto
buf generate --template buf.gen.gogo.yaml
buf generate --template buf.gen.pulsar.yaml
cd ..

cp -r autocctp.dev/* ./
cp -r api/noble/autocctp/* api/

rm -rf autocctp.dev
rm -rf api/noble
rm -rf noble
