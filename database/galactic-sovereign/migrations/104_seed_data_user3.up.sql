
-- another-test-user@another-provider.com -> very-nice-pseudo / oberon
INSERT INTO galactic_sovereign_schema.player ("id", "api_user", "universe", "name")
  VALUES (
    'e8db2006-3e35-49cd-8e1f-726491660a00',
    '4f26321f-d0ea-46a3-83dd-6aa1c6053aaf',
    '9682f17b-f5f0-4eda-a747-2537d2151837',
    'very-nice-pseudo'
  );

-- planet my-awesome-planet
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
