#!/bin/sh
tar -xzvf viacoin-0.16.3-x86_64-linux-gnu.tar.gz
cd viacoin-0.16.3/bin/
cp viacoind /usr/bin/viacoind
cp viacoin-cli /usr/bin/viacoin-cli

cd ~
if [ ! -d ".viacoin" ]; then
    mkdir .viacoin
  # Control will enter here if $DIRECTORY doesn't exist.
fi


cd ~/.viacoin
CONF=viacoin.conf
if [ ! -d $CONF ]; then
    rm $CONF
  # Control will enter here if $DIRECTORY doesn't exist.
fi
touch $CONF
echo "server=1" >> $CONF
echo "txindex=1" >> $CONF
echo "rpcallowip=*" >> $CONF
echo "rpcuser=via" >> $CONF
echo "rpcpassword=via" >> $CONF