START TRANSACTION;

--
-- Constraints for table `tblCustomerBankAccount`
--
ALTER TABLE `tblCustomerBankAccount`
  DROP FOREIGN KEY `FK_tblCustomerBankAccount_tblCustomerBasicInfo`;

--
-- Table structure for table `tblCustomerBankAccount`
--

DROP TABLE IF EXISTS `tblCustomerBankAccount`;

COMMIT;