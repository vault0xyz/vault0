package keystore

import (
	"context"
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"vault0/internal/config"
	coreCrypto "vault0/internal/core/crypto"
	"vault0/internal/core/keygen"
	"vault0/internal/db"
	"vault0/internal/errors"
	"vault0/internal/types"
)

// DBKeyStore implements the KeyStore interface using a local database
type DBKeyStore struct {
	db           *db.DB
	encryptor    coreCrypto.Encryptor
	keyGenerator keygen.KeyGenerator
	initialized  bool
}

// NewDBKeyStore creates a new DBKeyStore instance
func NewDBKeyStore(db *db.DB, cfg *config.Config) (*DBKeyStore, error) {
	if cfg.DBEncryptionKey == "" {
		return nil, errors.NewInvalidEncryptionKeyError("DB_ENCRYPTION_KEY environment variable is required")
	}

	// Create the encryptor
	encryptor, err := coreCrypto.NewAESEncryptorFromBase64(cfg.DBEncryptionKey)
	if err != nil {
		return nil, errors.NewEncryptionError(err)
	}

	return &DBKeyStore{
		db:           db,
		encryptor:    encryptor,
		keyGenerator: keygen.NewKeyGenerator(),
		initialized:  true,
	}, nil
}

// curveByName returns the elliptic.Curve instance for a given curve name
func curveByName(name string) (elliptic.Curve, error) {
	switch name {
	case types.CurveNameP256:
		return elliptic.P256(), nil
	case types.CurveNameSecp256k1:
		return coreCrypto.Secp256k1Curve, nil
	default:
		return nil, errors.NewInvalidCurveError("P-256 or secp256k1", name)
	}
}

// ECDSASignature represents an ECDSA signature's R and S values
type ECDSASignature struct {
	R, S *big.Int
}

// Create creates a new key with the given name and type
func (ks *DBKeyStore) Create(ctx context.Context, name string, keyType types.KeyType, curve elliptic.Curve, tags map[string]string) (*Key, error) {
	// Validate curve for ECDSA keys
	if keyType == types.KeyTypeECDSA {
		if curve == nil {
			curve = elliptic.P256() // Default to P-256 if no curve is specified
		}
	}

	// Convert tags to JSON
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return nil, errors.NewInvalidKeyError("failed to marshal tags", err)
	}

	// Generate a Snowflake ID for the key
	snowflakeID, err := ks.db.GenerateID()
	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}
	keyID := strconv.FormatInt(snowflakeID, 10)

	// Create the key
	key := &Key{
		ID:        keyID,
		Name:      name,
		Type:      keyType,
		Tags:      tags,
		CreatedAt: time.Now(),
		Curve:     curve,
	}

	// Generate cryptographic key material based on key type
	privateKey, publicKey, err := ks.keyGenerator.GenerateKeyPair(keyType, curve)
	if err != nil {
		return nil, err // Propagate error from keygen package
	}

	// Encrypt the private key before storing
	encryptedPrivateKey, err := ks.encryptor.Encrypt(privateKey)
	if err != nil {
		return nil, err // Propagate error from crypto package
	}

	// Set the key material
	key.PrivateKey = encryptedPrivateKey
	key.PublicKey = publicKey

	// Get curve name if applicable
	var curveName string
	if curve != nil {
		curveName = curve.Params().Name
	}

	// Insert the key into the database
	_, err = ks.db.ExecuteStatementContext(
		ctx,
		"INSERT INTO keys (id, name, key_type, curve, tags, created_at, private_key, public_key) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		key.ID,
		key.Name,
		string(key.Type),
		curveName,
		string(tagsJSON),
		key.CreatedAt.Unix(),
		key.PrivateKey,
		key.PublicKey,
	)
	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	return key, nil
}

// Import imports an existing key
func (ks *DBKeyStore) Import(ctx context.Context, name string, keyType types.KeyType, curve elliptic.Curve, privateKey, publicKey []byte, tags map[string]string) (*Key, error) {
	// Convert tags to JSON
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return nil, errors.NewInvalidKeyError("failed to marshal tags", err)
	}

	// Encrypt the private key
	encryptedPrivateKey, err := ks.encryptor.Encrypt(privateKey)
	if err != nil {
		return nil, err // Propagate error from crypto package
	}

	// Get curve name if applicable
	var curveName string
	if keyType == types.KeyTypeECDSA {
		if curve == nil {
			curve = elliptic.P256() // Default to P-256 if no curve is specified
		}
		curveName = curve.Params().Name
	} else if curve != nil {
		curveName = curve.Params().Name
	}

	// Generate a Snowflake ID for the key
	snowflakeID, err := ks.db.GenerateID()
	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}
	keyID := strconv.FormatInt(snowflakeID, 10)

	// Create the key
	key := &Key{
		ID:         keyID,
		Name:       name,
		Type:       keyType,
		Tags:       tags,
		CreatedAt:  time.Now(),
		PrivateKey: encryptedPrivateKey,
		PublicKey:  publicKey,
		Curve:      curve,
	}

	// Insert the key into the database
	_, err = ks.db.ExecuteStatementContext(
		ctx,
		"INSERT INTO keys (id, name, key_type, curve, tags, created_at, private_key, public_key) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		key.ID,
		key.Name,
		string(key.Type),
		curveName,
		string(tagsJSON),
		key.CreatedAt.Unix(),
		key.PrivateKey,
		key.PublicKey,
	)
	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	return key, nil
}

// GetPublicKey retrieves only the public part of a key by its ID
func (ks *DBKeyStore) GetPublicKey(ctx context.Context, id string) (*Key, error) {
	var (
		key       Key
		keyType   string
		tagsJSON  string
		curveName string
	)

	rows, err := ks.db.ExecuteQueryContext(
		ctx,
		"SELECT id, name, key_type, curve, tags, created_at, public_key FROM keys WHERE id = ?",
		id,
	)
	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.NewResourceNotFoundError("key", id)
	}

	err = rows.Scan(
		&key.ID,
		&key.Name,
		&keyType,
		&curveName,
		&tagsJSON,
		&key.CreatedAt,
		&key.PublicKey,
	)
	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	// Convert key type
	key.Type = types.KeyType(keyType)

	// Convert curve name to curve instance
	if curveName != "" {
		curve, err := curveByName(curveName)
		if err != nil {
			return nil, err
		}
		key.Curve = curve
	}

	// Parse tags JSON
	var tags map[string]string
	if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
		return nil, errors.NewInvalidKeyError("failed to unmarshal tags", err)
	}
	key.Tags = tags

	return &key, nil
}

// List retrieves keys in the keystore with pagination
func (ks *DBKeyStore) List(ctx context.Context, limit int, nextToken string) (*types.Page[*Key], error) {
	// Set default limit if not specified
	if limit <= 0 {
		limit = 50 // Default limit
	}

	// Default pagination column
	paginationColumn := "id"

	// Base query
	query := `SELECT id, name, key_type, curve, tags, created_at, public_key 
		FROM keys`

	args := []any{}

	// Handle token-based pagination
	if nextToken != "" {
		token, err := types.DecodeNextPageToken(nextToken, paginationColumn)
		if err != nil {
			return nil, err
		}

		if token != nil {
			query += " WHERE id > ?"
			args = append(args, token.Value)
		}
	}

	// Add ordering and limit
	query += " ORDER BY id ASC LIMIT ?"
	args = append(args, limit+1) // Fetch one extra to determine if there are more pages

	rows, err := ks.db.ExecuteQueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}
	defer rows.Close()

	var keys []*Key
	for rows.Next() {
		var (
			key       Key
			keyType   string
			tagsJSON  string
			curveName string
		)

		err := rows.Scan(
			&key.ID,
			&key.Name,
			&keyType,
			&curveName,
			&tagsJSON,
			&key.CreatedAt,
			&key.PublicKey,
		)
		if err != nil {
			return nil, errors.NewDatabaseError(err)
		}

		// Convert key type
		key.Type = types.KeyType(keyType)

		// Convert curve name to curve instance
		if curveName != "" {
			curve, err := curveByName(curveName)
			if err != nil {
				return nil, err
			}
			key.Curve = curve
		}

		// Parse tags JSON
		var tags map[string]string
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			return nil, errors.NewInvalidKeyError("failed to unmarshal tags", err)
		}
		key.Tags = tags

		keys = append(keys, &key)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	// Generate token function for pagination
	generateToken := func(key *Key) *types.NextPageToken {
		return &types.NextPageToken{
			Column: paginationColumn,
			Value:  key.ID,
		}
	}

	return types.NewPage(keys, limit, generateToken), nil
}

// Update modifies a key's metadata
func (ks *DBKeyStore) Update(ctx context.Context, id string, name string, tags map[string]string) (*Key, error) {
	// Convert tags to JSON
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return nil, errors.NewInvalidKeyError("failed to marshal tags", err)
	}

	// Update the key in the database
	result, err := ks.db.ExecuteStatementContext(
		ctx,
		"UPDATE keys SET name = ?, tags = ? WHERE id = ?",
		name,
		string(tagsJSON),
		id,
	)
	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	// Check if the key exists
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}
	if rowsAffected == 0 {
		return nil, errors.NewResourceNotFoundError("key", id)
	}

	// Return the updated key
	return ks.GetPublicKey(ctx, id)
}

// Delete removes a key from the keystore
func (ks *DBKeyStore) Delete(ctx context.Context, id string) error {
	result, err := ks.db.ExecuteStatementContext(ctx, "DELETE FROM keys WHERE id = ?", id)
	if err != nil {
		return errors.NewDatabaseError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError(err)
	}
	if rowsAffected == 0 {
		return errors.NewResourceNotFoundError("key", id)
	}

	return nil
}

// Sign performs a cryptographic signing operation using the specified key
func (ks *DBKeyStore) Sign(ctx context.Context, id string, data []byte, dataType DataType) ([]byte, error) {
	var (
		key       Key
		keyType   string
		curveName string
	)

	rows, err := ks.db.ExecuteQueryContext(
		ctx,
		"SELECT id, name, key_type, curve, private_key FROM keys WHERE id = ?",
		id,
	)
	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.NewResourceNotFoundError("key", id)
	}

	err = rows.Scan(
		&key.ID,
		&key.Name,
		&keyType,
		&curveName,
		&key.PrivateKey,
	)
	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	// Convert key type
	key.Type = types.KeyType(keyType)

	// Decrypt the private key
	privateKey, err := ks.encryptor.Decrypt(key.PrivateKey)
	if err != nil {
		return nil, err
	}

	// Sign the data
	return ks.signData(key.Type, privateKey, data, dataType, curveName)
}

// signData handles signing based on key type
func (ks *DBKeyStore) signData(keyType types.KeyType, privateKeyBytes, data []byte, dataType DataType, curveName string) ([]byte, error) {
	switch keyType {
	case types.KeyTypeECDSA:
		return ks.signWithECDSA(privateKeyBytes, data, dataType, curveName)
	case types.KeyTypeRSA:
		return ks.signWithRSA(privateKeyBytes, data, dataType)
	case types.KeyTypeEd25519:
		return ks.signWithEd25519(privateKeyBytes, data)
	case types.KeyTypeSymmetric:
		return ks.signWithHMAC(privateKeyBytes, data)
	default:
		return nil, errors.NewInvalidKeyTypeError(string(types.KeyTypeECDSA), string(keyType))
	}
}

// signWithECDSA signs data using an ECDSA private key
func (ks *DBKeyStore) signWithECDSA(privateKeyBytes, data []byte, dataType DataType, curveName string) ([]byte, error) {
	var privKey *ecdsa.PrivateKey
	var err error

	// Use custom unmarshal for secp256k1 curve
	if curveName == types.CurveNameSecp256k1 {
		privKey, err = coreCrypto.UnmarshalPrivateKey(privateKeyBytes)
		if err != nil {
			return nil, errors.NewInvalidKeyError("failed to parse secp256k1 private key", err)
		}
	} else {
		// For standard curves, use ParseECPrivateKey
		privKey, err = x509.ParseECPrivateKey(privateKeyBytes)
		if err != nil {
			return nil, errors.NewInvalidKeyError("failed to parse ECDSA private key", err)
		}
	}

	// Hash the data if needed
	var hash []byte
	if dataType == DataTypeRaw {
		h := sha256.Sum256(data)
		hash = h[:]
	} else {
		hash = data
	}

	// Sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash)
	if err != nil {
		return nil, errors.NewSigningError(err)
	}

	// Encode the signature
	signature := ECDSASignature{R: r, S: s}
	signatureBytes, err := asn1.Marshal(signature)
	if err != nil {
		return nil, errors.NewSigningError(err)
	}
	return signatureBytes, nil
}

// signWithRSA signs data using an RSA private key
func (ks *DBKeyStore) signWithRSA(privateKeyBytes, data []byte, dataType DataType) ([]byte, error) {
	// Parse the private key
	privKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		return nil, errors.NewInvalidKeyError("failed to parse RSA private key", err)
	}

	// Hash the data if needed
	var hash []byte
	if dataType == DataTypeRaw {
		h := sha256.Sum256(data)
		hash = h[:]
	} else {
		hash = data
	}

	// Sign the hash
	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hash)
	if err != nil {
		return nil, errors.NewSigningError(err)
	}
	return signature, nil
}

// signWithEd25519 signs data using an Ed25519 private key
func (ks *DBKeyStore) signWithEd25519(privateKeyBytes, data []byte) ([]byte, error) {
	// Parse the PKCS#8 private key
	key, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
	if err != nil {
		return nil, errors.NewInvalidKeyError("failed to parse Ed25519 private key", err)
	}

	// Convert to Ed25519 private key
	privKey, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.NewInvalidKeyError("invalid Ed25519 private key", nil)
	}

	// Sign the data (Ed25519 performs its own hashing)
	signature := ed25519.Sign(privKey, data)
	if signature == nil {
		return nil, errors.NewSigningError(fmt.Errorf("failed to sign with Ed25519"))
	}

	return signature, nil
}

// signWithHMAC signs data using an HMAC key
func (ks *DBKeyStore) signWithHMAC(keyBytes, data []byte) ([]byte, error) {
	// Create a new HMAC hasher
	h := hmac.New(sha256.New, keyBytes)

	// Write data to the hasher
	_, err := h.Write(data)
	if err != nil {
		return nil, errors.NewSigningError(err)
	}

	// Return the HMAC
	return h.Sum(nil), nil
}
