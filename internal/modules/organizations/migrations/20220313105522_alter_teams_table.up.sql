ALTER TABLE `teams`
    ADD COLUMN `country_code` varchar(2) not null,
    ADD COLUMN `division_name` varchar(25),
    ADD COLUMN `kind` varchar(25);
