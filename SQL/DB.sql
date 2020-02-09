/*

CREATE DATABASE phone_book;
USE phone_book;

CREATE TABLE `contacts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `phone_number` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `contactsID` INT NOT NULL REFERENCES contacts(id),
  `number` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `contacts` (`name`) VALUES
("Мама"),
("Брат"),
("Геннадий Петрович");

INSERT INTO `phone_number` (`contactsID`, `number`) VALUES
(2,	"89610079089"),
(1,	"+79621433323"),
(3, "+79623320909"),
(3, "84832543412");
