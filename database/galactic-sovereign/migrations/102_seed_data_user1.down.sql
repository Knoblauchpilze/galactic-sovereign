
-- planet colony
DELETE FROM planet_building WHERE planet = '110cdf6f-2103-4e34-924f-fd57eb87ea3e';
DELETE FROM planet_resource_storage WHERE planet = '110cdf6f-2103-4e34-924f-fd57eb87ea3e';
DELETE FROM planet_resource_production WHERE planet = '110cdf6f-2103-4e34-924f-fd57eb87ea3e';
DELETE FROM planet_resource WHERE planet = '110cdf6f-2103-4e34-924f-fd57eb87ea3e';

DELETE FROM homeworld WHERE planet = '110cdf6f-2103-4e34-924f-fd57eb87ea3e';
DELETE FROM planet WHERE id = '110cdf6f-2103-4e34-924f-fd57eb87ea3e';

-- planet homeworld
DELETE FROM planet_building WHERE planet = '167bd268-6ae7-4cf4-a359-9534beabfeff';
DELETE FROM planet_resource_storage WHERE planet = '167bd268-6ae7-4cf4-a359-9534beabfeff';
DELETE FROM planet_resource_production WHERE planet = '167bd268-6ae7-4cf4-a359-9534beabfeff';
DELETE FROM planet_resource WHERE planet = '167bd268-6ae7-4cf4-a359-9534beabfeff';

DELETE FROM homeworld WHERE planet = '167bd268-6ae7-4cf4-a359-9534beabfeff';
DELETE FROM planet WHERE id = '167bd268-6ae7-4cf4-a359-9534beabfeff';

-- user1 -> haha
DELETE FROM player WHERE id = '92a686c0-9a0a-4bc3-aa1b-9a57ed7f09d5';
