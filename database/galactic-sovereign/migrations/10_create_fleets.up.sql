
CREATE TABLE fleet(
  id UUID NOT NULL,
  player UUID NOT NULL,
  source UUID NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  arrival_at TIMESTAMP WITH TIME ZONE NOT NULL,
  return_at TIMESTAMP WITH TIME ZONE,
  version INTEGER DEFAULT 0,
  PRIMARY KEY (id),
  FOREIGN KEY (player) REFERENCES player(id),
  FOREIGN KEY (source) REFERENCES planet(id)
);

CREATE TABLE fleet_ship(
  fleet UUID NOT NULL,
  ship UUID NOT NULL,
  count INTEGER NOT NULL,
  FOREIGN KEY (fleet) REFERENCES fleet(id),
  FOREIGN KEY (ship) REFERENCES ship(id)
);

CREATE INDEX fleet_ship_fleet_index ON fleet_ship(fleet);
