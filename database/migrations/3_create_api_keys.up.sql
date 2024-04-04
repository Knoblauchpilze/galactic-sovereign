
CREATE TABLE api_key (
  id UUID NOT NULL,
  key UUID NOT NULL,
  api_user UUID NOT NULL,
  enabled boolean DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (api_user) REFERENCES api_user(id)
);
