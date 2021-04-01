#!/bin/sh
repeat_until_success () {
    echo Running command - "$@"
    i=0
    until $@
    do
        echo Command failed with exit code - $?
        if [ $i -gt 10 ]; then
            echo Giving up running command - "$@"
            return
        fi
        sleep 1
        echo Retrying
        i=`expr $i + 1`
    done
    echo Command finished successfully
}

#import private keys and then prefund them
repeat_until_success qtum-cli -rpcuser=qtum -rpcpassword=testpasswd importprivkey "cMbgxCJrTYUqgcmiC1berh5DFrtY1KeU4PXZ6NZxgenniF1mXCRk" # addr=qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW hdkeypath=m/88'/0'/1'
repeat_until_success qtum-cli -rpcuser=qtum -rpcpassword=testpasswd importprivkey "cRcG1jizfBzHxfwu68aMjhy78CpnzD9gJYZ5ggDbzfYD3EQfGUDZ" # addr=qLn9vqbr2Gx3TsVR9QyTVB5mrMoh4x43Uf hdkeypath=m/88'/0'/2'
repeat_until_success qtum-cli -rpcuser=qtum -rpcpassword=testpasswd generatetoaddress 50000 qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW
repeat_until_success qtum-cli -rpcuser=qtum -rpcpassword=testpasswd generatetoaddress 50000 qLn9vqbr2Gx3TsVR9QyTVB5mrMoh4x43Uf
echo Finished importing and seeding accounts