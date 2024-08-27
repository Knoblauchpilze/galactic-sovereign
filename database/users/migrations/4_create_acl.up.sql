
CREATE TABLE acl (
  id UUID NOT NULL,
  api_user UUID NOT NULL,
  resource TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (api_user) REFERENCES api_user(id),
  UNIQUE (api_user, resource)
);

CREATE TRIGGER trigger_acl_updated_at
  BEFORE UPDATE OR INSERT ON acl
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at();

CREATE INDEX acl_api_user_index ON acl (api_user);

CREATE TABLE acl_permission (
  acl UUID NOT NULL,
  permission TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (acl, permission),
  FOREIGN KEY (acl) REFERENCES acl(id)
);

CREATE TRIGGER trigger_acl_permission_updated_at
  BEFORE UPDATE OR INSERT ON acl_permission
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at();

CREATE TABLE user_limit (
  id UUID NOT NULL,
  name TEXT NOT NULL,
  api_user UUID NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (api_user) REFERENCES api_user(id),
  UNIQUE (name, api_user)
);

CREATE TRIGGER trigger_user_limit_updated_at
  BEFORE UPDATE OR INSERT ON user_limit
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at();

CREATE TABLE limits (
  id UUID NOT NULL,
  name TEXT NOT NULL,
  value TEXT NOT NULL,
  user_limit UUID NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (user_limit) REFERENCES user_limit(id),
  UNIQUE (name, user_limit)
);

CREATE TRIGGER trigger_limits_updated_at
  BEFORE UPDATE OR INSERT ON limits
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at();

CREATE INDEX limits_name_index ON limits (name);
