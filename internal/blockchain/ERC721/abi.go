package ERC721

const ERC721ABI = `[
	{
		"inputs": [
			{"internalType": "address", "name": "from", "type": "address"},
			{"internalType": "address", "name": "to", "type": "address"},
			{"internalType": "uint256", "name": "tokenId", "type": "uint256"}
		],
		"name": "Transfer",
		"type": "event"
	},
	{
		"inputs": [{"internalType": "uint256", "name": "tokenId", "type": "uint256"}],
		"name": "tokenURI",
		"outputs": [{"internalType": "string", "name": "", "type": "string"}],
		"stateMutability": "view",
		"type": "function"
	}
]`
