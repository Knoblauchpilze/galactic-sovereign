CREATE DATABASE db_stellar_dominion OWNER stellar_dominion_admin;
REVOKE ALL ON DATABASE db_stellar_dominion FROM public;

GRANT CONNECT ON DATABASE db_stellar_dominion TO stellar_dominion_user;

\connect db_stellar_dominion

CREATE SCHEMA stellar_dominion_schema AUTHORIZATION stellar_dominion_admin;

SET search_path = stellar_dominion_schema;

ALTER ROLE stellar_dominion_admin IN DATABASE db_stellar_dominion SET search_path = stellar_dominion_schema;
ALTER ROLE stellar_dominion_manager IN DATABASE db_stellar_dominion SET search_path = stellar_dominion_schema;
ALTER ROLE stellar_dominion_user IN DATABASE db_stellar_dominion SET search_path = stellar_dominion_schema;

GRANT USAGE  ON SCHEMA stellar_dominion_schema TO stellar_dominion_user;
GRANT CREATE ON SCHEMA stellar_dominion_schema TO stellar_dominion_admin;

ALTER DEFAULT PRIVILEGES FOR ROLE stellar_dominion_admin
GRANT SELECT ON TABLES TO stellar_dominion_user;

ALTER DEFAULT PRIVILEGES FOR ROLE stellar_dominion_admin
GRANT INSERT, UPDATE, DELETE ON TABLES TO stellar_dominion_manager;
