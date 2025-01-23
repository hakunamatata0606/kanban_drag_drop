CREATE TABLE `tasks` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(255) UNIQUE NOT NULL,
  `status_id` int NOT NULL,
  `parent_id` int,
  `title` varchar(255) NOT NULL,
  `description` varchar(255) NOT NULL
);

CREATE TABLE `status` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(255) UNIQUE NOT NULL
);

ALTER TABLE `tasks` ADD FOREIGN KEY (`status_id`) REFERENCES `status` (`id`);

ALTER TABLE `tasks` ADD FOREIGN KEY (`parent_id`) REFERENCES `tasks` (`id`);

INSERT INTO `status` (`name`) values ("idea"), ("todo"), ("inprogress"), ("inreview"), ("done");

INSERT INTO `tasks`(`name`, `status_id`, `title`, `description`) values
("name_1", 1, "title_1", "description_1"),
("name_2", 1, "title_2", "description_2"),
("name_3", 1, "title_3", "description_3"),
("name_4", 1, "title_4", "description_4");