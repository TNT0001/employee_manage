ALTER TABLE `personal_infos`
    MODIFY COLUMN `phone_number` varchar(11) unique ,
    ADD COLUMN `user_id` integer(11) unique not null ,
    ADD CONSTRAINT `personal_infos_fk1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
