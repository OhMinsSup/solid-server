CREATE TABLE {{.prefix}}categories (
    id varchar(36) NOT NULL,
    name varchar(100) NOT NULL,
    create_at BIGINT,
    update_at BIGINT,
    delete_at BIGINT,
    PRIMARY KEY (id)
    ) {{if .mysql}}DEFAULT CHARACTER SET utf8mb4{{end}};
