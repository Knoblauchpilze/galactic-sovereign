
CREATE TABLE resource (
  id uuid NOT NULL,
  name text,
  start_amount INTEGER NOT NULL,
  start_production INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE (name)
);

CREATE TRIGGER trigger_resource_updated_at
  BEFORE UPDATE OR INSERT ON resource
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at();
