
-- https://stackoverflow.com/questions/72985242/securely-create-role-in-postgres-using-sql-script-and-environment-variables
CREATE USER user_service_admin WITH CREATEDB PASSWORD :'admin_password';
CREATE USER user_service_manager WITH PASSWORD :'manager_password';
CREATE USER user_service_user WITH PASSWORD :'user_password';

GRANT user_service_user TO user_service_manager;
GRANT user_service_manager TO user_service_admin;
