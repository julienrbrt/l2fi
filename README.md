# L2 Force Inclusion UI

A simple UI to force inclusion of transactions on L2 networks when the sequencer is down or censoring transactions.

## Tech Stack

- Go with Chi router
- Tailwind CSS for styling
- [HTMX](https://htmx.org) for dynamic interactions - [Source](https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js)
- [Ethers.js](https://docs.ethers.io) for wallet integration - [Source](https://cdnjs.cloudflare.com/ajax/libs/ethers/6.13.5/ethers.umd.min.js)

## Setup

1. **Install dependencies:**

   ```bash
   # Install Go dependencies
   go mod download
   
   # Install Node.js dependencies for Tailwind CSS
   npm install
   ```

2. **Set up environment variables:**

   ```bash
   cp .env.example .env
   # Edit .env and add your Ethereum mainnet RPC URL
   ```

3. **Build CSS:**

   ```bash
   make build-css
   ```

4. **Build and run:**

   ```bash
   make build
   ./l2fi
   ```

   Or run directly:

   ```bash
   go run .
   ```

## Usage

1. Open <http://localhost:8080> in your browser
2. Select the target L2 network (Optimism or Arbitrum)
3. Fill in the transaction details:
   - **L2 Recipient Address**: The destination address on L2
   - **Value**: Amount in Wei to send
   - **Data**: Optional transaction data (hex format)
   - **L2 Gas Limit**: Gas limit for the L2 transaction
4. Click "Build Force Inclusion Transaction"
5. Review the generated transaction
6. Click "Sign with MetaMask" to submit the transaction

### CSS Development

For live CSS rebuilding during development:

```bash
make build-css-watch
```

### Project Structure

```
├── main.go              # HTTP server and routes
├── l2/                  # L2-specific force inclusion logic
│   ├── l2.go           # Common interface
│   ├── optimism.go     # Optimism implementation
│   └── arbitrum.go     # Arbitrum implementation
├── templates/          # HTML templates
├── static/             # Static assets (CSS, JS)
└── Makefile           # Build targets
```

## How Force Inclusion Works

### Optimism

- Uses the L1CrossDomainMessenger contract on Ethereum mainnet
- Calls `sendMessage()` to queue a transaction that will be included in the next L2 block
- The transaction bypasses the sequencer and is forced to be included

### Arbitrum  

- Uses the Delayed Inbox contract on Ethereum mainnet
- Calls `sendL2Message()` to submit a transaction directly to L2
- The transaction is processed after a delay period, bypassing the sequencer

Both methods require paying L1 gas fees but guarantee transaction inclusion when sequencers are down or censoring.

## Reference Links

- [Optimism: Bypassing the Sequencer](https://docs.optimism.io/stack/rollup/outages#bypassing-the-sequencer)
- [Optimism: Forced Transactions](https://docs.optimism.io/stack/transactions/forced-transaction)
- [Optimism: OptimismPortal Contract](https://github.com/ethereum-optimism/optimism/blob/111f3f3a3a2881899662e53e0f1b2f845b188a38/packages/contracts-bedrock/src/L1/OptimismPortal.sol#L209)
- [Optimism L2 Registry](https://github.com/ethereum-optimism/superchain-registry/tree/main/superchain/configs)
