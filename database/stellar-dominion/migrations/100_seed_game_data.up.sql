
-- Resources
INSERT INTO stellar_dominion_schema.resource ("id", "name", "start_amount")
  VALUES ('b4419b6b-b3bf-4576-aa92-055283addbc8', 'metal', 500);

INSERT INTO stellar_dominion_schema.resource ("id", "name", "start_amount")
  VALUES ('cd2ac9aa-9968-4ff5-b746-88f1f810fbb3', 'crystal', 500);

-- Buildings
INSERT INTO stellar_dominion_schema.building ("id", "name")
  VALUES ('d176e82d-f2ca-4611-996b-c4804096caef', 'metal mine');

INSERT INTO stellar_dominion_schema.building_cost ("building", "resource", "cost")
  VALUES (
    'd176e82d-f2ca-4611-996b-c4804096caef',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    60
  );

INSERT INTO stellar_dominion_schema.building_cost ("building", "resource", "cost")
  VALUES (
    'd176e82d-f2ca-4611-996b-c4804096caef',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    15
  );


INSERT INTO stellar_dominion_schema.building ("id", "name")
  VALUES ('3904d34d-9a7e-47d4-a332-091700e2c5c3', 'crystal mine');

INSERT INTO stellar_dominion_schema.building_cost ("building", "resource", "cost")
  VALUES (
    '3904d34d-9a7e-47d4-a332-091700e2c5c3',
    'b4419b6b-b3bf-4576-aa92-055283addbc8',
    48
  );

INSERT INTO stellar_dominion_schema.building_cost ("building", "resource", "cost")
  VALUES (
    '3904d34d-9a7e-47d4-a332-091700e2c5c3',
    'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3',
    24
  );
