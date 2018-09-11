START TRANSACTION;

ALTER TABLE `tblCustomerLoan` DROP `source`;
ALTER TABLE `tblCustomerLoanApplication` DROP `source`;
ALTER TABLE `tblCustomerLoanHistory` DROP `source`;

ALTER TABLE `tblCustomerPayment` DROP `destination`;
ALTER TABLE `tblCustomerPaymentHistory` DROP `destination`;


DROP TRIGGER IF EXISTS `trg_after_tblCustomerLoanApplication_insert`; 
CREATE TRIGGER `trg_after_tblCustomerLoanApplication_insert` AFTER INSERT ON `tblCustomerLoanApplication` FOR EACH ROW BEGIN
    DECLARE _balance DOUBLE(10,2);
	  IF (NEW.status="APPROVED") THEN      
		INSERT INTO tblCustomerLoan (fk_customer_id,loan_application_id,epoint_transaction_id,loan_amount,interest_amount,fee_amount,total_amount, reference_code,due_date,loan_date,approved_by,approved_date) VALUES (NEW.fk_customer_id,NEW.id,NEW.epoint_transaction_id, NEW.loan_amount, NEW.interest_amount,NEW.fee_amount,NEW.total_amount,NEW.reference_code,NEW.due_date,NEW.loan_date,NEW.processed_by,NEW.processed_date);
		
		UPDATE `tblCustomerCreditLine` 
		SET available_credit = (available_credit - NEW.loan_amount)
		WHERE fk_customer_id = NEW.fk_customer_id;  
		
		UPDATE `tblCustomerLoanTotal` 
		SET total_principal_amount = (total_principal_amount + NEW.loan_amount),
		    total_fee_amount = (total_fee_amount + NEW.fee_amount),
		    due_date = IF(total_amount=0,NEW.due_date,due_date),
		    total_amount = (total_amount + NEW.total_amount)
		    
		WHERE fk_customer_id = NEW.fk_customer_id;  
		
		
		set _balance = (select total_amount from tblCustomerLoanTotal where fk_customer_id = NEW.fk_customer_id);
		insert into tblCustomerLoanTransaction (fk_customer_id,`type`,debit,balance) VALUES (NEW.fk_customer_id,"LOAN",NEW.total_amount,_balance);
          END IF;
          
          INSERT INTO tblCustomerLoanHistory (fk_customer_id,loan_amount,interest_amount,fee_amount,total_amount,reference_code,due_date,loan_date,`status`,epoint_transaction_id,processed_by,processed_date) 
          VALUES (NEW.fk_customer_id,NEW.loan_amount,NEW.interest_amount,NEW.fee_amount,NEW.total_amount,NEW.reference_code,NEW.due_date,NEW.loan_date,NEW.status,NEW.epoint_transaction_id,NEW.processed_by,NEW.processed_date);
    END;


DROP TRIGGER IF EXISTS `trg_after_tblCustomerLoanApplication_update`; 
CREATE TRIGGER `trg_after_tblCustomerLoanApplication_update` AFTER UPDATE ON `tblCustomerLoanApplication` FOR EACH ROW BEGIN
    DECLARE _balance DOUBLE(10,2);
	  IF (NEW.status="APPROVED") THEN      
		INSERT INTO tblCustomerLoan (fk_customer_id,loan_application_id,epoint_transaction_id,loan_amount,interest_amount,fee_amount,total_amount, reference_code,due_date,loan_date,approved_by,approved_date) VALUES (NEW.fk_customer_id,NEW.id,NEW.epoint_transaction_id, NEW.loan_amount, NEW.interest_amount,NEW.fee_amount,NEW.total_amount,NEW.reference_code,NEW.due_date,NEW.loan_date,NEW.processed_by,NEW.processed_date);
		
		UPDATE `tblCustomerCreditLine` 
		SET available_credit = (available_credit - NEW.loan_amount)
		WHERE fk_customer_id = NEW.fk_customer_id;  
		
		UPDATE `tblCustomerLoanTotal` 
		SET total_principal_amount = (total_principal_amount + NEW.loan_amount),
		    total_fee_amount = (total_fee_amount + NEW.fee_amount),
		    total_amount = (total_amount + NEW.total_amount)
		WHERE fk_customer_id = NEW.fk_customer_id;  
		
		
		SET _balance = (SELECT total_amount FROM tblCustomerLoanTotal WHERE fk_customer_id = NEW.fk_customer_id);
		INSERT INTO tblCustomerLoanTransaction (fk_customer_id,`type`,debit,balance) VALUES (NEW.fk_customer_id,"LOAN",NEW.total_amount,_balance);
          END IF;
          
          IF (NEW.status<>OLD.status) THEN
		INSERT INTO tblCustomerLoanHistory (fk_customer_id,loan_amount,interest_amount,fee_amount,total_amount,reference_code,due_date,loan_date,`status`,epoint_transaction_id,processed_by,processed_date) 
		VALUES (NEW.fk_customer_id,NEW.loan_amount,NEW.interest_amount,NEW.fee_amount,NEW.total_amount,NEW.reference_code,NEW.due_date,NEW.loan_date,NEW.status,NEW.epoint_transaction_id,NEW.processed_by,NEW.processed_date);
          END IF;
    END;


DROP TRIGGER IF EXISTS `trg_after_tblCustomerPayment_insert`; 
CREATE TRIGGER `trg_after_tblCustomerPayment_insert` AFTER INSERT ON `tblCustomerPayment` FOR EACH ROW BEGIN
    DECLARE _balance DOUBLE(10,2);
	INSERT INTO `tblCustomerPaymentHistory` 
	SET fk_customer_id = NEW.fk_customer_id,
	    payment_amount = NEW.payment_amount,
	    reference_code = NEW.reference_code,
	    epoint_transaction_id=NEW.epoint_transaction_id,
	    paid_by = NEW.paid_by,
	    date_paid= NEW.date_paid;
	
	SET _balance = (SELECT total_amount FROM tblCustomerLoanTotal WHERE fk_customer_id = NEW.fk_customer_id) - NEW.payment_amount;
	if _balance<0 then
		SET _balance = 0.00;
	end if;
	INSERT INTO tblCustomerLoanTransaction (fk_customer_id,`type`,credit,balance) VALUES (NEW.fk_customer_id,"PAYMENT",NEW.payment_amount,_balance);
    END;


ALTER ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `view_transaction_history`  AS  select `tblCustomerLoanHistory`.`fk_customer_id` AS `fk_customer_id`,`tblCustomerLoanHistory`.`loan_amount` AS `amount`,'LOAN' AS `t_type`,`tblCustomerLoanHistory`.`loan_date` AS `t_date` from `tblCustomerLoanHistory` union select `tblCustomerPaymentHistory`.`fk_customer_id` AS `fk_customer_id`,`tblCustomerPaymentHistory`.`payment_amount` AS `amount`,'PAYMENT' AS `t_type`,`tblCustomerPaymentHistory`.`date_paid` AS `t_date` from `tblCustomerPaymentHistory`;

COMMIT;