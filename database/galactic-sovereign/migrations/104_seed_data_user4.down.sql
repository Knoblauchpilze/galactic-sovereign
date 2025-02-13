
-- planet deut-factory
DELETE FROM planet_building WHERE planet = '717ffa52-89bd-42eb-b34d-0f994a032e35';
DELETE FROM planet_resource_storage WHERE planet = '717ffa52-89bd-42eb-b34d-0f994a032e35';
DELETE FROM planet_resource_production WHERE planet = '717ffa52-89bd-42eb-b34d-0f994a032e35';
DELETE FROM planet_resource WHERE planet = '717ffa52-89bd-42eb-b34d-0f994a032e35';

DELETE FROM homeworld WHERE planet = '717ffa52-89bd-42eb-b34d-0f994a032e35';
DELETE FROM planet WHERE id = '717ffa52-89bd-42eb-b34d-0f994a032e35';

-- better-test-user@mail-client.org -> vend-deut
DELETE FROM player WHERE id = '2bab9414-7972-4483-a8b4-fdd169d0b073';
