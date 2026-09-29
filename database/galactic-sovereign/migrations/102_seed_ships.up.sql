
-- small cargo ship
INSERT INTO galactic_sovereign_schema.ship ("id", "name")
  VALUES('c0978950-601e-4d35-9c7c-28df69d2cd0e', 'small cargo ship');

INSERT INTO galactic_sovereign_schema.ship_cost ("ship", "resource", "cost")
  VALUES('c0978950-601e-4d35-9c7c-28df69d2cd0e', 'b4419b6b-b3bf-4576-aa92-055283addbc8', 2000);
INSERT INTO galactic_sovereign_schema.ship_cost ("ship", "resource", "cost")
  VALUES('c0978950-601e-4d35-9c7c-28df69d2cd0e', 'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3', 2000);

INSERT INTO galactic_sovereign_schema.ship_building_requirement("ship", "building", "level")
  VALUES('c0978950-601e-4d35-9c7c-28df69d2cd0e', '58d75842-6dc0-4ac0-b36d-55f91b8d060d', 2);

-- light fighter
INSERT INTO galactic_sovereign_schema.ship ("id", "name")
  VALUES('a31de13b-5905-4468-99c5-d1d1e529b36e', 'light fighter');

INSERT INTO galactic_sovereign_schema.ship_cost ("ship", "resource", "cost")
  VALUES('a31de13b-5905-4468-99c5-d1d1e529b36e', 'b4419b6b-b3bf-4576-aa92-055283addbc8', 3000);
INSERT INTO galactic_sovereign_schema.ship_cost ("ship", "resource", "cost")
  VALUES('a31de13b-5905-4468-99c5-d1d1e529b36e', 'cd2ac9aa-9968-4ff5-b746-88f1f810fbb3', 1000);

INSERT INTO galactic_sovereign_schema.ship_building_requirement("ship", "building", "level")
  VALUES('a31de13b-5905-4468-99c5-d1d1e529b36e', '58d75842-6dc0-4ac0-b36d-55f91b8d060d', 1);
