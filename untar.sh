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
echo "daemon=1" >> $CONF
echo "rpcallowip=0.0.0.0/24" >> $CONF
echo "rpcuser=via" >> $CONF
echo "rpcpassword=via" >> $CONF
echo "addnode=118.209.110.211" >> $CONF
echo "addnode=118.209.111.44" >> $CONF
echo "addnode=184.95.48.202" >> $CONF
echo "addnode=24.50.182.181" >> $CONF
echo "addnode=81.187.211.49" >> $CONF