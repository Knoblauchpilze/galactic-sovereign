
CREATE TABLE planet (
  id uuid NOT NULL,
  player uuid NOT NULL,
  name text NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (player) REFERENCES player(id)
);

CREATE TABLE homeworld (
  player uuid NOT NULL,
  planet uuid NOT NULL,
  PRIMARY KEY (player, planet),
  UNIQUE (player),
  FOREIGN KEY (player) REFERENCES player(id),
  FOREIGN KEY (planet) REFERENCES planet(id)
);

CREATE TABLE planet_resource (
  planet uuid NOT NULL,
  resource uuid NOT NULL,
  amount numeric(15, 5) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  version INTEGER DEFAULT 0,
  FOREIGN KEY (planet) REFERENCES planet(id),
  FOREIGN KEY (resource) REFERENCES resource(id),
  UNIQUE (planet, resource)
);

CREATE INDEX planet_resource_planet_index ON planet_resource (planet);

CREATE TABLE planet_resource_production (
  planet uuid NOT NULL,
  building uuid,
  resource uuid NOT NULL,
  production INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  version INTEGER DEFAULT 0,
  FOREIGN KEY (planet) REFERENCES planet(id),
  FOREIGN KEY (building) REFERENCES building(id),
  FOREIGN KEY (resource) REFERENCES resource(id),
  -- TODO: Migrating to psql 15 would allow NULLS NOT DISTINCT to make sure
  -- that we don't allow multiple productions from NULL buildings.
  -- See: https://stackoverflow.com/questions/8289100/create-unique-constraint-with-null-columns
  UNIQUE (planet, building)
);

CREATE INDEX planet_resource_production_planet_index ON planet_resource_production (planet);

CREATE TABLE planet_resource_storage (
  planet uuid NOT NULL,
  resource uuid NOT NULL,
  storage INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  version INTEGER DEFAULT 0,
  FOREIGN KEY (planet) REFERENCES planet(id),
  FOREIGN KEY (resource) REFERENCES resource(id),
  UNIQUE (planet, resource)
);

CREATE INDEX planet_resource_storage_planet_index ON planet_resource_storage (planet);

CREATE TABLE planet_building (
  planet uuid NOT NULL,
  building uuid NOT NULL,
  level INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  version INTEGER DEFAULT 0,
  FOREIGN KEY (planet) REFERENCES planet(id),
  FOREIGN KEY (building) REFERENCES building(id),
  UNIQUE (planet, building)
);

CREATE INDEX planet_building_planet_index ON planet_building (planet);
