
-- user1 -> throwaway-account / aquarius
INSERT INTO galactic_sovereign_schema.player ("id", "api_user", "universe", "name")
  VALUES (
    '04a7477c-a66b-4c47-9c17-ac209183c7a4',
    '0463ed3d-bfc9-4c10-b6ee-c223bbca0fab',
    '0ac6c027-11d6-47e6-ab15-514cfac48200',
    'throwaway-account'
  );

-- planet a-new-beginning
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
