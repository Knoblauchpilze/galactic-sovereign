
CREATE TABLE building_action (
  id uuid NOT NULL,
  planet uuid NOT NULL,
  building uuid NOT NULL,
  current_level integer NOT NULL,
  desired_level integer NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  completed_at TIMESTAMP WITH TIME ZONE NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (planet) REFERENCES planet(id),
  FOREIGN KEY (building) REFERENCES building(id),
  UNIQUE (planet)
);