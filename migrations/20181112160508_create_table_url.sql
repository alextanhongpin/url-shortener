-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS url (
	id int unsigned AUTO_INCREMENT,
	url varchar(255) NOT NULL,
	url_crc int unsigned NOT NULL DEFAULT 0,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted_at DATETIME NOT NULL DEFAULT '1900-01-01 00:00:00',
	INDEX (url_crc),
	UNIQUE (url),
	PRIMARY KEY (id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE url;
