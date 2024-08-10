
-- user1 -> the-great-test-user-overlord / oberon
INSERT INTO stellar_dominion_schema.player ("id", "api_user", "universe", "name")
  VALUES (
    '92a686c0-9a0a-4bc3-aa1b-9a57ed7f09d5',
    '0463ed3d-bfc9-4c10-b6ee-c223bbca0fab',
    '9682f17b-f5f0-4eda-a747-2537d2151837',
    'haha'
  );

-- homeworld
INSERT INTO stellar_dominion_schema.planet ("id", "player", "name")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    '92a686c0-9a0a-4bc3-aa1b-9a57ed7f09d5',
    'homeworld'
  );
INSERT INTO stellar_dominion_schema.homeworld ("player", "planet")
  VALUES (
    '92a686c0-9a0a-4bc3-aa1b-9a57ed7f09d5',
    '167bd268-6ae7-4cf4-a359-9534beabfeff'
  );

INSERT INTO stellar_dominion_schema.planet_resource ("planet", "resource", "amount")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    514
  );
INSERT INTO stellar_dominion_schema.planet_resource ("planet", "resource", "amount")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    17
  );

INSERT INTO stellar_dominion_schema.planet_building ("planet", "building", "level")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    'd176e82d-f2ca-4611-996b-c4804096caef',
    0
  );
INSERT INTO stellar_dominion_schema.planet_building ("planet", "building", "level")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    '3904d34d-9a7e-47d4-a332-091700e2c5c3',
    1
  );

-- colony
INSERT INTO stellar_dominion_schema.planet ("id", "player", "name")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    '92a686c0-9a0a-4bc3-aa1b-9a57ed7f09d5',
    'colony'
  );

INSERT INTO stellar_dominion_schema.planet_resource ("planet", "resource", "amount")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    500
  );
INSERT INTO stellar_dominion_schema.planet_resource ("planet", "resource", "amount")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    500
  );

INSERT INTO stellar_dominion_schema.planet_building ("planet", "building", "level")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    'd176e82d-f2ca-4611-996b-c4804096caef',
    0
  );
INSERT INTO stellar_dominion_schema.planet_building ("planet", "building", "level")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    '3904d34d-9a7e-47d4-a332-091700e2c5c3',
    0
  );

-- test-user@provider.com -> throwaway-account / aquarius
INSERT INTO stellar_dominion_schema.player ("id", "api_user", "universe", "name")
  VALUES (
    '04a7477c-a66b-4c47-9c17-ac209183c7a4',
    '0463ed3d-bfc9-4c10-b6ee-c223bbca0fab',
    '0ac6c027-11d6-47e6-ab15-514cfac48200',
    'throwaway-account'
  );

-- a-new-beginning
INSERT INTO stellar_dominion_schema.planet ("id", "player", "name")
  VALUES (
    'fafd18e9-2db6-439a-aaf3-010771d694c9',
    '04a7477c-a66b-4c47-9c17-ac209183c7a4',
    'a-new-beginning'
  );
INSERT INTO stellar_dominion_schema.homeworld ("player", "planet")
  VALUES (
    '04a7477c-a66b-4c47-9c17-ac209183c7a4',
    'fafd18e9-2db6-439a-aaf3-010771d694c9'
  );

INSERT INTO stellar_dominion_schema.planet_building ("planet", "building", "level")
  VALUES (
    'fafd18e9-2db6-439a-aaf3-010771d694c9',
    'd176e82d-f2ca-4611-996b-c4804096caef',
    2
  );
INSERT INTO stellar_dominion_schema.planet_building ("planet", "building", "level")
  VALUES (
    'fafd18e9-2db6-439a-aaf3-010771d694c9',
    '3904d34d-9a7e-47d4-a332-091700e2c5c3',
    1
  );

INSERT INTO stellar_dominion_schema.planet_resource ("planet", "resource", "amount")
  VALUES (
    'fafd18e9-2db6-439a-aaf3-010771d694c9',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    50
  );
INSERT INTO stellar_dominion_schema.planet_resource ("planet", "resource", "amount")
  VALUES (
    'fafd18e9-2db6-439a-aaf3-010771d694c9',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    603
  );

-- another-test-user@another-provider.com -> very-nice-pseudo / oberon
INSERT INTO stellar_dominion_schema.player ("id", "api_user", "universe", "name")
  VALUES (
    'e8db2006-3e35-49cd-8e1f-726491660a00',
    '4f26321f-d0ea-46a3-83dd-6aa1c6053aaf',
    '9682f17b-f5f0-4eda-a747-2537d2151837',
    'very-nice-pseudo'
  );

-- my-awesome-planet
INSERT INTO stellar_dominion_schema.planet ("id", "player", "name")
  VALUES (
    '00058def-e81d-43bb-aacf-a8402115449d',
    'e8db2006-3e35-49cd-8e1f-726491660a00',
    'my-awesome-planet'
  );
INSERT INTO stellar_dominion_schema.homeworld ("player", "planet")
  VALUES (
    'e8db2006-3e35-49cd-8e1f-726491660a00',
    '00058def-e81d-43bb-aacf-a8402115449d'
  );

INSERT INTO stellar_dominion_schema.planet_resource ("planet", "resource", "amount")
  VALUES (
    '00058def-e81d-43bb-aacf-a8402115449d',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    887
  );
INSERT INTO stellar_dominion_schema.planet_resource ("planet", "resource", "amount")
  VALUES (
    '00058def-e81d-43bb-aacf-a8402115449d',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    332
  );

INSERT INTO stellar_dominion_schema.planet_building ("planet", "building", "level")
  VALUES (
    '00058def-e81d-43bb-aacf-a8402115449d',
    'd176e82d-f2ca-4611-996b-c4804096caef',
    1
  );
INSERT INTO stellar_dominion_schema.planet_building ("planet", "building", "level")
  VALUES (
    '00058def-e81d-43bb-aacf-a8402115449d',
    '3904d34d-9a7e-47d4-a332-091700e2c5c3',
    0
  );

-- better-test-user@mail-client.org -> vend-deut / aquarius
INSERT INTO stellar_dominion_schema.player ("id", "api_user", "universe", "name")
  VALUES (
    '2bab9414-7972-4483-a8b4-fdd169d0b073',
    '00b265e6-6638-4b1b-aeac-5898c7307eb8',
    '0ac6c027-11d6-47e6-ab15-514cfac48200',
    'vend-deut'
  );

-- deut-factory
INSERT INTO stellar_dominion_schema.planet ("id", "player", "name")
  VALUES (
    '717ffa52-89bd-42eb-b34d-0f994a032e35',
    '2bab9414-7972-4483-a8b4-fdd169d0b073',
    'deut-factory'
  );
INSERT INTO stellar_dominion_schema.homeworld ("player", "planet")
  VALUES (
    '2bab9414-7972-4483-a8b4-fdd169d0b073',
    '717ffa52-89bd-42eb-b34d-0f994a032e35'
  );

INSERT INTO stellar_dominion_schema.planet_resource ("planet", "resource", "amount")
  VALUES (
    '717ffa52-89bd-42eb-b34d-0f994a032e35',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    500
  );
INSERT INTO stellar_dominion_schema.planet_resource ("planet", "resource", "amount")
  VALUES (
    '717ffa52-89bd-42eb-b34d-0f994a032e35',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    499
  );

INSERT INTO stellar_dominion_schema.planet_building ("planet", "building", "level")
  VALUES (
    '717ffa52-89bd-42eb-b34d-0f994a032e35',
    'd176e82d-f2ca-4611-996b-c4804096caef',
    2
  );
INSERT INTO stellar_dominion_schema.planet_building ("planet", "building", "level")
  VALUES (
    '717ffa52-89bd-42eb-b34d-0f994a032e35',
    '3904d34d-9a7e-47d4-a332-091700e2c5c3',
    3
  );


