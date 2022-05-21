ALTER TABLE `users`
    DROP FOREIGN KEY `users_fk1`,
    DROP INDEX `personal_info_id`,
    DROP COLUMN `personal_info_id`;
