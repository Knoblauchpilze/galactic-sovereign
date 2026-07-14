
CREATE TABLE universe(
  id UUID NOT NULL,
  name TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  version INTEGER DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE (name)
);

CREATE TABLE universe_topology (
  universe UUID NOT NULL,
  galaxies INTEGER NOT NULL,
  solar_systems INTEGER NOT NULL,
  orbits INTEGER NOT NULL,
  UNIQUE (universe),
  FOREIGN KEY (universe) REFERENCES universe(id)
);

CREATE TRIGGER trigger_universe_updated_at
  BEFORE UPDATE OR INSERT ON universe
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at();
