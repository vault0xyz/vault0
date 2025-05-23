package signer

import (
	"context"

	"vault0/internal/errors"
	"vault0/internal/logger"
	"vault0/internal/types"
)

// Service defines the signer management operations interface
type Service interface {
	// CreateSigner creates a new signer with the given name and type
	// If it's an internal signer, userID must be provided
	// Returns an error if the operation fails
	CreateSigner(ctx context.Context, name string, signerType SignerType, userID *int64) (*Signer, error)

	// UpdateSigner modifies an existing signer's information
	// Returns an error if the signer doesn't exist or the operation fails
	UpdateSigner(ctx context.Context, id int64, name string, signerType SignerType, userID *int64) (*Signer, error)

	// DeleteSigner removes a signer by ID
	// Returns an error if the signer doesn't exist or the operation fails
	DeleteSigner(ctx context.Context, id int64) error

	// GetSignerByID retrieves a signer by ID
	// Returns an error if the signer doesn't exist or the operation fails
	GetSignerByID(ctx context.Context, id int64) (*Signer, error)

	// FindSignersByUserID retrieves all signers for a user
	// Returns an empty slice if no signers are found
	FindSignersByUserID(ctx context.Context, userID int64) ([]*Signer, error)

	// ListSigners returns a paginated collection of signers
	// Default limit is 10 if limit < 1
	// nextToken is used for token-based pagination (empty string for first page)
	ListSigners(ctx context.Context, limit int, nextToken string) (*types.Page[*Signer], error)

	// AddAddress creates a new address for a signer
	// Returns an error if the signer doesn't exist or the operation fails
	AddAddress(ctx context.Context, signerID int64, chainType, address string) (*Address, error)

	// DeleteAddress removes an address from a signer
	// Returns an error if the address doesn't exist or the operation fails
	DeleteAddress(ctx context.Context, id int64) error

	// GetAddresses retrieves all addresses for a signer
	// Returns an error if the signer doesn't exist or the operation fails
	GetAddresses(ctx context.Context, signerID int64) ([]*Address, error)
}

// service implements the Service interface
type service struct {
	log        logger.Logger
	repository Repository
}

// NewService creates a new signer service
func NewService(log logger.Logger, repository Repository) Service {
	return &service{
		log:        log,
		repository: repository,
	}
}

// CreateSigner creates a new signer
func (s *service) CreateSigner(ctx context.Context, name string, signerType SignerType, userID *int64) (*Signer, error) {
	// Validate input
	if name == "" {
		return nil, errors.NewInvalidInputError("Name is required", "name", "")
	}

	// For internal signers, userID must be provided
	if signerType == Internal && (userID == nil || *userID == 0) {
		return nil, errors.NewInvalidInputError("User ID is required for internal signers", "user_id", nil)
	}

	// Create signer entity
	signer := &Signer{
		Name:   name,
		Type:   signerType,
		UserID: userID,
	}

	// Save to repository
	err := s.repository.Create(ctx, signer)
	if err != nil {
		s.log.Error("Failed to create signer",
			logger.String("name", name),
			logger.String("type", string(signerType)),
			logger.Error(err),
		)
		return nil, err
	}

	s.log.Info("Created new signer",
		logger.Int64("id", signer.ID),
		logger.String("name", name),
		logger.String("type", string(signerType)),
	)

	return signer, nil
}

// UpdateSigner modifies an existing signer
func (s *service) UpdateSigner(ctx context.Context, id int64, name string, signerType SignerType, userID *int64) (*Signer, error) {
	// Validate input
	if id == 0 {
		return nil, errors.NewInvalidInputError("Signer ID is required", "id", 0)
	}

	if name == "" {
		return nil, errors.NewInvalidInputError("Name cannot be empty", "name", "")
	}

	// For internal signers, userID must be provided
	if signerType == Internal && (userID == nil || *userID == 0) {
		return nil, errors.NewInvalidInputError("User ID is required for internal signers", "user_id", nil)
	}

	// Get existing signer
	signer, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	signer.Name = name
	signer.Type = signerType
	signer.UserID = userID

	// Save to repository
	err = s.repository.Update(ctx, signer)
	if err != nil {
		s.log.Error("Failed to update signer",
			logger.Int64("id", id),
			logger.String("name", name),
			logger.Error(err),
		)
		return nil, err
	}

	s.log.Info("Updated signer",
		logger.Int64("id", id),
		logger.String("name", name),
		logger.String("type", string(signerType)),
	)

	return signer, nil
}

// DeleteSigner removes a signer
func (s *service) DeleteSigner(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.NewInvalidInputError("Signer ID is required", "id", 0)
	}

	err := s.repository.Delete(ctx, id)
	if err != nil {
		s.log.Error("Failed to delete signer",
			logger.Int64("id", id),
			logger.Error(err),
		)
		return err
	}

	s.log.Info("Deleted signer", logger.Int64("id", id))
	return nil
}

// GetSignerByID retrieves a signer by ID
func (s *service) GetSignerByID(ctx context.Context, id int64) (*Signer, error) {
	if id == 0 {
		return nil, errors.NewInvalidInputError("Signer ID is required", "id", 0)
	}

	signer, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return signer, nil
}

// FindSignersByUserID retrieves all signers for a user
func (s *service) FindSignersByUserID(ctx context.Context, userID int64) ([]*Signer, error) {
	if userID == 0 {
		return nil, errors.NewInvalidInputError("User ID is required", "user_id", 0)
	}

	signers, err := s.repository.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return signers, nil
}

// ListSigners returns a paginated collection of signers
func (s *service) ListSigners(ctx context.Context, limit int, nextToken string) (*types.Page[*Signer], error) {
	// Apply default pagination values if not specified
	// if limit <= 0, use default limit of 10
	if limit <= 0 {
		limit = 10
	}

	s.log.Debug("Listing signers",
		logger.Int("limit", limit),
		logger.String("next_token", nextToken),
	)

	return s.repository.List(ctx, limit, nextToken)
}

// AddAddress creates a new address for a signer
func (s *service) AddAddress(ctx context.Context, signerID int64, chainType, address string) (*Address, error) {
	// Validate input
	if signerID == 0 {
		return nil, errors.NewInvalidInputError("Signer ID is required", "signer_id", 0)
	}

	if chainType == "" {
		return nil, errors.NewInvalidInputError("Chain type is required", "chain_type", "")
	}

	if address == "" {
		return nil, errors.NewInvalidInputError("Address is required", "address", "")
	}

	// Verify signer exists
	_, err := s.repository.GetByID(ctx, signerID)
	if err != nil {
		return nil, err
	}

	// Create address entity
	newAddress := &Address{
		SignerID:  signerID,
		ChainType: chainType,
		Address:   address,
	}

	// Save to repository
	err = s.repository.AddAddress(ctx, newAddress)
	if err != nil {
		s.log.Error("Failed to add address",
			logger.Int64("signer_id", signerID),
			logger.String("chain_type", chainType),
			logger.String("address", address),
			logger.Error(err),
		)
		return nil, err
	}

	s.log.Info("Added address to signer",
		logger.Int64("signer_id", signerID),
		logger.String("chain_type", chainType),
		logger.String("address", address),
	)

	return newAddress, nil
}

// DeleteAddress removes an address from a signer
func (s *service) DeleteAddress(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.NewInvalidInputError("Address ID is required", "id", 0)
	}

	err := s.repository.DeleteAddress(ctx, id)
	if err != nil {
		s.log.Error("Failed to delete address",
			logger.Int64("id", id),
			logger.Error(err),
		)
		return err
	}

	s.log.Info("Deleted address", logger.Int64("id", id))
	return nil
}

// GetAddresses retrieves all addresses for a signer
func (s *service) GetAddresses(ctx context.Context, signerID int64) ([]*Address, error) {
	if signerID == 0 {
		return nil, errors.NewInvalidInputError("Signer ID is required", "signer_id", 0)
	}

	// Verify signer exists
	_, err := s.repository.GetByID(ctx, signerID)
	if err != nil {
		return nil, err
	}

	addresses, err := s.repository.GetAddresses(ctx, signerID)
	if err != nil {
		return nil, err
	}

	return addresses, nil
}
