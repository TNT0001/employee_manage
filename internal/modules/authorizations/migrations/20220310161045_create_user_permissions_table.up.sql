CREATE TABLE IF NOT EXISTS `user_permissions` (
    `id` integer(11) PRIMARY KEY AUTO_INCREMENT,
    `user_id` integer(11) not null,
    `permission_id` integer(11) not null,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT `user_permissions_fk1` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
