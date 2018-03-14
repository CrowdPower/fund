CREATE TABLE Users (
    username VARCHAR(64) NOT NULL,
    password VARCHAR(128) NOT NULL,
    email VARCHAR(256) NOT NULL,
    invalidatedtokens BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (username)
);

CREATE TABLE Deposits (
    id CHAR(36) NOT NULL,
    username VARCHAR(64) NOT NULL,
    amount LONG NOT NULL,
    time VARCHAR(32) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (username) REFERENCES Users(username)
);

CREATE TABLE Payments (
    id CHAR(36) NOT NULL,
    username VARCHAR(64) NOT NULL,
    amount LONG NOT NULL,
    time VARCHAR(32) NOT NULL,
    url VARCHAR(2048) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (username) REFERENCES Users(username)
);

CREATE VIEW Balances AS
SELECT Users.username, depositsum - paymentsum AS balance
FROM Users
LEFT JOIN 
    (SELECT username, IFNULL(SUM(amount), 0) AS depositsum
    FROM Deposits GROUP BY username) AS Deposits
ON Users.username = Deposits.username
LEFT JOIN 
    (SELECT username, IFNULL(SUM(amount), 0) AS paymentsum
    FROM Payments GROUP BY username) AS Payments
ON Users.username = Payments.username;

CREATE TRIGGER BalanceCheck
AFTER INSERT ON Payments
WHEN 0 > (SELECT balance FROM Balances WHERE Balances.username = NEW.username)
BEGIN
    SELECT RAISE(ROLLBACK, "Insufficient Funds");
END;
