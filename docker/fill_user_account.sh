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
repeat_until_success qtum-cli -rpcuser=qtum -rpcpassword=testpasswd importprivkey "cMbgxCJrTYUqgcmiC1berh5DFrtY1KeU4PXZ6NZxgenniF1mXCRk" address1 # addr=qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW hdkeypath=m/88'/0'/1'
repeat_until_success qtum-cli -rpcuser=qtum -rpcpassword=testpasswd importprivkey "cRcG1jizfBzHxfwu68aMjhy78CpnzD9gJYZ5ggDbzfYD3EQfGUDZ" address2 # addr=qLn9vqbr2Gx3TsVR9QyTVB5mrMoh4x43Uf hdkeypath=m/88'/0'/2'
repeat_until_success qtum-cli -rpcuser=qtum -rpcpassword=testpasswd importprivkey "cV79qBoCSA2NDrJz8S3T7J8f3zgkGfg4ua4hRRXfhbnq5VhXkukT" address3 # addr=qdbfjUGtAei3uU13mFiDUZz1e8CWYnawLh
repeat_until_success qtum-cli -rpcuser=qtum -rpcpassword=testpasswd importprivkey "cV93kaaV8hvNqZ711s2z9jVWLYEtwwsVpyFeEZCP6otiZgrCTiEW" address4 # addr=qUXNJHShdBHzNc53zW9KeGjPKYZbEYC3gB
repeat_until_success qtum-cli -rpcuser=qtum -rpcpassword=testpasswd importprivkey "cVPHpTvmv3UjQsZfsMRrW5RrGCyTSAZ3MWs1f8R1VeKJSYxy5uac" address5 # addr=qUXZFWC99kG1jYdJepgHyAwS97uAahbRPX
repeat_until_success qtum-cli -rpcuser=qtum -rpcpassword=testpasswd importprivkey "cTs5NqY4Ko9o6FESHGBDEG77qqz9me7cyYCoinHcWEiqMZgLC6XY" address6 # addr=qPA8MJH28xgbhV8fQaY1q7VNymp8zokzYv
echo Finished importing accounts
echo Seeding accounts
echo Seeding qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW
repeat_until_success qtum-cli -rpcuser=qtum -rpcpassword=testpasswd generatetoaddress 1000 qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW
echo Seeding qLn9vqbr2Gx3TsVR9QyTVB5mrMoh4x43Uf
repeat_until_success qtum-cli -rpcuser=qtum -rpcpassword=testpasswd generatetoaddress 1000 qLn9vqbr2Gx3TsVR9QyTVB5mrMoh4x43Uf
echo Seeding qdbfjUGtAei3uU13mFiDUZz1e8CWYnawLh
repeat_until_success qtum-cli -rpcuser=qtum -rpcpassword=testpasswd generatetoaddress 500 qdbfjUGtAei3uU13mFiDUZz1e8CWYnawLh
echo Seeding qUXNJHShdBHzNc53zW9KeGjPKYZbEYC3gB
repeat_until_success qtum-cli -rpcuser=qtum -rpcpassword=testpasswd generatetoaddress 500 qUXNJHShdBHzNc53zW9KeGjPKYZbEYC3gB
echo Seeding qUXZFWC99kG1jYdJepgHyAwS97uAahbRPX
repeat_until_success qtum-cli -rpcuser=qtum -rpcpassword=testpasswd generatetoaddress 500 qUXZFWC99kG1jYdJepgHyAwS97uAahbRPX
echo Seeding qPA8MJH28xgbhV8fQaY1q7VNymp8zokzYv
repeat_until_success qtum-cli -rpcuser=qtum -rpcpassword=testpasswd generatetoaddress 500 qPA8MJH28xgbhV8fQaY1q7VNymp8zokzYv
echo Finished importing and seeding accounts