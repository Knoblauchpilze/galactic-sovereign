
-- user1
INSERT INTO user_service_schema.api_user ("id", "email", "password")
  VALUES (
    '0463ed3d-bfc9-4c10-b6ee-c223bbca0fab',
    'user1',
    'pwd1'
  );

-- https://www.postgresql.org/docs/current/functions-datetime.html#FUNCTIONS-DATETIME-CURRENT
INSERT INTO user_service_schema.api_key ("id", "key", "api_user", "valid_until")
  VALUES (
    'a5eff7a9-9bd6-4f51-9b42-a7ca5ffd3f5e',
    '3e8d49a3-9220-4ea0-88eb-299520c6ab85',
    '0463ed3d-bfc9-4c10-b6ee-c223bbca0fab',
     current_timestamp + make_interval(hours => 6)
  );

-- another-test-user@another-provider.com
INSERT INTO user_service_schema.api_user ("id", "email", "password")
  VALUES (
    '4f26321f-d0ea-46a3-83dd-6aa1c6053aaf',
    'another-test-user@another-provider.com',
    'super-strong-password'
  );

INSERT INTO user_service_schema.api_key ("id", "key", "api_user", "valid_until")
  VALUES (
    'fd8136c4-c584-4bbf-a390-53d5c2548fb8',
    '2da3e9ec-7299-473a-be0f-d722d870f51a',
    '4f26321f-d0ea-46a3-83dd-6aa1c6053aaf',
     current_timestamp + make_interval(hours => 6)
  );

-- better-test-user@mail-client.org
INSERT INTO user_service_schema.api_user ("id", "email", "password")
  VALUES (
    '00b265e6-6638-4b1b-aeac-5898c7307eb8',
    'better-test-user@mail-client.org',
    'weakpassword'
  );

INSERT INTO user_service_schema.api_key ("id", "key", "api_user", "valid_until")
  VALUES (
    '42698272-5b8f-42db-a43c-8108eaad66e1',
    'e9c3ce0d-d6d6-45cb-ad93-c407d429469f',
    '00b265e6-6638-4b1b-aeac-5898c7307eb8',
     current_timestamp + make_interval(hours => 6)
  );


-- i-dont-care-about-@security.de
INSERT INTO user_service_schema.api_user ("id", "email", "password")
  VALUES (
    'beb2a2dc-2a9f-48d6-b2ca-fd3b5ca3249f',
    'i-dont-care-about-@security.de',
    'mycatismypassword'
  );

INSERT INTO user_service_schema.api_key ("id", "key", "api_user", "valid_until")
  VALUES (
    'a610adcb-d966-4617-9f15-caf6e48b6325',
    'c64f4da4-8bc5-4e19-a038-cd8755bd07d5',
    'beb2a2dc-2a9f-48d6-b2ca-fd3b5ca3249f',
     current_timestamp + make_interval(hours => 6)
  );
