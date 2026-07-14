
CREATE TABLE building_action(
  id UUID NOT NULL,
  planet UUID NOT NULL,
  building UUID NOT NULL,
  desired_level INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  completed_at TIMESTAMP WITH TIME ZONE NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (planet) REFERENCES planet(id),
  FOREIGN KEY (building) REFERENCES building(id),
  UNIQUE (planet)
);

CREATE INDEX building_action_planet_index ON building_action(planet);

CREATE TABLE building_action_cost(
  action UUID NOT NULL,
  resource UUID NOT NULL,
  amount INTEGER NOT NULL,
  FOREIGN KEY (action) REFERENCES building_action(id),
  FOREIGN KEY (resource) REFERENCES resource(id),
  UNIQUE (action, resource)
);

CREATE TABLE building_action_resource_production(
  action UUID NOT NULL,
  resource UUID NOT NULL,
  production INTEGER NOT NULL,
  FOREIGN KEY (action) REFERENCES building_action(id),
  FOREIGN KEY (resource) REFERENCES resource(id),
  UNIQUE (action, resource)
);

CREATE TABLE building_action_resource_storage(
  action UUID NOT NULL,
  resource UUID NOT NULL,
  storage INTEGER NOT NULL,
  FOREIGN KEY (action) REFERENCES building_action(id),
  FOREIGN KEY (resource) REFERENCES resource(id),
  UNIQUE (action, resource)
);
