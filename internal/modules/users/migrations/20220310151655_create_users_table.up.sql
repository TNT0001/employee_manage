create table if not exists `users` (
    `id` integer(11) PRIMARY KEY AUTO_INCREMENT,
    `keycloak_user_id` varchar(255) not null unique,
    `user_name` varchar(255) not null unique ,
    `join_date` date not null,
    `team_id` integer(11),
    `personal_info_id` integer(11) unique,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` int(11) UNSIGNED NOT NULL DEFAULT 0,
    CONSTRAINT `users_fk1` FOREIGN KEY (`personal_info_id`) REFERENCES `personal_infos` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
