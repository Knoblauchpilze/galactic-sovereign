
CREATE TABLE api_key (
  id UUID NOT NULL,
  key UUID NOT NULL,
  api_user UUID NOT NULL,
  enabled boolean DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE,
  version INTEGER DEFAULT 0,
  PRIMARY KEY (id),
  FOREIGN KEY (api_user) REFERENCES api_user(id)
);

CREATE TRIGGER trigger_api_key_updated_at
  BEFORE UPDATE OR INSERT ON api_key
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at();
