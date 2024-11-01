CREATE TABLE IF NOT EXISTS product(
    id bigserial PRIMARY KEY,
    prodName text NOT NULL,
    addedDate timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);