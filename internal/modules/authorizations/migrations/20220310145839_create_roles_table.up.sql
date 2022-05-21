CREATE TABLE IF NOT EXISTS `roles` (
    `id` integer(11) PRIMARY KEY AUTO_INCREMENT,
    `name` varchar(25) not null unique,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
