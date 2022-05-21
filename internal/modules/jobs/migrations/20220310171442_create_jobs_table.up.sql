CREATE TABLE IF NOT EXISTS `jobs` (
    `id` integer(11) PRIMARY KEY AUTO_INCREMENT,
    `user_id` integer(11) not null,
    `project_name` varchar(255) not null,
    `assign_type_id`  integer(11),
    `assign_percent`  integer(11),
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
