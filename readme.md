# Locksmith

[![Go Version](https://img.shields.io/badge/go-1.18+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Open Source Love](https://badges.frapsoft.com/os/v1/open-source.svg?v=103)](https://opensource.org/)

A simple, secure, and open-source password manager designed to run locally from your command line. All data is encrypted and stored on your local machine, ensuring you have full control over your sensitive information.

## About The Project

Locksmith is a command-line tool for securely storing and managing your passwords and other secrets. It uses strong, modern encryption standards to protect your data in a local "vault" file. With no cloud services involved, you can be confident that your data remains private.

**Core Features:**

*   **Strong Encryption:** Uses AES-256 for encryption and Argon2 for key derivation, protecting your vault against brute-force attacks.
*   **100% Local:** Your vault is stored on your local filesystem. No data is ever sent over the network.
*   **CLI Based:** Manage your secrets directly from the terminal for a fast and efficient workflow.
*   **Password Recovery:** A secure security question-based mechanism helps you regain access to your vault if you forget your master password.
*   **Clipboard Integration:** Automatically copies passwords to the clipboard for a short period to prevent them from being snooped.

## Getting Started

Follow these simple steps to get Locksmith up and running.

### Prerequisites

*   **Go:** You need to have Go version 1.18 or higher installed. You can download it from the [official Go website](https://go.dev/dl/).

### Installation & Building

1.  **Clone the repository:**
    ```sh
    git clone <your-repo-url>
    cd locksmith
    ```

2.  **Build the executable:**
    ```sh
    go build -o locksmith
    ```
    This will create a `locksmith` executable in the project directory. For convenience, you may want to move this to a directory in your system's PATH (e.g., `/usr/local/bin`).

## Usage

Locksmith is operated via a series of simple commands.

### 1. Initialize Your Vault
First, create your encrypted vault and set your master password.
```sh
./locksmith init
```
You will be prompted to create a master password and set up security questions for recovery.

### 2. Add a New Entry
Add a new password entry to your vault.
```sh
./locksmith add
```
You will be prompted for the site name, username, and password.

### 3. List All Entries
View a list of all sites you have stored.
```sh
./locksmith list
```

### 4. Get an Entry
Retrieve and view the details for a specific site. The password will be automatically copied to your clipboard for 30 seconds.
```sh
./locksmith get <site-name>
```

### 5. Delete an Entry
Permanently remove an entry from your vault.
```sh
./locksmith delete <site-name>
```

### 6. Change Master Password
Change your master password and re-encrypt the vault.
```sh
./locksmith change-master
```

### 7. Recover Your Vault
If you have forgotten your master password, you can use the recovery feature.
```sh
./locksmith recover
```

## Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".

1.  Fork the Project
2.  Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3.  Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4.  Push to the Branch (`git push origin feature/AmazingFeature`)
5.  Open a Pull Request

## License

Distributed under the MIT License. See `LICENSE` for more information.