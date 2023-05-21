# vault
Vault is a password store which stores your passwords in a secure way using aes encryption with GCM(Galois/Counter Mode). It also makes things easy to backup. All you have to do is to copy the `$HOME/.vault/.passwdstore` file to the same location in another system and everything works as expected.

## Installation
You need to have go toolchain to install this module. With proper go tool chain installation run the command `go install -v github.com/231tr0n/vault@latest`.

## Usage
Use the flag `-help` for knowing all the options of the vault.
```
  -change
    	Changes the vault password
  -clear
    	Clears all the passwords in the vault
  -delete string
    	Deletes the password in the vault
  -generate int
    	Generates a new random password of length given
  -get string
    	Gets the password from the vault
  -list
    	Lists all the passwords in the vault
  -put string
    	Puts the password in the vault
```

**Note:** You have to set the password using the `-change` argument initially since it is not set. Give an empty password when prompted for vault password.
