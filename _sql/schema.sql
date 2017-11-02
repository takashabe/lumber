use lumber;

CREATE TABLE IF NOT EXISTS entries (
  `id`         int          NOT NULL AUTO_INCREMENT,
  `title`      varchar(256) NOT NULL,
  `content`    text         NOT NULL,
  `status`     int          NOT NULL,
  `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS entryStatus (
  `id` int NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
