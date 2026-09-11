
CREATE TABLE ship(
  id UUID NOT NULL,
  name text NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);

CREATE TRIGGER trigger_ship_updated_at
  BEFORE UPDATE OR INSERT ON ship
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at();

CREATE TABLE ship_cost(
  ship UUID NOT NULL,
  resource UUID NOT NULL,
  cost INTEGER NOT NULL,
  FOREIGN KEY (ship) REFERENCES ship(id),
  FOREIGN KEY (resource) REFERENCES resource(id),
  UNIQUE (ship, resource)
);

CREATE TABLE ship_building_requirement(
  ship UUID NOT NULL,
  building UUID NOT NULL,
  level INTEGER NOT NULL,
  FOREIGN KEY (ship) REFERENCES ship(id),
  FOREIGN KEY (building) REFERENCES building(id),
  UNIQUE (ship, building)
);
