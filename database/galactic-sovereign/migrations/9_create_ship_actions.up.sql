
CREATE TABLE ship_action(
  id UUID NOT NULL,
  planet UUID NOT NULL,
  ship UUID NOT NULL,
  count INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  next_completion_at TIMESTAMP WITH TIME ZONE NOT NULL,
  completed_at TIMESTAMP WITH TIME ZONE NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (planet) REFERENCES planet(id),
  FOREIGN KEY (ship) REFERENCES ship(id)
);

CREATE INDEX ship_action_planet_index ON ship_action(planet);
