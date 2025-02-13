
-- user1 -> haha / oberon
INSERT INTO galactic_sovereign_schema.player ("id", "api_user", "universe", "name")
  VALUES (
    '92a686c0-9a0a-4bc3-aa1b-9a57ed7f09d5',
    '0463ed3d-bfc9-4c10-b6ee-c223bbca0fab',
    '9682f17b-f5f0-4eda-a747-2537d2151837',
    'haha'
  );

-- homeworld
INSERT INTO galactic_sovereign_schema.planet ("id", "player", "name", "created_at", "updated_at")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    '92a686c0-9a0a-4bc3-aa1b-9a57ed7f09d5',
    'homeworld',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.homeworld ("player", "planet")
  VALUES (
    '92a686c0-9a0a-4bc3-aa1b-9a57ed7f09d5',
    '167bd268-6ae7-4cf4-a359-9534beabfeff'
  );

INSERT INTO galactic_sovereign_schema.planet_resource ("planet", "resource", "amount", "created_at", "updated_at")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    1514,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_resource ("planet", "resource", "amount", "created_at", "updated_at")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    517,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

INSERT INTO galactic_sovereign_schema.planet_resource_production ("planet", "building", "resource", "production", "created_at", "updated_at")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    NULL,
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    30,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_resource_production ("planet", "building", "resource", "production", "created_at", "updated_at")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    NULL,
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    15,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

INSERT INTO galactic_sovereign_schema.planet_resource_storage ("planet", "resource", "storage", "created_at", "updated_at")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    10000,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_resource_storage ("planet", "resource", "storage", "created_at", "updated_at")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    9000,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

INSERT INTO galactic_sovereign_schema.planet_building ("planet", "building", "level", "created_at", "updated_at")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    'd176e82d-f2ca-4611-996b-c4804096caef',
    0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_building ("planet", "building", "level", "created_at", "updated_at")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    '3904d34d-9a7e-47d4-a332-091700e2c5c3',
    1,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_building ("planet", "building", "level", "created_at", "updated_at")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    '22b4c0c3-c8e5-4493-89fc-522fdbb0beee',
    0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_building ("planet", "building", "level", "created_at", "updated_at")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    'd9c8df28-bb71-4be4-8702-ce2bea8bd943',
    0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

-- colony
INSERT INTO galactic_sovereign_schema.planet ("id", "player", "name", "created_at", "updated_at")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    '92a686c0-9a0a-4bc3-aa1b-9a57ed7f09d5',
    'colony',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

INSERT INTO galactic_sovereign_schema.planet_resource ("planet", "resource", "amount", "created_at", "updated_at")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    500,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_resource ("planet", "resource", "amount", "created_at", "updated_at")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    500,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

INSERT INTO galactic_sovereign_schema.planet_resource_production ("planet", "building", "resource", "production", "created_at", "updated_at")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    NULL,
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    30,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_resource_production ("planet", "building", "resource", "production", "created_at", "updated_at")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    NULL,
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    15,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

INSERT INTO galactic_sovereign_schema.planet_resource_storage ("planet", "resource", "storage", "created_at", "updated_at")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    10000,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_resource_storage ("planet", "resource", "storage", "created_at", "updated_at")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    9000,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

INSERT INTO galactic_sovereign_schema.planet_building ("planet", "building", "level", "created_at", "updated_at")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    'd176e82d-f2ca-4611-996b-c4804096caef',
    0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_building ("planet", "building", "level", "created_at", "updated_at")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    '3904d34d-9a7e-47d4-a332-091700e2c5c3',
    0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_building ("planet", "building", "level", "created_at", "updated_at")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    '22b4c0c3-c8e5-4493-89fc-522fdbb0beee',
    0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_building ("planet", "building", "level", "created_at", "updated_at")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    'd9c8df28-bb71-4be4-8702-ce2bea8bd943',
    0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

-- user1 -> throwaway-account / aquarius
INSERT INTO galactic_sovereign_schema.player ("id", "api_user", "universe", "name")
  VALUES (
    '04a7477c-a66b-4c47-9c17-ac209183c7a4',
    '0463ed3d-bfc9-4c10-b6ee-c223bbca0fab',
    '0ac6c027-11d6-47e6-ab15-514cfac48200',
    'throwaway-account'
  );

-- a-new-beginning
INSERT INTO galactic_sovereign_schema.planet ("id", "player", "name", "created_at", "updated_at")
  VALUES (
    'fafd18e9-2db6-439a-aaf3-010771d694c9',
    '04a7477c-a66b-4c47-9c17-ac209183c7a4',
    'a-new-beginning',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.homeworld ("player", "planet")
  VALUES (
    '04a7477c-a66b-4c47-9c17-ac209183c7a4',
    'fafd18e9-2db6-439a-aaf3-010771d694c9'
  );

INSERT INTO galactic_sovereign_schema.planet_resource ("planet", "resource", "amount", "created_at", "updated_at")
  VALUES (
    'fafd18e9-2db6-439a-aaf3-010771d694c9',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    50,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_resource ("planet", "resource", "amount", "created_at", "updated_at")
  VALUES (
    'fafd18e9-2db6-439a-aaf3-010771d694c9',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    603,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

INSERT INTO galactic_sovereign_schema.planet_resource_production ("planet", "building", "resource", "production", "created_at", "updated_at")
  VALUES (
    'fafd18e9-2db6-439a-aaf3-010771d694c9',
    NULL,
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    30,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_resource_production ("planet", "building", "resource", "production", "created_at", "updated_at")
  VALUES (
    'fafd18e9-2db6-439a-aaf3-010771d694c9',
    NULL,
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    15,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

INSERT INTO galactic_sovereign_schema.planet_resource_storage ("planet", "resource", "storage", "created_at", "updated_at")
  VALUES (
    'fafd18e9-2db6-439a-aaf3-010771d694c9',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    10000,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_resource_storage ("planet", "resource", "storage", "created_at", "updated_at")
  VALUES (
    'fafd18e9-2db6-439a-aaf3-010771d694c9',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    9000,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

INSERT INTO galactic_sovereign_schema.planet_building ("planet", "building", "level", "created_at", "updated_at")
  VALUES (
    'fafd18e9-2db6-439a-aaf3-010771d694c9',
    'd176e82d-f2ca-4611-996b-c4804096caef',
    2,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_building ("planet", "building", "level", "created_at", "updated_at")
  VALUES (
    'fafd18e9-2db6-439a-aaf3-010771d694c9',
    '3904d34d-9a7e-47d4-a332-091700e2c5c3',
    1,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_building ("planet", "building", "level", "created_at", "updated_at")
  VALUES (
    'fafd18e9-2db6-439a-aaf3-010771d694c9',
    '22b4c0c3-c8e5-4493-89fc-522fdbb0beee',
    0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_building ("planet", "building", "level", "created_at", "updated_at")
  VALUES (
    'fafd18e9-2db6-439a-aaf3-010771d694c9',
    'd9c8df28-bb71-4be4-8702-ce2bea8bd943',
    0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

-- another-test-user@another-provider.com -> very-nice-pseudo / oberon
INSERT INTO galactic_sovereign_schema.player ("id", "api_user", "universe", "name")
  VALUES (
    'e8db2006-3e35-49cd-8e1f-726491660a00',
    '4f26321f-d0ea-46a3-83dd-6aa1c6053aaf',
    '9682f17b-f5f0-4eda-a747-2537d2151837',
    'very-nice-pseudo'
  );

-- my-awesome-planet
INSERT INTO galactic_sovereign_schema.planet ("id", "player", "name", "created_at", "updated_at")
  VALUES (
    '00058def-e81d-43bb-aacf-a8402115449d',
    'e8db2006-3e35-49cd-8e1f-726491660a00',
    'my-awesome-planet',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.homeworld ("player", "planet")
  VALUES (
    'e8db2006-3e35-49cd-8e1f-726491660a00',
    '00058def-e81d-43bb-aacf-a8402115449d'
  );

INSERT INTO galactic_sovereign_schema.planet_resource ("planet", "resource", "amount", "created_at", "updated_at")
  VALUES (
    '00058def-e81d-43bb-aacf-a8402115449d',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    887,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_resource ("planet", "resource", "amount", "created_at", "updated_at")
  VALUES (
    '00058def-e81d-43bb-aacf-a8402115449d',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    332,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

INSERT INTO galactic_sovereign_schema.planet_resource_production ("planet", "building", "resource", "production", "created_at", "updated_at")
  VALUES (
    '00058def-e81d-43bb-aacf-a8402115449d',
    NULL,
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    30,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_resource_production ("planet", "building", "resource", "production", "created_at", "updated_at")
  VALUES (
    '00058def-e81d-43bb-aacf-a8402115449d',
    NULL,
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    15,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

INSERT INTO galactic_sovereign_schema.planet_resource_storage ("planet", "resource", "storage", "created_at", "updated_at")
  VALUES (
    '00058def-e81d-43bb-aacf-a8402115449d',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    10000,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_resource_storage ("planet", "resource", "storage", "created_at", "updated_at")
  VALUES (
    '00058def-e81d-43bb-aacf-a8402115449d',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    9000,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

INSERT INTO galactic_sovereign_schema.planet_building ("planet", "building", "level", "created_at", "updated_at")
  VALUES (
    '00058def-e81d-43bb-aacf-a8402115449d',
    'd176e82d-f2ca-4611-996b-c4804096caef',
    1,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_building ("planet", "building", "level", "created_at", "updated_at")
  VALUES (
    '00058def-e81d-43bb-aacf-a8402115449d',
    '3904d34d-9a7e-47d4-a332-091700e2c5c3',
    0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_building ("planet", "building", "level", "created_at", "updated_at")
  VALUES (
    '00058def-e81d-43bb-aacf-a8402115449d',
    '22b4c0c3-c8e5-4493-89fc-522fdbb0beee',
    0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_building ("planet", "building", "level", "created_at", "updated_at")
  VALUES (
    '00058def-e81d-43bb-aacf-a8402115449d',
    'd9c8df28-bb71-4be4-8702-ce2bea8bd943',
    0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

-- better-test-user@mail-client.org -> vend-deut / aquarius
INSERT INTO galactic_sovereign_schema.player ("id", "api_user", "universe", "name")
  VALUES (
    '2bab9414-7972-4483-a8b4-fdd169d0b073',
    '00b265e6-6638-4b1b-aeac-5898c7307eb8',
    '0ac6c027-11d6-47e6-ab15-514cfac48200',
    'vend-deut'
  );

-- deut-factory
INSERT INTO galactic_sovereign_schema.planet ("id", "player", "name", "created_at", "updated_at")
  VALUES (
    '717ffa52-89bd-42eb-b34d-0f994a032e35',
    '2bab9414-7972-4483-a8b4-fdd169d0b073',
    'deut-factory',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.homeworld ("player", "planet")
  VALUES (
    '2bab9414-7972-4483-a8b4-fdd169d0b073',
    '717ffa52-89bd-42eb-b34d-0f994a032e35'
  );

INSERT INTO galactic_sovereign_schema.planet_resource ("planet", "resource", "amount", "created_at", "updated_at")
  VALUES (
    '717ffa52-89bd-42eb-b34d-0f994a032e35',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    500,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_resource ("planet", "resource", "amount", "created_at", "updated_at")
  VALUES (
    '717ffa52-89bd-42eb-b34d-0f994a032e35',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    499,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

INSERT INTO galactic_sovereign_schema.planet_resource_production ("planet", "building", "resource", "production", "created_at", "updated_at")
  VALUES (
    '717ffa52-89bd-42eb-b34d-0f994a032e35',
    NULL,
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    30,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_resource_production ("planet", "building", "resource", "production", "created_at", "updated_at")
  VALUES (
    '717ffa52-89bd-42eb-b34d-0f994a032e35',
    NULL,
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    15,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

INSERT INTO galactic_sovereign_schema.planet_resource_storage ("planet", "resource", "storage", "created_at", "updated_at")
  VALUES (
    '717ffa52-89bd-42eb-b34d-0f994a032e35',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    10000,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_resource_storage ("planet", "resource", "storage", "created_at", "updated_at")
  VALUES (
    '717ffa52-89bd-42eb-b34d-0f994a032e35',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    9000,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

INSERT INTO galactic_sovereign_schema.planet_building ("planet", "building", "level", "created_at", "updated_at")
  VALUES (
    '717ffa52-89bd-42eb-b34d-0f994a032e35',
    'd176e82d-f2ca-4611-996b-c4804096caef',
    2,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_building ("planet", "building", "level", "created_at", "updated_at")
  VALUES (
    '717ffa52-89bd-42eb-b34d-0f994a032e35',
    '3904d34d-9a7e-47d4-a332-091700e2c5c3',
    3,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_building ("planet", "building", "level", "created_at", "updated_at")
  VALUES (
    '717ffa52-89bd-42eb-b34d-0f994a032e35',
    '22b4c0c3-c8e5-4493-89fc-522fdbb0beee',
    0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_building ("planet", "building", "level", "created_at", "updated_at")
  VALUES (
    '717ffa52-89bd-42eb-b34d-0f994a032e35',
    'd9c8df28-bb71-4be4-8702-ce2bea8bd943',
    0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );


