
-- Metal
INSERT INTO galactic_sovereign_schema.resource("id", "name", "start_amount", "start_production", "start_storage")
  VALUES ('b4419b6b-b3bf-4576-aa92-055283addbc8', 'metal', 500, 30, 10000);

INSERT INTO galactic_sovereign_schema.resource_metabolization_rate_building("resource", "hours_per_unit")
  VALUES ('b4419b6b-b3bf-4576-aa92-055283addbc8', 1.0/2500.0);
INSERT INTO galactic_sovereign_schema.resource_metabolization_rate_shipyard("resource", "hours_per_unit")
  VALUES ('b4419b6b-b3bf-4576-aa92-055283addbc8', 1.0/2500.0);

-- Crystal
INSERT INTO galactic_sovereign_schema.resource("id", "name", "start_amount", "start_production", "start_storage")
  VALUES ('cd2ac9aa-9968-4ff5-b746-88f1f810fbb3', 'crystal', 500, 15, 10000);

INSERT INTO galactic_sovereign_schema.resource_metabolization_rate_building("resource", "hours_per_unit")
  VALUES ('cd2ac9aa-9968-4ff5-b746-88f1f810fbb3', 1.0/2500.0);
INSERT INTO galactic_sovereign_schema.resource_metabolization_rate_shipyard("resource", "hours_per_unit")
  VALUES ('cd2ac9aa-9968-4ff5-b746-88f1f810fbb3', 1.0/2500.0);

-- Deuterium
INSERT INTO galactic_sovereign_schema.resource("id", "name", "start_amount", "start_production", "start_storage")
  VALUES ('9665303f-d37f-41e3-ad12-70f8ba8edd14', 'deuterium', 0, 0, 10000);

INSERT INTO galactic_sovereign_schema.resource_metabolization_rate_building("resource", "hours_per_unit")
  VALUES ('9665303f-d37f-41e3-ad12-70f8ba8edd14', 0);
INSERT INTO galactic_sovereign_schema.resource_metabolization_rate_shipyard("resource", "hours_per_unit")
  VALUES ('9665303f-d37f-41e3-ad12-70f8ba8edd14', 0);
