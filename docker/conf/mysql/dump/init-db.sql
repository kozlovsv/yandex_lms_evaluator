-- --------------------------------------------------------
-- Хост:                         127.0.0.1
-- Версия сервера:               10.6.16-MariaDB-1:10.6.16+maria~ubu2004 - mariadb.org binary distribution
-- Операционная система:         debian-linux-gnu
-- HeidiSQL Версия:              12.6.0.6765
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

-- Дамп структуры для таблица evaluator.agent
CREATE TABLE IF NOT EXISTS `agent` (
  `name` varchar(255) NOT NULL,
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '0 - доступен, 1 - недоступен',
  `last_ping` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `current_op` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`name`),
  KEY `last_ping` (`last_ping`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Дамп данных таблицы evaluator.agent: ~2 rows (приблизительно)
DELETE FROM `agent`;
INSERT INTO `agent` (`name`, `status`, `last_ping`, `current_op`) VALUES
	('Agent1', 0, '2024-02-15 12:43:12', ''),
	('Agent2', 0, '2024-02-15 12:43:12', '');

-- Дамп структуры для таблица evaluator.expression
CREATE TABLE `expression` (
	`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
	`user_id` INT(11) NOT NULL,
	`value` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_general_ci',
	`status` TINYINT(3) UNSIGNED NOT NULL DEFAULT '0' COMMENT '0 - новый\r\n1 - в процессе вычисления\r\n2 - ошибка\r\n3 - вычислен',
	`result` VARCHAR(50) NOT NULL DEFAULT '0' COLLATE 'utf8mb4_general_ci',
	`idempotency_key` VARCHAR(36) NOT NULL COLLATE 'utf8mb4_general_ci',
	`updated_at` TIMESTAMP NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	`created_at` TIMESTAMP NOT NULL DEFAULT current_timestamp(),
	`result_text` VARCHAR(50) NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `idempotency_key` (`idempotency_key`) USING BTREE
)COLLATE='utf8mb4_general_ci' ENGINE=InnoDB;


-- Дамп структуры для таблица evaluator.settings
CREATE TABLE IF NOT EXISTS `settings` (
  `type` enum('Плюс','Минус','Умножение','Деление','Таймаут проверки агента','Таймаут удаления агента') NOT NULL,
  `value` int(11) NOT NULL,
  PRIMARY KEY (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Таблица настроек';

-- Дамп данных таблицы evaluator.settings: ~6 rows (приблизительно)
DELETE FROM `settings`;
INSERT INTO `settings` (`type`, `value`) VALUES
	('Плюс', 1000),
	('Минус', 2000),
	('Умножение', 3000),
	('Деление', 4000),
	('Таймаут проверки агента', 120000),
	('Таймаут удаления агента', 240000);


CREATE TABLE IF NOT EXISTS `users` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`email` VARCHAR(50) NOT NULL COLLATE 'utf8mb4_general_ci',
	`pass_hash` BLOB NOT NULL,
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `email` (`email`) USING BTREE
)
COMMENT='Пользователи'
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB;


/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;