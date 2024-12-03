CREATE TABLE `blog_api`.`related_articles` (
  `parent_article_id` INT NOT NULL,
  `related_article_id` INT NOT NULL,
  PRIMARY KEY (`parent_article_id`, `related_article_id`),
  FOREIGN KEY (`parent_article_id`) REFERENCES `article`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`related_article_id`) REFERENCES `article`(`id`) ON DELETE CASCADE
);