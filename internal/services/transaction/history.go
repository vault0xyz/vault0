package transaction

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"
	"vault0/internal/config"
	"vault0/internal/core/blockchain"
	"vault0/internal/core/blockexplorer"
	"vault0/internal/core/tokenstore"
	"vault0/internal/errors"
	"vault0/internal/logger"
	"vault0/internal/types"
)

// HistoryService interface defines methods for transaction history synchronization
type HistoryService interface {
	// MonitorAddress adds an address to be monitored for transaction history
	MonitorAddress(address types.Address, startBlockNumber *big.Int) error

	// UnmonitorAddress removes an address from monitoring
	UnmonitorAddress(address types.Address) error

	// StartTransactionSyncing starts the background synchronization process
	StartTransactionSyncing(ctx context.Context) error

	// StopTransactionSyncing stops the synchronization process
	StopTransactionSyncing()

	// HistoryEvents returns a channel that emits processed historical transactions
	HistoryEvents() <-chan *TransactionEvent
}

// NewHistoryService creates a new transaction history service
func NewHistoryService(
	config *config.Config,
	log logger.Logger,
	blockExplorerFactory blockexplorer.Factory,
	blockchainFactory blockchain.Factory,
	transformer TransformerService,
	repository Repository,
	tokenStore tokenstore.TokenStore,
) HistoryService {
	service := &historyService{
		config:               config,
		log:                  log,
		blockExplorerFactory: blockExplorerFactory,
		blockchainFactory:    blockchainFactory,
		transformerService:   transformer,
		repository:           repository,
		tokenStore:           tokenStore,
		syncMutex:            sync.RWMutex{},
		syncAddresses:        make(map[string]addressSyncInfo),
		historyEventsChan:    make(chan *TransactionEvent, 100),
	}

	return service
}

type historyService struct {
	// Synchronization lifecycle management
	syncCtx       context.Context
	syncCancel    context.CancelFunc
	syncMutex     sync.RWMutex
	syncAddresses map[string]addressSyncInfo

	// Channel for emitting history events
	historyEventsChan chan *TransactionEvent

	// Dependencies
	config               *config.Config
	log                  logger.Logger
	blockExplorerFactory blockexplorer.Factory
	blockchainFactory    blockchain.Factory
	transformerService   TransformerService
	repository           Repository
	tokenStore           tokenstore.TokenStore
}

// addressSyncInfo holds the address and start block number for syncing
type addressSyncInfo struct {
	Address    types.Address
	StartBlock *big.Int
}

// deriveAddressKey creates a unique key for the address map combining chain and address
func (s *historyService) deriveAddressKey(address types.Address) string {
	return fmt.Sprintf("%s:%s", string(address.ChainType), strings.ToLower(address.ToChecksum()))
}

// MonitorAddress adds an address to be monitored for transaction history
func (s *historyService) MonitorAddress(address types.Address, startBlockNumber *big.Int) error {
	s.syncMutex.Lock()
	defer s.syncMutex.Unlock()

	// Use helper method to create composite key
	addrKey := s.deriveAddressKey(address)
	if _, exists := s.syncAddresses[addrKey]; exists {
		return nil // Already monitoring this address
	}

	s.syncAddresses[addrKey] = addressSyncInfo{
		Address:    address,
		StartBlock: startBlockNumber,
	}

	s.log.Info("Added address for transaction history syncing",
		logger.String("chain", string(address.ChainType)),
		logger.String("address", address.String()),
		logger.Int64("start_block_number", startBlockNumber.Int64()))
	return nil
}

// UnmonitorAddress removes an address from monitoring
func (s *historyService) UnmonitorAddress(address types.Address) error {
	s.syncMutex.Lock()
	defer s.syncMutex.Unlock()

	// Use helper method to create composite key
	addrKey := s.deriveAddressKey(address)
	if _, exists := s.syncAddresses[addrKey]; !exists {
		return nil // Address not being monitored
	}

	delete(s.syncAddresses, addrKey)
	s.log.Info("Removed address from transaction history syncing",
		logger.String("chain", string(address.ChainType)),
		logger.String("address", address.String()))
	return nil
}

// HistoryEvents returns a channel that emits processed historical transactions
func (s *historyService) HistoryEvents() <-chan *TransactionEvent {
	return s.historyEventsChan
}

// StartTransactionSyncing starts the background synchronization process
func (s *historyService) StartTransactionSyncing(ctx context.Context) error {
	s.syncMutex.Lock()

	// Check if already syncing
	if s.syncCtx != nil {
		s.syncMutex.Unlock()
		s.log.Info("Transaction history syncing is already active")
		return nil
	}

	// Create cancellable context
	s.syncCtx, s.syncCancel = context.WithCancel(ctx)
	s.syncMutex.Unlock()

	// Get interval from config with fallback to default
	interval := 300 // Default to 5 minutes if not specified
	if s.config.Transaction.HistorySynchInterval > 0 {
		interval = s.config.Transaction.HistorySynchInterval
	}

	s.log.Info("Starting transaction history syncing",
		logger.Int("interval_seconds", interval))

	// Start the scheduler goroutine
	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		// Initial sync immediately after starting
		s.syncTransactions(s.syncCtx)

		for {
			select {
			case <-s.syncCtx.Done():
				s.log.Info("Transaction history syncing stopped")
				return
			case <-ticker.C:
				s.syncTransactions(s.syncCtx)
			}
		}
	}()

	return nil
}

// StopTransactionSyncing stops the synchronization process
func (s *historyService) StopTransactionSyncing() {
	s.syncMutex.Lock()
	defer s.syncMutex.Unlock()

	if s.syncCtx == nil {
		return
	}

	s.log.Info("Stopping transaction history syncing")

	// Cancel the context to signal goroutines to stop
	if s.syncCancel != nil {
		s.syncCancel()
	}

	// Reset context and cancel function
	s.syncCtx = nil
	s.syncCancel = nil

	s.log.Info("Transaction history syncing stopped")
}

// syncTransactions syncs transaction history for all monitored addresses
func (s *historyService) syncTransactions(ctx context.Context) {
	s.syncMutex.RLock()
	// Create a copy of addresses to avoid holding the lock during processing
	addresses := make([]types.Address, 0, len(s.syncAddresses))
	for _, addr := range s.syncAddresses {
		addresses = append(addresses, addr.Address)
	}
	s.syncMutex.RUnlock()

	if len(addresses) == 0 {
		s.log.Debug("No addresses configured for transaction history syncing")
		return
	}

	s.log.Info("Starting transaction history sync cycle",
		logger.Int("address_count", len(addresses)))

	for _, address := range addresses {
		// Check if context is cancelled
		if ctx.Err() != nil {
			s.log.Info("Transaction history sync cycle interrupted",
				logger.Error(ctx.Err()))
			return
		}

		err := s.syncTransactionsForAddress(ctx, address)
		if err != nil {
			s.log.Error("Failed to sync history for address",
				logger.String("address", address.String()),
				logger.Error(err))
			continue
		}
	}

	s.log.Info("Completed transaction history sync cycle")
}

// syncTransactionsForAddressByType fetches and processes transactions of a specific type for an address
func (s *historyService) syncTransactionsForAddressByType(ctx context.Context, address types.Address, txType blockexplorer.TransactionType) error {
	// Get explorer for this chain
	explorer, err := s.blockExplorerFactory.NewExplorer(address.ChainType)
	if err != nil {
		s.log.Error("Failed to get explorer for chain",
			logger.String("chain", string(address.ChainType)),
			logger.Error(err))
		return errors.NewOperationFailedError("get explorer", err)
	}

	// Retrieve start block number using getStartBlock
	startBlock := s.getStartBlock(address)

	// Configure options for transaction history
	options := blockexplorer.TransactionHistoryOptions{
		TransactionType: txType,
		StartBlock:      startBlock.Int64(),
		Limit:           9000,
	}

	s.log.Info("Fetching transaction history",
		logger.String("address", address.String()),
		logger.String("chain", string(address.ChainType)),
		logger.String("tx_type", string(txType)),
		logger.Int64("start_block", startBlock.Int64()))

	// Fetch transaction history from explorer
	page, err := explorer.GetTransactionHistory(ctx, address.ToChecksum(), options, "")
	if err != nil {
		s.log.Error("Failed to get transaction history",
			logger.String("address", address.String()),
			logger.String("tx_type", string(txType)),
			logger.Error(err))
		return errors.NewOperationFailedError("get transaction history", err)
	}

	if len(page.Items) == 0 {
		s.log.Info("No transactions found for address",
			logger.String("address", address.String()),
			logger.String("tx_type", string(txType)))
		return nil
	}

	s.log.Info("Found transactions for address",
		logger.String("address", address.String()),
		logger.String("tx_type", string(txType)),
		logger.Int("count", len(page.Items)))

	// Process and save transactions
	for _, item := range page.Items {
		// Get the core transaction
		rawTx := item.GetTransaction()

		// Apply transformers
		transformedTx := s.transformerService.Apply(ctx, rawTx)
		if transformedTx == nil {
			s.log.Warn("Transaction is nil after transformation, skipping",
				logger.String("tx_hash", rawTx.Hash))
			continue
		}

		// For ERC20 transactions, verify that the token address exists in the token store
		if txType == blockexplorer.TxTypeERC20 {
			if !s.isValidERC20Token(ctx, transformedTx) {
				continue
			}
		}

		// Convert to service transaction
		serviceTx := FromCoreTransaction(transformedTx)
		if serviceTx == nil {
			s.log.Error("Failed to convert transaction to service model",
				logger.String("tx_hash", transformedTx.Hash))
			continue
		}

		// Save or update transaction
		existingTx, err := s.repository.GetByHash(ctx, serviceTx.Hash)
		if err != nil && !errors.IsError(err, errors.ErrCodeTransactionNotFound) {
			s.log.Error("Error checking for existing transaction",
				logger.String("tx_hash", serviceTx.Hash),
				logger.Error(err))
			continue
		}

		var isNewTransaction bool
		if existingTx != nil {
			// Update existing transaction with new data
			serviceTx.ID = existingTx.ID
			if err := s.repository.Update(ctx, serviceTx); err != nil {
				s.log.Error("Failed to update transaction",
					logger.String("tx_hash", serviceTx.Hash),
					logger.Error(err))
				continue
			}
			s.log.Debug("Updated existing transaction",
				logger.String("tx_hash", serviceTx.Hash))
			isNewTransaction = false
		} else {
			// Create new transaction
			if err := s.repository.Create(ctx, serviceTx); err != nil {
				s.log.Error("Failed to create transaction",
					logger.String("tx_hash", serviceTx.Hash),
					logger.Error(err))
				continue
			}
			s.log.Debug("Created new transaction",
				logger.String("tx_hash", serviceTx.Hash))
			isNewTransaction = true
		}

		// Emit transaction event
		event := &TransactionEvent{
			Transaction: transformedTx,
			IsNew:       isNewTransaction,
		}
		select {
		case s.historyEventsChan <- event:
			s.log.Debug("Emitted history transaction event",
				logger.String("tx_hash", transformedTx.Hash),
				logger.Bool("is_new", isNewTransaction))
		default:
			s.log.Warn("History events channel is full, dropping event",
				logger.String("tx_hash", transformedTx.Hash))
		}
	}

	// Update start block number to the latest block processed using setStartBlock
	latestBlockNumber := page.Items[len(page.Items)-1].GetTransaction().BlockNumber
	s.setStartBlock(address, new(big.Int).Add(latestBlockNumber, big.NewInt(1)))

	return nil
}

// isValidERC20Token checks if the ERC20 token in the transaction exists in the token store
// Returns true if the token is valid, false otherwise
func (s *historyService) isValidERC20Token(ctx context.Context, tx *types.Transaction) bool {
	tokenAddress, exists := tx.Metadata.GetAddress(types.ERC20TokenAddressMetadataKey)
	if !exists {
		s.log.Warn("ERC20 token address not found in metadata, skipping",
			logger.String("tx_hash", tx.Hash))
		return false
	}

	// Check if the token exists in the token store
	token, err := s.tokenStore.GetToken(ctx, tokenAddress.Hex())
	if token != nil && err == nil {
		s.log.Debug("Found token in token store",
			logger.String("tx_hash", tx.Hash),
			logger.String("token_address", tokenAddress.Hex()),
			logger.String("token_symbol", token.Symbol))
		return true
	}

	if errors.IsError(err, errors.ErrCodeResourceNotFound) {
		s.log.Warn("Token not found in token store, skipping",
			logger.String("tx_hash", tx.Hash),
			logger.String("token_address", tokenAddress.Hex()))
	} else {
		s.log.Error("Failed to get token from token store",
			logger.String("tx_hash", tx.Hash),
			logger.String("token_address", tokenAddress.Hex()),
			logger.Error(err))
	}

	return false
}

// syncTransactionsForAddress immediately syncs transaction history for a specific address
func (s *historyService) syncTransactionsForAddress(ctx context.Context, address types.Address) error {
	// First sync normal transactions
	if err := s.syncTransactionsForAddressByType(ctx, address, blockexplorer.TxTypeNormal); err != nil {
		s.log.Error("Failed to sync normal transactions for address",
			logger.String("address", address.String()),
			logger.Error(err))
		// Continue with other transaction types despite error
	}

	// Then sync ERC20 transactions
	if err := s.syncTransactionsForAddressByType(ctx, address, blockexplorer.TxTypeERC20); err != nil {
		s.log.Error("Failed to sync ERC20 transactions for address",
			logger.String("address", address.String()),
			logger.Error(err))
		// Continue with other transaction types despite error
	}

	// Could add support for other transaction types in the future
	// such as blockexplorer.TxTypeInternal, blockexplorer.TxTypeERC721

	return nil
}

// getStartBlock retrieves the start block number for a given address
func (s *historyService) getStartBlock(address types.Address) *big.Int {
	s.syncMutex.RLock()
	defer s.syncMutex.RUnlock()

	addrKey := s.deriveAddressKey(address)
	info, exists := s.syncAddresses[addrKey]
	if !exists {
		s.log.Warn("Address not found in syncAddresses",
			logger.String("address", address.String()))
		return big.NewInt(0) // Default to block 0 if not found
	}

	return info.StartBlock
}

// setStartBlock updates the start block number for a given address
func (s *historyService) setStartBlock(address types.Address, blockNumber *big.Int) {
	s.syncMutex.Lock()
	defer s.syncMutex.Unlock()

	addrKey := s.deriveAddressKey(address)
	info, exists := s.syncAddresses[addrKey]
	if !exists {
		s.log.Warn("Address not found in syncAddresses",
			logger.String("address", address.String()))
		return
	}

	info.StartBlock = blockNumber
	s.syncAddresses[addrKey] = info
}
