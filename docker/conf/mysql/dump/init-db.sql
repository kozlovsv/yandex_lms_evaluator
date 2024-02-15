/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

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
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Дамп данных таблицы evaluator.expression: ~3 rows (приблизительно)
DELETE FROM `expression`;
INSERT INTO `expression` (`id`, `value`, `status`, `result`, `idempotency_key`, `updated_at`, `created_at`, `result_text`) VALUES
	(1, '((23+587)*(4+5)*(6+7))/(5-3+4)+45+78', 0, '0', '786706b8-ed80-443a-80f6-ea1fa8cc1b57', '2024-02-09 10:55:01', '2024-02-08 07:44:18', ''),
	(2, '2-2**2-2', 0, '0', '786706b8-ed80-443a-80f6-ea1fa8cc1b58', '2024-02-09 10:55:01', '2024-02-08 07:44:36', ''),
	(3, '(2+5)/(11+2)', 0, '0', '786706b8-ed80-443a-80f6-ea1fa8cc1b54', '2024-02-09 10:55:01', '2024-02-08 07:44:47', '');

-- Дамп структуры для таблица evaluator.node
CREATE TABLE IF NOT EXISTS `node` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `status` enum('Доступен','Недоступен') NOT NULL DEFAULT 'Доступен',
  `last_ping` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Дамп данных таблицы evaluator.node: ~0 rows (приблизительно)
DELETE FROM `node`;

-- Дамп структуры для таблица evaluator.settings
CREATE TABLE IF NOT EXISTS `settings` (
  `id` tinyint(4) NOT NULL AUTO_INCREMENT,
  `type` enum('Плюс','Минус','Умножение','Деление') NOT NULL,
  `value` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Таблица настроек';

-- Дамп данных таблицы evaluator.settings: ~0 rows (приблизительно)
DELETE FROM `settings`;

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
