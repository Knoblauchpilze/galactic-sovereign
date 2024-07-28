
CREATE TABLE planet (
  id uuid NOT NULL,
  player uuid NOT NULL,
  name text NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (player) REFERENCES player(id)
);

CREATE TRIGGER trigger_planet_updated_at
  BEFORE UPDATE OR INSERT ON planet
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at();

CREATE TABLE homeworld (
  player uuid NOT NULL,
  planet uuid NOT NULL,
  PRIMARY KEY (player, planet),
  UNIQUE (player, planet),
  FOREIGN KEY (player) REFERENCES player(id),
  FOREIGN KEY (planet) REFERENCES planet(id)
);

CREATE TABLE resource (
  id uuid NOT NULL,
  name text,
  start_amount INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE (name)
);

CREATE TRIGGER trigger_resource_updated_at
  BEFORE UPDATE OR INSERT ON resource
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at();

CREATE TABLE planet_resource (
  planet uuid NOT NULL,
  resource uuid NOT NULL,
  amount numeric(15, 5) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  version INTEGER DEFAULT 0,
  FOREIGN KEY (planet) REFERENCES planet(id),
  FOREIGN KEY (resource) REFERENCES resource(id),
  UNIQUE (planet, resource)
);
