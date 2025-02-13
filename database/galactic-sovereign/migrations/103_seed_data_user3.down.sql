
-- planet my-awesome-planet
DELETE FROM planet_building WHERE planet = '00058def-e81d-43bb-aacf-a8402115449d';
DELETE FROM planet_resource_storage WHERE planet = '00058def-e81d-43bb-aacf-a8402115449d';
DELETE FROM planet_resource_production WHERE planet = '00058def-e81d-43bb-aacf-a8402115449d';
DELETE FROM planet_resource WHERE planet = '00058def-e81d-43bb-aacf-a8402115449d';

DELETE FROM homeworld WHERE planet = '00058def-e81d-43bb-aacf-a8402115449d';
DELETE FROM planet WHERE id = '00058def-e81d-43bb-aacf-a8402115449d';

-- another-test-user@another-provider.com -> very-nice-pseudo
DELETE FROM player WHERE id = 'e8db2006-3e35-49cd-8e1f-726491660a00';
