// main.go
package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Deklarasi variabel global
var db *sql.DB
var currentAccountID int
var currentAccountNumber string

// Struct Account untuk menyimpan data akun user
type Account struct {
	ID            int
	AccountNumber string
	Name          string
	Balance       float64
}

// Fungsi utama program
func main() {
	initDB()             // Koneksi ke database dan membuat tabel jika belum ada
	defer db.Close()     // Pastikan koneksi ditutup saat program selesai
	rand.Seed(time.Now().UnixNano()) // Untuk keperluan generate nomor rekening

	// Loop utama program (menu utama)
	for {
		clearScreen()
		fmt.Println("=== ATM SYSTEM ===")
		fmt.Println("1. Login")
		fmt.Println("2. Buat Akun Baru")
		fmt.Println("3. Keluar")
		fmt.Print("Pilihan: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			login()
		case 2:
			createAccount()
		case 3:
			fmt.Println("Terima kasih telah menggunakan ATM kami!")
			os.Exit(0)
		default:
			fmt.Println("Pilihan tidak valid!")
			waitForEnter()
		}
	}
}

// Inisialisasi koneksi database dan pembuatan tabel
func initDB() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/atm_db")
	if err != nil {
		logFatal("Gagal konek ke database", err)
	}

	err = db.Ping()
	if err != nil {
		logFatal("Database tidak merespon", err)
	}

	// Buat tabel accounts jika belum ada
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
			id INT AUTO_INCREMENT PRIMARY KEY,
			account_number VARCHAR(10) NOT NULL UNIQUE,
			name VARCHAR(100) NOT NULL,
			password VARCHAR(255) NOT NULL,
			pin CHAR(4) NOT NULL,
			balance DECIMAL(15,2) DEFAULT 0.00,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`)

	if err != nil {
		logFatal("Gagal membuat tabel accounts", err)
	}

	// Buat tabel transactions jika belum ada
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS transactions (
			id INT AUTO_INCREMENT PRIMARY KEY,
			account_id INT NOT NULL,
			type ENUM('deposit', 'withdraw', 'transfer') NOT NULL,
			amount DECIMAL(15,2) NOT NULL,
			recipient_id INT NULL,
			recipient_account VARCHAR(10) NULL,
			description VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (account_id) REFERENCES accounts(id)
		)`)

	if err != nil {
		logFatal("Gagal membuat tabel transactions", err)
	}
}

// Fungsi untuk generate nomor rekening unik
func generateAccountNumber() string {
	for {
		num := fmt.Sprintf("%08d", rand.Intn(100000000)) // 8 digit
		var exists bool
		db.QueryRow("SELECT EXISTS(SELECT 1 FROM accounts WHERE account_number = ?)", num).Scan(&exists)
		if !exists {
			return num
		}
	}
}

// Membuat akun baru
func createAccount() {
	clearScreen()
	fmt.Println("=== BUAT AKUN BARU ===")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Print("PIN (4 digit angka): ")
	pin, _ := reader.ReadString('\n')
	pin = strings.TrimSpace(pin)

	if len(pin) != 4 {
		fmt.Println("PIN harus 4 digit angka!")
		waitForEnter()
		return
	}

	accountNumber := generateAccountNumber()

	_, err := db.Exec(
		"INSERT INTO accounts (account_number, name, password, pin, balance) VALUES (?, ?, ?, ?, ?)",
		accountNumber, name, password, pin, 0,
	)

	if err != nil {
		fmt.Println("Gagal membuat akun:", err)
	} else {
		fmt.Printf("Akun berhasil dibuat!\nNomor Rekening Anda: %s\n", accountNumber)
	}
	waitForEnter()
}

// Fungsi login user ke sistem
func login() {
	clearScreen()
	fmt.Println("=== LOGIN ===")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Print("PIN (4 digit): ")
	pin, _ := reader.ReadString('\n')
	pin = strings.TrimSpace(pin)

	var account Account
	var dbPassword, dbPIN string

	// Ambil data dari database
	err := db.QueryRow(
		"SELECT id, account_number, name, password, pin, balance FROM accounts WHERE name = ?",
		name,
	).Scan(&account.ID, &account.AccountNumber, &account.Name, &dbPassword, &dbPIN, &account.Balance)

	if err != nil {
		fmt.Println("Login gagal: Username tidak ditemukan")
		waitForEnter()
		return
	}

	// Validasi password dan PIN
	if password != dbPassword || pin != dbPIN {
		fmt.Println("Login gagal: Password atau PIN salah")
		waitForEnter()
		return
	}

	currentAccountID = account.ID
	currentAccountNumber = account.AccountNumber
	fmt.Printf("Login berhasil!\nNomor Rekening Anda: %s\n", account.AccountNumber)
	waitForEnter()
	accountMenu(&account)
}

// Ambil saldo terbaru dari database
func refreshBalance(account *Account) {
	db.QueryRow("SELECT balance FROM accounts WHERE id = ?", account.ID).Scan(&account.Balance)
}

// Menu utama setelah login
func accountMenu(account *Account) {
	for {
		refreshBalance(account)
		clearScreen()
		fmt.Printf("=== HALO %s ===\n", account.Name)
		fmt.Printf("Nomor Rekening: %s\n", account.AccountNumber)
		fmt.Println("1. Cek Saldo")
		fmt.Println("2. Setor Tunai")
		fmt.Println("3. Tarik Tunai")
		fmt.Println("4. Transfer")
		fmt.Println("5. Riwayat Transaksi")
		fmt.Println("6. Logout")
		fmt.Print("Pilihan: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			checkBalance(account)
		case 2:
			deposit(account)
		case 3:
			withdraw(account)
		case 4:
			transferMenu(account)
		case 5:
			transactionHistory(account)
		case 6:
			currentAccountID = 0
			currentAccountNumber = ""
			return
		default:
			fmt.Println("Pilihan tidak valid!")
		}
		waitForEnter()
	}
}

// Tampilkan saldo akun
func checkBalance(account *Account) {
	clearScreen()
	refreshBalance(account)
	fmt.Println("=== CEK SALDO ===")
	fmt.Printf("Saldo Anda: Rp %.2f\n", account.Balance)
}

// Setor uang ke akun
func deposit(account *Account) {
	clearScreen()
	fmt.Println("=== SETOR TUNAI ===")
	fmt.Print("Jumlah setoran: Rp ")
	var amount float64
	fmt.Scanln(&amount)

	if amount <= 0 {
		fmt.Println("Jumlah harus lebih dari 0")
		return
	}

	_, err := db.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, account.ID)
	if err != nil {
		fmt.Println("Gagal menyetor:", err)
		return
	}

	// Simpan riwayat transaksi
	_, err = db.Exec(
		"INSERT INTO transactions (account_id, type, amount, description) VALUES (?, 'deposit', ?, ?)",
		account.ID, amount, "Setor tunai",
	)

	if err != nil {
		fmt.Println("Gagal mencatat transaksi:", err)
		return
	}

	fmt.Printf("Berhasil menyetor Rp %.2f\n", amount)
	refreshBalance(account)
}

// Tarik uang dari akun
func withdraw(account *Account) {
	clearScreen()
	fmt.Println("=== TARIK TUNAI ===")
	fmt.Print("Jumlah penarikan: Rp ")
	var amount float64
	fmt.Scanln(&amount)

	if amount <= 0 {
		fmt.Println("Jumlah harus lebih dari 0")
		return
	}

	refreshBalance(account)
	if amount > account.Balance {
		fmt.Println("Saldo tidak mencukupi")
		return
	}

	_, err := db.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, account.ID)
	if err != nil {
		fmt.Println("Gagal menarik uang:", err)
		return
	}

	// Simpan riwayat transaksi
	_, err = db.Exec(
		"INSERT INTO transactions (account_id, type, amount, description) VALUES (?, 'withdraw', ?, ?)",
		account.ID, amount, "Tarik tunai",
	)

	if err != nil {
		fmt.Println("Gagal mencatat transaksi:", err)
		return
	}

	fmt.Printf("Berhasil menarik Rp %.2f\n", amount)
	refreshBalance(account)
}

// Menu pilihan transfer
func transferMenu(account *Account) {
	clearScreen()
	fmt.Println("=== MENU TRANSFER ===")
	fmt.Println("1. Transfer dengan Nama Penerima")
	fmt.Println("2. Transfer dengan Nomor Rekening")
	fmt.Println("3. Kembali")
	fmt.Print("Pilihan: ")

	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		transferByName(account)
	case 2:
		transferByAccountNumber(account)
	case 3:
		return
	default:
		fmt.Println("Pilihan tidak valid!")
	}
}

// Transfer berdasarkan nama penerima
func transferByName(account *Account) {
	clearScreen()
	fmt.Println("=== TRANSFER (NAMA PENERIMA) ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Nama penerima: ")
	recipientName, _ := reader.ReadString('\n')
	recipientName = strings.TrimSpace(recipientName)

	if recipientName == account.Name {
		fmt.Println("Tidak bisa transfer ke diri sendiri")
		return
	}

	var recipientID int
	var recipientAccountNumber string
	err := db.QueryRow("SELECT id, account_number FROM accounts WHERE name = ?", recipientName).Scan(&recipientID, &recipientAccountNumber)
	if err != nil {
		fmt.Println("Penerima tidak ditemukan")
		return
	}

	processTransfer(account, recipientID, recipientAccountNumber, recipientName)
}

// Transfer berdasarkan nomor rekening
func transferByAccountNumber(account *Account) {
	clearScreen()
	fmt.Println("=== TRANSFER (NOMOR REKENING) ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Nomor rekening tujuan: ")
	accountNumber, _ := reader.ReadString('\n')
	accountNumber = strings.TrimSpace(accountNumber)

	if accountNumber == account.AccountNumber {
		fmt.Println("Tidak bisa transfer ke diri sendiri")
		return
	}

	var recipientID int
	var recipientName string
	err := db.QueryRow("SELECT id, name FROM accounts WHERE account_number = ?", accountNumber).Scan(&recipientID, &recipientName)
	if err != nil {
		fmt.Println("Nomor rekening tidak valid")
		return
	}

	processTransfer(account, recipientID, accountNumber, recipientName)
}

// Proses transfer dana (digunakan oleh 2 metode di atas)
func processTransfer(sender *Account, recipientID int, recipientAccountNumber, recipientName string) {
	fmt.Print("Jumlah transfer: Rp ")
	var amount float64
	fmt.Scanln(&amount)

	refreshBalance(sender)
	if amount <= 0 || amount > sender.Balance {
		fmt.Println("Jumlah tidak valid atau saldo tidak mencukupi")
		return
	}

	tx, err := db.Begin() // Mulai transaksi database
	if err != nil {
		fmt.Println("Gagal memulai transaksi:", err)
		return
	}

	// Kurangi saldo pengirim
	_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, sender.ID)
	if err != nil {
		tx.Rollback()
		fmt.Println("Gagal transfer:", err)
		return
	}

	// Tambah saldo penerima
	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, recipientID)
	if err != nil {
		tx.Rollback()
		fmt.Println("Gagal transfer:", err)
		return
	}

	// Simpan riwayat transaksi pengirim
	_, err = tx.Exec(
		"INSERT INTO transactions (account_id, type, amount, recipient_id, recipient_account, description) VALUES (?, 'transfer', ?, ?, ?, ?)",
		sender.ID, amount, recipientID, recipientAccountNumber,
		fmt.Sprintf("Transfer ke %s (%s)", recipientName, recipientAccountNumber),
	)
	if err != nil {
		tx.Rollback()
		fmt.Println("Gagal mencatat transaksi pengirim:", err)
		return
	}

	// Simpan riwayat transaksi penerima
	_, err = tx.Exec(
		"INSERT INTO transactions (account_id, type, amount, recipient_id, recipient_account, description) VALUES (?, 'transfer', ?, ?, ?, ?)",
		recipientID, amount, sender.ID, sender.AccountNumber,
		fmt.Sprintf("Transfer dari %s (%s)", sender.Name, sender.AccountNumber),
	)
	if err != nil {
		tx.Rollback()
		fmt.Println("Gagal mencatat transaksi penerima:", err)
		return
	}

	tx.Commit()
	fmt.Printf("Berhasil transfer Rp %.2f ke %s (%s)\n", amount, recipientName, recipientAccountNumber)
	refreshBalance(sender)
}

// Tampilkan 10 riwayat transaksi terakhir
func transactionHistory(account *Account) {
	clearScreen()
	fmt.Println("=== RIWAYAT TRANSAKSI ===")

	rows, err := db.Query(`
		SELECT type, amount, description, created_at
		FROM transactions
		WHERE account_id = ?
		ORDER BY created_at DESC LIMIT 10`, account.ID)
	if err != nil {
		fmt.Println("Gagal mengambil riwayat:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var txType, description, createdAt string
		var amount float64

		err = rows.Scan(&txType, &amount, &description, &createdAt)
		if err != nil {
			fmt.Println("Error membaca data:", err)
			continue
		}

		fmt.Printf("[%s] %s: Rp %.2f\n   %s\n", createdAt, strings.ToUpper(txType), amount, description)
	}
}

// Bersihkan layar terminal
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// Tunggu input Enter untuk melanjutkan
func waitForEnter() {
	fmt.Println("\nTekan Enter untuk melanjutkan...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// Tampilkan error fatal dan keluar program
func logFatal(message string, err error) {
	fmt.Printf("%s: %v\n", message, err)
	waitForEnter()
	os.Exit(1)
}
