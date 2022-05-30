CREATE TABLE IF NOT EXISTS {{.prefix}}posts (
    id VARCHAR(100),
    slug VARCHAR(100) NOT NULL,
    sub_title VARCHAR(100),
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    cover_image VARCHAR(100),
    disabled_comment BOOLEAN NOT NULL DEFAULT FALSE,
    publishing_at BIGINT,
    create_at     BIGINT,
    update_at     BIGINT,
    delete_at     BIGINT,
    user_id varchar(36) NOT NULL,
    PRIMARY KEY (id)
) {{if .mysql}}DEFAULT CHARACTER SET utf8mb4{{end}};
