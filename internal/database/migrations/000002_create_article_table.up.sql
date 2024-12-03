CREATE TABLE `blog_api`.`article` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `category_id` INT NOT NULL,
  `company_uuid` VARCHAR(36) NOT NULL,
  `language` ENUM('uk', 'en') NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  `user_uuid` VARCHAR(36) NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`category_id`) REFERENCES `category`(`id`) ON DELETE CASCADE
);
