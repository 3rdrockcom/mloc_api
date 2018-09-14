START TRANSACTION;

UPDATE
    `tblSystemSettings`
SET
    `sms_message` = 'Welcome to MLOC! Your account has been activated and was approved for {amount} {currency} credit line.',
    `email_message` = 'Hello {firstname}, <br><br>Welcome to MLOC! Your account has been activated and was approved for {amount} {currency} credit line. <br><br>Regards, <br>MLOC Team<br>'
WHERE
    `code` = 'CANM';

UPDATE
    `tblSystemSettings`
SET
    `sms_message` = 'Welcome to MLOC! Your application for {amount} {currency} credit line is being reviewed. ',
    `email_message` = 'Hello {firstname}, <br><br>Welcome to MLOC! Your application for {amount} {currency} credit line is being reviewed. <br><br>Regards, <br>MLOC Team<br>'
WHERE
    `code` = 'CPNM';

UPDATE
    `tblSystemSettings`
SET
    `sms_message` = 'Your loan request of {amount} {currency} has been approved and  credited to your mobile wallet. Thank you for your business.',
    `email_message` = 'Hello {firstname}, <br><br>Your loan request of {amount} {currency} has been approved and  credited to your mobile wallet. Thank you for your business. <br><br>Regards, <br>MLOC Team<br>'
WHERE
    `code` = 'LANM';

UPDATE
    `tblSystemSettings`
SET
    `sms_message` = 'Your loan request of {amount} {currency} is being reviewed.',
    `email_message` = 'Hello {firstname}, <br><br>Your loan request of {amount} {currency} is being reviewed. <br><br>Regards, <br>MLOC Team<br>'
WHERE
    `code` = 'LPNM';

UPDATE
    `tblSystemSettings`
SET
    `sms_message` = 'Your payment amounting of {amount} {currency} has been received. Thank you for your business. ',
    `email_message` = 'Hello {firstname}, <br><br>Your payment amounting of {amount} {currency} has been received. Thank you for your business. <br><br>Regards, <br>MLOC Team<br>'
WHERE
    `code` = 'LPANM';

COMMIT;