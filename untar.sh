#!/bin/sh
tar -xzvf viacoin-0.16.3-x86_64-linux-gnu.tar.gz
cd viacoin-0.16.3/bin/
cp viacoind /usr/bin/viacoind
cp viacoin-cli /usr/bin/viacoin-cli