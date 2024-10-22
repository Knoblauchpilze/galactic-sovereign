
CREATE USER galactic_sovereign_admin WITH CREATEDB PASSWORD :'admin_password';
CREATE USER galactic_sovereign_manager WITH PASSWORD :'manager_password';
CREATE USER galactic_sovereign_user WITH PASSWORD :'user_password';

GRANT galactic_sovereign_user TO galactic_sovereign_manager;
GRANT galactic_sovereign_manager TO galactic_sovereign_admin;
