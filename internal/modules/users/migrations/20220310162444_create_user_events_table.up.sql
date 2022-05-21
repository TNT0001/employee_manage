create table if not exists `user_events` (
    `id` integer(11) PRIMARY KEY AUTO_INCREMENT,
    `user_id` integer(11) not null,
    `time_line` date not null,
    `description` text not null,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` int(11) UNSIGNED NOT NULL DEFAULT 0,
    CONSTRAINT `user_events_fk1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
