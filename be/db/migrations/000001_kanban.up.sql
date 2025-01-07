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