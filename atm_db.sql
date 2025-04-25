-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Apr 25, 2025 at 08:04 PM
-- Server version: 10.4.32-MariaDB
-- PHP Version: 8.2.12

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `atm_db`
--

-- --------------------------------------------------------

--
-- Table structure for table `accounts`
--

CREATE TABLE `accounts` (
  `id` int(11) NOT NULL,
  `account_number` varchar(10) DEFAULT NULL,
  `name` varchar(100) NOT NULL,
  `password` varchar(255) NOT NULL,
  `pin` char(4) NOT NULL,
  `balance` decimal(15,2) DEFAULT 0.00,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ;

--
-- Dumping data for table `accounts`
--

INSERT INTO `accounts` (`id`, `account_number`, `name`, `password`, `pin`, `balance`, `created_at`, `updated_at`) VALUES
(1, NULL, 'John Doe', 'hashed_password_1', '1234', 998500.00, '2025-04-25 14:12:47', '2025-04-25 14:45:48'),
(2, NULL, 'Jane Smith', 'hashed_password_2', '5678', 500500.00, '2025-04-25 14:12:47', '2025-04-25 14:45:48'),
(3, '89043075', 'stevanu', '1111', '1111', 0.00, '2025-04-25 16:57:17', '2025-04-25 17:53:22'),
(4, '63705527', 'dika', '1111', '1111', 2000000.00, '2025-04-25 17:14:00', '2025-04-25 17:14:33');

-- --------------------------------------------------------

--
-- Table structure for table `transactions`
--

CREATE TABLE `transactions` (
  `id` int(11) NOT NULL,
  `account_id` int(11) NOT NULL,
  `type` enum('deposit','withdraw','transfer_out','transfer_in') NOT NULL,
  `amount` decimal(15,2) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `related_account_id` int(11) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ;

--
-- Dumping data for table `transactions`
--

INSERT INTO `transactions` (`id`, `account_id`, `type`, `amount`, `description`, `related_account_id`, `created_at`) VALUES
(1, 1, 'deposit', 500000.00, 'Initial deposit', NULL, '2025-04-25 14:12:47'),
(2, 2, 'deposit', 300000.00, 'Initial deposit', NULL, '2025-04-25 14:12:47'),
(3, 1, 'withdraw', 200000.00, 'ATM withdrawal', NULL, '2025-04-25 14:12:47'),
(4, 1, 'transfer_out', 100000.00, 'Transfer to Jane Smith', 2, '2025-04-25 14:12:47'),
(5, 2, 'transfer_in', 100000.00, 'Transfer from John Doe', 1, '2025-04-25 14:12:47'),
(6, 3, 'deposit', 10000000.00, 'Setor tunai', NULL, '2025-04-25 16:58:01'),
(7, 3, 'deposit', 10.00, 'Setor tunai', NULL, '2025-04-25 16:58:17'),
(8, 3, 'deposit', 200000.00, 'Setor tunai', NULL, '2025-04-25 17:03:36'),
(9, 3, 'withdraw', 10000.00, 'Tarik tunai', NULL, '2025-04-25 17:03:54'),
(10, 3, 'withdraw', 1000000.00, 'Tarik tunai', NULL, '2025-04-25 17:04:24'),
(11, 3, 'withdraw', 100000.00, 'Tarik tunai', NULL, '2025-04-25 17:11:08'),
(12, 3, 'withdraw', 1000000.00, 'Tarik tunai', NULL, '2025-04-25 17:11:15'),
(13, 4, 'deposit', 2000000.00, 'Setor tunai', NULL, '2025-04-25 17:14:33'),
(14, 3, 'withdraw', 1000000.00, 'Tarik tunai', NULL, '2025-04-25 17:52:29'),
(15, 3, 'deposit', 100000.00, 'Setor tunai', NULL, '2025-04-25 17:52:52'),
(16, 3, 'withdraw', 7190010.00, 'Tarik tunai', NULL, '2025-04-25 17:53:22');

-- --------------------------------------------------------

--
-- Table structure for table `transfer_logs`
--

CREATE TABLE `transfer_logs` (
  `id` int(11) NOT NULL,
  `sender_id` int(11) NOT NULL,
  `recipient_id` int(11) NOT NULL,
  `amount` decimal(15,2) NOT NULL,
  `status` enum('pending','completed','failed') DEFAULT 'completed',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ;

--
-- Dumping data for table `transfer_logs`
--

INSERT INTO `transfer_logs` (`id`, `sender_id`, `recipient_id`, `amount`, `status`, `created_at`) VALUES
(1, 1, 2, 100000.00, 'completed', '2025-04-25 14:12:47');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `accounts`
--
ALTER TABLE `accounts`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`),
  ADD UNIQUE KEY `account_number` (`account_number`);

--
-- Indexes for table `transactions`
--
ALTER TABLE `transactions`
  ADD PRIMARY KEY (`id`),
  ADD KEY `account_id` (`account_id`),
  ADD KEY `related_account_id` (`related_account_id`);

--
-- Indexes for table `transfer_logs`
--
ALTER TABLE `transfer_logs`
  ADD PRIMARY KEY (`id`),
  ADD KEY `sender_id` (`sender_id`),
  ADD KEY `recipient_id` (`recipient_id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `accounts`
--
ALTER TABLE `accounts`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `transactions`
--
ALTER TABLE `transactions`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `transfer_logs`
--
ALTER TABLE `transfer_logs`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `transactions`
--
ALTER TABLE `transactions`
  ADD CONSTRAINT `transactions_ibfk_1` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`),
  ADD CONSTRAINT `transactions_ibfk_2` FOREIGN KEY (`related_account_id`) REFERENCES `accounts` (`id`);

--
-- Constraints for table `transfer_logs`
--
ALTER TABLE `transfer_logs`
  ADD CONSTRAINT `transfer_logs_ibfk_1` FOREIGN KEY (`sender_id`) REFERENCES `accounts` (`id`),
  ADD CONSTRAINT `transfer_logs_ibfk_2` FOREIGN KEY (`recipient_id`) REFERENCES `accounts` (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
