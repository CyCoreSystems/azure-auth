#!/bin/bash
set -e

pushd cmd/azure-token
go install
popd


pushd cmd/send-azure
go install
popd
