# Pass

A simple command-line password manager written in Go.

## Features

- AES-256-GCM encryption for all stored credentials  
- Argon2id key derivation from a user-provided master password  
- Encrypted JSON-based vault storage  
- Random salt and nonce generation  
- Command-line interface  

## Security Model

### Master Password

The master password is **Automatically generated** and is required to access the vault.
For convenience, the application can optionally persist the derived key in Master.key. This is intended for personal use

### Key Derivation

Encryption keys are derived using **Argon2id**, a memory-hard password-based key derivation function designed to resist brute-force attacks.

Current parameters:

- Iterations: 5  
- Memory: 2048 MB  
- Parallelism: 4  
- Output key length: 32 bytes (256 bits)

### Encryption

Vault entries are encrypted using **AES-256 in Galois/Counter Mode (GCM)** 

## Technologies

- Go  
- AES-256-GCM  
- Argon2id  
