# ATM-CLI

![Go Version](https://img.shields.io/badge/Go-1.21+-blue?logo=go)
![License](https://img.shields.io/github/license/stevanu/atm-CLI)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen)

🚀 **ATM-CLI** is a simple Command Line Interface (CLI) simulation of an ATM system, built with **Golang** and **MySQL**.  
This project demonstrates CRUD operations, database transactions, user authentication, and real-time balance updates.

---

## ✨ Features

- ✅ Register a new account with name, PIN (bcrypt hashed), and initial balance
- ✅ Login with name and PIN
- ✅ Check current balance
- ✅ Deposit money
- ✅ Withdraw money
- ✅ Transfer to another account
- ✅ Transaction history stored in MySQL
- ✅ All transactions instantly update balances in database

---

## ⚙️ Tech Stack

- **Language:** Go (Golang)
- **Database:** MySQL
- **Libraries:**
  - [`github.com/go-sql-driver/mysql`](https://pkg.go.dev/github.com/go-sql-driver/mysql) - MySQL driver
  - [`golang.org/x/crypto/bcrypt`](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - PIN hashing

---

## 📂 Project Structure

atm-CLI/
│
├── main.go # Entry point
├── database.go # DB connection & queries
├── models.go # Data models & structs
├── utils.go # Helpers (validation, bcrypt)
├── go.mod # Go module
└── go.sum # Checksums


---

## 🔄 Basic Flow

```plaintext
+-------------+         +----------+         +------------+
|   Register  |  --->   |   Login  |  --->   | Operations |
|   / Login   |         |          |         | (Balance,  |
|             |         |          |         | Deposit,   |
|             |         |          |         | Withdraw,  |
|             |         |          |         | Transfer)  |
+-------------+         +----------+         +------------+
```

🛠 Installation & Setup
1️⃣ Clone this repository
bash
Copy
Edit
git clone https://github.com/stevanu/atm-CLI.git
cd atm-CLI
2️⃣ Setup MySQL database
sql
Copy
Edit
CREATE DATABASE atm_db;
USE atm_db;

```
CREATE TABLE users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE,
  pin_hash VARCHAR(255) NOT NULL,
  balance DECIMAL(15,2) NOT NULL DEFAULT 0
);

CREATE TABLE transactions (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT,
  type ENUM('deposit','withdraw','transfer') NOT NULL,
  amount DECIMAL(15,2) NOT NULL,
  description VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id)
);
```
Update connection string in database.go:
```
db, err := sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/atm_db")
```

3️⃣ Run the app
```
go run main.go
```

🚀 How to Use
Register your account by entering a name, PIN, and initial balance.

Login with your name & PIN.

Access features: check balance, deposit, withdraw, transfer, view transaction history.

Logout when finished.
