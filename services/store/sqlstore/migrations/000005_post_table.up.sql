CREATE TABLE IF NOT EXISTS {{.prefix}}posts (
    id VARCHAR(100),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    create_at    BIGINT,
    update_at    BIGINT,
    delete_at    BIGINT,
    PRIMARY KEY (id)
) {{if .mysql}}DEFAULT CHARACTER SET utf8mb4{{end}};
