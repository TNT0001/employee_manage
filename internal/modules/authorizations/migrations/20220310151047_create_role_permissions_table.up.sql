CREATE TABLE IF NOT EXISTS `role_permissions` (
    `id` integer(11) PRIMARY KEY AUTO_INCREMENT,
    `role_id` integer(11) not null,
    `permission_id` integer(11) not null,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT `role_permissions_fk1` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`),
    CONSTRAINT `role_permissions_fk2` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
