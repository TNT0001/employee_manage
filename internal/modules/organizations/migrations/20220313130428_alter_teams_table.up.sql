ALTER TABLE `teams`
    ADD CONSTRAINT `teams_UC1` UNIQUE (`country_code`, `team_name`, `division_name`);