
-- Resources
-- https://ogame.fandom.com/wiki/Buildings
INSERT INTO galactic_sovereign_schema.resource("id", "name", "start_amount", "start_production", "start_storage")
  VALUES ('b4419b6b-b3bf-4576-aa92-055283addbc8', 'metal', 500, 30, 10000);

INSERT INTO galactic_sovereign_schema.resource_metabolization_rate_building("resource", "hours_per_unit")
  VALUES ('b4419b6b-b3bf-4576-aa92-055283addbc8', 1.0/2500.0);
INSERT INTO galactic_sovereign_schema.resource_metabolization_rate_shipyard("resource", "hours_per_unit")
  VALUES ('b4419b6b-b3bf-4576-aa92-055283addbc8', 1.0/2500.0);

INSERT INTO galactic_sovereign_schema.resource("id", "name", "start_amount", "start_production", "start_storage")
  VALUES ('cd2ac9aa-9968-4ff5-b746-88f1f810fbb3', 'crystal', 500, 15, 10000);

INSERT INTO galactic_sovereign_schema.resource_metabolization_rate_building("resource", "hours_per_unit")
  VALUES ('cd2ac9aa-9968-4ff5-b746-88f1f810fbb3', 1.0/2500.0);
INSERT INTO galactic_sovereign_schema.resource_metabolization_rate_shipyard("resource", "hours_per_unit")
  VALUES ('cd2ac9aa-9968-4ff5-b746-88f1f810fbb3', 1.0/2500.0);

INSERT INTO galactic_sovereign_schema.resource("id", "name", "start_amount", "start_production", "start_storage")
  VALUES ('9665303f-d37f-41e3-ad12-70f8ba8edd14', 'deuterium', 0, 0, 10000);

INSERT INTO galactic_sovereign_schema.resource_metabolization_rate_building("resource", "hours_per_unit")
  VALUES ('9665303f-d37f-41e3-ad12-70f8ba8edd14', 0);
INSERT INTO galactic_sovereign_schema.resource_metabolization_rate_shipyard("resource", "hours_per_unit")
  VALUES ('9665303f-d37f-41e3-ad12-70f8ba8edd14', 0);

-- Buildings
-- metal mine
INSERT INTO galactic_sovereign_schema.building("id", "name")
  VALUES ('d176e82d-f2ca-4611-996b-c4804096caef', 'metal mine');

INSERT INTO galactic_sovereign_schema.building_cost("building", "resource", "cost", "progress")
  VALUES (
    'd176e82d-f2ca-4611-996b-c4804096caef',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    60,
    1.5
  );
INSERT INTO galactic_sovereign_schema.building_cost("building", "resource", "cost", "progress")
  VALUES (
    'd176e82d-f2ca-4611-996b-c4804096caef',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    15,
    1.5
  );

INSERT INTO galactic_sovereign_schema.building_resource_production("building", "resource", "base", "progress")
  VALUES (
    'd176e82d-f2ca-4611-996b-c4804096caef',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    30,
    1.1
  );

-- crystal mine
INSERT INTO galactic_sovereign_schema.building("id", "name")
  VALUES ('3904d34d-9a7e-47d4-a332-091700e2c5c3', 'crystal mine');

INSERT INTO galactic_sovereign_schema.building_cost("building", "resource", "cost", "progress")
  VALUES (
    '3904d34d-9a7e-47d4-a332-091700e2c5c3',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    48,
    1.6
  );
INSERT INTO galactic_sovereign_schema.building_cost("building", "resource", "cost", "progress")
  VALUES (
    '3904d34d-9a7e-47d4-a332-091700e2c5c3',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    24,
    1.6
  );

INSERT INTO galactic_sovereign_schema.building_resource_production("building", "resource", "base", "progress")
  VALUES (
    '3904d34d-9a7e-47d4-a332-091700e2c5c3',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    20,
    1.1
  );

-- deuterium synthetizer
INSERT INTO galactic_sovereign_schema.building("id", "name")
  VALUES ('54a0ce97-bf8b-4fae-ba6e-caa9ae96265f', 'deuterium synthetizer');

INSERT INTO galactic_sovereign_schema.building_cost("building", "resource", "cost", "progress")
  VALUES (
    '54a0ce97-bf8b-4fae-ba6e-caa9ae96265f',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    225,
    1.5
  );
INSERT INTO galactic_sovereign_schema.building_cost("building", "resource", "cost", "progress")
  VALUES (
    '54a0ce97-bf8b-4fae-ba6e-caa9ae96265f',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    75,
    1.5
  );

INSERT INTO galactic_sovereign_schema.building_resource_production("building", "resource", "base", "progress")
  VALUES (
    '54a0ce97-bf8b-4fae-ba6e-caa9ae96265f',
    '9665303f-d37f-41e3-ad12-70f8ba8edd14',
    10,
    1.1
  );

-- metal storage
INSERT INTO galactic_sovereign_schema.building("id", "name")
  VALUES ('22b4c0c3-c8e5-4493-89fc-522fdbb0beee', 'metal storage');

INSERT INTO galactic_sovereign_schema.building_cost("building", "resource", "cost", "progress")
  VALUES (
    '22b4c0c3-c8e5-4493-89fc-522fdbb0beee',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    1000,
    2.0
  );

INSERT INTO galactic_sovereign_schema.building_resource_storage("building", "resource", "base", "scale", "progress")
  VALUES (
    '22b4c0c3-c8e5-4493-89fc-522fdbb0beee',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    5000,
    2.5,
    1.833195476
  );

-- crystal storage
INSERT INTO galactic_sovereign_schema.building("id", "name")
  VALUES ('d9c8df28-bb71-4be4-8702-ce2bea8bd943', 'crystal storage');

INSERT INTO galactic_sovereign_schema.building_cost("building", "resource", "cost", "progress")
  VALUES (
    'd9c8df28-bb71-4be4-8702-ce2bea8bd943',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    1000,
    2.0
  );
INSERT INTO galactic_sovereign_schema.building_cost("building", "resource", "cost", "progress")
  VALUES (
    'd9c8df28-bb71-4be4-8702-ce2bea8bd943',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    500,
    2.0
  );

INSERT INTO galactic_sovereign_schema.building_resource_storage("building", "resource", "base", "scale", "progress")
  VALUES (
    'd9c8df28-bb71-4be4-8702-ce2bea8bd943',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    5000,
    2.5,
    1.833195476
  );

-- deuterium tank
INSERT INTO galactic_sovereign_schema.building("id", "name")
  VALUES ('6b81a99f-d826-475b-8dd5-d066b501b1df', 'deuterium tank');

INSERT INTO galactic_sovereign_schema.building_cost("building", "resource", "cost", "progress")
  VALUES (
    '6b81a99f-d826-475b-8dd5-d066b501b1df',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    1000,
    2.0
  );
INSERT INTO galactic_sovereign_schema.building_cost("building", "resource", "cost", "progress")
  VALUES (
    '6b81a99f-d826-475b-8dd5-d066b501b1df',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    1000,
    2.0
  );

INSERT INTO galactic_sovereign_schema.building_resource_storage("building", "resource", "base", "scale", "progress")
  VALUES (
    '6b81a99f-d826-475b-8dd5-d066b501b1df',
    '9665303f-d37f-41e3-ad12-70f8ba8edd14',
    5000,
    2.5,
    1.833195476
  );

-- shipyard
INSERT INTO galactic_sovereign_schema.building("id", "name")
  VALUES ('58d75842-6dc0-4ac0-b36d-55f91b8d060d', 'shipyard');

INSERT INTO galactic_sovereign_schema.building_cost("building", "resource", "cost", "progress")
  VALUES (
    '58d75842-6dc0-4ac0-b36d-55f91b8d060d',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    400,
    2.0
  );
INSERT INTO galactic_sovereign_schema.building_cost("building", "resource", "cost", "progress")
  VALUES (
    '58d75842-6dc0-4ac0-b36d-55f91b8d060d',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    200,
    2.0
  );
INSERT INTO galactic_sovereign_schema.building_cost("building", "resource", "cost", "progress")
  VALUES (
    '58d75842-6dc0-4ac0-b36d-55f91b8d060d',
    '9665303f-d37f-41e3-ad12-70f8ba8edd14',
    100,
    2.0
  );

-- Ships
-- small cargo ship
INSERT INTO galactic_sovereign_schema.ship ("id", "name")
  VALUES('c0978950-601e-4d35-9c7c-28df69d2cd0e', 'small cargo ship');

INSERT INTO galactic_sovereign_schema.ship_cost ("ship", "resource", "cost")
  VALUES('c0978950-601e-4d35-9c7c-28df69d2cd0e', 'b4419b6b-b3bf-4576-aa92-055283addbc8', 2000);
INSERT INTO galactic_sovereign_schema.ship_cost ("ship", "resource", "cost")
  VALUES('c0978950-601e-4d35-9c7c-28df69d2cd0e', 'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3', 2000);

-- light fighter
INSERT INTO galactic_sovereign_schema.ship ("id", "name")
  VALUES('a31de13b-5905-4468-99c5-d1d1e529b36e', 'light fighter');

INSERT INTO galactic_sovereign_schema.ship_cost ("ship", "resource", "cost")
  VALUES('a31de13b-5905-4468-99c5-d1d1e529b36e', 'b4419b6b-b3bf-4576-aa92-055283addbc8', 3000);
INSERT INTO galactic_sovereign_schema.ship_cost ("ship", "resource", "cost")
  VALUES('a31de13b-5905-4468-99c5-d1d1e529b36e', 'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3', 1000);
