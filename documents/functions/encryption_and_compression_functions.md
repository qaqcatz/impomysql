# MySQL Document

https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html

# Overview

```shell
[1]: only tidb features
[2]: supported features of impo
[3]: unsupported features of impo

[2]AesDecrypt               = "aes_decrypt"
[2]AesEncrypt               = "aes_encrypt"
[2]Compress                 = "compress"
[1]Decode                   = "decode"
[1]DesDecrypt               = "des_decrypt"
[1]DesEncrypt               = "des_encrypt"
[1]Encode                   = "encode"
[1]Encrypt                  = "encrypt"
[2]MD5                      = "md5"
[1]OldPassword              = "old_password"
[1]PasswordFunc             = "password_func"
[3]RandomBytes              = "random_bytes"
[2]SHA1                     = "sha1"
[2]SHA                      = "sha"
[2]SHA2                     = "sha2"
[2]Uncompress               = "uncompress"
[2]UncompressedLength       = "uncompressed_length"
[2]ValidatePasswordStrength = "validate_password_strength"
```

# aes_decrypt

`(crypt_str,key_str[,init_vector][,kdf_name][,salt][,info | iterations])`

This function decrypts data using the official AES (Advanced Encryption Standard) algorithm. For more information, see the description of [`AES_ENCRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-encrypt).

Statements that use [`AES_DECRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-decrypt) are unsafe for statement-based replication.

see aes_encrypt.

# aes_encrypt

`(str,key_str[,init_vector][,kdf_name][,salt][,info | iterations])`

[`AES_ENCRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-encrypt) and [`AES_DECRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-decrypt) implement encryption and decryption of data using the official AES (Advanced Encryption Standard) algorithm, previously known as “Rijndael.” The AES standard permits various key lengths. By default these functions implement AES with a 128-bit key length. Key lengths of 196 or 256 bits can be used, as described later. The key length is a trade off between performance and security.

[`AES_ENCRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-encrypt) encrypts the string *`str`* using the key string *`key_str`*, and returns a binary string containing the encrypted output. [`AES_DECRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-decrypt) decrypts the encrypted string *`crypt_str`* using the key string *`key_str`*, and returns the original plaintext string. If either function argument is `NULL`, the function returns `NULL`. If [`AES_DECRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-decrypt) detects invalid data or incorrect padding, it returns `NULL`. However, it is possible for [`AES_DECRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-decrypt) to return a non-`NULL` value (possibly garbage) if the input data or the key is invalid.

From MySQL 8.0.30, the functions support the use of a key derivation function (KDF) to create a cryptographically strong secret key from the information passed in *`key_str`*. The derived key is used to encrypt and decrypt the data, and it remains in the MySQL Server instance and is not accessible to users. Using a KDF is highly recommended, as it provides better security than specifying your own premade key or deriving it by a simpler method as you use the function. The functions support HKDF (available from OpenSSL 1.1.0), for which you can specify an optional salt and context-specific information to include in the keying material, and PBKDF2 (available from OpenSSL 1.0.2), for which you can specify an optional salt and set the number of iterations used to produce the key.

[`AES_ENCRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-encrypt) and [`AES_DECRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-decrypt) permit control of the block encryption mode. The [`block_encryption_mode`](https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html#sysvar_block_encryption_mode) system variable controls the mode for block-based encryption algorithms. Its default value is `aes-128-ecb`, which signifies encryption using a key length of 128 bits and ECB mode. For a description of the permitted values of this variable, see [Section 5.1.8, “Server System Variables”](https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html). The optional *`init_vector`* argument is used to provide an initialization vector for block encryption modes that require it.

Statements that use [`AES_ENCRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-encrypt) or [`AES_DECRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-decrypt) are unsafe for statement-based replication.

If [`AES_ENCRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-encrypt) is invoked from within the [**mysql**](https://dev.mysql.com/doc/refman/8.0/en/mysql.html) client, binary strings display using hexadecimal notation, depending on the value of the [`--binary-as-hex`](https://dev.mysql.com/doc/refman/8.0/en/mysql-command-options.html#option_mysql_binary-as-hex). For more information about that option, see [Section 4.5.1, “mysql — The MySQL Command-Line Client”](https://dev.mysql.com/doc/refman/8.0/en/mysql.html).

The arguments for the [`AES_ENCRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-encrypt) and [`AES_DECRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-decrypt) functions are as follows:

- **str**

  The string for [`AES_ENCRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-encrypt) to encrypt using the key string *`key_str`*, or (from MySQL 8.0.30) the key derived from it by the specified KDF. The string can be any length. Padding is automatically added to *`str`* so it is a multiple of a block as required by block-based algorithms such as AES. This padding is automatically removed by the [`AES_DECRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-decrypt) function.

- **crypt_str**

  The encrypted string for [`AES_DECRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-decrypt) to decrypt using the key string *`key_str`*, or (from MySQL 8.0.30) the key derived from it by the specified KDF. The string can be any length. The length of *`crypt_str`* can be calculated from the length of the original string using this formula:`16 * (trunc(*string_length* / 16) + 1)`

- **key_str**

  The encryption key, or the input keying material that is used as the basis for deriving a key using a key derivation function (KDF). For the same instance of data, use the same value of *`key_str`* for encryption with [`AES_ENCRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-encrypt) and decryption with [`AES_DECRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-decrypt).If you are using a KDF, which you can from MySQL 8.0.30, *`key_str`* can be any arbitrary information such as a password or passphrase. In the further arguments for the function, you specify the KDF name, then add further options to increase the security as appropriate for the KDF.When you use a KDF, the function creates a cryptographically strong secret key from the information passed in *`key_str`* and any salt or additional information that you provide in the other arguments. The derived key is used to encrypt and decrypt the data, and it remains in the MySQL Server instance and is not accessible to users. Using a KDF is highly recommended, as it provides better security than specifying your own premade key or deriving it by a simpler method as you use the function.If you are not using a KDF, for a key length of 128 bits, the most secure way to pass a key to the *`key_str`* argument is to create a truly random 128-bit value and pass it as a binary value. For example:`INSERT INTO t VALUES (1,AES_ENCRYPT('text',UNHEX('F3229A0B371ED2D9441B830D21A390C3')));`A passphrase can be used to generate an AES key by hashing the passphrase. For example:`INSERT INTO t VALUES (1,AES_ENCRYPT('text', UNHEX(SHA2('My secret passphrase',512))));`If you exceed the maximum key length of 128 bits, a warning is returned. If you are not using a KDF, do not pass a password or passphrase directly to *`key_str`*, hash it first. Previous versions of this documentation suggested the former approach, but it is no longer recommended as the examples shown here are more secure.

- **init_vector**

  An initialization vector, for block encryption modes that require it. The [`block_encryption_mode`](https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html#sysvar_block_encryption_mode) system variable controls the mode. For the same instance of data, use the same value of *`init_vector`* for encryption with [`AES_ENCRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-encrypt) and decryption with [`AES_DECRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-decrypt).NoteIf you are using a KDF, you must specify an initialization vector or a null string for this argument, in order to access the later arguments to define the KDF.For modes that require an initialization vector, it must be 16 bytes or longer (bytes in excess of 16 are ignored). An error occurs if *`init_vector`* is missing. For modes that do not require an initialization vector, it is ignored and a warning is generated if *`init_vector`* is specified, unless you are using a KDF.The default value for the [`block_encryption_mode`](https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html#sysvar_block_encryption_mode) system variable is `aes-128-ecb`, or ECB mode, which does not require an initialization vector. The alternative permitted block encryption modes CBC, CFB1, CFB8, CFB128, and OFB all require an initialization vector.A random string of bytes to use for the initialization vector can be produced by calling [`RANDOM_BYTES(16)`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_random-bytes).

- **kdf_name**

  The name of the key derivation function (KDF) to create a key from the input keying material passed in *`key_str`*, and other arguments as appropriate for the KDF. This optional argument is available from MySQL 8.0.30.For the same instance of data, use the same value of *`kdf_name`* for encryption with [`AES_ENCRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-encrypt) and decryption with [`AES_DECRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-decrypt). When you specify *`kdf_name`*, you must specify *`init_vector`*, using either a valid initialization vector, or a null string if the encryption mode does not require an initialization vector.The following values are supported:**`hkdf`**HKDF, which is available from OpenSSL 1.1.0. HKDF extracts a pseudorandom key from the keying material then expands it into additional keys. With HKDF, you can specify an optional salt (*`salt`*) and context-specific information such as application details (*`info`*) to include in the keying material.**`pbkdf2_hmac`**PBKDF2, which is available from OpenSSL 1.0.2. PBKDF2 applies a pseudorandom function to the keying material, and repeats this process a large number of times to produce the key. With PBKDF2, you can specify an optional salt (*`salt`*) to include in the keying material, and set the number of iterations used to produce the key (*`iterations`*).In this example, HKDF is specified as the key derivation function, and a salt and context information are provided. The argument for the initialization vector is included but is the empty string:`SELECT AES_ENCRYPT('mytext','mykeystring', '', 'hkdf', 'salt', 'info');`In this example, PBKDF2 is specified as the key derivation function, a salt is provided, and the number of iterations is doubled from the recommended minimum:`SELECT AES_ENCRYPT('mytext','mykeystring', '', 'pbkdf2_hmac','salt', '2000');`

- **salt**

  A salt to be passed to the key derivation function (KDF). This optional argument is available from MySQL 8.0.30. Both HKDF and PBKDF2 can use salts, and their use is recommended to help prevent attacks based on dictionaries of common passwords or rainbow tables.A salt consists of random data, which for security must be different for each encryption operation. A random string of bytes to use for the salt can be produced by calling [`RANDOM_BYTES()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_random-bytes). This example produces a 64-bit salt:`SET @salt = RANDOM_BYTES(8);`For the same instance of data, use the same value of *`salt`* for encryption with [`AES_ENCRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-encrypt) and decryption with [`AES_DECRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-decrypt). The salt can safely be stored along with the encrypted data.

- **info**

  Context-specific information for HKDF to include in the keying material, such as information about the application. This optional argument is available from MySQL 8.0.30 when you specify `hkdf` as the KDF name. HKDF adds this information to the keying material specified in *`key_str`* and the salt specified in *`salt`* to produce the key.For the same instance of data, use the same value of *`info`* for encryption with [`AES_ENCRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-encrypt) and decryption with [`AES_DECRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-decrypt).

- **iterations**

  The iteration count for PBKDF2 to use when producing the key. This optional argument is available from MySQL 8.0.30 when you specify `pbkdf2_hmac` as the KDF name. A higher count gives greater resistance to brute-force attacks because it has a greater computational cost for the attacker, but the same is necessarily true for the key derivation process. The default if you do not specify this argument is 1000, which is the minimum recommended by the OpenSSL standard.For the same instance of data, use the same value of *`iterations`* for encryption with [`AES_ENCRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-encrypt) and decryption with [`AES_DECRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-decrypt).

For the same instance of data, use the same value of *`iterations`* for encryption with [`AES_ENCRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-encrypt) and decryption with [`AES_DECRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-decrypt).

```sql
mysql> SET block_encryption_mode = 'aes-256-cbc';
mysql> SET @key_str = SHA2('My secret passphrase',512);
mysql> SET @init_vector = RANDOM_BYTES(16);
mysql> SET @crypt_str = AES_ENCRYPT('text',@key_str,@init_vector);
mysql> SELECT AES_DECRYPT(@crypt_str,@key_str,@init_vector);
+-----------------------------------------------+
| AES_DECRYPT(@crypt_str,@key_str,@init_vector) |
+-----------------------------------------------+
| text                                          |
+-----------------------------------------------+
```

# compress

Compresses a string and returns the result as a binary string. This function requires MySQL to have been compiled with a compression library such as `zlib`. Otherwise, the return value is always `NULL`. The return value is also `NULL` if *`string_to_compress`* is `NULL`. The compressed string can be uncompressed with [`UNCOMPRESS()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_uncompress).

```sql
mysql> SELECT COMPRESS('hello world');
+--------------------------------------------------+
| COMPRESS('hello world')                          |
+--------------------------------------------------+
| 0x0B000000789CCB48CDC9C95728CF2FCA4901001A0B045D |
+--------------------------------------------------+
1 row in set (0.00 sec)

mysql> SELECT COMPRESS(123);
+----------------------------------+
| COMPRESS(123)                    |
+----------------------------------+
| 0x03000000789C3334320600012D0097 |
+----------------------------------+
1 row in set (0.00 sec)

mysql> SELECT COMPRESS(NULL);
+--------------------------------+
| COMPRESS(NULL)                 |
+--------------------------------+
| NULL                           |
+--------------------------------+
1 row in set (0.00 sec)
```

# decode

> not found

# des_decrypt

> not found

# des_encrypt

> not found

# encode

> not found

# encrypt

> not found

# md5

Calculates an MD5 128-bit checksum for the string. The value is returned as a string of 32 hexadecimal digits, or `NULL` if the argument was `NULL`. The return value can, for example, be used as a hash key. See the notes at the beginning of this section about storing hash values efficiently.

The return value is a string in the connection character set.

```sql
mysql> SELECT MD5('testing');
+----------------------------------+
| MD5('testing')                   |
+----------------------------------+
| ae2b1fca515949e5d54fb22b8ed95575 |
+----------------------------------+
1 row in set (0.00 sec)

mysql> SELECT MD5(123);
+----------------------------------+
| MD5(123)                         |
+----------------------------------+
| 202cb962ac59075b964b07152d234b70 |
+----------------------------------+
1 row in set (0.00 sec)

mysql> SELECT MD5(NULL);
+-----------+
| MD5(NULL) |
+-----------+
| NULL      |
+-----------+
1 row in set (0.01 sec)
```

# old_password

> not found

# password_func

> not found

# random_bytes

This function returns a binary string of *`len`* random bytes generated using the random number generator of the SSL library. Permitted values of *`len`* range from 1 to 1024. For values outside that range, an error occurs. Returns `NULL` if *`len`* is `NULL`.

[`RANDOM_BYTES()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_random-bytes) can be used to provide the initialization vector for the [`AES_DECRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-decrypt) and [`AES_ENCRYPT()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_aes-encrypt) functions. For use in that context, *`len`* must be at least 16. Larger values are permitted, but bytes in excess of 16 are ignored.

[`RANDOM_BYTES()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_random-bytes) generates a random value, which makes its result nondeterministic. Consequently, statements that use this function are unsafe for statement-based replication.

If [`RANDOM_BYTES()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_random-bytes) is invoked from within the [**mysql**](https://dev.mysql.com/doc/refman/8.0/en/mysql.html) client, binary strings display using hexadecimal notation, depending on the value of the [`--binary-as-hex`](https://dev.mysql.com/doc/refman/8.0/en/mysql-command-options.html#option_mysql_binary-as-hex). For more information about that option, see [Section 4.5.1, “mysql — The MySQL Command-Line Client”](https://dev.mysql.com/doc/refman/8.0/en/mysql.html).

```sql
mysql> SELECT random_bytes(123);
+----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| random_bytes(123)                                                                                                                                                                                                                                        |
+----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| 0xDF454B234C44ED1E4CBA4238724A3D0CD26F9FD4AC91BC77B381A91A68AE7D9FE2CD7C2E6EBD4921684F839B6266B3BA3B3D33659647FBD5006E81606A93F337B6E7C6B051BBD98F9CB3F3A777108CCCC3FEEAB41D9D139FC426B2B3965F177CF62149BDA0979CEB68E6F86FDDF33E2452085FD892A71BE047E925 |
+----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
1 row in set (0.00 sec)

mysql> SELECT random_bytes(123);
+----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| random_bytes(123)                                                                                                                                                                                                                                        |
+----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| 0xC938CB5FFA7D4704F07CAA8FCC4530954366C2DDF43063FB3655C97B3CF41504B6758445C30B708E27183F6BCAEA76D1556C3784A7B469EA2347576A5E296D7785A82F617BC9CE2DCA50DF9BE5E6FA78CD96DCD52F6F59D049011971461244846FECDD43E116C31A03419D2C1F6B1B12160CC58DD4BCE2F8C4E997 |
+----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
1 row in set (0.00 sec)
```

> uncertain, do not support

# sha1

Calculates an SHA-1 160-bit checksum for the string, as described in RFC 3174 (Secure Hash Algorithm). The value is returned as a string of 40 hexadecimal digits, or `NULL` if the argument is `NULL`. One of the possible uses for this function is as a hash key. See the notes at the beginning of this section about storing hash values efficiently. [`SHA()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_sha1) is synonymous with [`SHA1()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_sha1).

The return value is a string in the connection character set.

```sql
mysql>  SELECT SHA1('abc');
+------------------------------------------+
| SHA1('abc')                              |
+------------------------------------------+
| a9993e364706816aba3e25717850c26c9cd0d89d |
+------------------------------------------+
1 row in set (0.00 sec)

mysql>  SELECT SHA1(123);
+------------------------------------------+
| SHA1(123)                                |
+------------------------------------------+
| 40bd001563085fc35165329ea1ff5c5ecbdbbeef |
+------------------------------------------+
1 row in set (0.00 sec)

mysql>  SELECT SHA1(NULL);
+------------+
| SHA1(NULL) |
+------------+
| NULL       |
+------------+
1 row in set (0.00 sec)
```

# sha

see sha1

# sha2

Calculates the SHA-2 family of hash functions (SHA-224, SHA-256, SHA-384, and SHA-512). The first argument is the plaintext string to be hashed. The second argument indicates the desired bit length of the result, which must have a value of 224, 256, 384, 512, or 0 (which is equivalent to 256). If either argument is `NULL` or the hash length is not one of the permitted values, the return value is `NULL`. Otherwise, the function result is a hash value containing the desired number of bits. See the notes at the beginning of this section about storing hash values efficiently.

The return value is a string in the connection character set.

```sql
mysql> SELECT SHA2('abc', 224);
+----------------------------------------------------------+
| SHA2('abc', 224)                                         |
+----------------------------------------------------------+
| 23097d223405d8228642a477bda255b32aadbce4bda0b3f7e36c9da7 |
+----------------------------------------------------------+
1 row in set (0.00 sec)

mysql> SELECT SHA2('abc', 0);
+------------------------------------------------------------------+
| SHA2('abc', 0)                                                   |
+------------------------------------------------------------------+
| ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad |
+------------------------------------------------------------------+
1 row in set (0.00 sec)

mysql> SELECT SHA2('abc', -1);
+-----------------+
| SHA2('abc', -1) |
+-----------------+
| NULL            |
+-----------------+
1 row in set, 1 warning (0.00 sec)

mysql> SELECT SHA2('abc', 512);
+----------------------------------------------------------------------------------------------------------------------------------+
| SHA2('abc', 512)                                                                                                                 |
+----------------------------------------------------------------------------------------------------------------------------------+
| ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a2192992a274fc1a836ba3c23a3feebbd454d4423643ce80e2a9ac94fa54ca49f |
+----------------------------------------------------------------------------------------------------------------------------------+
1 row in set (0.00 sec)
```

# uncompress

Uncompresses a string compressed by the [`COMPRESS()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_compress) function. If the argument is not a compressed value, the result is `NULL`; if *`string_to_uncompress`* is `NULL`, the result is also `NULL`. This function requires MySQL to have been compiled with a compression library such as `zlib`. Otherwise, the return value is always `NULL`.

```sql
mysql> SELECT UNCOMPRESSED_LENGTH('HELLO WORLD');
+------------------------------------+
| UNCOMPRESSED_LENGTH('HELLO WORLD') |
+------------------------------------+
|                          206325064 |
+------------------------------------+
1 row in set (0.01 sec)

mysql> SELECT UNCOMPRESSED_LENGTH('123456');
+-------------------------------+
| UNCOMPRESSED_LENGTH('123456') |
+-------------------------------+
|                     875770417 |
+-------------------------------+
1 row in set (0.00 sec)
```

# uncompressed_length

Returns the length that the compressed string had before being compressed. Returns `NULL` if *`compressed_string`* is `NULL`.

```sql
mysql> SELECT UNCOMPRESSED_LENGTH('HELLO WORLD');
+------------------------------------+
| UNCOMPRESSED_LENGTH('HELLO WORLD') |
+------------------------------------+
|                          206325064 |
+------------------------------------+
1 row in set (0.00 sec)

mysql> SELECT UNCOMPRESSED_LENGTH(123456);
+-----------------------------+
| UNCOMPRESSED_LENGTH(123456) |
+-----------------------------+
|                   875770417 |
+-----------------------------+
1 row in set (0.00 sec)
```

# validate_password_strength

Given an argument representing a plaintext password, this function returns an integer to indicate how strong the password is, or `NULL` if the argument is `NULL`. The return value ranges from 0 (weak) to 100 (strong).

Password assessment by [`VALIDATE_PASSWORD_STRENGTH()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_validate-password-strength) is done by the `validate_password` component. If that component is not installed, the function always returns 0. For information about installing `validate_password`, see [Section 6.4.3, “The Password Validation Component”](https://dev.mysql.com/doc/refman/8.0/en/validate-password.html). To examine or configure the parameters that affect password testing, check or set the system variables implemented by `validate_password`. See [Section 6.4.3.2, “Password Validation Options and Variables”](https://dev.mysql.com/doc/refman/8.0/en/validate-password-options-variables.html).

The password is subjected to increasingly strict tests and the return value reflects which tests were satisfied, as shown in the following table. In addition, if the [`validate_password.check_user_name`](https://dev.mysql.com/doc/refman/8.0/en/validate-password-options-variables.html#sysvar_validate_password.check_user_name) system variable is enabled and the password matches the user name, [`VALIDATE_PASSWORD_STRENGTH()`](https://dev.mysql.com/doc/refman/8.0/en/encryption-functions.html#function_validate-password-strength) returns 0 regardless of how other `validate_password` system variables are set.

| Password Test                                                | Return Value |
| :----------------------------------------------------------- | :----------- |
| Length < 4                                                   | 0            |
| Length ≥ 4 and < [`validate_password.length`](https://dev.mysql.com/doc/refman/8.0/en/validate-password-options-variables.html#sysvar_validate_password.length) | 25           |
| Satisfies policy 1 (`LOW`)                                   | 50           |
| Satisfies policy 2 (`MEDIUM`)                                | 75           |
| Satisfies policy 3 (`STRONG`)                                | 100          |

```sql
mysql> SELECT validate_password_strength('hklSD3#@$123sdfsdfsd5343#!@#OISOIUDOUOLKJSD');
+---------------------------------------------------------------------------+
| validate_password_strength('hklSD3#@$123sdfsdfsd5343#!@#OISOIUDOUOLKJSD') |
+---------------------------------------------------------------------------+
|                                                                         0 |
+---------------------------------------------------------------------------+
1 row in set (0.00 sec)

```





