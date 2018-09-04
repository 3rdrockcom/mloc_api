START TRANSACTION;

--
-- Procedures
--
CREATE PROCEDURE `prc_get_SystemCategory` (`_where` TEXT)  BEGIN
SET @QUERY = CONCAT('SELECT * FROM tblSystemCategory a ',_where);
  PREPARE stmt FROM @QUERY;
  EXECUTE stmt;
  DEALLOCATE PREPARE stmt;
END;

CREATE PROCEDURE `prc_get_SystemMenu` (`_where` TEXT)  BEGIN
SET @QUERY = CONCAT('SELECT a.menuid,a.link,a.title,a.status,a.arranged,a.comments,a.icon,a.fk_template_id FROM `tblSystemMenu` a ',_where);
  PREPARE stmt FROM @QUERY;
  EXECUTE stmt;
  DEALLOCATE PREPARE stmt;
END;

CREATE PROCEDURE `prc_get_SystemRole` (`_where` TEXT)  BEGIN
SET @QUERY = CONCAT('SELECT a.roleid,a.code,a.description,a.fixed,a.icon,a.created_by,a.date_created FROM tblSystemRole a ',_where);
  PREPARE stmt FROM @QUERY;
  EXECUTE stmt;
  DEALLOCATE PREPARE stmt;
END;

CREATE PROCEDURE `prc_get_SystemRoleMenu` (IN `_syscatid` VARCHAR(100), IN `_role` VARCHAR(100), IN `_userid` VARCHAR(100))  BEGIN
	
	if IFNULL(_syscatid,'')='' then
		SELECT DISTINCT fk_syscatid,c.description AS title,'#' AS `link`,c.icon,c.cat_label as label FROM tblSystemRoleMenu a
		INNER JOIN tblSystemMenuCategory b ON b.menucat_id=a.fk_menucat_id
		INNER JOIN tblSystemCategory c ON c.syscat_id=b.fk_syscatid
		INNER JOIN tblSystemMenu d ON d.menuid=b.fk_menuid
		WHERE FIND_IN_SET(fk_roleid,_role) AND d.status='SHOW' ORDER BY c.arranged;
	else
		SELECT DISTINCT fk_menucat_id,fk_syscatid,title,d.link,d.icon,can_read,can_add,can_edit,can_delete,d.view_folder,d.menu_label AS label FROM tblSystemRoleMenu a
		INNER JOIN tblSystemMenuCategory b ON b.menucat_id=a.fk_menucat_id
		INNER JOIN tblSystemCategory c ON c.syscat_id=b.fk_syscatid
		INNER JOIN tblSystemMenu d ON d.menuid=b.fk_menuid
		WHERE FIND_IN_SET(fk_roleid,_role) AND fk_syscatid=_syscatid AND d.status='SHOW' ORDER BY d.arranged;
	end if;
END;

CREATE PROCEDURE `prc_get_UserInfo` (`_where` TEXT)  BEGIN
SET @QUERY = CONCAT('SELECT a.fk_userid,a.username,a.status,a.ipadd,CONCAT(b.lname,\', \',b.fname,\' \',IFNULL(mname,"")) AS fullname,b.lname,b.fname,b.mname,b.userno,b.gender,b.email,
GROUP_CONCAT(d.description) AS role_assigned
FROM tblUserInfo a 
INNER JOIN tblUserInformation b ON b.userid=a.fk_userid
LEFT JOIN tblUserRole c ON c.fk_userid=a.fk_userid
LEFT JOIN tblSystemRole d ON d.roleid=c.fk_roleid ',_where);
  PREPARE stmt FROM @QUERY;
  EXECUTE stmt;
  DEALLOCATE PREPARE stmt;
	
END;

CREATE PROCEDURE `prc_get_UserInformation` (`_where` TEXT)  BEGIN
SET @QUERY = CONCAT('SELECT * FROM tblUserInformation  ',_where);
  PREPARE stmt FROM @QUERY;
  EXECUTE stmt;
  DEALLOCATE PREPARE stmt;
	
END;

CREATE PROCEDURE `prc_get_UserInfo_login` (`_where` TEXT)  BEGIN
SET @QUERY = CONCAT('SELECT a.fk_userid,a.username,a.status,a.ipadd,CONCAT(b.lname,\', \',b.fname) AS FULLNAME,b.userno,b.gender FROM tblUserInfo a INNER JOIN tblUserInformation b ON b.userid=a.fk_userid ',_where);
  PREPARE stmt FROM @QUERY;
  EXECUTE stmt;
  DEALLOCATE PREPARE stmt;
END;

CREATE PROCEDURE `prc_get_UserRole` (`_userid` VARCHAR(11))  BEGIN
	declare _count int(11);
	SET _count = (SELECT COUNT(fk_roleid) FROM tblUserRole a 
		      INNER JOIN tblSystemRole b ON b.roleid=a.fk_roleid
		      WHERE a.fk_userid=_userid);
	if _count>1 then
		SELECT GROUP_CONCAT(fk_roleid) as fk_roleid,'All' as `code`, 'All' as description, 'fa fa-star' as icon
		FROM tblUserRole a 
		INNER JOIN tblSystemRole b ON b.roleid=a.fk_roleid
		WHERE a.fk_userid=_userid
		UNION
		SELECT fk_roleid,b.code,b.description,b.icon FROM tblUserRole a 
		INNER JOIN tblSystemRole b ON b.roleid=a.fk_roleid
		WHERE a.fk_userid=_userid;
	else
		SELECT fk_roleid,b.code,b.description,b.icon FROM tblUserRole a 
		INNER JOIN tblSystemRole b ON b.roleid=a.fk_roleid
		WHERE a.fk_userid=_userid;
	end if;
	
    END;

CREATE PROCEDURE `prc_set_Creditlimit` (`_action` VARCHAR(1), `_clid` VARCHAR(11), `_cl_code` VARCHAR(150), `_cl_desc` VARCHAR(150), `_amount` DOUBLE(10,2), `_no_of_days` INT(11), `_active` VARCHAR(30))  BEGIN
IF _cl_code='' THEN SET _cl_code=NULL; END IF;
IF _cl_desc='' THEN SET _cl_desc=NULL; END IF;
  CASE _action
    WHEN 'S' THEN
	
	IF EXISTS (SELECT id FROM `tblLoanCreditLimit` WHERE id = _clid)
	  THEN
	    UPDATE tblLoanCreditLimit 
	    SET `code` = _cl_code,
		description=_cl_desc,
		amount=_amount,
		no_of_days=_no_of_days,
		active=_active
	    WHERE id = _clid;
	  ELSE 
	    INSERT INTO tblLoanCreditLimit (`code`,description,amount,no_of_days,active) 
	    VALUES (_cl_code,_cl_desc,_amount,_no_of_days,_active); 
	END IF;
    WHEN 'D' THEN 
	DELETE FROM `tblLoanCreditLimit` WHERE id = _clid;
  END CASE;
END;

CREATE PROCEDURE `prc_set_Fee` (`_action` VARCHAR(1), `_feeid` VARCHAR(11), `_fee_code` VARCHAR(150), `_fee_desc` VARCHAR(150), `_percentage` DOUBLE(10,2), `_fixed` DOUBLE(10,2), `_active` VARCHAR(30))  BEGIN
IF _fee_code='' THEN SET _fee_code=NULL; END IF;
IF _fee_desc='' THEN SET _fee_desc=NULL; END IF;
  CASE _action
    WHEN 'S' THEN
	
	IF EXISTS (SELECT id FROM tblFee WHERE id = _feeid)
	  THEN
	    UPDATE tblFee 
	    SET `code` = _fee_code,
		description=_fee_desc,
		percentage=_percentage,
		`fixed`=_fixed,
		active=_active
	    WHERE id = _feeid;
	  ELSE 
	    INSERT INTO tblFee (`code`,description,percentage,`fixed`,active) 
	    VALUES (_fee_code,_fee_desc,_percentage,_fixed,_active); 
	END IF;
    WHEN 'D' THEN 
	DELETE FROM `tblFee` WHERE id = _feeid;
  END CASE;
END;

CREATE PROCEDURE `prc_set_Interest` (`_action` VARCHAR(1), `_interestid` VARCHAR(11), `_interest_code` VARCHAR(150), `_interest_desc` VARCHAR(150), `_percentage` DOUBLE(10,2), `_fixed` DOUBLE(10,2), `_active` VARCHAR(30))  BEGIN
IF _interest_code='' THEN SET _interest_code=NULL; END IF;
IF _interest_desc='' THEN SET _interest_desc=NULL; END IF;
  CASE _action
    WHEN 'S' THEN
	
	IF EXISTS (SELECT id FROM tblInterest WHERE id = _interestid)
	  THEN
	    UPDATE tblInterest 
	    SET `code` = _interest_code,
		description=_interest_desc,
		percentage=_percentage,
		`fixed`=_fixed,
		active=_active
	    WHERE id = _interestid;
	  ELSE 
	    INSERT INTO tblInterest (`code`,description,percentage,`fixed`,active) 
	    VALUES (_interest_code,_interest_desc,_percentage,_fixed,_active); 
	END IF;
    WHEN 'D' THEN 
	DELETE FROM `tblInterest` WHERE id = _interestid;
  END CASE;
END;

CREATE PROCEDURE `prc_set_Interval` (`_action` VARCHAR(1), `_interval_id` VARCHAR(11), `_interval_desc` VARCHAR(150), `_no_of_days` INT(11), `_active` VARCHAR(30))  BEGIN
IF _no_of_days='' THEN SET _no_of_days=NULL; END IF;
IF _interval_desc='' THEN SET _interval_desc=NULL; END IF;
  CASE _action
    WHEN 'S' THEN
	
	IF EXISTS (SELECT id FROM `tblLoanInterval` WHERE id = _interval_id)
	  THEN
	    UPDATE tblLoanInterval 
	    SET description=_interval_desc,
		no_of_days=_no_of_days,
		active=_active
	    WHERE id = _interval_id;
	  ELSE 
	    INSERT INTO tblLoanInterval (description,no_of_days,active) 
	    VALUES (_interval_desc,_no_of_days,_active); 
	END IF;
    WHEN 'D' THEN 
	DELETE FROM `tblLoanInterval` WHERE id = _interval_id;
  END CASE;
END;

CREATE PROCEDURE `prc_set_SystemCategory` (`_action` VARCHAR(1), `_syscat_id` VARCHAR(11), `_description` VARCHAR(150), `_arrange` VARCHAR(11), `_icon` VARCHAR(50), OUT `syscat_id_return` INT(11))  BEGIN
IF _description='' THEN SET _description=NULL; END IF;
IF _arrange='' THEN SET _arrange=NULL; END IF;
IF _icon='' THEN SET _icon=NULL; END IF;
SET syscat_id_return = _syscat_id;
  CASE _action
    WHEN 'S' THEN
	
	IF EXISTS (SELECT syscat_id FROM tblSystemCategory WHERE syscat_id = _syscat_id)
	  THEN
	    UPDATE tblSystemCategory 
	    SET description = _description,
		arranged=_arrange,
		icon=_icon
	    WHERE syscat_id = _syscat_id;
	  ELSE 
	    INSERT INTO tblSystemCategory (description,arranged,icon) 
	    VALUES (_description,_arrange,_icon); 
	    SET syscat_id_return = LAST_INSERT_ID();
	END IF;
    WHEN 'D' THEN 
	DELETE FROM tblSystemCategory WHERE syscat_id = _syscat_id;
  END CASE;
END;

CREATE PROCEDURE `prc_set_SystemMenu` (`_menuid` VARCHAR(11), `_title` VARCHAR(150), `_icon` VARCHAR(50), `_status` VARCHAR(10), `_arranged` VARCHAR(11), `_assigned_syscat_id` VARCHAR(11))  BEGIN
IF _title='' THEN SET _title=NULL; END IF;
IF _icon='' THEN SET _icon=NULL; END IF;
IF _status='' THEN SET _status=NULL; END IF;
IF _arranged='' THEN SET _arranged=NULL; END IF;
IF _assigned_syscat_id='' THEN SET _assigned_syscat_id=NULL; END IF;
UPDATE tblSystemMenu 
SET title = _title,
    icon=_icon,
    `status`=_status,
    arranged=_arranged
WHERE menuid = _menuid;
IF EXISTS (SELECT fk_syscatid FROM tblSystemMenuCategory WHERE fk_menuid = _menuid)
then
	UPDATE tblSystemMenuCategory SET fk_syscatid = _assigned_syscat_id WHERE fk_menuid = _menuid;
else
	INSERT INTO tblSystemMenuCategory(fk_syscatid,fk_menuid) VALUES (_assigned_syscat_id,_menuid);
end if;
END;

CREATE PROCEDURE `prc_set_SystemMenuCategory` (`_action` VARCHAR(1), `_syscat_id` VARCHAR(255), `_menuid` VARCHAR(255))  BEGIN
  CASE _action
    WHEN 'S' THEN
	
	IF NOT EXISTS (SELECT fk_syscatid FROM tblSystemMenuCategory WHERE fk_syscatid = _syscat_id AND fk_menuid=_menuid)
	  THEN
	    INSERT INTO tblSystemMenuCategory (fk_syscatid,fk_menuid) 
	    VALUES (_syscat_id,_menuid); 
	END IF;
    WHEN 'D' THEN 
	DELETE FROM tblSystemMenuCategory WHERE fk_syscatid = _syscat_id AND NOT FIND_IN_SET(fk_menuid,_menuid);
  END CASE;
END;

CREATE PROCEDURE `prc_set_SystemRole` (`_action` VARCHAR(1), `_roleid` VARCHAR(11), `_code` VARCHAR(150), `_description` VARCHAR(100), `_icon` VARCHAR(100), `_createdby` INT(11), OUT `_roleid_return` INT(11))  BEGIN
IF _code='' THEN SET _code=NULL; END IF;
IF _description='' THEN SET _description=NULL; END IF;
SET _roleid_return = _roleid;
  CASE _action
    WHEN 'S' THEN
	
	IF EXISTS (SELECT roleid FROM tblSystemRole WHERE roleid = _roleid)
	  THEN
	    UPDATE tblSystemRole 
	    SET `code`= _code,
		description=_description,
		icon=_icon
	    WHERE roleid = _roleid;
	    
	  ELSE 
	    INSERT INTO tblSystemRole (`code`,description,icon,created_by,date_created) 
	    VALUES (_code,_description,_icon,_createdby,CURRENT_TIMESTAMP); 
	    
	    SET _roleid_return = LAST_INSERT_ID();
	END IF;
    WHEN 'D' THEN 
	DELETE FROM tblSystemRole WHERE roleid = _roleid ;
  END CASE;
END;

CREATE PROCEDURE `prc_set_SystemRoleMenu` (`_action` VARCHAR(1), `_roleid` VARCHAR(11), `_menucat_id` VARCHAR(11), `_canread` VARCHAR(1), `_canadd` VARCHAR(1), `_canedit` VARCHAR(1), `_candelete` VARCHAR(1))  BEGIN
  CASE _action
    WHEN 'S' THEN
	    INSERT INTO tblSystemRoleMenu (fk_roleid,fk_menucat_id,can_read,can_add,can_edit,can_delete,datecreated)
	    VALUES (_roleid,_menucat_id,_canread,_canadd,_canedit,_candelete,CURRENT_TIMESTAMP); 
    WHEN 'D' THEN 
	DELETE FROM tblSystemRoleMenu WHERE fk_roleid=_roleid;
  END CASE;
END;

CREATE PROCEDURE `prc_set_UserInformation` (`_action` VARCHAR(1), `_userid` VARCHAR(11), `_userno` VARCHAR(150), `_lname` VARCHAR(100), `_fname` VARCHAR(100), `_mname` VARCHAR(100), `_email` VARCHAR(150), `_password` TEXT, `_createdby` INT(11), OUT `_userid_return` INT(11))  BEGIN
IF _userid='' THEN SET _userid=NULL; END IF;
IF _userno='' THEN SET _userno=NULL; END IF;
IF _email='' THEN SET _email=NULL; END IF;
SET _userid_return = _userid;
  CASE _action
    WHEN 'S' THEN
	
	IF EXISTS (SELECT userid FROM tblUserInformation WHERE userid = _userid)
	  THEN
	    UPDATE tblUserInformation 
	    SET userno = _userno,
		lname=_lname,
		fname=_fname,
		mname=_mname,
		email=_email
	    WHERE userid = _userid;
	    
	    UPDATE tblUserInfo SET username=_userno WHERE fk_userid = _userid;
	    IF _password<>'' THEN 
	    UPDATE tblUserInfo SET `password`=_password WHERE fk_userid = _userid;
	    END IF;
	    
	  ELSE 
	    INSERT INTO tblUserInformation (userno,lname,fname,mname,email,createdby,datecreated) 
	    VALUES (_userno,_lname,_fname,_mname,_email,_createdby,current_timestamp); 
	    
	    set _userid_return = last_insert_id();
	    set _userid = _userid_return;
	    
	    if _password<>'' then 
	    update tblUserInfo SET `password`=_password WHERE fk_userid = _userid;
	    end if;
	END IF;
    WHEN 'D' THEN 
	UPDATE tblUserInfo SET `status`='INACTIVE' WHERE fk_userid = _userid;
    WHEN 'A' THEN 
	UPDATE tblUserInfo SET `status`='ACTIVE' WHERE fk_userid = _userid;
  END CASE;
END;

CREATE PROCEDURE `prc_set_UserRole` (`_action` VARCHAR(1), `_userid` VARCHAR(11), `_roleid` VARCHAR(11), `_type` VARCHAR(11))  BEGIN
  CASE _action
    WHEN 'S' THEN
	
	IF NOT EXISTS (SELECT fk_userid FROM tblUserRole WHERE fk_userid = _userid AND fk_roleid=_roleid)
	  THEN
	    INSERT INTO tblUserRole (fk_userid,fk_roleid) 
	    VALUES (_userid,_roleid); 
	END IF;
    WHEN 'D' THEN 
	if _type='role' then
		DELETE FROM tblUserRole WHERE fk_userid = _userid AND NOT FIND_IN_SET(fk_roleid,_roleid);
	else
		DELETE FROM tblUserRole WHERE fk_roleid = _roleid AND NOT FIND_IN_SET(fk_userid,_userid);
	end if;
  END CASE;
END;

--
-- Functions
--
CREATE FUNCTION `fn_date_format` (`_date` DATETIME) RETURNS VARCHAR(255) CHARSET utf8 BEGIN
	DECLARE _return VARCHAR(150);
	SET _return = DATE_FORMAT(_date,'%M %d, %Y %h:%i %p');
	RETURN _return;
    END;

CREATE FUNCTION `fn_fullname_get` (`_customer_id` INT) RETURNS VARCHAR(150) CHARSET latin1 BEGIN
    DECLARE fullname VARCHAR(150);
    SET fullname = (SELECT CONCAT(last_name,', ',first_name,' ',IFNULL(middle_name,'')) FROM `tblCustomerBasicInfo` WHERE id=_customer_id);
    RETURN fullname;
    END;

CREATE FUNCTION `fn_systemrole_get` (`_role` INT) RETURNS VARCHAR(150) CHARSET latin1 BEGIN
    DECLARE utype VARCHAR(150);
    SET utype = (SELECT description FROM tblSystemRole WHERE moduleid=_role);
    RETURN utype;
    END;

CREATE FUNCTION `fn_user_fullname_get` (`_userid` INT) RETURNS VARCHAR(150) CHARSET latin1 BEGIN
    DECLARE fullname VARCHAR(150);
    SET fullname = (SELECT CONCAT(lname,', ',fname,' ',IFNULL(mname,'')) FROM tblUserInformation WHERE userid=_userid);
    RETURN fullname;
    END;


-- --------------------------------------------------------

--
-- Table structure for table `tblApiAccess`
--
CREATE TABLE IF NOT EXISTS `tblApiAccess` (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `key` varchar(40) NOT NULL DEFAULT '',
  `all_access` tinyint(1) NOT NULL DEFAULT '1',
  `controller` varchar(50) NOT NULL DEFAULT '',
  `date_created` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `date_modified` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `tblApiKey`
--
CREATE TABLE IF NOT EXISTS `tblApiKey` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `fk_customer_id` int(11) DEFAULT NULL,
  `key` blob NOT NULL,
  `level` int(2) NOT NULL,
  `ignore_limits` tinyint(1) NOT NULL DEFAULT '0',
  `is_private_key` tinyint(1) NOT NULL DEFAULT '0',
  `ip_addresses` text,
  `date_created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `FK_tblApiKey_tblCustomerBasicInfo` (`fk_customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Triggers `tblApiKey`
--
CREATE TRIGGER `trg_after_tblApiKey_delete` AFTER DELETE ON `tblApiKey` FOR EACH ROW BEGIN
	DELETE FROM tblApiAccess 
	WHERE `key` = OLD.key;
    END;

CREATE TRIGGER `trg_after_tblApiKey_insert` AFTER INSERT ON `tblApiKey` FOR EACH ROW BEGIN
	INSERT INTO tblApiAccess 
	SET `key` = NEW.key;
    END;

CREATE TRIGGER `trg_after_tblApiKey_update` AFTER UPDATE ON `tblApiKey` FOR EACH ROW BEGIN
	UPDATE tblApiAccess 
	SET `key` = NEW.key WHERE `key`=OLD.key;
    END;

-- --------------------------------------------------------

--
-- Table structure for table `tblApiLimit`
--
CREATE TABLE IF NOT EXISTS `tblApiLimit` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uri` varchar(255) NOT NULL,
  `count` int(10) NOT NULL,
  `hour_started` int(11) NOT NULL,
  `api_key` varchar(40) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `tblApiLogs`
--
CREATE TABLE IF NOT EXISTS `tblApiLogs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uri` varchar(255) NOT NULL,
  `method` varchar(6) NOT NULL,
  `params` text,
  `api_key` varchar(40) NOT NULL,
  `ip_address` varchar(45) NOT NULL,
  `time` int(11) NOT NULL,
  `rtime` float DEFAULT NULL,
  `authorized` varchar(1) NOT NULL,
  `response_code` smallint(3) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `tblCity`
--
CREATE TABLE IF NOT EXISTS `tblCity` (
  `city_id` int(11) NOT NULL AUTO_INCREMENT,
  `city` varchar(50) NOT NULL,
  `state_code` char(2) NOT NULL,
  PRIMARY KEY (`city_id`),
  KEY `idx_state_code` (`state_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `tblCountry`
--
CREATE TABLE IF NOT EXISTS `tblCountry` (
  `country_id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL,
  `iso_code_2` varchar(2) NOT NULL,
  `iso_code_3` varchar(3) NOT NULL,
  `address_format` text NOT NULL,
  `postcode_required` tinyint(1) NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '1',
  `mobile_prefix` varchar(6) DEFAULT NULL,
  PRIMARY KEY (`country_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerAgreement`
--
CREATE TABLE IF NOT EXISTS `tblCustomerAgreement` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `fk_customer_id` int(11) DEFAULT NULL,
  `mloc_access` int(1) DEFAULT '0',
  `registration` int(1) DEFAULT '0',
  `term_and_condition` int(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `FK_tblCustomerAgreement_tblCustomerBasicInfo` (`fk_customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerBasicInfo`
--
CREATE TABLE IF NOT EXISTS `tblCustomerBasicInfo` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `first_name` varchar(50) DEFAULT NULL,
  `middle_name` varchar(50) DEFAULT NULL,
  `last_name` varchar(50) DEFAULT NULL,
  `suffix` varchar(30) DEFAULT NULL,
  `birth_date` date DEFAULT NULL,
  `address1` varchar(255) DEFAULT NULL,
  `address2` varchar(255) DEFAULT NULL,
  `city` int(11) DEFAULT NULL,
  `state` int(11) DEFAULT NULL,
  `country` int(11) DEFAULT NULL,
  `zipcode` varchar(20) DEFAULT NULL,
  `home_number` varchar(20) DEFAULT NULL,
  `mobile_number` varchar(20) DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  `gender` enum('MALE','FEMALE') DEFAULT NULL,
  `program_id` int(11) DEFAULT '1',
  `program_customer_id` int(11) DEFAULT NULL,
  `program_customer_mobile` varchar(50) DEFAULT NULL,
  `cust_unique_id` varchar(100) DEFAULT NULL COMMENT 'this is for unique identity of customer from different programs',
  `created_by` varchar(50) DEFAULT NULL,
  `created_date` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Triggers `tblCustomerBasicInfo`
--
CREATE TRIGGER `trg_after_tblCustomerBasicInfo_insert` AFTER INSERT ON `tblCustomerBasicInfo` FOR EACH ROW BEGIN
	INSERT INTO `tblCustomerOtherInfo` SET `fk_customer_id` = NEW.id;
	INSERT INTO `tblCustomerAgreement` SET `fk_customer_id` = NEW.id, `mloc_access` = 1;
	INSERT INTO `tblCustomerCreditLine` SET `fk_customer_id` = NEW.id;
	INSERT INTO `tblCustomerLoanTotal` SET `fk_customer_id` = NEW.id;
	
    END;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerCreditLine`
--
CREATE TABLE IF NOT EXISTS `tblCustomerCreditLine` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `fk_customer_id` int(11) DEFAULT NULL COMMENT 'tblCustomerBasicInfo',
  `credit_line_id` int(11) DEFAULT NULL COMMENT 'tblCreditLimit',
  `credit_limit` double(10,2) DEFAULT '0.00',
  `available_credit` double(10,2) DEFAULT '0.00',
  `is_suspended` enum('YES','NO') DEFAULT 'NO',
  `approved_by` varchar(100) DEFAULT NULL,
  `approved_date` datetime DEFAULT NULL,
  `created_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `FK_tblCustomerCreditLine_tblCustomerBasicInfo` (`fk_customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerCreditLineApplication`
--
CREATE TABLE IF NOT EXISTS `tblCustomerCreditLineApplication` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `fk_customer_id` int(11) DEFAULT NULL COMMENT 'tblCustomerBasicInfo',
  `credit_line_id` int(11) DEFAULT NULL,
  `credit_line_amount` double(10,2) DEFAULT '0.00',
  `reference_code` varchar(100) DEFAULT NULL COMMENT 'reference for every application',
  `status` enum('PENDING','APPROVED','REJECTED') DEFAULT 'PENDING',
  `processed_by` varchar(100) DEFAULT NULL COMMENT 'if number means tblUserInfo else AUTOMATIC',
  `processed_date` datetime DEFAULT NULL,
  `created_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `FK_tblCustomerCreditLineApplication_tblCustomerBasicInfo` (`fk_customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Triggers `tblCustomerCreditLineApplication`
--
CREATE TRIGGER `trg_after_tblCustomerCreditLineApplication_insert` AFTER INSERT ON `tblCustomerCreditLineApplication` FOR EACH ROW BEGIN
	  IF (NEW.status="APPROVED") THEN
		UPDATE `tblCustomerCreditLine` 
		SET credit_limit = (credit_limit + NEW.credit_line_amount), 
	            available_credit = (available_credit + NEW.credit_line_amount),
	            credit_line_id = NEW.credit_line_id,
	            approved_by = NEW.processed_by,
	            approved_date = NEW.processed_date
		WHERE fk_customer_id = NEW.fk_customer_id;        
          END IF;
          
          INSERT INTO tblCustomerCreditLineHistory (fk_customer_id,credit_line_id,credit_line_amount,reference_code,`status`,processed_by,processed_date) 
          VALUES (NEW.fk_customer_id,NEW.credit_line_id,NEW.credit_line_amount,NEW.reference_code,NEW.status,NEW.processed_by,NEW.processed_date);
    END;

CREATE TRIGGER `trg_after_tblCustomerCreditLineApplication_update` AFTER UPDATE ON `tblCustomerCreditLineApplication` FOR EACH ROW BEGIN
	  IF (NEW.status="APPROVED") THEN
		UPDATE `tblCustomerCreditLine` 
		SET credit_limit = (credit_limit + NEW.credit_line_amount), 
	            available_credit = (available_credit + NEW.credit_line_amount),
	            credit_line_id = NEW.credit_line_id,
	            approved_by = NEW.processed_by,
	            approved_date = NEW.processed_date
		WHERE fk_customer_id = NEW.fk_customer_id;        
          END IF;
          
          IF (NEW.status<>OLD.status) THEN
		  INSERT INTO tblCustomerCreditLineHistory (fk_customer_id,credit_line_id,credit_line_amount,reference_code,`status`,processed_by,processed_date) 
		  VALUES (NEW.fk_customer_id,NEW.credit_line_id,NEW.credit_line_amount,NEW.reference_code,NEW.status,NEW.processed_by,NEW.processed_date);
          END IF;
    END;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerCreditLineHistory`
--
CREATE TABLE IF NOT EXISTS `tblCustomerCreditLineHistory` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `fk_customer_id` int(11) DEFAULT NULL COMMENT 'tblCustomerBasicInfo',
  `credit_line_id` int(11) DEFAULT NULL COMMENT 'tblCreditLimit',
  `credit_line_amount` double(10,2) DEFAULT NULL,
  `reference_code` varchar(100) DEFAULT NULL COMMENT 'reference for every application',
  `status` varchar(50) DEFAULT NULL,
  `processed_by` varchar(100) DEFAULT NULL COMMENT 'tblUserInformation',
  `processed_date` datetime DEFAULT NULL,
  `created_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `FK_tblCustomerCreditLineHistory_tblCustomerBasicInfo` (`fk_customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerLoan`
--
CREATE TABLE IF NOT EXISTS `tblCustomerLoan` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `fk_customer_id` int(11) DEFAULT NULL COMMENT 'tblCustomerBasicInfo',
  `loan_application_id` int(11) DEFAULT NULL,
  `epoint_transaction_id` varchar(50) DEFAULT NULL,
  `loan_interval_id` int(11) DEFAULT NULL COMMENT 'tblLoanInterval',
  `loan_term_id` int(11) DEFAULT NULL COMMENT 'tblLoanTerm',
  `loan_amount` double(10,2) DEFAULT '0.00',
  `interest_amount` double(10,2) DEFAULT '0.00' COMMENT 'depends if interest_fixed or interest_percentage',
  `fee_amount` double(10,2) DEFAULT '0.00' COMMENT 'depends if fee_fixed or fee_percentage',
  `total_amount` double(10,2) DEFAULT '0.00' COMMENT 'loan_amount + interest_amount + fee_amount',
  `total_paid_principal` double(10,2) DEFAULT '0.00',
  `total_paid_fee` double(10,2) DEFAULT '0.00',
  `total_paid_amount` double(10,2) DEFAULT '0.00',
  `reference_code` varchar(50) DEFAULT NULL,
  `is_paid` int(1) DEFAULT '0',
  `due_date` datetime DEFAULT NULL,
  `loan_date` datetime DEFAULT NULL,
  `approved_by` varchar(100) DEFAULT NULL,
  `approved_date` datetime DEFAULT NULL,
  `created_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `FK_tblCustomerLoan_tblCustomerBasicInfo` (`fk_customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Triggers `tblCustomerLoan`
--
CREATE TRIGGER `trg_after_tblCustomerLoan_update` AFTER UPDATE ON `tblCustomerLoan` FOR EACH ROW BEGIN
	declare _total_principal double(10,2);
	DECLARE _total_fee DOUBLE(10,2);
	IF (NEW.total_paid_principal<>OLD.total_paid_principal) then
	set _total_principal = NEW.total_paid_principal-OLD.total_paid_principal;
		UPDATE `tblCustomerLoanTotal`
		set total_principal_amount = total_principal_amount - _total_principal,
		    total_amount = total_amount - _total_principal,
		    due_date = IF(total_amount=0,NULL,due_date)
		WHERE fk_customer_id = NEW.fk_customer_id; 
		
		UPDATE `tblCustomerCreditLine`
		SET available_credit = available_credit + _total_principal
		WHERE fk_customer_id = NEW.fk_customer_id;           
	end if;
	
	IF (NEW.total_paid_fee<>OLD.total_paid_fee) THEN
	SET _total_fee = NEW.total_paid_fee-OLD.total_paid_fee;
		UPDATE `tblCustomerLoanTotal`
		SET total_fee_amount = total_fee_amount - _total_fee,
		    total_amount = total_amount - _total_fee,
		    due_date = IF(total_amount=0,NULL,due_date)
		WHERE fk_customer_id = NEW.fk_customer_id;   
		
	END IF;
    END;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerLoanApplication`
--
CREATE TABLE IF NOT EXISTS `tblCustomerLoanApplication` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `fk_customer_id` int(11) DEFAULT NULL COMMENT 'tblCustomerBasicInfo',
  `loan_interval_id` int(11) DEFAULT NULL COMMENT 'tblLoanInterval',
  `loan_term_id` int(11) DEFAULT NULL COMMENT 'tblLoanTerm',
  `loan_amount` double(10,2) DEFAULT '0.00',
  `interest_amount` double(10,2) DEFAULT '0.00' COMMENT 'depends if interest_fixed or interest_percentage',
  `fee_amount` double(10,2) DEFAULT '0.00' COMMENT 'depends if fee_fixed or fee_percentage',
  `total_amount` double(10,2) DEFAULT '0.00' COMMENT 'loan_amount + interest_amount + fee_amount',
  `reference_code` varchar(50) DEFAULT NULL,
  `due_date` datetime DEFAULT NULL,
  `loan_date` datetime DEFAULT NULL,
  `status` enum('PENDING','APPROVED','REJECTED') DEFAULT 'PENDING',
  `epoint_transaction_id` varchar(50) DEFAULT NULL,
  `processed_by` varchar(100) DEFAULT NULL,
  `processed_date` datetime DEFAULT NULL,
  `created_by` varchar(100) DEFAULT NULL COMMENT 'it could be system generated or tblUserInfo',
  `created_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `FK_tblCustomerLoanApplication_tblCustomerBasicInfo` (`fk_customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Triggers `tblCustomerLoanApplication`
--
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

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerLoanHistory`
--
CREATE TABLE IF NOT EXISTS `tblCustomerLoanHistory` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `fk_customer_id` int(11) DEFAULT NULL COMMENT 'tblCustomerBasicInfo',
  `loan_interval_id` int(11) DEFAULT NULL COMMENT 'tblLoanInterval',
  `loan_term_id` int(11) DEFAULT NULL COMMENT 'tblLoanTerm',
  `loan_amount` double(10,2) DEFAULT '0.00',
  `interest_amount` double(10,2) DEFAULT '0.00' COMMENT 'depends if interest_fixed or interest_percentage',
  `fee_amount` double(10,2) DEFAULT '0.00' COMMENT 'depends if fee_fixed or fee_percentage',
  `total_amount` double(10,2) DEFAULT '0.00' COMMENT 'loan_amount + interest_amount + fee_amount',
  `reference_code` varchar(50) DEFAULT NULL,
  `due_date` datetime DEFAULT NULL,
  `loan_date` datetime DEFAULT NULL,
  `status` varchar(100) DEFAULT NULL,
  `epoint_transaction_id` varchar(50) DEFAULT NULL,
  `processed_by` varchar(100) DEFAULT NULL,
  `processed_date` datetime DEFAULT NULL,
  `created_by` varchar(100) DEFAULT NULL,
  `created_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `FK_tblCustomerLoanHistory_tblCustomerBasicInfo` (`fk_customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerLoanTotal`
--
CREATE TABLE IF NOT EXISTS `tblCustomerLoanTotal` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `fk_customer_id` int(11) DEFAULT NULL,
  `total_principal_amount` double(10,2) DEFAULT '0.00',
  `total_fee_amount` double(10,2) DEFAULT '0.00',
  `total_amount` double(10,2) DEFAULT '0.00',
  `due_date` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `FK_tblCustomerLoanTotal_tblCustomerBasicInfo` (`fk_customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerLoanTransaction`
--
CREATE TABLE IF NOT EXISTS `tblCustomerLoanTransaction` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `fk_customer_id` int(11) DEFAULT NULL,
  `type` enum('LOAN','PAYMENT') DEFAULT NULL,
  `debit` double(10,2) DEFAULT '0.00' COMMENT 'loan',
  `credit` double(10,2) DEFAULT '0.00' COMMENT 'payment',
  `balance` double(10,2) DEFAULT '0.00',
  `created_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `FK_tblCustomerLoanTransaction_tblCustomerBasicInfo` (`fk_customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerOtherInfo`
--
CREATE TABLE IF NOT EXISTS `tblCustomerOtherInfo` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `fk_customer_id` int(11) DEFAULT NULL COMMENT 'tblCustomerBasicInfo',
  `company_name` varchar(100) DEFAULT NULL,
  `phone_number` varchar(50) DEFAULT NULL,
  `net_pay_percheck` double(10,2) DEFAULT '0.00',
  `income_source` int(11) DEFAULT NULL COMMENT 'tblIncomeSource',
  `pay_frequency` int(11) DEFAULT NULL COMMENT 'tblPayFrequency',
  `next_paydate` date DEFAULT NULL,
  `following_paydate` date DEFAULT NULL,
  `created_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `FK_tblCustomerOtherInfo_tblCustomerBasicInfo` (`fk_customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Triggers `tblCustomerOtherInfo`
--
CREATE TRIGGER `trg_after_tblCustomerOtherInfo_update` AFTER UPDATE ON `tblCustomerOtherInfo` FOR EACH ROW BEGIN
	UPDATE `tblCustomerAgreement` 
	SET `registration` = '1' WHERE `fk_customer_id`=OLD.fk_customer_id;
    END;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerPayment`
--
CREATE TABLE IF NOT EXISTS `tblCustomerPayment` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `fk_customer_id` int(11) DEFAULT NULL COMMENT 'tblCustomerBasicInfo',
  `reference_code` varchar(50) DEFAULT NULL,
  `epoint_transaction_id` varchar(50) DEFAULT NULL,
  `payment_amount` double(10,2) DEFAULT NULL,
  `date_paid` datetime DEFAULT NULL,
  `paid_by` varchar(20) DEFAULT NULL COMMENT 'tblCustomerBasicInfo',
  PRIMARY KEY (`id`),
  KEY `FK_tblCustomerPayment_tblCustomerBasicInfo` (`fk_customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Triggers `tblCustomerPayment`
--
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

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerPaymentHistory`
--
CREATE TABLE IF NOT EXISTS `tblCustomerPaymentHistory` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `fk_customer_id` int(11) DEFAULT NULL COMMENT 'tblCustomerBasicInfo',
  `reference_code` varchar(50) DEFAULT NULL,
  `epoint_transaction_id` varchar(50) DEFAULT NULL,
  `payment_amount` double(10,2) DEFAULT NULL,
  `date_paid` datetime DEFAULT NULL,
  `paid_by` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `FK_tblCustomerPaymentHistory_tblCustomerBasicInfo` (`fk_customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tblCustomerSettlement`
--
CREATE TABLE IF NOT EXISTS `tblCustomerSettlement` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `fk_customer_id` int(11) DEFAULT NULL COMMENT 'tblCustomerBasicInfo',
  `customer_loan_id` int(11) DEFAULT NULL,
  `customer_payment_id` int(11) DEFAULT NULL COMMENT 'tblCustomerPayment',
  `settlement_amount` double(10,2) DEFAULT '0.00',
  `principal_amount` double(10,2) DEFAULT '0.00',
  `fee_amount` double(10,2) DEFAULT '0.00',
  `is_settled` enum('YES','NO') DEFAULT 'NO',
  `created_date` datetime DEFAULT NULL,
  `settled_date` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `FK_tblCustomerSettlement_tblCustomerBasicInfo` (`fk_customer_id`),
  KEY `FK_tblCustomerSettlement_tblCustomerPayment` (`customer_payment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tblFee`
--
CREATE TABLE IF NOT EXISTS `tblFee` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(50) DEFAULT NULL,
  `description` varchar(150) DEFAULT NULL,
  `percentage` double(10,2) DEFAULT '0.00',
  `fixed` double(10,2) DEFAULT '0.00',
  `active` enum('YES','NO') DEFAULT 'NO',
  `created_date` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tblIncomeSource`
--
CREATE TABLE IF NOT EXISTS `tblIncomeSource` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `description` varchar(150) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tblInterest`
--
CREATE TABLE IF NOT EXISTS `tblInterest` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(50) DEFAULT NULL,
  `description` varchar(100) DEFAULT NULL,
  `percentage` double(10,2) DEFAULT '0.00',
  `fixed` double(10,2) DEFAULT '0.00',
  `active` enum('YES','NO') DEFAULT 'NO',
  `created_date` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tblLoanCreditLimit`
--
CREATE TABLE IF NOT EXISTS `tblLoanCreditLimit` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `tier` int(11) DEFAULT NULL,
  `code` varchar(50) DEFAULT NULL,
  `description` varchar(100) DEFAULT NULL,
  `amount` double(10,2) DEFAULT NULL,
  `no_of_days` int(11) DEFAULT NULL,
  `active` enum('YES','NO') DEFAULT 'NO',
  `created_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tblLoanInterval`
--
CREATE TABLE IF NOT EXISTS `tblLoanInterval` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `description` varchar(150) DEFAULT NULL,
  `no_of_days` int(11) DEFAULT NULL,
  `active` enum('YES','NO') DEFAULT 'NO',
  `created_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tblLoanPaymentCycle`
--
CREATE TABLE IF NOT EXISTS `tblLoanPaymentCycle` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `description` varchar(100) DEFAULT NULL,
  `created_date` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tblLoanTerm`
--
CREATE TABLE IF NOT EXISTS `tblLoanTerm` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `term` varchar(10) DEFAULT NULL,
  `created_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tblPayFrequency`
--
CREATE TABLE IF NOT EXISTS `tblPayFrequency` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `description` varchar(150) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tblProgram`
--
CREATE TABLE IF NOT EXISTS `tblProgram` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `program_code` varchar(30) DEFAULT NULL,
  `program_name` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tblSession`
--
CREATE TABLE IF NOT EXISTS `tblSession` (
  `id` varchar(60) NOT NULL,
  `fk_userid` int(11) DEFAULT NULL COMMENT 'tblUserInfo',
  `fullname` varchar(50) DEFAULT NULL,
  `ip_address` varchar(45) NOT NULL,
  `timestamp` int(10) UNSIGNED NOT NULL DEFAULT '0',
  `data` blob NOT NULL,
  PRIMARY KEY (`id`),
  KEY `ci_sessions_timestamp` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tblState`
--
CREATE TABLE IF NOT EXISTS `tblState` (
  `state_id` int(11) NOT NULL AUTO_INCREMENT,
  `state` varchar(22) NOT NULL,
  `state_code` char(2) NOT NULL,
  `country_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`state_code`),
  KEY `id` (`state_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `tblSystemCategory`
--
CREATE TABLE IF NOT EXISTS `tblSystemCategory` (
  `syscat_id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `description` text,
  `arranged` int(11) DEFAULT NULL,
  `icon` varchar(150) DEFAULT NULL,
  `cat_label` text COMMENT 'tblLanguageToken',
  PRIMARY KEY (`syscat_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `tblSystemMenu`
--
CREATE TABLE IF NOT EXISTS `tblSystemMenu` (
  `menuid` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `link` text,
  `title` text COMMENT 'maximum of 18 character',
  `status` enum('SHOW','HIDDEN') CHARACTER SET latin1 DEFAULT NULL,
  `arranged` int(11) DEFAULT NULL,
  `comments` text,
  `icon` text,
  `view_folder` text COMMENT '''~'' is = to ''/''',
  `menu_label` text COMMENT 'tblLanguageToken',
  PRIMARY KEY (`menuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `tblSystemMenuCategory`
--
CREATE TABLE IF NOT EXISTS `tblSystemMenuCategory` (
  `menucat_id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `fk_syscatid` int(11) UNSIGNED DEFAULT NULL COMMENT 'tblSystemCategory',
  `fk_menuid` int(11) UNSIGNED DEFAULT NULL COMMENT 'tblSystemMenu',
  PRIMARY KEY (`menucat_id`),
  KEY `FK_tblSystemMenuCategory_tblSystemCategory` (`fk_syscatid`),
  KEY `FK_tblSystemMenuCategory_tblSystemMenu` (`fk_menuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `tblSystemMenuSub`
--
CREATE TABLE IF NOT EXISTS `tblSystemMenuSub` (
  `sub_menuid` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `fk_menuid` int(11) UNSIGNED DEFAULT NULL COMMENT 'tblSystemMenu',
  `link` text,
  `title` text COMMENT 'maximum of 18 character',
  `status` enum('SHOW','HIDDEN') CHARACTER SET latin1 DEFAULT NULL,
  `arranged` int(11) DEFAULT NULL,
  `comments` text,
  `icon` text,
  `fk_template_id` int(11) UNSIGNED DEFAULT NULL COMMENT 'tblFieldTemplate',
  `view_folder` text COMMENT '''~'' is = to ''/''',
  `menu_label` text COMMENT 'tblLanguageToken',
  PRIMARY KEY (`sub_menuid`),
  KEY `FK_tblSystemMenu_tblFieldTemplate` (`fk_template_id`),
  KEY `FK_tblSystemMenuSub_tblSystemMenu` (`fk_menuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `tblSystemRole`
--
CREATE TABLE IF NOT EXISTS `tblSystemRole` (
  `roleid` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` varchar(30) CHARACTER SET latin1 DEFAULT NULL,
  `description` varchar(100) CHARACTER SET latin1 DEFAULT NULL,
  `fixed` enum('YES','NO') CHARACTER SET latin1 DEFAULT 'NO',
  `icon` varchar(150) CHARACTER SET latin1 DEFAULT NULL,
  `created_by` int(11) DEFAULT NULL,
  `date_created` datetime DEFAULT NULL,
  PRIMARY KEY (`roleid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `tblSystemRoleMenu`
--
CREATE TABLE IF NOT EXISTS `tblSystemRoleMenu` (
  `role_menu_id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `fk_roleid` int(11) UNSIGNED DEFAULT NULL COMMENT 'tblSystemRole',
  `fk_menucat_id` int(11) UNSIGNED DEFAULT NULL COMMENT 'tblMenuCategory',
  `can_read` enum('1','0') CHARACTER SET latin1 NOT NULL DEFAULT '1',
  `can_add` enum('1','0') CHARACTER SET latin1 NOT NULL DEFAULT '1',
  `can_edit` enum('1','0') CHARACTER SET latin1 NOT NULL DEFAULT '1',
  `can_delete` enum('1','0') CHARACTER SET latin1 NOT NULL DEFAULT '1',
  `datecreated` datetime DEFAULT NULL,
  `created_by` int(11) DEFAULT NULL,
  PRIMARY KEY (`role_menu_id`),
  KEY `FK_tblsystem_module_menu_menuid` (`fk_menucat_id`),
  KEY `FK_tblsystem_module_menu_moduleid` (`fk_roleid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `tblSystemSettings`
--
CREATE TABLE IF NOT EXISTS `tblSystemSettings` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(150) DEFAULT NULL,
  `code` varchar(30) DEFAULT NULL,
  `description` varchar(200) DEFAULT NULL,
  `value` varchar(30) DEFAULT NULL,
  `setting_type` enum('OTHERS','NOTIFICATION','SCHEDULE') DEFAULT 'OTHERS',
  `is_active` enum('YES','NO') DEFAULT 'YES',
  `sms_message` blob,
  `email_message` blob,
  `subject` varchar(200) DEFAULT NULL,
  `from` text,
  `to` text,
  `cc` text,
  `bcc` text,
  `updated_by` int(11) DEFAULT NULL COMMENT 'tblUserInfo',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tblUserInfo`
--
CREATE TABLE IF NOT EXISTS `tblUserInfo` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `fk_userid` int(11) UNSIGNED DEFAULT NULL COMMENT 'tblUserInformation',
  `username` text,
  `password` text CHARACTER SET latin1,
  `status` enum('ACTIVE','INACTIVE') CHARACTER SET latin1 DEFAULT 'ACTIVE',
  `ipadd` varchar(100) CHARACTER SET latin1 DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `FK_tbluser_info_userid` (`fk_userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `tblUserInformation`
--
CREATE TABLE IF NOT EXISTS `tblUserInformation` (
  `userid` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `userno` text,
  `lname` text,
  `fname` text,
  `mname` text,
  `bdate` date DEFAULT NULL,
  `bplace` text,
  `gender` enum('MALE','FEMALE') CHARACTER SET latin1 DEFAULT NULL,
  `email` varchar(150) CHARACTER SET latin1 DEFAULT NULL,
  `mobile` varchar(150) CHARACTER SET latin1 DEFAULT NULL,
  `createdby` int(11) DEFAULT NULL,
  `datecreated` datetime DEFAULT NULL,
  PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Triggers `tblUserInformation`
--
CREATE TRIGGER `trg_after_tblUserInformation_insert` AFTER INSERT ON `tblUserInformation` FOR EACH ROW BEGIN
	/** INSERT IN tbluser_info*/
	INSERT INTO tblUserInfo (fk_userid,username,`password`) VALUES (NEW.userid,NEW.userno,md5(NEW.lname));
    END;

-- --------------------------------------------------------

--
-- Table structure for table `tblUserRole`
--

CREATE TABLE IF NOT EXISTS `tblUserRole` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `fk_userid` int(11) UNSIGNED DEFAULT NULL COMMENT 'tblUserInformation',
  `fk_roleid` int(11) UNSIGNED DEFAULT NULL COMMENT 'tblSystemRole',
  PRIMARY KEY (`id`),
  KEY `FK_tblUserRole_tblSystemRole` (`fk_roleid`),
  KEY `FK_tblUserRole_tblUserInformation` (`fk_userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Stand-in structure for view `view_customer_info`
-- (See below for the actual view)
--
-- CREATE TABLE IF NOT EXISTS `view_customer_info` (
-- `cust_unique_id` varchar(100)
-- ,`customer_id` int(11)
-- ,`first_name` varchar(50)
-- ,`middle_name` varchar(50)
-- ,`last_name` varchar(50)
-- ,`suffix` varchar(30)
-- ,`birth_date` date
-- ,`address1` varchar(255)
-- ,`address2` varchar(255)
-- ,`country_id` int(11)
-- ,`country_desc` varchar(128)
-- ,`state_id` int(11)
-- ,`state_code` char(2)
-- ,`state_desc` varchar(22)
-- ,`city_id` int(11)
-- ,`city_desc` varchar(50)
-- ,`zipcode` varchar(20)
-- ,`home_number` varchar(20)
-- ,`mobile_number` varchar(20)
-- ,`email` varchar(100)
-- ,`company_name` varchar(100)
-- ,`phone_number` varchar(50)
-- ,`net_pay_percheck` double(10,2)
-- ,`income_source_id` int(11)
-- ,`mloc_access` int(1)
-- ,`registration` int(1)
-- ,`term_and_condition` int(1)
-- ,`income_source_desc` varchar(150)
-- ,`pay_frequency_id` int(11)
-- ,`pay_frequency_desc` varchar(150)
-- ,`next_paydate` date
-- ,`following_paydate` date
-- ,`key` blob
-- ,`credit_limit` double(10,2)
-- ,`available_credit` double(10,2)
-- ,`is_suspended` enum('YES','NO')
-- ,`credit_line_id` int(11)
-- ,`credit_approved_by` varchar(100)
-- ,`loan_total_principal_amount` double(10,2)
-- ,`loan_total_fee_amount` double(10,2)
-- ,`loan_total_amount` double(10,2)
-- ,`program_customer_id` int(11)
-- ,`program_customer_mobile` varchar(50)
-- );

-- -- --------------------------------------------------------

-- --
-- -- Stand-in structure for view `view_transaction_history`
-- -- (See below for the actual view)
-- --
-- CREATE TABLE IF NOT EXISTS `view_transaction_history` (
-- `fk_customer_id` int(11)
-- ,`amount` double(10,2)
-- ,`t_type` varchar(7)
-- ,`t_date` datetime
-- );

-- --------------------------------------------------------

--
-- Structure for view `view_customer_info`
--
CREATE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `view_customer_info`  AS  (select `a`.`cust_unique_id` AS `cust_unique_id`,`a`.`id` AS `customer_id`,`a`.`first_name` AS `first_name`,`a`.`middle_name` AS `middle_name`,`a`.`last_name` AS `last_name`,`a`.`suffix` AS `suffix`,`a`.`birth_date` AS `birth_date`,`a`.`address1` AS `address1`,`a`.`address2` AS `address2`,`c`.`country_id` AS `country_id`,`c`.`name` AS `country_desc`,`d`.`state_id` AS `state_id`,`d`.`state_code` AS `state_code`,`d`.`state` AS `state_desc`,`e`.`city_id` AS `city_id`,`e`.`city` AS `city_desc`,`a`.`zipcode` AS `zipcode`,`a`.`home_number` AS `home_number`,`a`.`mobile_number` AS `mobile_number`,`a`.`email` AS `email`,`b`.`company_name` AS `company_name`,`b`.`phone_number` AS `phone_number`,`b`.`net_pay_percheck` AS `net_pay_percheck`,`b`.`income_source` AS `income_source_id`,`h`.`mloc_access` AS `mloc_access`,`h`.`registration` AS `registration`,`h`.`term_and_condition` AS `term_and_condition`,`f`.`description` AS `income_source_desc`,`b`.`pay_frequency` AS `pay_frequency_id`,`g`.`description` AS `pay_frequency_desc`,`b`.`next_paydate` AS `next_paydate`,`b`.`following_paydate` AS `following_paydate`,`i`.`key` AS `key`,`j`.`credit_limit` AS `credit_limit`,`j`.`available_credit` AS `available_credit`,`j`.`is_suspended` AS `is_suspended`,`j`.`credit_line_id` AS `credit_line_id`,`j`.`approved_by` AS `credit_approved_by`,`k`.`total_principal_amount` AS `loan_total_principal_amount`,`k`.`total_fee_amount` AS `loan_total_fee_amount`,`k`.`total_amount` AS `loan_total_amount`,`a`.`program_customer_id` AS `program_customer_id`,`a`.`program_customer_mobile` AS `program_customer_mobile` from ((((((((((`tblCustomerBasicInfo` `a` left join `tblCustomerOtherInfo` `b` on((`b`.`fk_customer_id` = `a`.`id`))) left join `tblCountry` `c` on((`c`.`country_id` = `a`.`country`))) left join `tblState` `d` on((`d`.`state_id` = `a`.`state`))) left join `tblCity` `e` on((`e`.`city_id` = `a`.`city`))) left join `tblIncomeSource` `f` on((`f`.`id` = `b`.`income_source`))) left join `tblPayFrequency` `g` on((`g`.`id` = `b`.`pay_frequency`))) left join `tblCustomerAgreement` `h` on((`h`.`fk_customer_id` = `a`.`id`))) left join `tblApiKey` `i` on((`i`.`fk_customer_id` = `a`.`id`))) left join `tblCustomerCreditLine` `j` on((`j`.`fk_customer_id` = `a`.`id`))) left join `tblCustomerLoanTotal` `k` on((`k`.`fk_customer_id` = `a`.`id`)))) ;

-- --------------------------------------------------------

--
-- Structure for view `view_transaction_history`
--
CREATE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `view_transaction_history`  AS  select `tblCustomerLoanHistory`.`fk_customer_id` AS `fk_customer_id`,`tblCustomerLoanHistory`.`loan_amount` AS `amount`,'LOAN' AS `t_type`,`tblCustomerLoanHistory`.`loan_date` AS `t_date` from `tblCustomerLoanHistory` union select `tblCustomerPaymentHistory`.`fk_customer_id` AS `fk_customer_id`,`tblCustomerPaymentHistory`.`payment_amount` AS `amount`,'PAYMENT' AS `t_type`,`tblCustomerPaymentHistory`.`date_paid` AS `t_date` from `tblCustomerPaymentHistory` ;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `tblApiKey`
--
ALTER TABLE `tblApiKey`
  ADD CONSTRAINT `FK_tblApiKey_tblCustomerBasicInfo` FOREIGN KEY (`fk_customer_id`) REFERENCES `tblCustomerBasicInfo` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `tblCustomerAgreement`
--
ALTER TABLE `tblCustomerAgreement`
  ADD CONSTRAINT `FK_tblCustomerAgreement_tblCustomerBasicInfo` FOREIGN KEY (`fk_customer_id`) REFERENCES `tblCustomerBasicInfo` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `tblCustomerCreditLine`
--
ALTER TABLE `tblCustomerCreditLine`
  ADD CONSTRAINT `FK_tblCustomerCreditLine_tblCustomerBasicInfo` FOREIGN KEY (`fk_customer_id`) REFERENCES `tblCustomerBasicInfo` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `tblCustomerCreditLineApplication`
--
ALTER TABLE `tblCustomerCreditLineApplication`
  ADD CONSTRAINT `FK_tblCustomerCreditLineApplication_tblCustomerBasicInfo` FOREIGN KEY (`fk_customer_id`) REFERENCES `tblCustomerBasicInfo` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `tblCustomerCreditLineHistory`
--
ALTER TABLE `tblCustomerCreditLineHistory`
  ADD CONSTRAINT `FK_tblCustomerCreditLineHistory_tblCustomerBasicInfo` FOREIGN KEY (`fk_customer_id`) REFERENCES `tblCustomerBasicInfo` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `tblCustomerLoan`
--
ALTER TABLE `tblCustomerLoan`
  ADD CONSTRAINT `FK_tblCustomerLoan_tblCustomerBasicInfo` FOREIGN KEY (`fk_customer_id`) REFERENCES `tblCustomerBasicInfo` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `tblCustomerLoanApplication`
--
ALTER TABLE `tblCustomerLoanApplication`
  ADD CONSTRAINT `FK_tblCustomerLoanApplication_tblCustomerBasicInfo` FOREIGN KEY (`fk_customer_id`) REFERENCES `tblCustomerBasicInfo` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `tblCustomerLoanHistory`
--
ALTER TABLE `tblCustomerLoanHistory`
  ADD CONSTRAINT `FK_tblCustomerLoanHistory_tblCustomerBasicInfo` FOREIGN KEY (`fk_customer_id`) REFERENCES `tblCustomerBasicInfo` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `tblCustomerLoanTotal`
--
ALTER TABLE `tblCustomerLoanTotal`
  ADD CONSTRAINT `FK_tblCustomerLoanTotal_tblCustomerBasicInfo` FOREIGN KEY (`fk_customer_id`) REFERENCES `tblCustomerBasicInfo` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `tblCustomerLoanTransaction`
--
ALTER TABLE `tblCustomerLoanTransaction`
  ADD CONSTRAINT `FK_tblCustomerLoanTransaction_tblCustomerBasicInfo` FOREIGN KEY (`fk_customer_id`) REFERENCES `tblCustomerBasicInfo` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `tblCustomerOtherInfo`
--
ALTER TABLE `tblCustomerOtherInfo`
  ADD CONSTRAINT `FK_tblCustomerOtherInfo_tblCustomerBasicInfo` FOREIGN KEY (`fk_customer_id`) REFERENCES `tblCustomerBasicInfo` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `tblCustomerPayment`
--
ALTER TABLE `tblCustomerPayment`
  ADD CONSTRAINT `FK_tblCustomerPayment_tblCustomerBasicInfo` FOREIGN KEY (`fk_customer_id`) REFERENCES `tblCustomerBasicInfo` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `tblCustomerPaymentHistory`
--
ALTER TABLE `tblCustomerPaymentHistory`
  ADD CONSTRAINT `FK_tblCustomerPaymentHistory_tblCustomerBasicInfo` FOREIGN KEY (`fk_customer_id`) REFERENCES `tblCustomerBasicInfo` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `tblCustomerSettlement`
--
ALTER TABLE `tblCustomerSettlement`
  ADD CONSTRAINT `FK_tblCustomerSettlement_tblCustomerBasicInfo` FOREIGN KEY (`fk_customer_id`) REFERENCES `tblCustomerBasicInfo` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `FK_tblCustomerSettlement_tblCustomerPayment` FOREIGN KEY (`customer_payment_id`) REFERENCES `tblCustomerPayment` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `tblSystemMenuCategory`
--
ALTER TABLE `tblSystemMenuCategory`
  ADD CONSTRAINT `FK_tblSystemMenuCategory_tblSystemCategory` FOREIGN KEY (`fk_syscatid`) REFERENCES `tblSystemCategory` (`syscat_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `FK_tblSystemMenuCategory_tblSystemMenu` FOREIGN KEY (`fk_menuid`) REFERENCES `tblSystemMenu` (`menuid`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `tblSystemMenuSub`
--
ALTER TABLE `tblSystemMenuSub`
  ADD CONSTRAINT `FK_tblSystemMenuSub_tblSystemMenu` FOREIGN KEY (`fk_menuid`) REFERENCES `tblSystemMenu` (`menuid`) ON DELETE SET NULL ON UPDATE CASCADE;

--
-- Constraints for table `tblSystemRoleMenu`
--
ALTER TABLE `tblSystemRoleMenu`
  ADD CONSTRAINT `FK_tblSystemRoleMenu_tblSystemMenuCategory` FOREIGN KEY (`fk_menucat_id`) REFERENCES `tblSystemMenuCategory` (`menucat_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `FK_tblSystemRoleMenu_tblSystemRole` FOREIGN KEY (`fk_roleid`) REFERENCES `tblSystemRole` (`roleid`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `tblUserInfo`
--
ALTER TABLE `tblUserInfo`
  ADD CONSTRAINT `FK_tblUserInfo_tblUserInformation` FOREIGN KEY (`fk_userid`) REFERENCES `tblUserInformation` (`userid`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `tblUserRole`
--
ALTER TABLE `tblUserRole`
  ADD CONSTRAINT `FK_tblUserRole_tblSystemRole` FOREIGN KEY (`fk_roleid`) REFERENCES `tblSystemRole` (`roleid`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `FK_tblUserRole_tblUserInformation` FOREIGN KEY (`fk_userid`) REFERENCES `tblUserInformation` (`userid`) ON DELETE CASCADE ON UPDATE CASCADE;

COMMIT;