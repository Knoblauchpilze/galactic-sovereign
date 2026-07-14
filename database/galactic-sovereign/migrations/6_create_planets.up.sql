
CREATE TABLE planet(
  id UUID NOT NULL,
  player UUID NOT NULL,
  name text NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  version INTEGER DEFAULT 0,
  PRIMARY KEY (id),
  FOREIGN KEY (player) REFERENCES player(id)
);

CREATE TABLE homeworld(
  player UUID NOT NULL,
  planet UUID NOT NULL,
  PRIMARY KEY (player, planet),
  UNIQUE (player),
  FOREIGN KEY (player) REFERENCES player(id),
  FOREIGN KEY (planet) REFERENCES planet(id)
);

CREATE TABLE planet_resource(
  planet UUID NOT NULL,
  resource UUID NOT NULL,
  amount NUMERIC(15, 5) NOT NULL,
  FOREIGN KEY (planet) REFERENCES planet(id),
  FOREIGN KEY (resource) REFERENCES resource(id),
  UNIQUE (planet, resource)
);

CREATE INDEX planet_resource_planet_index ON planet_resource(planet);

CREATE TABLE planet_resource_production(
  planet UUID NOT NULL,
  building UUID,
  resource UUID NOT NULL,
  production INTEGER NOT NULL,
  FOREIGN KEY (planet) REFERENCES planet(id),
  FOREIGN KEY (building) REFERENCES building(id),
  FOREIGN KEY (resource) REFERENCES resource(id),
  -- https://stackoverflow.com/questions/8289100/create-unique-constraint-with-null-columns
  UNIQUE NULLS NOT DISTINCT (planet, building, resource)
);

CREATE INDEX planet_resource_production_planet_index ON planet_resource_production(planet);

CREATE TABLE planet_resource_storage(
  planet UUID NOT NULL,
  resource UUID NOT NULL,
  storage INTEGER NOT NULL,
  FOREIGN KEY (planet) REFERENCES planet(id),
  FOREIGN KEY (resource) REFERENCES resource(id),
  UNIQUE (planet, resource)
);

CREATE INDEX planet_resource_storage_planet_index ON planet_resource_storage(planet);

CREATE TABLE planet_building(
  planet UUID NOT NULL,
  building UUID NOT NULL,
  level INTEGER NOT NULL,
  FOREIGN KEY (planet) REFERENCES planet(id),
  FOREIGN KEY (building) REFERENCES building(id),
  UNIQUE (planet, building)
);

CREATE INDEX planet_building_planet_index ON planet_building(planet);
