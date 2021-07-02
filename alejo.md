## latest Janus fix

Add `Logs: []eth.Log{}` attribute in  dummy tx receipt for non vm tx in eth_getTransactionReceipt

```javascript
		return &eth.GetTransactionReceiptResponse{
			TransactionHash:   ethTx.Hash,
			TransactionIndex:  ethTx.TransactionIndex,
			BlockHash:         ethTx.BlockHash,
			BlockNumber:       ethTx.BlockNumber,
			CumulativeGasUsed: "0x0",
			GasUsed:           "0x0",
			From:              ethTx.From,
			To:                ethTx.To,
			Logs:		   []eth.Log{},
			LogsBloom:         eth.EmptyLogsBloom,
			Status:            "0x0",
```


## current status
-After the latest Janus fix, the graph starts ingesting blocks succesfully
-A new error `invalid address` is displayed by theGraph:

```bash
2021-07-02T23:47:23.167671781Z Jul 02 23:47:23.165 INFO Graph Node version: 0.22.0 (2021-02-24)
2021-07-02T23:47:23.167705680Z Jul 02 23:47:23.167 INFO Generating configuration from command line arguments
2021-07-02T23:47:23.221343142Z Jul 02 23:47:23.221 INFO Starting up
2021-07-02T23:47:23.224479940Z Jul 02 23:47:23.224 INFO Trying IPFS node at: http://ipfs:5001/
2021-07-02T23:47:23.252864957Z Jul 02 23:47:23.252 INFO Creating transport, capabilities: archive, trace, url: http://172.25.0.1:23889, network: mainnet
2021-07-02T23:47:23.255988028Z Jul 02 23:47:23.255 INFO Successfully connected to IPFS node at: http://ipfs:5001/
2021-07-02T23:47:23.290956812Z Jul 02 23:47:23.290 INFO Connecting to Postgres, weight: 1, conn_pool_size: 10, url: postgresql://graph-node:HIDDEN_PASSWORD@postgres:5432/graph-node, pool: main, shard: primary
2021-07-02T23:47:23.306665922Z Jul 02 23:47:23.306 INFO Pool successfully connected to Postgres, pool: main, shard: primary, component: Store
2021-07-02T23:47:23.310327295Z Jul 02 23:47:23.310 INFO Waiting for other graph-node instances to finish migrating, shard: primary, component: Store
2021-07-02T23:47:23.312815895Z Jul 02 23:47:23.312 INFO Running migrations, shard: primary, component: Store
2021-07-02T23:47:24.021201226Z Jul 02 23:47:24.021 INFO Migrations finished, shard: primary, component: Store
2021-07-02T23:47:24.022567651Z Jul 02 23:47:24.022 INFO Connecting to Ethereum..., capabilities: archive, trace, network: mainnet
2021-07-02T23:47:24.218315886Z Jul 02 23:47:24.218 INFO Connected to Ethereum, capabilities: archive, trace, network_version: 81, network: mainnet
2021-07-02T23:47:24.248076361Z Jul 02 23:47:24.247 INFO Creating LoadManager in disabled mode, component: LoadManager
2021-07-02T23:47:24.248087228Z Jul 02 23:47:24.247 INFO Starting block ingestors
2021-07-02T23:47:24.248091208Z Jul 02 23:47:24.248 INFO Starting block ingestor for network, network_name: mainnet
2021-07-02T23:47:24.249219458Z Jul 02 23:47:24.249 INFO Starting JSON-RPC admin server at: http://localhost:8020, component: JsonRpcServer
2021-07-02T23:47:24.249839017Z Jul 02 23:47:24.249 INFO Started all subgraphs, component: SubgraphRegistrar
2021-07-02T23:47:24.250401222Z Jul 02 23:47:24.250 INFO Starting GraphQL HTTP server at: http://localhost:8000, component: GraphQLServer
2021-07-02T23:47:24.250413790Z Jul 02 23:47:24.250 INFO Starting index node server at: http://localhost:8030, component: IndexNodeServer
2021-07-02T23:47:24.252613982Z Jul 02 23:47:24.252 INFO Starting metrics server at: http://localhost:8040, component: MetricsServer
2021-07-02T23:47:24.252621949Z Jul 02 23:47:24.252 INFO Starting GraphQL WebSocket server at: ws://localhost:8001, component: SubscriptionServer
2021-07-02T23:47:24.369661683Z Jul 02 23:47:24.369 INFO Downloading latest blocks from Ethereum. This may take a few minutes..., network_name: mainnet, component: BlockIngestor
2021-07-02T23:47:30.074965312Z Jul 02 23:47:30.074 INFO Syncing 50 blocks from Ethereum., code: BlockIngestionLagging, blocks_needed: 50, blocks_behind: 1531, latest_block_head: 10063, current_block_head: 8532, network_name: mainnet, component: BlockInge
stor
2021-07-02T23:47:33.661108479Z Jul 02 23:47:33.660 INFO Syncing 50 blocks from Ethereum., code: BlockIngestionLagging, blocks_needed: 50, blocks_behind: 1120, latest_block_head: 11183, current_block_head: 10063, network_name: mainnet, component: BlockIng
estor
2021-07-02T23:47:37.921647268Z Jul 02 23:47:37.921 INFO Syncing 50 blocks from Ethereum., code: BlockIngestionLagging, blocks_needed: 50, blocks_behind: 1379, latest_block_head: 12562, current_block_head: 11183, network_name: mainnet, component: BlockIng
estor
2021-07-02T23:47:41.603372334Z Jul 02 23:47:41.603 INFO Syncing 50 blocks from Ethereum., code: BlockIngestionLagging, blocks_needed: 50, blocks_behind: 1326, latest_block_head: 13888, current_block_head: 12562, network_name: mainnet, component: BlockIng
estor
2021-07-02T23:47:45.039407600Z Jul 02 23:47:45.039 INFO Syncing 50 blocks from Ethereum., code: BlockIngestionLagging, blocks_needed: 50, blocks_behind: 1186, latest_block_head: 15074, current_block_head: 13888, network_name: mainnet, component: BlockIng
estor
2021-07-02T23:47:48.689598108Z Jul 02 23:47:48.689 INFO Syncing 50 blocks from Ethereum., code: BlockIngestionLagging, blocks_needed: 50, blocks_behind: 1282, latest_block_head: 16356, current_block_head: 15074, network_name: mainnet, component: BlockIng
estor
2021-07-02T23:47:52.193212873Z Jul 02 23:47:52.193 INFO Syncing 50 blocks from Ethereum., code: BlockIngestionLagging, blocks_needed: 50, blocks_behind: 1227, latest_block_head: 17583, current_block_head: 16356, network_name: mainnet, component: BlockIng
estor
2021-07-02T23:47:55.832417460Z Jul 02 23:47:55.832 INFO Syncing 50 blocks from Ethereum., code: BlockIngestionLagging, blocks_needed: 50, blocks_behind: 1589, latest_block_head: 19172, current_block_head: 17583, network_name: mainnet, component: BlockIng
estor
2021-07-02T23:48:00.105888902Z Jul 02 23:48:00.105 INFO Syncing 50 blocks from Ethereum., code: BlockIngestionLagging, blocks_needed: 50, blocks_behind: 2012, latest_block_head: 21184, current_block_head: 19172, network_name: mainnet, component: BlockIng
estor
2021-07-02T23:48:04.286256921Z Jul 02 23:48:04.286 INFO Syncing 50 blocks from Ethereum., code: BlockIngestionLagging, blocks_needed: 50, blocks_behind: 2082, latest_block_head: 23266, current_block_head: 21184, network_name: mainnet, component: BlockIng
estor
2021-07-02T23:48:08.376502075Z Jul 02 23:48:08.376 INFO Syncing 50 blocks from Ethereum., code: BlockIngestionLagging, blocks_needed: 50, blocks_behind: 1958, latest_block_head: 25224, current_block_head: 23266, network_name: mainnet, component: BlockIng
estor
2021-07-02T23:48:12.591625800Z Jul 02 23:48:12.591 INFO Syncing 50 blocks from Ethereum., code: BlockIngestionLagging, blocks_needed: 50, blocks_behind: 2037, latest_block_head: 27261, current_block_head: 25224, network_name: mainnet, component: BlockIng
estor
2021-07-02T23:48:33.672341124Z Jul 02 23:48:33.672 WARN Trying again after block polling failed: RPC error: Error { code: ServerError(100), message: "invalid address", data: None }, network_name: mainnet, component: BlockIngestor
2021-07-02T23:48:34.696259203Z Jul 02 23:48:34.696 INFO Syncing 50 blocks from Ethereum., code: BlockIngestionLagging, blocks_needed: 50, blocks_behind: 12539, latest_block_head: 37763, current_block_head: 25224, network_name: mainnet, component: BlockIn
gestor
2021-07-02T23:48:38.626620047Z Jul 02 23:48:38.626 INFO Syncing 50 blocks from Ethereum., code: BlockIngestionLagging, blocks_needed: 50, blocks_behind: 1851, latest_block_head: 39614, current_block_head: 37763, network_name: mainnet, component: BlockIng
estor
2021-07-02T23:48:42.833452892Z Jul 02 23:48:42.833 INFO Syncing 50 blocks from Ethereum., code: BlockIngestionLagging, blocks_needed: 50, blocks_behind: 1975, latest_block_head: 41589, current_block_head: 39614, network_name: mainnet, component: BlockIng
estor
2021-07-02T23:49:18.616046204Z Jul 02 23:49:18.615 WARN Trying again after block polling failed: RPC error: Error { code: ServerError(100), message: "invalid address", data: None }, network_name: mainnet, component: BlockIngestor
2021-07-02T23:49:19.649545543Z Jul 02 23:49:19.649 INFO Syncing 50 blocks from Ethereum., code: BlockIngestionLagging, blocks_needed: 50, blocks_behind: 18929, latest_block_head: 58543, current_block_head: 39614, network_name: mainnet, component: BlockIn
gestor
2021-07-02T23:49:58.203958696Z Jul 02 23:49:58.203 WARN Trying again after block polling failed: RPC error: Error { code: ServerError(100), message: "invalid address", data: None }, network_name: mainnet, component: BlockIngestor
2021-07-02T23:49:59.243803337Z Jul 02 23:49:59.243 INFO Syncing 50 blocks from Ethereum., code: BlockIngestionLagging, blocks_needed: 50, blocks_behind: 37339, latest_block_head: 76953, current_block_head: 39614, network_name: mainnet, component: BlockIn
gestor
```

## updated Janus and theGraph logs

https://docs.google.com/spreadsheets/d/1gi9B8VnIFgbhwmZbc4BH5PX4miyS_LZphyyR5nwEmg4/edit?usp=sharing

