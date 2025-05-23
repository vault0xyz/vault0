// Package wallet provides blockchain wallet functionality for managing
// cryptocurrency accounts, signing transactions, and interacting with various
// blockchain networks.
//
// The wallet package is part of the Core/Infrastructure Layer and provides:
//   - Multi-chain wallet support (Ethereum, Polygon, Base)
//   - Secure key management through keystore integration
//   - Address derivation
//   - Transaction creation and signing
//   - Support for native and token transactions
//
// Key Features:
//   - EVM compatibility for multiple chains
//   - Secure private key handling
//   - Standardized transaction interface
//   - Chain-specific implementations
package wallet

import (
	"context"
	"math/big"
	"vault0/internal/types"
)

// WalletInfo contains the wallet identification and address information.
// This structure is used to store and transmit wallet metadata without
// exposing sensitive key material.
type WalletInfo struct {
	// KeyID is the unique identifier for the key in the keystore.
	// This ID is used to retrieve the key for signing operations.
	KeyID string

	// Address is the blockchain-specific address derived from the public key.
	// The format of this address depends on the blockchain network
	// (e.g., 0x-prefixed for EVM chains).
	Address string

	// ChainType identifies the blockchain network this wallet is for.
	// This determines which implementation and parameters are used
	// for blockchain interactions.
	ChainType types.ChainType
}

// WalletManager defines the interface for blockchain wallet operations.
// Implementations must ensure secure handling of private keys and
// proper transaction signing according to chain-specific requirements.
type WalletManager interface {
	// Chain returns the blockchain network information for this wallet.
	// This includes network parameters, chain ID, and other chain-specific
	// configuration needed for transaction creation and signing.
	Chain() types.Chain

	// DeriveAddress derives the blockchain address for this wallet.
	// The derivation process depends on the blockchain network:
	//   - For EVM chains: Keccak256(public key) -> 0x-prefixed address
	//   - For other chains: Implementation-specific derivation
	//
	// The context can be used to cancel long-running operations.
	DeriveAddress(ctx context.Context) (string, error)

	// CreateNativeTransaction creates an unsigned transaction for the chain's
	// native currency (e.g., ETH for Ethereum, MATIC for Polygon).
	//
	// Parameters:
	//   - ctx: Context for the operation
	//   - toAddress: Recipient's blockchain address
	//   - amount: Transaction amount in the smallest unit (e.g., wei)
	//   - options: Chain-specific transaction parameters (gas price, nonce, etc.)
	//
	// Returns:
	//   - *types.Transaction: Unsigned transaction ready for signing
	//   - error: Any error during transaction creation
	CreateNativeTransaction(ctx context.Context, toAddress string, amount *big.Int, options types.TransactionOptions) (*types.Transaction, error)

	// CreateTokenTransaction creates an unsigned transaction for ERC20 or similar
	// tokens on the blockchain.
	//
	// Parameters:
	//   - ctx: Context for the operation
	//   - tokenAddress: Contract address of the token
	//   - toAddress: Recipient's blockchain address
	//   - amount: Token amount in the smallest unit (depends on token decimals)
	//   - options: Chain-specific transaction parameters (gas price, nonce, etc.)
	//
	// Returns:
	//   - *types.Transaction: Unsigned transaction ready for signing
	//   - error: Any error during transaction creation
	CreateTokenTransaction(ctx context.Context, tokenAddress, toAddress string, amount *big.Int, options types.TransactionOptions) (*types.Transaction, error)

	// SignTransaction signs a transaction using the wallet's private key.
	// The signing process follows chain-specific requirements:
	//   - For EVM chains: Signs according to EIP-155
	//   - For other chains: Implementation-specific signing
	//
	// Parameters:
	//   - ctx: Context for the operation
	//   - tx: Transaction to sign
	//
	// Returns:
	//   - []byte: Signed transaction bytes ready for broadcasting
	//   - error: Any error during signing
	SignTransaction(ctx context.Context, tx *types.Transaction) ([]byte, error)

	// CreateContractCallTransaction creates an unsigned transaction to call a
	// method on a smart contract. It uses the provided ABI to encode the
	// method call data.
	//
	// Parameters:
	//   - ctx: Context for the operation.
	//   - contractAddress: The address of the smart contract to interact with.
	//   - value: The amount of native currency to send with the call (for payable functions).
	//   - abiString: The JSON ABI string of the contract.
	//   - method: The name of the contract method to call.
	//   - args: A slice of arguments for the contract method call.
	//   - options: Chain-specific transaction parameters (gas price, nonce, etc.).
	//
	// Returns:
	//   - *types.Transaction: An unsigned transaction object ready for signing.
	//                       The 'Data' field will contain the ABI-encoded call.
	//                       The 'Value' field will contain the native currency amount.
	//   - error: Any error encountered during ABI encoding or transaction creation.
	CreateContractCallTransaction(ctx context.Context, contractAddress string, value *big.Int, abiString string, method string, args []any, options types.TransactionOptions) (*types.Transaction, error)
}
