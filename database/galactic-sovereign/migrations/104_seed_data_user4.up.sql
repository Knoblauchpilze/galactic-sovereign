
-- better-test-user@mail-client.org -> vend-deut / aquarius
INSERT INTO galactic_sovereign_schema.player ("id", "api_user", "universe", "name")
  VALUES (
    '2bab9414-7972-4483-a8b4-fdd169d0b073',
    '00b265e6-6638-4b1b-aeac-5898c7307eb8',
    '0ac6c027-11d6-47e6-ab15-514cfac48200',
    'vend-deut'
  );

-- planet deut-factory
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
