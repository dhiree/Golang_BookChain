# Library Blockchain

This project is a simple blockchain implementation using Go, designed to manage book checkouts in a library system. It uses cryptographic hashing to ensure the integrity of the data stored in the blockchain.

## Features

- **Blockchain:** A linked list of blocks, each containing a checkout record.
- **Book Checkout:** Allows users to check out books and records the transactions on the blockchain.
- **Book Management:** Provides the ability to create new book records with unique IDs generated using MD5 hashing.
- **Genesis Block:** The first block in the blockchain is automatically generated.
- **API Endpoints:**
  - `GET /`: Retrieve the current state of the blockchain.
  - `POST /`: Add a new block to the blockchain by checking out a book.
  - `POST /new`: Create a new book record.

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/dhiree/Golang_BookChain.git
