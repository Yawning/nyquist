### nyquist - A Noise Protocol Framework implementation
#### Yawning Angel (yawning at schwanenlied dot me)

This package implements the [Noise Protocol Framework][1].

#### Why?

> Yeah, well, I'm gonna go build my own theme park with blackjack and
> hookers.  In fact, forget the park.

#### Notes

It is assumed that developers using this package are familiar with the Noise
Protocol Framework specification.

While this package attempts to be as complete of an implementation of the
specification as possible, certain features are currently unimplemented,
primarily due to lack of time on the part of the author.  It is belived
that it is possible to implement the missing functionality via the public
APIs.

 * 10.2. The `fallback` modifier

Care is taken to attempt to sanitize private key material from memory where
possible, however due to limitations in `crypto/cipher.AEAD`, `crypto/hkdf`,
`crypto/hmac`, and all of the hash functions, this is not particularly
comprehensive.

This package will `panic` only if invariants are violated.  Under normal
and correct use this situation should not occur ("correct" being defined as,
"Yes, it will panic if an invalid configuration is provided when initializing
a handshake").

Several "non-standard" cryptography libraries are used in lieu of runtime and
`x/crypto` equivalents.  If more "standard" implementations are desired it is
possible to implement the relevant cryptography functions using the external
interface.  The libraries and rationale are as follows:

 * [bsaes][2] Provides a constant time AES256-GCM.  The runtime library's
   implementation of both AES and GHASH is insecure on systems without
   hardware support and dedicated assembly language implementations.

 * [ed25519][3] Provides a significantly faster X25519 scalar basepoint
   multiply on supported platforms.

Several non-standard protocol extensions are supported by this implementation:

 * The maximum message size can be set to an arbitrary value or entirely
   disabled, on a per-session basis.  The implementation will default to
   the value in the specification.

 * AEAD algorithms with authentication tags that are not 128 bits (16 bytes)
   in size should be supported.

 * Non-standard DH, Cipher and Hash functions are trivial to support by
   implementing the appropriate interface, as long as the following
   constraints are met:

    * For any given DH scheme, all public keys must be DHLEN bytes in size.

    * HASHLEN must be at least 256 bits (32 bytes) in size.

    * AEAD implementations must be able to tollerate always being passed
      a key that is 256 bits (32 bytes) in size.

 * Non-standard (or unimplemented) patterns are trivial to support by
   implementing the appropriate interface.

#### TODO

 * Add a pattern validity checker.

 * Improve tests.

[1]: https://noiseprotocol.org/
[2]: https://gitlab.com/yawning/bsaes
[3]: https://github.com/oasislabs/ed25519
