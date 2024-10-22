
CREATE TABLE player (
  id UUID NOT NULL,
  api_user UUID NOT NULL,
  universe UUID NOT NULL,
  name TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  version INTEGER DEFAULT 0,
  PRIMARY KEY (id),
  FOREIGN KEY (universe) REFERENCES universe(id),
  UNIQUE (api_user, universe),
  UNIQUE (universe, name)
);

CREATE TRIGGER trigger_player_updated_at
  BEFORE UPDATE OR INSERT ON player
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at();

CREATE INDEX player_api_user_index ON player (api_user);
