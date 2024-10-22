CREATE DATABASE db_galactic_sovereign OWNER galactic_sovereign_admin;
REVOKE ALL ON DATABASE db_galactic_sovereign FROM public;

GRANT CONNECT ON DATABASE db_galactic_sovereign TO galactic_sovereign_user;

\connect db_galactic_sovereign

CREATE SCHEMA galactic_sovereign_schema AUTHORIZATION galactic_sovereign_admin;

SET search_path = galactic_sovereign_schema;

ALTER ROLE galactic_sovereign_admin IN DATABASE db_galactic_sovereign SET search_path = galactic_sovereign_schema;
ALTER ROLE galactic_sovereign_manager IN DATABASE db_galactic_sovereign SET search_path = galactic_sovereign_schema;
ALTER ROLE galactic_sovereign_user IN DATABASE db_galactic_sovereign SET search_path = galactic_sovereign_schema;

GRANT USAGE  ON SCHEMA galactic_sovereign_schema TO galactic_sovereign_user;
GRANT CREATE ON SCHEMA galactic_sovereign_schema TO galactic_sovereign_admin;

ALTER DEFAULT PRIVILEGES FOR ROLE galactic_sovereign_admin
GRANT SELECT ON TABLES TO galactic_sovereign_user;

ALTER DEFAULT PRIVILEGES FOR ROLE galactic_sovereign_admin
GRANT INSERT, UPDATE, DELETE ON TABLES TO galactic_sovereign_manager;
