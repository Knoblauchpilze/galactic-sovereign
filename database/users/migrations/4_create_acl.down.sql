
DROP TRIGGER trigger_limits_updated_at ON limits;
DROP TRIGGER trigger_user_limits_updated_at ON user_limits;
DROP TRIGGER trigger_acl_permissions_updated_at ON acl_permissions;
DROP TRIGGER trigger_acl_updated_at ON acl;

DROP TABLE limits;
DROP TABLE user_limits;
DROP TABLE acl_permissions;
DROP TABLE acl;
