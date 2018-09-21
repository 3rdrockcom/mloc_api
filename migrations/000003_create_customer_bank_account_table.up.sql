START TRANSACTION;

--
-- Table structure for table `tblCustomerBankAccount`
--

CREATE TABLE IF NOT EXISTS `tblCustomerBankAccount` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `customer_id` int(11) NOT NULL,
  `alias` varchar(80) NOT NULL DEFAULT '',
  `account_type` varchar(32) NOT NULL,
  `bank_code` varchar(16) NOT NULL DEFAULT '',
  `account_number` varchar(32) NOT NULL,
  `kms_id` int(11) NOT NULL,
  `evault_id` int(11) NOT NULL,
  `date_created` datetime NOT NULL DEFAULT current_timestamp(),
  `date_updated` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `customer_id` (`customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `tblCustomerBankAccount`
--
ALTER TABLE `tblCustomerBankAccount`
  ADD CONSTRAINT `FK_tblCustomerBankAccount_tblCustomerBasicInfo` FOREIGN KEY (`customer_id`) REFERENCES `tblCustomerBasicInfo` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

COMMIT;