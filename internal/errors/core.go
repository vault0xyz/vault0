package errors

import (
	"fmt"
	"math/big"
	// Removed: "vault0/internal/types"
)

// Core module error codes
const (
	// Database errors
	ErrCodeDatabaseError    = "database_error"
	ErrCodeDatabaseNotFound = "database_not_found"

	// Configuration errors
	ErrCodeConfiguration = "configuration_error"

	// Resource errors
	ErrCodeResourceNotFound = "resource_not_found"
	ErrCodeResourceExists   = "resource_exists"

	// Token errors
	ErrCodeInvalidToken = "invalid_token"

	// Blockchain errors
	ErrCodeBlockchainError            = "blockchain_error"
	ErrCodeInvalidBlockchainConfig    = "invalid_blockchain_config"
	ErrCodeChainNotSupported          = "chain_not_supported"
	ErrCodeInsufficientFunds          = "insufficient_funds"
	ErrCodeInvalidTransaction         = "invalid_transaction"
	ErrCodeTransactionNotFound        = "transaction_not_found"
	ErrCodeBlockNotFound              = "block_not_found"
	ErrCodeInvalidBlockIdentifier     = "invalid_block_identifier"
	ErrCodeRPCError                   = "rpc_error"
	ErrCodeInvalidAddress             = "invalid_address"
	ErrCodeTransactionFailed          = "transaction_failed"
	ErrCodeInvalidContract            = "invalid_contract"
	ErrCodeContractNotFound           = "contract_not_found"
	ErrCodeMethodNotFound             = "method_not_found"
	ErrCodeInvalidEventSignature      = "invalid_event_signature"
	ErrCodeInvalidEventArgs           = "invalid_event_args"
	ErrCodeUnsupportedEventArgType    = "unsupported_event_arg_type"
	ErrCodeTransactionCreationFailed  = "transaction_creation_failed"
	ErrCodeTransactionSigningFailed   = "transaction_signing_failed"
	ErrCodeTransactionBroadcastFailed = "transaction_broadcast_failed"

	// Keystore errors
	ErrCodeKeystoreError   = "keystore_error"
	ErrCodeKeyNotFound     = "key_not_found"
	ErrCodeKeyExists       = "key_exists"
	ErrCodeInvalidKey      = "invalid_key"
	ErrCodeSigningError    = "signing_error"
	ErrCodeInvalidKeystore = "invalid_keystore"

	// Wallet errors
	ErrCodeWalletError         = "wallet_error"
	ErrCodeInvalidWalletConfig = "invalid_wallet_config"
	ErrCodeInvalidKeyType      = "invalid_key_type"
	ErrCodeInvalidCurve        = "invalid_curve"
	ErrCodeInvalidSignature    = "invalid_signature"
	ErrCodeSignatureRecovery   = "signature_recovery_failed"
	ErrCodeAddressMismatch     = "address_mismatch"

	// Crypto errors
	ErrCodeCryptoError          = "crypto_error"
	ErrCodeEncryptionError      = "encryption_error"
	ErrCodeDecryptionError      = "decryption_error"
	ErrCodeInvalidEncryptionKey = "invalid_encryption_key"

	// Block explorer errors
	ErrCodeExplorerError           = "explorer_error"
	ErrCodeRateLimitExceeded       = "rate_limit_exceeded"
	ErrCodeInvalidAPIKey           = "invalid_api_key"
	ErrCodeInvalidExplorerResponse = "invalid_explorer_response"
	ErrCodeExplorerRequestFailed   = "explorer_request_failed"
	ErrCodeMissingAPIKey           = "missing_api_key"

	// Blockchain transaction errors
	ErrCodeInvalidNonce        = "invalid_nonce"
	ErrCodeInvalidGasPrice     = "invalid_gas_price"
	ErrCodeInvalidGasLimit     = "invalid_gas_limit"
	ErrCodeInvalidContractCall = "invalid_contract_call"
	ErrCodeInvalidAmount       = "invalid_amount"

	// New token errors
	ErrCodeInvalidTokenBalance = "invalid_token_balance"

	// Price Feed errors
	ErrCodePriceFeedRequestFailed        = "price_feed_request_failed"
	ErrCodeInvalidPriceFeedResponse      = "invalid_price_feed_response"
	ErrCodePriceFeedProviderNotSupported = "price_feed_provider_not_supported"

	// Log parsing errors
	ErrCodeLogTopicIndexOutOfBounds = "log_topic_index_out_of_bounds"
	ErrCodeLogTopicInvalidFormat    = "log_topic_invalid_format"

	// Pagination errors
	ErrCodeInvalidPaginationToken = "invalid_pagination_token"
	ErrCodeTokenEncodingFailed    = "token_encoding_failed"
	ErrCodeTokenDecodingFailed    = "token_decoding_failed"

	// Blockchain node errors
	ErrCodeBlockchainNodeUnreachable = "blockchain_node_unreachable"

	// ABI errors
	ErrCodeABIError     = "abi_error"
	ErrCodeMappingError = "mapping_error"

	// ABI Proxy errors
	ErrCodeABIProxyParseFailed             = "abi_proxy_parse_failed"
	ErrCodeABIProxyMethodNotFound          = "abi_proxy_method_not_found"
	ErrCodeABIProxyMethodSignatureInvalid  = "abi_proxy_method_signature_invalid"
	ErrCodeABIProxyPackFailed              = "abi_proxy_pack_failed"
	ErrCodeABIProxyCallFailed              = "abi_proxy_call_failed"
	ErrCodeABIProxyEmptyResult             = "abi_proxy_empty_result"
	ErrCodeABIProxyUnpackFailed            = "abi_proxy_unpack_failed"
	ErrCodeABIProxyAddressConversionFailed = "abi_proxy_address_conversion_failed"

	// General ABI processing errors (extending existing ErrCodeABIError context)
	ErrCodeABIParseFailed              = "abi_parse_failed"
	ErrCodeABIMethodNotFound           = "abi_method_not_found" // Can be used for both by name and by ID
	ErrCodeABIMethodSelectorMismatch   = "abi_method_selector_mismatch"
	ErrCodeABIInputDataTooShort        = "abi_input_data_too_short"
	ErrCodeABIInputDataEmpty           = "abi_input_data_empty"
	ErrCodeABIInputDataInvalidLength   = "abi_input_data_invalid_length"
	ErrCodeABIUnpackFailed             = "abi_unpack_failed"
	ErrCodeABIPackFailed               = "abi_pack_failed"
	ErrCodeABIArgumentNotFound         = "abi_argument_not_found"
	ErrCodeABIArgumentConversionFailed = "abi_argument_conversion_failed"
	ErrCodeABIArgumentInvalidType      = "abi_argument_invalid_type"
	ErrCodeABIArgumentNilValue         = "abi_argument_nil_value"

	// ABI availability errors
	ErrCodeABIUnavailableOrUnverified = "abi_unavailable_or_unverified"
)

// NewConfigurationError creates an error for configuration issues.
func NewConfigurationError(message string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeConfiguration,
		Message: fmt.Sprintf("Configuration error: %s", message),
	}
}

// NewInvalidBlockchainConfigError creates an error for invalid blockchain configuration
func NewInvalidBlockchainConfigError(chain string, key string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidBlockchainConfig,
		Message: fmt.Sprintf("Invalid blockchain configuration for %s: missing %s", chain, key),
		Details: map[string]any{
			"chain": chain,
			"key":   key,
		},
	}
}

// NewDatabaseError creates an error for database failures
func NewDatabaseError(err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeDatabaseError,
		Message: "Database operation failed",
		Err:     err,
	}
}

// NewDatabaseNotFoundError creates an error for database record not found
func NewDatabaseNotFoundError(entity string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeDatabaseNotFound,
		Message: entity + " not found",
	}
}

// NewBlockchainError creates a new error for blockchain-related issues
func NewBlockchainError(err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeBlockchainError,
		Message: fmt.Sprintf("Blockchain error: %v", err),
		Details: map[string]any{
			"error": err.Error(),
		},
	}
}

// NewChainNotSupportedError creates a new error for unsupported chains
func NewChainNotSupportedError(chain string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeChainNotSupported,
		Message: fmt.Sprintf("Chain not supported: %s", chain),
		Details: map[string]any{
			"chain": chain,
		},
	}
}

// NewInsufficientFundsError creates a new error for insufficient funds
func NewInsufficientFundsError(balance string, required string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInsufficientFunds,
		Message: fmt.Sprintf("Insufficient funds: have %s, need %s", balance, required),
		Details: map[string]any{
			"balance":  balance,
			"required": required,
		},
	}
}

// NewKeystoreError creates an error for keystore operations
func NewKeystoreError(err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeKeystoreError,
		Message: "Keystore operation failed",
		Err:     err,
	}
}

// NewKeyNotFoundError creates an error for missing keys
func NewKeyNotFoundError(keyID string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeKeyNotFound,
		Message: "Key not found: " + keyID,
	}
}

// NewInvalidKeyError creates an error for invalid keys
func NewInvalidKeyError(msg string, err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidKey,
		Message: msg,
		Err:     err,
	}
}

// NewCryptoError creates an error for cryptographic operations
func NewCryptoError(err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeCryptoError,
		Message: "Cryptographic operation failed",
		Err:     err,
	}
}

// NewEncryptionError creates an error for encryption failures
func NewEncryptionError(err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeEncryptionError,
		Message: "Encryption failed",
		Err:     err,
	}
}

// NewDecryptionError creates an error for decryption failures
func NewDecryptionError(err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeDecryptionError,
		Message: "Decryption failed",
		Err:     err,
	}
}

// NewExplorerError creates a new error for block explorer issues
func NewExplorerError(err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeExplorerError,
		Message: fmt.Sprintf("Block explorer error: %v", err),
		Err:     err,
	}
}

// NewRateLimitExceededError creates a new error for rate limit issues
func NewRateLimitExceededError() *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeRateLimitExceeded,
		Message: "Rate limit exceeded",
	}
}

// NewInvalidAddressError creates a new error for invalid addresses
func NewInvalidAddressError(address string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidAddress,
		Message: fmt.Sprintf("Invalid address: %s", address),
		Details: map[string]any{
			"address": address,
		},
	}
}

// NewInvalidAmountError creates a new error for invalid amounts
func NewInvalidAmountError(amount string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidAmount,
		Message: fmt.Sprintf("Invalid amount: %s", amount),
		Details: map[string]any{
			"amount": amount,
		},
	}
}

// NewRPCError creates a new error for RPC-related issues
func NewRPCError(err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeRPCError,
		Message: fmt.Sprintf("RPC error: %v", err),
		Err:     err,
	}
}

// NewInvalidTransactionError creates a new error for invalid transactions
func NewInvalidTransactionError(err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidTransaction,
		Message: fmt.Sprintf("Invalid transaction: %v", err),
		Err:     err,
	}
}

// NewInvalidAPIKeyError creates a new error for invalid API keys
func NewInvalidAPIKeyError() *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidAPIKey,
		Message: "Invalid API key",
	}
}

// NewInvalidExplorerResponseError creates a new error for invalid responses from block explorer
func NewInvalidExplorerResponseError(err error, response string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidExplorerResponse,
		Message: fmt.Sprintf("Invalid response from block explorer: %v", err),
		Err:     err,
		Details: map[string]any{
			"response": response,
		},
	}
}

// NewExplorerRequestFailedError creates a new error for failed block explorer requests
func NewExplorerRequestFailedError(err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeExplorerRequestFailed,
		Message: fmt.Sprintf("Block explorer request failed: %v", err),
		Err:     err,
	}
}

// NewMissingAPIKeyError creates a new error for missing API key
func NewMissingAPIKeyError() *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeMissingAPIKey,
		Message: "Missing API key for block explorer",
	}
}

// NewInvalidNonceError creates a new error for invalid nonce
func NewInvalidNonceError(address string, nonce uint64) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidNonce,
		Message: fmt.Sprintf("Invalid nonce for address %s: %d", address, nonce),
		Details: map[string]any{
			"address": address,
			"nonce":   nonce,
		},
	}
}

// NewInvalidGasPriceError creates a new error for invalid gas price
func NewInvalidGasPriceError(gasPrice *big.Int) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidGasPrice,
		Message: fmt.Sprintf("Invalid gas price: %s", gasPrice.String()),
		Details: map[string]any{
			"gas_price": gasPrice.String(),
		},
	}
}

// NewInvalidGasLimitError creates a new error for invalid gas limit
func NewInvalidGasLimitError(gasLimit uint64) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidGasLimit,
		Message: fmt.Sprintf("Invalid gas limit: %d", gasLimit),
		Details: map[string]any{
			"gas_limit": gasLimit,
		},
	}
}

// NewInvalidContractCallError creates a new error for invalid contract calls
func NewInvalidContractCallError(contract string, err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidContractCall,
		Message: "Contract call failed",
		Details: map[string]any{
			"contract": contract,
		},
		Err: err,
	}
}

// NewInvalidContractError creates a new error for invalid contracts
func NewInvalidContractError(contract string, err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidContract,
		Message: fmt.Sprintf("Invalid contract: %s", contract),
		Err:     err,
		Details: map[string]any{
			"contract": contract,
			"error":    err.Error(),
		},
	}
}

// NewWalletError creates a new error for general wallet operations
func NewWalletError(msg string, err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeWalletError,
		Message: fmt.Sprintf("Wallet error: %s", msg),
		Err:     err,
		Details: map[string]any{
			"error": err.Error(),
		},
	}
}

// NewInvalidWalletConfigError creates a new error for invalid wallet configuration
func NewInvalidWalletConfigError(msg string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidWalletConfig,
		Message: fmt.Sprintf("Invalid wallet configuration: %s", msg),
	}
}

// NewInvalidKeyTypeError creates a new error for invalid key types
func NewInvalidKeyTypeError(expected, got string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidKeyType,
		Message: fmt.Sprintf("Invalid key type: expected %s, got %s", expected, got),
		Details: map[string]any{
			"expected": expected,
			"got":      got,
		},
	}
}

// NewInvalidCurveError creates a new error for invalid curves
func NewInvalidCurveError(expected, got string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidCurve,
		Message: fmt.Sprintf("Invalid curve: expected %s, got %s", expected, got),
		Details: map[string]any{
			"expected": expected,
			"got":      got,
		},
	}
}

// NewInvalidSignatureError creates a new error for invalid signatures
func NewInvalidSignatureError(err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidSignature,
		Message: "Invalid signature",
		Err:     err,
		Details: map[string]any{
			"error": err.Error(),
		},
	}
}

// NewSignatureRecoveryError creates a new error for signature recovery failures
func NewSignatureRecoveryError(err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeSignatureRecovery,
		Message: "Failed to recover signature",
		Err:     err,
	}
}

// NewAddressMismatchError creates a new error for address mismatches
func NewAddressMismatchError(expected, got string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeAddressMismatch,
		Message: fmt.Sprintf("Address mismatch: expected %s, got %s", expected, got),
		Details: map[string]any{
			"expected": expected,
			"got":      got,
		},
	}
}

// NewTransactionNotFoundError creates a new error for transaction not found
func NewTransactionNotFoundError(hash string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeTransactionNotFound,
		Message: fmt.Sprintf("Transaction not found: %s", hash),
		Details: map[string]any{
			"hash": hash,
		},
	}
}

// NewTransactionFailedError creates an error for failed transaction
func NewTransactionFailedError(err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeTransactionFailed,
		Message: "Transaction failed",
		Err:     err,
	}
}

// NewInvalidEncryptionKeyError creates a new error for invalid encryption keys
func NewInvalidEncryptionKeyError(key string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidEncryptionKey,
		Message: fmt.Sprintf("Invalid encryption key: %s", key),
	}
}

// NewResourceNotFoundError creates an error for when a resource is not found
func NewResourceNotFoundError(resource, id string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeResourceNotFound,
		Message: fmt.Sprintf("%s not found: %s", resource, id),
		Details: map[string]any{
			"resource": resource,
			"id":       id,
		},
	}
}

// NewResourceAlreadyExistsError creates an error for when a resource already exists
func NewResourceAlreadyExistsError(resource, attribute, value string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeResourceExists,
		Message: fmt.Sprintf("%s with %s '%s' already exists", resource, attribute, value),
		Details: map[string]any{
			"resource":  resource,
			"attribute": attribute,
			"value":     value,
		},
	}
}

// NewSigningError creates an error for signing operations
func NewSigningError(err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeSigningError,
		Message: "Signing operation failed",
		Err:     err,
		Details: map[string]any{
			"error": err.Error(),
		},
	}
}

// NewInvalidKeystoreError creates an error for invalid or uninitialized keystore
func NewInvalidKeystoreError(keystore string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidKeystore,
		Message: fmt.Sprintf("Invalid keystore: %s", keystore),
	}
}

// NewInvalidTokenError creates an error for invalid token data
func NewInvalidTokenError(msg string, err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidToken,
		Message: fmt.Sprintf("Invalid token: %s", msg),
		Err:     err,
	}
}

// NewContractNotFoundError creates a new error for when a contract is not found at the specified address
func NewContractNotFoundError(address string, chainType string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeContractNotFound,
		Message: fmt.Sprintf("Contract not found at address %s on chain %s", address, chainType),
		Details: map[string]any{
			"address": address,
			"chain":   chainType,
		},
	}
}

// NewInvalidEventSignatureError creates a new error for invalid event signatures
func NewInvalidEventSignatureError(signature string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidEventSignature,
		Message: fmt.Sprintf("Invalid event signature format: %s", signature),
		Details: map[string]any{
			"signature": signature,
		},
	}
}

// NewInvalidEventArgsError creates a new error for invalid event arguments
func NewInvalidEventArgsError(msg string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidEventArgs,
		Message: msg,
		Details: map[string]any{
			"error": msg,
		},
	}
}

// NewUnsupportedEventArgTypeError creates a new error for unsupported event argument types
func NewUnsupportedEventArgTypeError(paramIndex int) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeUnsupportedEventArgType,
		Message: fmt.Sprintf("Unsupported argument type for parameter %d", paramIndex),
		Details: map[string]any{
			"parameter_index": paramIndex,
		},
	}
}

// NewBlockNotFoundError creates a new error for when a block cannot be found
func NewBlockNotFoundError(identifier string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeBlockNotFound,
		Message: fmt.Sprintf("Block not found: %s", identifier),
		Details: map[string]any{
			"identifier": identifier,
		},
	}
}

// NewInvalidBlockIdentifierError creates a new error for invalid block identifiers
func NewInvalidBlockIdentifierError(identifier string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidBlockIdentifier,
		Message: fmt.Sprintf("Invalid block identifier: %s", identifier),
		Details: map[string]any{
			"identifier": identifier,
		},
	}
}

// NewInvalidTokenBalanceError creates a new error for failed token balance requests
func NewInvalidTokenBalanceError(token string, err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidTokenBalance,
		Message: "Failed to get token balance",
		Details: map[string]any{
			"token": token,
		},
		Err: err,
	}
}

// NewPriceFeedRequestFailed creates a new error for failed price feed API requests.
func NewPriceFeedRequestFailed(err error, details string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodePriceFeedRequestFailed,
		Message: "Price feed API request failed",
		Details: map[string]any{"details": details},
		Err:     err,
	}
}

// NewInvalidPriceFeedResponse creates a new error for invalid price feed API responses.
func NewInvalidPriceFeedResponse(err error, details string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidPriceFeedResponse,
		Message: "Invalid price feed API response",
		Details: map[string]any{"details": details},
		Err:     err,
	}
}

// NewPriceFeedProviderNotSupported creates an error for unsupported price feed providers.
func NewPriceFeedProviderNotSupported(provider string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodePriceFeedProviderNotSupported,
		Message: fmt.Sprintf("Price feed provider '%s' not supported", provider),
		Err:     nil, // No underlying Go error
	}
}

// NewLogTopicIndexOutOfBoundsError creates an error for when a topic index is out of bounds.
func NewLogTopicIndexOutOfBoundsError(index, count int) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeLogTopicIndexOutOfBounds,
		Message: fmt.Sprintf("Log topic index %d is out of bounds for %d topics", index, count),
		Details: map[string]any{
			"index": index,
			"count": count,
		},
	}
}

// NewLogTopicInvalidFormatError creates an error for when a log topic has an invalid format.
func NewLogTopicInvalidFormatError(index int, topicValue, reason string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeLogTopicInvalidFormat,
		Message: fmt.Sprintf("Log topic[%d] has an invalid format: %s", index, reason),
		Details: map[string]any{
			"index":       index,
			"topic_value": topicValue,
			"reason":      reason,
		},
	}
}

// NewInvalidPaginationTokenError creates an error for invalid pagination tokens
func NewInvalidPaginationTokenError(token string, err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeInvalidPaginationToken,
		Message: "Invalid pagination token format",
		Err:     err,
		Details: map[string]any{
			"token": token,
		},
	}
}

// NewTokenEncodingFailedError creates an error for token encoding failures
func NewTokenEncodingFailedError(err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeTokenEncodingFailed,
		Message: "Failed to encode pagination token",
		Err:     err,
	}
}

// NewTokenDecodingFailedError creates an error for token decoding failures
func NewTokenDecodingFailedError(token string, err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeTokenDecodingFailed,
		Message: "Failed to decode pagination token",
		Err:     err,
		Details: map[string]any{
			"token": token,
		},
	}
}

// NewMethodNotFoundError creates an error when a method is not found in a contract's ABI
func NewMethodNotFoundError(methodName, contractAddress string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeMethodNotFound,
		Message: fmt.Sprintf("Method '%s' not found in ABI for contract %s", methodName, contractAddress),
		Details: map[string]any{
			"method_name":      methodName,
			"contract_address": contractAddress,
		},
	}
}

// NewTransactionCreationError creates an error for failures during transaction creation
func NewTransactionCreationError(context string, err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeTransactionCreationFailed,
		Message: fmt.Sprintf("Failed to create transaction (%s)", context),
		Err:     err,
		Details: map[string]any{
			"context": context,
		},
	}
}

// NewTransactionSigningError creates an error for failures during transaction signing
func NewTransactionSigningError(err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeTransactionSigningFailed,
		Message: "Failed to sign transaction",
		Err:     err,
	}
}

// NewTransactionBroadcastError creates an error for failures during transaction broadcast
func NewTransactionBroadcastError(err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeTransactionBroadcastFailed,
		Message: "Failed to broadcast transaction",
		Err:     err,
	}
}

// NewBlockchainNodeUnreachableError creates a new error for unreachable blockchain nodes.
func NewBlockchainNodeUnreachableError(node string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeBlockchainNodeUnreachable,
		Message: fmt.Sprintf("Blockchain node '%s' is unreachable", node),
		Details: map[string]any{
			"node": node,
		},
	}
}

// NewABIError creates a new error for ABI processing issues.
func NewABIError(err error, context string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIError,
		Message: fmt.Sprintf("ABI processing failed: %s", context),
		Details: map[string]any{"context": context},
		Err:     err,
	}
}

// NewMappingError creates a new error for transaction mapping failures.
func NewMappingError(txHashStr string, reason string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeMappingError,
		Message: fmt.Sprintf("Failed to map transaction %s: %s", txHashStr, reason),
		Details: map[string]any{
			"transaction_hash": txHashStr,
			"reason":           reason,
		},
	}
}

// NewABIProxyParseError creates an error for proxy ABI parsing failures.
func NewABIProxyParseError(err error, proxyAddress string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIProxyParseFailed,
		Message: fmt.Sprintf("Failed to parse ABI for proxy contract %s", proxyAddress),
		Err:     err,
		Details: map[string]any{
			"proxy_address": proxyAddress,
		},
	}
}

// NewABIProxyMethodNotFoundError creates an error when the proxy implementation method is not found.
func NewABIProxyMethodNotFoundError(proxyAddress string, methodName string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIProxyMethodNotFound,
		Message: fmt.Sprintf("Method '%s' not found in ABI for proxy contract %s", methodName, proxyAddress),
		Details: map[string]any{
			"proxy_address": proxyAddress,
			"method_name":   methodName,
		},
	}
}

// NewABIProxyMethodSignatureInvalidError creates an error for invalid proxy method signatures.
func NewABIProxyMethodSignatureInvalidError(proxyAddress string, methodName string, expectedSignature string, actualSignature string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIProxyMethodSignatureInvalid,
		Message: fmt.Sprintf("Proxy contract %s method '%s' has invalid signature. Expected '%s', got '%s'", proxyAddress, methodName, expectedSignature, actualSignature),
		Details: map[string]any{
			"proxy_address":      proxyAddress,
			"method_name":        methodName,
			"expected_signature": expectedSignature,
			"actual_signature":   actualSignature,
		},
	}
}

// NewABIProxyPackError creates an error for failures when packing data for a proxy method call.
func NewABIProxyPackError(err error, proxyAddress string, methodName string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIProxyPackFailed,
		Message: fmt.Sprintf("Failed to pack call data for method '%s' on proxy %s", methodName, proxyAddress),
		Err:     err,
		Details: map[string]any{
			"proxy_address": proxyAddress,
			"method_name":   methodName,
		},
	}
}

// NewABIProxyCallError creates an error for failures when calling a proxy contract method.
func NewABIProxyCallError(err error, proxyAddress string, methodName string, chainType string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIProxyCallFailed,
		Message: fmt.Sprintf("Failed to call method '%s' on proxy %s (chain: %s)", methodName, proxyAddress, chainType),
		Err:     err,
		Details: map[string]any{
			"proxy_address": proxyAddress,
			"method_name":   methodName,
			"chain_type":    chainType,
		},
	}
}

// NewABIProxyEmptyResultError creates an error when a proxy method call returns empty data.
func NewABIProxyEmptyResultError(proxyAddress string, methodName string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIProxyEmptyResult,
		Message: fmt.Sprintf("Call to method '%s' on proxy %s returned empty data", methodName, proxyAddress),
		Details: map[string]any{
			"proxy_address": proxyAddress,
			"method_name":   methodName,
		},
	}
}

// NewABIProxyUnpackError creates an error for failures when unpacking results from a proxy method call.
func NewABIProxyUnpackError(err error, proxyAddress string, methodName string, detail string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIProxyUnpackFailed,
		Message: fmt.Sprintf("Failed to unpack result from '%s' call on proxy %s: %s", methodName, proxyAddress, detail),
		Err:     err,
		Details: map[string]any{
			"proxy_address": proxyAddress,
			"method_name":   methodName,
			"detail":        detail,
		},
	}
}

// NewABIProxyAddressConversionError creates an error for failures converting an address from a proxy call.
func NewABIProxyAddressConversionError(err error, invalidAddress string, chainType string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIProxyAddressConversionFailed,
		Message: fmt.Sprintf("Failed to convert implementation address '%s' to types.Address for chain %s", invalidAddress, chainType),
		Err:     err,
		Details: map[string]any{
			"invalid_address": invalidAddress,
			"chain_type":      chainType,
		},
	}
}

// NewABIParseError creates an error for ABI string parsing failures.
func NewABIParseError(err error) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIParseFailed,
		Message: "Failed to parse ABI string",
		Err:     err,
	}
}

// NewABIMethodNotFoundError creates an error when a method is not found in an ABI.
// Can specify if lookup was by name or ID.
func NewABIMethodNotFoundError(identifier string, byName bool) *Vault0Error {
	lookupType := "ID"
	if byName {
		lookupType = "name"
	}
	return &Vault0Error{
		Code:    ErrCodeABIMethodNotFound,
		Message: fmt.Sprintf("Method with %s '%s' not found in ABI", lookupType, identifier),
		Details: map[string]any{
			"identifier":  identifier,
			"lookup_type": lookupType,
		},
	}
}

// NewABIMethodSelectorMismatchError creates an error for method selector mismatches.
func NewABIMethodSelectorMismatchError(methodName string, expectedSelector []byte, actualSelector []byte) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIMethodSelectorMismatch,
		Message: fmt.Sprintf("Method selector mismatch for method '%s': expected %x, got %x", methodName, expectedSelector, actualSelector),
		Details: map[string]any{
			"method_name":       methodName,
			"expected_selector": fmt.Sprintf("%x", expectedSelector),
			"actual_selector":   fmt.Sprintf("%x", actualSelector),
		},
	}
}

// NewABIInputDataTooShortError creates an error when input data is too short.
func NewABIInputDataTooShortError(length int, requiredLength int) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIInputDataTooShort,
		Message: fmt.Sprintf("Input data is too short (length %d) to contain a method selector (requires %d bytes)", length, requiredLength),
		Details: map[string]any{
			"length":          length,
			"required_length": requiredLength,
		},
	}
}

// NewABIInputDataEmptyError creates an error for empty input data when arguments are required.
func NewABIInputDataEmptyError(methodName string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIInputDataEmpty,
		Message: fmt.Sprintf("Input data is empty for method %s which requires arguments", methodName),
		Details: map[string]any{
			"method_name": methodName,
		},
	}
}

// NewABIInputDataInvalidLengthError creates an error for input data with invalid length.
func NewABIInputDataInvalidLengthError(methodName string, length int) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIInputDataInvalidLength,
		Message: fmt.Sprintf("Input data length %d is not a multiple of 32 bytes for method %s", length, methodName),
		Details: map[string]any{
			"method_name": methodName,
			"length":      length,
		},
	}
}

// NewABIUnpackFailedError creates an error for failures when unpacking ABI data.
func NewABIUnpackFailedError(err error, methodName string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIUnpackFailed,
		Message: fmt.Sprintf("Failed to unpack input/output for method %s", methodName),
		Err:     err,
		Details: map[string]any{
			"method_name": methodName,
		},
	}
}

// NewABIPackFailedError creates an error for failures when packing ABI data.
func NewABIPackFailedError(err error, methodName string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIPackFailed,
		Message: fmt.Sprintf("Failed to pack arguments for method '%s'", methodName),
		Err:     err,
		Details: map[string]any{
			"method_name": methodName,
		},
	}
}

// NewABIArgumentNotFoundError creates an error when an argument is not found.
func NewABIArgumentNotFoundError(argName string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIArgumentNotFound,
		Message: fmt.Sprintf("Argument '%s' not found", argName),
		Details: map[string]any{
			"argument_name": argName,
		},
	}
}

// NewABIArgumentConversionError creates an error for argument type conversion failures.
func NewABIArgumentConversionError(err error, argName string, targetType string, valueAttempted any) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIArgumentConversionFailed,
		Message: fmt.Sprintf("Failed to convert argument '%s' to type '%s'", argName, targetType),
		Err:     err,
		Details: map[string]any{
			"argument_name":   argName,
			"target_type":     targetType,
			"value_attempted": fmt.Sprintf("%v", valueAttempted), // Be cautious with logging raw values
		},
	}
}

// NewABIArgumentInvalidTypeError creates an error for invalid argument types.
func NewABIArgumentInvalidTypeError(argName string, expectedType string, actualType string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIArgumentInvalidType,
		Message: fmt.Sprintf("Argument '%s' has invalid type: expected %s, got %s", argName, expectedType, actualType),
		Details: map[string]any{
			"argument_name": argName,
			"expected_type": expectedType,
			"actual_type":   actualType,
		},
	}
}

// NewABIArgumentNilValueError creates an error when an argument expected to be non-nil is nil.
func NewABIArgumentNilValueError(argName string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIArgumentNilValue,
		Message: fmt.Sprintf("Argument '%s' resolved to a nil value where a non-nil value was expected", argName),
		Details: map[string]any{
			"argument_name": argName,
		},
	}
}

// NewABIUnavailableOrUnverifiedError creates an error for when an ABI is not available from the explorer or the contract is not verified.
func NewABIUnavailableOrUnverifiedError(contractAddress string, chainType string) *Vault0Error {
	return &Vault0Error{
		Code:    ErrCodeABIUnavailableOrUnverified,
		Message: fmt.Sprintf("ABI for contract %s on chain %s is unavailable or contract is not verified", contractAddress, chainType),
		Details: map[string]any{
			"contract_address": contractAddress,
			"chain_type":       chainType,
		},
	}
}
