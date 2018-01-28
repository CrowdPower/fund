CREATE TABLE `Users` (
    `username` VARCHAR(64) NOT NULL,
    `password` VARCHAR(128) NOT NULL,
    `email` VARCHAR(256) NOT NULL,
    `invalidatedtokens` BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (`username`)
);

CREATE TABLE `Deposits` (
    `id` CHAR(36) NOT NULL,
    `username` VARCHAR(64) NOT NULL,
    `amount` LONG NOT NULL,
    `time` VARCHAR(32) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`username`) REFERENCES Users(`username`)
);

CREATE TABLE `Payments` (
    `id` CHAR(36) NOT NULL,
    `username` VARCHAR(64) NOT NULL,
    `amount` LONG NOT NULL,
    `time` VARCHAR(32) NOT NULL,
    `url` VARCHAR(2048) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`username`) REFERENCES Users(`username`)
);
