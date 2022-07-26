CREATE TABLE {{.prefix}}post_categories (
    id varchar(36) NOT NULL,
    category_id varchar(36) NOT NULL,
    post_id varchar(36) NOT NULL,
    PRIMARY KEY (id)
    ) {{if .mysql}}DEFAULT CHARACTER SET utf8mb4{{end}};

CREATE INDEX idx_post_categories_post_id_category_id ON {{.prefix}}post_categories(post_id, category_id);
