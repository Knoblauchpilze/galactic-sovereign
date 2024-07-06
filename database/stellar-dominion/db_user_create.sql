
CREATE USER stellar_dominion_admin WITH CREATEDB PASSWORD :'admin_password';
CREATE USER stellar_dominion_manager WITH PASSWORD :'manager_password';
CREATE USER stellar_dominion_user WITH PASSWORD :'user_password';

GRANT stellar_dominion_user TO stellar_dominion_manager;
GRANT stellar_dominion_manager TO stellar_dominion_admin;
