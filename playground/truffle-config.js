module.exports = {
  // See <http://truffleframework.com/docs/advanced/configuration>
  // to customize your Truffle configuration!
  development: {
    host: '127.0.0.1',
    port: 23889, // janus QTUM-ETH RPC bridge
    network_id: '*', // eslint-disable-line camelcase
    from: '0x7926223070547d2d15b2ef5e7383e541c338ffe9',
    gasPrice: '0x64', // minimal gas for qtum
  },
};
