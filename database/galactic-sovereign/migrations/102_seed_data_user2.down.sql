
-- planet a-new-beginning
DELETE FROM planet_building WHERE planet = 'fafd18e9-2db6-439a-aaf3-010771d694c9';
DELETE FROM planet_resource_storage WHERE planet = 'fafd18e9-2db6-439a-aaf3-010771d694c9';
DELETE FROM planet_resource_production WHERE planet = 'fafd18e9-2db6-439a-aaf3-010771d694c9';
DELETE FROM planet_resource WHERE planet = 'fafd18e9-2db6-439a-aaf3-010771d694c9';

DELETE FROM homeworld WHERE planet = 'fafd18e9-2db6-439a-aaf3-010771d694c9';
DELETE FROM planet WHERE id = 'fafd18e9-2db6-439a-aaf3-010771d694c9';

-- user1 -> throwaway-account
DELETE FROM player WHERE id = '04a7477c-a66b-4c47-9c17-ac209183c7a4';
