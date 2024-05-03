-- https://dba.stackexchange.com/questions/117109/how-to-manage-default-privileges-for-users-on-a-database-vs-schema/117661#117661
CREATE DATABASE db_user_service OWNER user_service_admin;
REVOKE ALL ON DATABASE db_user_service FROM public;

GRANT CONNECT ON DATABASE db_user_service TO user_service_user;

\connect db_user_service

CREATE SCHEMA user_service_schema AUTHORIZATION user_service_admin;

SET search_path = user_service_schema;

ALTER ROLE user_service_admin IN DATABASE db_user_service SET search_path = user_service_schema;
ALTER ROLE user_service_manager IN DATABASE db_user_service SET search_path = user_service_schema;
ALTER ROLE user_service_user IN DATABASE db_user_service SET search_path = user_service_schema;

GRANT USAGE  ON SCHEMA user_service_schema TO user_service_user;
GRANT CREATE ON SCHEMA user_service_schema TO user_service_admin;

ALTER DEFAULT PRIVILEGES FOR ROLE user_service_admin
GRANT SELECT ON TABLES TO user_service_user;

ALTER DEFAULT PRIVILEGES FOR ROLE user_service_admin
GRANT INSERT, UPDATE, DELETE ON TABLES TO user_service_manager;
