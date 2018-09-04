START TRANSACTION;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `tblApiKey`
--
ALTER TABLE `tblApiKey`
  DROP FOREIGN KEY `FK_tblApiKey_tblCustomerBasicInfo`;

--
-- Constraints for table `tblCustomerAgreement`
--
ALTER TABLE `tblCustomerAgreement`
  DROP FOREIGN KEY `FK_tblCustomerAgreement_tblCustomerBasicInfo`;

--
-- Constraints for table `tblCustomerCreditLine`
--
ALTER TABLE `tblCustomerCreditLine`
  DROP FOREIGN KEY `FK_tblCustomerCreditLine_tblCustomerBasicInfo`;

--
-- Constraints for table `tblCustomerCreditLineApplication`
--
ALTER TABLE `tblCustomerCreditLineApplication`
  DROP FOREIGN KEY `FK_tblCustomerCreditLineApplication_tblCustomerBasicInfo`;

--
-- Constraints for table `tblCustomerCreditLineHistory`
--
ALTER TABLE `tblCustomerCreditLineHistory`
  DROP FOREIGN KEY `FK_tblCustomerCreditLineHistory_tblCustomerBasicInfo`;

--
-- Constraints for table `tblCustomerLoan`
--
ALTER TABLE `tblCustomerLoan`
  DROP FOREIGN KEY `FK_tblCustomerLoan_tblCustomerBasicInfo`;

--
-- Constraints for table `tblCustomerLoanApplication`
--
ALTER TABLE `tblCustomerLoanApplication`
  DROP FOREIGN KEY `FK_tblCustomerLoanApplication_tblCustomerBasicInfo`;

--
-- Constraints for table `tblCustomerLoanHistory`
--
ALTER TABLE `tblCustomerLoanHistory`
  DROP FOREIGN KEY `FK_tblCustomerLoanHistory_tblCustomerBasicInfo`;

--
-- Constraints for table `tblCustomerLoanTotal`
--
ALTER TABLE `tblCustomerLoanTotal`
  DROP FOREIGN KEY `FK_tblCustomerLoanTotal_tblCustomerBasicInfo`;

--
-- Constraints for table `tblCustomerLoanTransaction`
--
ALTER TABLE `tblCustomerLoanTransaction`
  DROP FOREIGN KEY `FK_tblCustomerLoanTransaction_tblCustomerBasicInfo`;

--
-- Constraints for table `tblCustomerOtherInfo`
--
ALTER TABLE `tblCustomerOtherInfo`
  DROP FOREIGN KEY `FK_tblCustomerOtherInfo_tblCustomerBasicInfo`;

--
-- Constraints for table `tblCustomerPayment`
--
ALTER TABLE `tblCustomerPayment`
  DROP FOREIGN KEY `FK_tblCustomerPayment_tblCustomerBasicInfo`;

--
-- Constraints for table `tblCustomerPaymentHistory`
--
ALTER TABLE `tblCustomerPaymentHistory`
  DROP FOREIGN KEY `FK_tblCustomerPaymentHistory_tblCustomerBasicInfo`;

--
-- Constraints for table `tblCustomerSettlement`
--
ALTER TABLE `tblCustomerSettlement`
  DROP FOREIGN KEY `FK_tblCustomerSettlement_tblCustomerBasicInfo`,
  DROP FOREIGN KEY `FK_tblCustomerSettlement_tblCustomerPayment`;

--
-- Constraints for table `tblSystemMenuCategory`
--
ALTER TABLE `tblSystemMenuCategory`
  DROP FOREIGN KEY `FK_tblSystemMenuCategory_tblSystemCategory`,
  DROP FOREIGN KEY `FK_tblSystemMenuCategory_tblSystemMenu`;

--
-- Constraints for table `tblSystemMenuSub`
--
ALTER TABLE `tblSystemMenuSub`
  DROP FOREIGN KEY `FK_tblSystemMenuSub_tblSystemMenu`;

--
-- Constraints for table `tblSystemRoleMenu`
--
ALTER TABLE `tblSystemRoleMenu`
  DROP FOREIGN KEY `FK_tblSystemRoleMenu_tblSystemMenuCategory`,
  DROP FOREIGN KEY `FK_tblSystemRoleMenu_tblSystemRole`;

--
-- Constraints for table `tblUserInfo`
--
ALTER TABLE `tblUserInfo`
  DROP FOREIGN KEY `FK_tblUserInfo_tblUserInformation`;

--
-- Constraints for table `tblUserRole`
--
ALTER TABLE `tblUserRole`
  DROP FOREIGN KEY `FK_tblUserRole_tblSystemRole`,
  DROP FOREIGN KEY `FK_tblUserRole_tblUserInformation`;

--
-- Procedures
--
DROP PROCEDURE IF EXISTS `prc_get_SystemCategory`;


DROP PROCEDURE IF EXISTS `prc_get_SystemMenu`;


DROP PROCEDURE IF EXISTS `prc_get_SystemRole`;


DROP PROCEDURE IF EXISTS `prc_get_SystemRoleMenu`;


DROP PROCEDURE IF EXISTS `prc_get_UserInfo`;


DROP PROCEDURE IF EXISTS `prc_get_UserInformation`;


DROP PROCEDURE IF EXISTS `prc_get_UserInfo_login`;


DROP PROCEDURE IF EXISTS `prc_get_UserRole`;


DROP PROCEDURE IF EXISTS `prc_set_Creditlimit`;


DROP PROCEDURE IF EXISTS `prc_set_Fee`;


DROP PROCEDURE IF EXISTS `prc_set_Interest`;


DROP PROCEDURE IF EXISTS `prc_set_Interval`;


DROP PROCEDURE IF EXISTS `prc_set_SystemCategory`;


DROP PROCEDURE IF EXISTS `prc_set_SystemMenu`;


DROP PROCEDURE IF EXISTS `prc_set_SystemMenuCategory`;


DROP PROCEDURE IF EXISTS `prc_set_SystemRole`;


DROP PROCEDURE IF EXISTS `prc_set_SystemRoleMenu`;


DROP PROCEDURE IF EXISTS `prc_set_UserInformation`;


DROP PROCEDURE IF EXISTS `prc_set_UserRole`;

--
-- Functions
--
DROP FUNCTION IF EXISTS `fn_date_format`;


DROP FUNCTION IF EXISTS `fn_fullname_get`;


DROP FUNCTION IF EXISTS `fn_systemrole_get`;


DROP FUNCTION IF EXISTS `fn_user_fullname_get`;



-- --------------------------------------------------------

--
-- Table structure for table `tblApiAccess`
--

DROP TABLE IF EXISTS `tblApiAccess`;


-- --------------------------------------------------------

--
-- Table structure for table `tblApiKey`
--

DROP TABLE IF EXISTS `tblApiKey`;


--
-- Triggers `tblApiKey`
--
DROP TRIGGER IF EXISTS `trg_after_tblApiKey_delete`;


DROP TRIGGER IF EXISTS `trg_after_tblApiKey_insert`;


DROP TRIGGER IF EXISTS `trg_after_tblApiKey_update`;


-- --------------------------------------------------------

--
-- Table structure for table `tblApiLimit`
--

DROP TABLE IF EXISTS `tblApiLimit`;


-- --------------------------------------------------------

--
-- Table structure for table `tblApiLogs`
--

DROP TABLE IF EXISTS `tblApiLogs`;


-- --------------------------------------------------------

--
-- Table structure for table `tblCity`
--

DROP TABLE IF EXISTS `tblCity`;

-- --------------------------------------------------------

--
-- Table structure for table `tblCountry`
--

DROP TABLE IF EXISTS `tblCountry`;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerAgreement`
--

DROP TABLE IF EXISTS `tblCustomerAgreement`;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerBasicInfo`
--

DROP TABLE IF EXISTS `tblCustomerBasicInfo`;

--
-- Triggers `tblCustomerBasicInfo`
--
DROP TRIGGER IF EXISTS `trg_after_tblCustomerBasicInfo_insert`;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerCreditLine`
--

DROP TABLE IF EXISTS `tblCustomerCreditLine`;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerCreditLineApplication`
--

DROP TABLE IF EXISTS `tblCustomerCreditLineApplication`;

--
-- Triggers `tblCustomerCreditLineApplication`
--
DROP TRIGGER IF EXISTS `trg_after_tblCustomerCreditLineApplication_insert`;

DROP TRIGGER IF EXISTS `trg_after_tblCustomerCreditLineApplication_update`;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerCreditLineHistory`
--

DROP TABLE IF EXISTS `tblCustomerCreditLineHistory`;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerLoan`
--

DROP TABLE IF EXISTS `tblCustomerLoan`;

--
-- Triggers `tblCustomerLoan`
--
DROP TRIGGER IF EXISTS `trg_after_tblCustomerLoan_update`;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerLoanApplication`
--

DROP TABLE IF EXISTS `tblCustomerLoanApplication`;

--
-- Triggers `tblCustomerLoanApplication`
--
DROP TRIGGER IF EXISTS `trg_after_tblCustomerLoanApplication_insert`;

DROP TRIGGER IF EXISTS `trg_after_tblCustomerLoanApplication_update`;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerLoanHistory`
--

DROP TABLE IF EXISTS `tblCustomerLoanHistory`;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerLoanTotal`
--

DROP TABLE IF EXISTS `tblCustomerLoanTotal`;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerLoanTransaction`
--

DROP TABLE IF EXISTS `tblCustomerLoanTransaction`;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerOtherInfo`
--

DROP TABLE IF EXISTS `tblCustomerOtherInfo`;

--
-- Triggers `tblCustomerOtherInfo`
--
DROP TRIGGER IF EXISTS `trg_after_tblCustomerOtherInfo_update`;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerPayment`
--

DROP TABLE IF EXISTS `tblCustomerPayment`;

--
-- Triggers `tblCustomerPayment`
--
DROP TRIGGER IF EXISTS `trg_after_tblCustomerPayment_insert`;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerPaymentHistory`
--

DROP TABLE IF EXISTS `tblCustomerPaymentHistory`;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerSettlement`
--

DROP TABLE IF EXISTS `tblCustomerSettlement`;

-- --------------------------------------------------------

--
-- Table structure for table `tblFee`
--

DROP TABLE IF EXISTS `tblFee`;

-- --------------------------------------------------------

--
-- Table structure for table `tblIncomeSource`
--

DROP TABLE IF EXISTS `tblIncomeSource`;

-- --------------------------------------------------------

--
-- Table structure for table `tblInterest`
--

DROP TABLE IF EXISTS `tblInterest`;

-- --------------------------------------------------------

--
-- Table structure for table `tblLoanCreditLimit`
--

DROP TABLE IF EXISTS `tblLoanCreditLimit`;

-- --------------------------------------------------------

--
-- Table structure for table `tblLoanInterval`
--

DROP TABLE IF EXISTS `tblLoanInterval`;

-- --------------------------------------------------------

--
-- Table structure for table `tblLoanPaymentCycle`
--

DROP TABLE IF EXISTS `tblLoanPaymentCycle`;
-- --------------------------------------------------------

--
-- Table structure for table `tblLoanTerm`
--

DROP TABLE IF EXISTS `tblLoanTerm`;

-- --------------------------------------------------------

--
-- Table structure for table `tblPayFrequency`
--

DROP TABLE IF EXISTS `tblPayFrequency`;

-- --------------------------------------------------------

--
-- Table structure for table `tblProgram`
--

DROP TABLE IF EXISTS `tblProgram`;

-- --------------------------------------------------------

--
-- Table structure for table `tblSession`
--

DROP TABLE IF EXISTS `tblSession`;

-- --------------------------------------------------------

--
-- Table structure for table `tblState`
--

DROP TABLE IF EXISTS `tblState`;

-- --------------------------------------------------------

--
-- Table structure for table `tblSystemCategory`
--

DROP TABLE IF EXISTS `tblSystemCategory`;

-- --------------------------------------------------------

--
-- Table structure for table `tblSystemMenu`
--

DROP TABLE IF EXISTS `tblSystemMenu`;

-- --------------------------------------------------------

--
-- Table structure for table `tblSystemMenuCategory`
--

DROP TABLE IF EXISTS `tblSystemMenuCategory`;

-- --------------------------------------------------------

--
-- Table structure for table `tblSystemMenuSub`
--

DROP TABLE IF EXISTS `tblSystemMenuSub`;

-- --------------------------------------------------------

--
-- Table structure for table `tblSystemRole`
--

DROP TABLE IF EXISTS `tblSystemRole`;

-- --------------------------------------------------------

--
-- Table structure for table `tblSystemRoleMenu`
--

DROP TABLE IF EXISTS `tblSystemRoleMenu`;

-- --------------------------------------------------------

--
-- Table structure for table `tblSystemSettings`
--

DROP TABLE IF EXISTS `tblSystemSettings`;

-- --------------------------------------------------------

--
-- Table structure for table `tblUserInfo`
--

DROP TABLE IF EXISTS `tblUserInfo`;

-- --------------------------------------------------------

--
-- Table structure for table `tblUserInformation`
--

DROP TABLE IF EXISTS `tblUserInformation`;

--
-- Triggers `tblUserInformation`
--
DROP TRIGGER IF EXISTS `trg_after_tblUserInformation_insert`;

-- --------------------------------------------------------

--
-- Table structure for table `tblUserRole`
--

DROP TABLE IF EXISTS `tblUserRole`;

-- --------------------------------------------------------

--
-- Structure for view `view_customer_info`
--
DROP VIEW IF EXISTS `view_customer_info`;
-- --------------------------------------------------------

--
-- Structure for view `view_transaction_history`
--
DROP VIEW IF EXISTS `view_transaction_history`;

COMMIT;