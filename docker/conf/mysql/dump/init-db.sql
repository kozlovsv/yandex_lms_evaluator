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
CREATE TABLE IF NOT EXISTS `expression` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `value` varchar(255) NOT NULL,
  `status` tinyint(3) unsigned NOT NULL DEFAULT 0 COMMENT '0 - новый\r\n1 - в процессе вычисления\r\n2 - ошибка\r\n3 - вычислен',
  `result` varchar(50) NOT NULL DEFAULT '0',
  `idempotency_key` varchar(36) NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `result_text` varchar(50) DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idempotency_key` (`idempotency_key`)
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Дамп данных таблицы evaluator.expression: ~16 rows (приблизительно)
DELETE FROM `expression`;
INSERT INTO `expression` (`id`, `value`, `status`, `result`, `idempotency_key`, `updated_at`, `created_at`, `result_text`) VALUES
	(1, '((23+587)*(4+5)*(6+7))/(5-3+4)+45+78', 0, '0', '786706b8-ed80-443a-80f6-ea1fa8cc1b57', '2024-02-15 14:36:25', '2024-02-14 16:33:00', ''),
	(2, '2-2**2-2', 0, '0', '786706b8-ed80-443a-80f6-ea1fa8cc1b58', '2024-02-15 14:36:25', '2024-02-14 16:32:48', ''),
	(3, '(2+5)/(11+2)', 0, '0', '786706b8-ed80-443a-80f6-ea1fa8cc1b54', '2024-02-15 14:36:25', '2024-02-14 16:32:48', ''),
	(6, '2+5*8', 0, '0', 'dd820bf9-be7e-4947-bac6-1b442929afae', '2024-02-15 14:36:25', '2024-02-14 16:17:00', ''),
	(10, '12+2+3', 0, '0', 'e3e7b159-c731-41e4-97b3-ebf2559ef072', '2024-02-15 14:36:25', '2024-02-14 16:17:01', ''),
	(13, '2+2*8', 0, '0', 'cae15d2b-5aaa-44b6-980b-27d35db0b722', '2024-02-15 14:36:25', '2024-02-14 16:32:47', ''),
	(14, '2+4+5+6', 0, '0', '2acfce29-84bc-486c-8f64-dbcb94f9fec9', '2024-02-15 14:36:25', '2024-02-14 16:33:06', ''),
	(15, '2*3', 0, '0', 'be22b645-60fe-46be-8163-2ac50c4c8806', '2024-02-15 14:36:25', '2024-02-14 16:33:12', ''),
	(16, '2*3*4/5+5', 0, '0', '3958beb8-7103-42a5-927e-915b05af7436', '2024-02-15 14:36:25', '2024-02-14 16:33:23', ''),
	(17, '2*3*4', 0, '0', '711822c7-edba-4f29-9ade-1811fdd194b8', '2024-02-15 14:36:25', '2024-02-14 16:33:18', ''),
	(18, '2*2*2*2*2*2*2', 0, '0', 'ef6656aa-ef21-40f9-8804-cb95250ccb30', '2024-02-15 14:36:25', '2024-02-14 16:33:40', ''),
	(19, '2*2*2*2', 0, '0', '9ff01b85-78cb-45a3-b6c6-db7510902db2', '2024-02-15 14:36:25', '2024-02-14 16:33:30', ''),
	(20, '2*3*3*4', 0, '0', '888b78b3-8785-42d6-a4f0-711d8a19e82c', '2024-02-15 14:36:25', '2024-02-14 16:33:29', ''),
	(21, '2*3*4*5', 0, '0', '3b1f4c6c-a36d-4a90-8a8e-16f5298562b9', '2024-02-15 14:36:25', '2024-02-14 16:33:48', ''),
	(23, '2*2*2*2*2*2', 0, '0', 'e4e168fc-8e63-4f48-9807-56d151cd2bde', '2024-02-15 14:36:25', '2024-02-14 16:34:06', ''),
	(24, '2+2', 0, '0', '98ab7712-84de-4f04-a420-c74076ea6109', '2024-02-15 14:36:25', '2024-02-14 16:33:58', ''),
	(25, '8*8', 0, '0', 'f5068745-4195-4e24-9308-4b3271b07f0c', '2024-02-15 14:36:25', '2024-02-14 16:34:15', ''),
	(26, '5*5', 0, '0', '4134db9c-42a3-491f-9a57-675e492f9c3a', '2024-02-15 14:36:25', '2024-02-14 16:34:21', ''),
	(27, '5*5*5*5', 0, '0', '8c52abae-6bc3-4a6b-ba0a-25056ef4cb93', '2024-02-15 14:36:25', '2024-02-14 16:34:28', '');

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

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
