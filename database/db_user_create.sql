
CREATE USER user_service_admin WITH CREATEDB PASSWORD 'VCNNGJjsLSmoU5nxnSBBs';
CREATE USER user_service_manager WITH PASSWORD 'Zyj94bZzKzCr7uG4QvwRB';
CREATE USER user_service_user WITH PASSWORD 'ksM5Vuj32XRWcqv2FMCJz';

GRANT user_service_user TO user_service_manager;
GRANT user_service_manager TO user_service_admin;
