
CREATE TABLE api_key (
  id UUID NOT NULL,
  key UUID NOT NULL,
  api_user UUID NOT NULL,
  valid_until TIMESTAMP WITH TIME ZONE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (api_user) REFERENCES api_user(id),
  UNIQUE (api_user)
);

CREATE INDEX api_key_api_user_index ON api_key (api_user);
CREATE INDEX api_key_key_index ON api_key (key);
