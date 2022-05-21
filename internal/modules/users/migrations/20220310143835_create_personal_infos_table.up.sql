create table if not exists `personal_infos` (
    `id` integer(11) PRIMARY KEY AUTO_INCREMENT,
    `name` varchar(255),
    `sur_name` varchar(255),
    `full_name` varchar(255) not null,
    `email` varchar(255) unique,
    `facebook` varchar(255) unique,
    `linken` varchar(255) unique,
    `phone_number` bigint unique,
    `address` varchar(255),
    `on_probationary_period` bool not null,
    `start_probationary_date` date,
    `end_probationary_date` date,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` int(11) UNSIGNED NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
