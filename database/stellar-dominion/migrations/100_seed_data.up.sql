
-- Universes
INSERT INTO stellar_dominion_schema.universe ("id", "name")
  VALUES ('9682f17b-f5f0-4eda-a747-2537d2151837', 'oberon');

INSERT INTO stellar_dominion_schema.universe ("id", "name")
  VALUES ('0ac6c027-11d6-47e6-ab15-514cfac48200', 'aquarius');

-- test-user@provider.com
INSERT INTO stellar_dominion_schema.player ("id", "api_user", "universe", "name")
  VALUES (
    '92a686c0-9a0a-4bc3-aa1b-9a57ed7f09d5',
    '0463ed3d-bfc9-4c10-b6ee-c223bbca0fab',
    '9682f17b-f5f0-4eda-a747-2537d2151837',
    'the-great-test-user-overlord'
  );

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

INSERT INTO stellar_dominion_schema.planet ("id", "player", "name")
  VALUES (
    '110cdf6f-2103-4e34-924f-fd57eb87ea3e',
    '92a686c0-9a0a-4bc3-aa1b-9a57ed7f09d5',
    'colony'
  );

INSERT INTO stellar_dominion_schema.player ("id", "api_user", "universe", "name")
  VALUES (
    '04a7477c-a66b-4c47-9c17-ac209183c7a4',
    '0463ed3d-bfc9-4c10-b6ee-c223bbca0fab',
    '0ac6c027-11d6-47e6-ab15-514cfac48200',
    'throwaway-account'
  );

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

-- another-test-user@another-provider.com
INSERT INTO stellar_dominion_schema.player ("id", "api_user", "universe", "name")
  VALUES (
    'e8db2006-3e35-49cd-8e1f-726491660a00',
    '4f26321f-d0ea-46a3-83dd-6aa1c6053aaf',
    '9682f17b-f5f0-4eda-a747-2537d2151837',
    'very-nice-pseudo'
  );

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

-- better-test-user@mail-client.org
INSERT INTO stellar_dominion_schema.player ("id", "api_user", "universe", "name")
  VALUES (
    '2bab9414-7972-4483-a8b4-fdd169d0b073',
    '00b265e6-6638-4b1b-aeac-5898c7307eb8',
    '0ac6c027-11d6-47e6-ab15-514cfac48200',
    'vend-deut'
  );

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


