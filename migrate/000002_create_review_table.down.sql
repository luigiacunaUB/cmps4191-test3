CREATE TABLE IF NOT EXISTS review(
    id bigserial PRIMARY KEY,
    prodID bigserial,
    reviewText text NOT NULL,
    addedDate timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);