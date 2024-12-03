CREATE TABLE `blog_api`.`translations` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `type` ENUM('category', 'article') NOT NULL,
  `object_id` INT NOT NULL,
  `field` VARCHAR(255) NOT NULL,
  `language` ENUM('uk', 'en') NOT NULL,
  `content` LONGTEXT NOT NULL,
  PRIMARY KEY (`id`)
);