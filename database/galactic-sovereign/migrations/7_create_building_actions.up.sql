
CREATE TABLE building_action (
  id uuid NOT NULL,
  planet uuid NOT NULL,
  building uuid NOT NULL,
  current_level INTEGER NOT NULL,
  desired_level INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  completed_at TIMESTAMP WITH TIME ZONE NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (planet) REFERENCES planet(id),
  FOREIGN KEY (building) REFERENCES building(id),
  UNIQUE (planet)
);

CREATE INDEX building_action_planet_index ON building_action (planet);

CREATE TABLE building_action_cost (
  action uuid NOT NULL,
  resource uuid NOT NULL,
  amount INTEGER NOT NULL,
  FOREIGN KEY (action) REFERENCES building_action(id),
  FOREIGN KEY (resource) REFERENCES resource(id),
  UNIQUE (action, resource)
);

CREATE TABLE building_action_resource_production (
  action uuid NOT NULL,
  resource uuid NOT NULL,
  production INTEGER NOT NULL,
  FOREIGN KEY (action) REFERENCES building_action(id),
  FOREIGN KEY (resource) REFERENCES resource(id),
  UNIQUE (action, resource)
);

CREATE TABLE building_action_resource_storage (
  action uuid NOT NULL,
  resource uuid NOT NULL,
  storage INTEGER NOT NULL,
  FOREIGN KEY (action) REFERENCES building_action(id),
  FOREIGN KEY (resource) REFERENCES resource(id),
  UNIQUE (action, resource)
);
