
CREATE TABLE resource(
  id UUID NOT NULL,
  name text,
  start_amount INTEGER NOT NULL,
  start_production INTEGER NOT NULL,
  start_storage INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE (name)
);

CREATE TRIGGER trigger_resource_updated_at
  BEFORE UPDATE OR INSERT ON resource
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at();

CREATE TABLE resource_metabolization_rate_building(
  resource UUID NOT NULL,
  hours_per_unit NUMERIC(15, 5) NOT NULL,
  FOREIGN KEY (resource) REFERENCES resource(id),
  UNIQUE (resource)
);

CREATE TABLE resource_metabolization_rate_shipyard(
  resource UUID NOT NULL,
  hours_per_unit NUMERIC(15, 5) NOT NULL,
  FOREIGN KEY (resource) REFERENCES resource(id),
  UNIQUE (resource)
);
