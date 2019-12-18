
typedef unsigned hyper uint64;
typedef opaque uint256[32];

enum CryptoKeyType
{
    KEY_TYPE_ED25519 = 0,
    KEY_TYPE_MUXED_ED25519 = 256,
    KEY_TYPE_PRE_AUTH_TX = 1,
    KEY_TYPE_HASH_X = 2
};

enum PublicKeyType
{
    PUBLIC_KEY_TYPE_ED25519 = KEY_TYPE_ED25519
};

enum SignerKeyType
{
    SIGNER_KEY_TYPE_ED25519 = KEY_TYPE_ED25519,
    SIGNER_KEY_TYPE_PRE_AUTH_TX = KEY_TYPE_PRE_AUTH_TX,
    SIGNER_KEY_TYPE_HASH_X = KEY_TYPE_HASH_X
};

union PublicKey switch (PublicKeyType type)
{
case PUBLIC_KEY_TYPE_ED25519:
    uint256 ed25519;
};

union SignerKey switch (SignerKeyType type)
{
case SIGNER_KEY_TYPE_ED25519:
    uint256 ed25519;
case SIGNER_KEY_TYPE_PRE_AUTH_TX:
    /* SHA-256 Hash of TransactionSignaturePayload structure */
    uint256 preAuthTx;
case SIGNER_KEY_TYPE_HASH_X:
    /* Hash of random 256 bit preimage X */
    uint256 hashX;
};

// Source or destination of a payment operation
union MuxedAccount switch (CryptoKeyType type) {
 case KEY_TYPE_ED25519:
     uint256 ed25519;
 case KEY_TYPE_MUXED_ED25519:
     struct {
         uint64 id;
         uint256 ed25519;
     } med25519;
};

