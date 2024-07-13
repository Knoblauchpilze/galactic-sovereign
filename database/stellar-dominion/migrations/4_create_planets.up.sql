
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
