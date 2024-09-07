
CREATE TABLE building (
  id uuid NOT NULL,
  name text NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);

CREATE TRIGGER trigger_building_updated_at
  BEFORE UPDATE OR INSERT ON building
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at();

CREATE TABLE building_cost (
  building uuid NOT NULL,
  resource uuid NOT NULL,
  cost integer NOT NULL,
  progress numeric(15, 5) NOT NULL,
  FOREIGN KEY (building) REFERENCES building(id),
  FOREIGN KEY (resource) REFERENCES resource(id),
  UNIQUE (building, resource)
);