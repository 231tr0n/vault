# Vault

[![Go Reference](https://pkg.go.dev/badge/github.com/231tr0n/vault/cmd/vault.svg)](https://pkg.go.dev/github.com/231tr0n/vault/cmd/vault)

Vault is a simple password manager with a single master password. It uses AES encryption with GCM(Galois/Counter Mode) to encrypt the passwords.

## Installation
With proper go installation, run the command `go install -v github.com/231tr0n/vault/cmd/vault@latest`.

## Usage
Use the flag `-help` for knowing all the options of the vault.
**Note:** You have to set the password using the `-change` argument initially since it is not set. Give an empty password when prompted for vault's old password.

## Backup
All you have to do is to copy the `$HOME/.vault/.passwdstore` file to the same location in another system and everything works as expected.

