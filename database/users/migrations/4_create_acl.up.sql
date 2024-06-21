
CREATE TABLE acl (
  id UUID NOT NULL,
  resource TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE,
  PRIMARY KEY (id)
);

CREATE TRIGGER trigger_acl_updated_at
  BEFORE UPDATE OR INSERT ON acl
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at();

CREATE TABLE acl_permissions (
  acl UUID NOT NULL,
  permission TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE,
  PRIMARY KEY (acl, permission),
  FOREIGN KEY (acl) REFERENCES acl(id)
);

CREATE TRIGGER trigger_acl_permissions_updated_at
  BEFORE UPDATE OR INSERT ON acl_permissions
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at();

CREATE TABLE user_limits (
  id UUID NOT NULL,
  api_user UUID NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE,
  PRIMARY KEY (id),
  FOREIGN KEY (api_user) REFERENCES api_user(id)
);

CREATE TRIGGER trigger_user_limits_updated_at
  BEFORE UPDATE OR INSERT ON user_limits
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at();

CREATE TABLE limits (
  id UUID NOT NULL,
  name TEXT NOT NULL,
  value TEXT NOT NULL,
  user_limit UUID NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE,
  PRIMARY KEY (id),
  FOREIGN KEY (user_limit) REFERENCES user_limits(id),
  UNIQUE (name, user_limit)
);

CREATE TRIGGER trigger_limits_updated_at
  BEFORE UPDATE OR INSERT ON limits
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at();

CREATE INDEX limits_name_index ON limits (name);
