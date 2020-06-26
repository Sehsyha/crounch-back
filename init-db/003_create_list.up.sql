BEGIN;

CREATE TABLE IF NOT EXISTS "list"(
  id UUID PRIMARY KEY,
  name VARCHAR(60) NOT NULL,
  creation_date TIMESTAMP WITH TIME ZONE NOT NULL,
  user_id UUID NOT NULL,
  CONSTRAINT FK_list_user_id FOREIGN KEY (user_id) REFERENCES "user"(id)
);

COMMIT;
