
-- user1 -> haha / oberon
INSERT INTO galactic_sovereign_schema.player("id", "api_user", "universe", "name")
  VALUES (
    '92a686c0-9a0a-4bc3-aa1b-9a57ed7f09d5',
    '0463ed3d-bfc9-4c10-b6ee-c223bbca0fab',
    '9682f17b-f5f0-4eda-a747-2537d2151837',
    'haha'
  );

-- planet homeworld
INSERT INTO galactic_sovereign_schema.planet("id", "player", "name", "fields", "created_at", "updated_at")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    '92a686c0-9a0a-4bc3-aa1b-9a57ed7f09d5',
    'homeworld',
    163,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.homeworld("player", "planet")
  VALUES (
    '92a686c0-9a0a-4bc3-aa1b-9a57ed7f09d5',
    '167bd268-6ae7-4cf4-a359-9534beabfeff'
  );
INSERT INTO galactic_sovereign_schema.planet_coordinate("planet", "universe", "galaxy", "solar_system", "position")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    '9682f17b-f5f0-4eda-a747-2537d2151837',
    1, 1, 1
  );

INSERT INTO galactic_sovereign_schema.planet_resource("planet", "resource", "amount")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    1514
  );
INSERT INTO galactic_sovereign_schema.planet_resource("planet", "resource", "amount")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    517
  );

INSERT INTO galactic_sovereign_schema.planet_resource_production("planet", "building", "resource", "production")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    NULL,
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    30
  );
INSERT INTO galactic_sovereign_schema.planet_resource_production("planet", "building", "resource", "production")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    NULL,
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    15
  );

INSERT INTO galactic_sovereign_schema.planet_resource_storage("planet", "resource", "storage")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    10000
  );
INSERT INTO galactic_sovereign_schema.planet_resource_storage("planet", "resource", "storage")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    9000
  );
INSERT INTO galactic_sovereign_schema.planet_resource_storage("planet", "resource", "storage")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    '9665303f-d37f-41e3-ad12-70f8ba8edd14',
    10000
  );

INSERT INTO galactic_sovereign_schema.planet_building("planet", "building", "level")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    'd176e82d-f2ca-4611-996b-c4804096caef',
    0
  );
INSERT INTO galactic_sovereign_schema.planet_building("planet", "building", "level")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    '3904d34d-9a7e-47d4-a332-091700e2c5c3',
    1
  );
INSERT INTO galactic_sovereign_schema.planet_building("planet", "building", "level")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    '22b4c0c3-c8e5-4493-89fc-522fdbb0beee',
    0
  );
INSERT INTO galactic_sovereign_schema.planet_building("planet", "building", "level")
  VALUES (
    '167bd268-6ae7-4cf4-a359-9534beabfeff',
    'd9c8df28-bb71-4be4-8702-ce2bea8bd943',
    0
  );

-- planet colony
INSERT INTO galactic_sovereign_schema.planet("id", "player", "name", "fields", "created_at", "updated_at")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    '92a686c0-9a0a-4bc3-aa1b-9a57ed7f09d5',
    'colony',
    95,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );
INSERT INTO galactic_sovereign_schema.planet_coordinate("planet", "universe", "galaxy", "solar_system", "position")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    '9682f17b-f5f0-4eda-a747-2537d2151837',
    1, 2, 3
  );

INSERT INTO galactic_sovereign_schema.planet_resource("planet", "resource", "amount")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    500
  );
INSERT INTO galactic_sovereign_schema.planet_resource("planet", "resource", "amount")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    500
  );

INSERT INTO galactic_sovereign_schema.planet_resource_production("planet", "building", "resource", "production")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    NULL,
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    30
  );
INSERT INTO galactic_sovereign_schema.planet_resource_production("planet", "building", "resource", "production")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    NULL,
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    15
  );

INSERT INTO galactic_sovereign_schema.planet_resource_storage("planet", "resource", "storage")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    10000
  );
INSERT INTO galactic_sovereign_schema.planet_resource_storage("planet", "resource", "storage")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    9000
  );
INSERT INTO galactic_sovereign_schema.planet_resource_storage("planet", "resource", "storage")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    '9665303f-d37f-41e3-ad12-70f8ba8edd14',
    10000
  );

INSERT INTO galactic_sovereign_schema.planet_building("planet", "building", "level")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    'd176e82d-f2ca-4611-996b-c4804096caef',
    0
  );
INSERT INTO galactic_sovereign_schema.planet_building("planet", "building", "level")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    '3904d34d-9a7e-47d4-a332-091700e2c5c3',
    0
  );
INSERT INTO galactic_sovereign_schema.planet_building("planet", "building", "level")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    '22b4c0c3-c8e5-4493-89fc-522fdbb0beee',
    0
  );
INSERT INTO galactic_sovereign_schema.planet_building("planet", "building", "level")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    'd9c8df28-bb71-4be4-8702-ce2bea8bd943',
    0
  );
