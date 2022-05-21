CREATE TABLE IF NOT EXISTS `user_roles` (
    `id` integer(11) PRIMARY KEY AUTO_INCREMENT,
    `user_id` integer(11) not null,
    `role_id` integer(11) not null,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT `user_roles_fk1` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
