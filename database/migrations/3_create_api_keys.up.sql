
CREATE TABLE api_keys (
  id UUID NOT NULL,
  key UUID NOT NULL,
  api_user UUID NOT NULL,
  enabled boolean DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);
