package nyquist

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/yawning/nyquist.git/cipher"
)

func TestCipherStateMalformedKey(t *testing.T) {
	require := require.New(t)
	cs := newCipherState(cipher.ChaChaPoly, DefaultMaxMessageSize)

	var oversizedKey [33]byte
	err := cs.setKey(oversizedKey[:])
	require.Error(err, "cs.setKey(oversized)")
	require.Panics(func() {
		cs.InitializeKey(oversizedKey[:])
	}, "cs.InitializeKey(oversized)")

	var undersizedKey [16]byte
	err = cs.setKey(undersizedKey[:])
	require.Error(err, "cs.setKey(undersized)")
	require.Panics(func() {
		cs.InitializeKey(undersizedKey[:])
	}, "cs.InitializeKey(undersized)")
}

func TestCipherStateExhaustedNonce(t *testing.T) {
	require := require.New(t)
	cs := newCipherState(cipher.ChaChaPoly, DefaultMaxMessageSize)

	var testKey [32]byte
	cs.InitializeKey(testKey[:])
	cs.SetNonce(maxnonce)

	ciphertext, err := cs.EncryptWithAd(nil, nil, []byte("exhausted nonce plaintext"))
	require.Equal(ErrNonceExhausted, err, "cs.EncryptWithAd() - exhauted nonce")
	require.Nil(ciphertext, "cs.EncryptWithAd() - exhauted nonce")

	plaintext, err := cs.DecryptWithAd(nil, nil, []byte("exhausted nonce ciphertext"))
	require.Equal(ErrNonceExhausted, err, "cs.DecryptWithAd() - exhauted nonce")
	require.Nil(plaintext, "cs.DecryptWithAd() - exhausted nonce")
}

func TestCipherStateMaxMessageSize(t *testing.T) {
	require := require.New(t)
	cs := newCipherState(cipher.DeoxysII, DefaultMaxMessageSize)

	var testKey [32]byte
	cs.InitializeKey(testKey[:])

	// The max message size includes the tag.
	ciphertext, err := cs.EncryptWithAd(nil, nil, make([]byte, DefaultMaxMessageSize-15))
	require.Equal(ErrMessageSize, err, "cs.EncryptWithAd(oversized)")
	require.Nil(ciphertext, "cs.EncryptWithAd(oversized)")

	maxPlaintext := make([]byte, DefaultMaxMessageSize-16)
	ciphertext, err = cs.EncryptWithAd(nil, nil, maxPlaintext)
	require.NoError(err, "cs.EncryptWithAd(maxMessageSize-tagLen)")
	require.NotNil(ciphertext, "cs.EncryptWithAd(maxMessageSize-tagLen)")

	plaintext, err := cs.DecryptWithAd(nil, nil, make([]byte, DefaultMaxMessageSize+1))
	require.Equal(ErrMessageSize, err, "cs.DecryptWithAd(oversized)")
	require.Nil(plaintext, "cs.DecryptWithAd(oversized")

	cs.SetNonce(0)
	plaintext, err = cs.DecryptWithAd(nil, nil, ciphertext)
	require.NoError(err, "cs.DecryptWithAd(maxMessageSize)")
	require.Equal(maxPlaintext, plaintext, "cs.DecryptWithAd(maxMessageSize)")
}

func TestCipherStateRekey(t *testing.T) {
	require := require.New(t)
	cs := newCipherState(cipher.ChaChaPoly, DefaultMaxMessageSize)

	err := cs.Rekey()
	require.Error(err, "cs.Rekey() - no key")

	testPlaintext := []byte("rekey test plaintext")

	var testKey [32]byte
	cs.InitializeKey(testKey[:])
	ciphertext, err := cs.EncryptWithAd(nil, nil, testPlaintext)
	require.NoError(err, "cs.EncryptWithAd()")

	err = cs.Rekey()
	require.NoError(err, "cs.Rekey()")

	cs.SetNonce(0)
	newCiphertext, err := cs.EncryptWithAd(nil, nil, testPlaintext)
	require.NoError(err, "cs.EncryptWithAd() - rekeyed")
	require.NotEqual(ciphertext, newCiphertext, "rekey actually changed key")
}

func TestCipherStateReset(t *testing.T) {
	// The main purpose of this test is to exercise the code that invokes
	// cipher.Resetable().
	require := require.New(t)
	cs := newCipherState(cipher.DeoxysII, DefaultMaxMessageSize)

	var testKey [32]byte
	cs.InitializeKey(testKey[:])
	cs.Reset()

	require.Nil(cs.aead, "cs.Reset()")
}

func TestCipherStateAuth(t *testing.T) {
	require := require.New(t)
	cs := newCipherState(cipher.DeoxysII, DefaultMaxMessageSize)

	testPlaintext := []byte("auth test plaintext")

	var testKey [32]byte
	cs.InitializeKey(testKey[:])
	ciphertext, err := cs.EncryptWithAd(nil, nil, testPlaintext)
	require.NoError(err, "cs.EncryptWithAd()")

	cs.SetNonce(0)
	_, err = cs.DecryptWithAd(nil, []byte("bogus ad"), ciphertext)
	require.Equal(ErrOpen, err, "cs.DecryptWithAd(bogus ad)")

	ciphertext[0] ^= 0xa5
	_, err = cs.DecryptWithAd(nil, nil, ciphertext)
	require.Equal(ErrOpen, err, "cs.DecryptWithAd(tampered ciphertext)")

	ciphertext[0] ^= 0xa5
	plaintext, err := cs.DecryptWithAd(nil, nil, ciphertext)
	require.NoError(err, "cs.DecryptWithAd()")
	require.Equal(testPlaintext, plaintext, "cs.DecryptWithAd()")
}
