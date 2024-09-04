# Hush: Secure CLI Password Manager

Hush is a command-line interface (CLI) password manager designed for security and ease of use.

## Security Features

- **Master Password**: Single point of access for all stored passwords
- **Argon2 Key Derivation**: Robust key derivation from the master password
- **AES-GCM Encryption**: State-of-the-art encryption for stored passwords
- **Secure Storage**: Encrypted passwords stored locally with restricted access

## Commands

### Initialize Hush
```bash
hush init
```
Set up Hush and create your master password.

### Add a Password
```bash
hush add <password-name>
```
Store a new password entry.

### List Passwords
```bash
hush list
```
Display all stored password names.

### Retrieve a Password
```bash
hush get <account-name> [flags]
```
Fetch a stored password.

Flags:
- `-d, --display`: Display the password instead of copying to clipboard

By default, the password is copied to the clipboard for security.

### Remove a Password
```bash
hush remove <password-name>
```
Delete a stored password.

### Generate a Password
```bash
hush generate [flags]
```
Create a strong, random password.

Flags:
- `-l, --length <number>`: Specify the length of the generated password (default: 16)
- `-d, --display`: Display the generated password instead of copying to clipboard

By default, the generated password is copied to the clipboard for convenience.

### Delete All Data
```bash
hush implode
```
Remove all stored data and the Hush directory.

## Usage Example

```bash
$ hush init
Enter your master password: ********
Hush initialized successfully!

$ hush generate -l 16 -d
Generated password: P@ssw0rd!2345678
Password copied to clipboard.
Do you want to save this password? (y/N): y
Enter a name for this password: example_password

$ hush list
example_password

$ hush get example_password
Enter your master password: ********
Password copied to clipboard.

$ hush remove example_password
Enter your master password: ********
Password removed successfully.
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

If you find Hush useful, please consider giving this repository a star ⭐. Your support helps to increase the visibility of this project and is greatly appreciated!
