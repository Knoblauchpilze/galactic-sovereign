
DROP TRIGGER trigger_limits_updated_at ON limits;
DROP TRIGGER trigger_user_limit_updated_at ON user_limit;
DROP TRIGGER trigger_acl_permission_updated_at ON acl_permission;
DROP TRIGGER trigger_acl_updated_at ON acl;

DROP TABLE limits;
DROP TABLE user_limit;
DROP TABLE acl_permission;
DROP TABLE acl;
