import * as ethers from "ethers";
import pairABI from "./contracts/pair.abi.json";
import { TxSwapSide } from "./types/enums";

const wsRPC = process.env.WS_RPC || "wss://bsc-ws-node.nariox.org:443";

const pairs = {
	bnbbusd: {
		addresses: [
			"0x58f876857a02d6762e0101bb5c46a8c1ed44dc16",
			"0xaCAac9311b0096E04Dfe96b6D87dec867d3883Dc",
		],
		isToken0: true,
	},
};

const main = async () => {
	const wsprovider = new ethers.providers.WebSocketProvider(wsRPC);

	//Parse ABI
	const PairABI = new ethers.utils.Interface(pairABI);

	// Event Filter
	const query: ethers.EventFilter = {
		topics: [PairABI.getEventTopic("Swap")],
		address: ethers.utils.getAddress(pairs.bnbbusd.addresses[0]),
	};

	console.log(`== Tx Logs Ticker ==`);

	wsprovider.addListener(query, (txLog) => {
		const swapData = PairABI.decodeEventLog("Swap", txLog.data, txLog.topics);

		if (pairs.bnbbusd.isToken0) {
			const side = swapData.amount0Out > 0 ? TxSwapSide.BUY : TxSwapSide.SELL;
			const [amount0, amount1] =
				side === TxSwapSide.BUY
					? [swapData.amount0Out, swapData.amount1In]
					: [swapData.amount0In, swapData.amount1Out];
			console.log(`${side}: ${amount1 / amount0} BUSD`);
		} else {
			const side = swapData.amount1Out > 0 ? "BUY" : "SELL";
			const [amount0, amount1] =
				side === TxSwapSide.BUY
					? [swapData.amount0In, swapData.amount1Out]
					: [swapData.amount0Out, swapData.amount1In];
			console.log(`${side}: ${amount0 / amount1} BUSD`);
		}
	});
};

main();
