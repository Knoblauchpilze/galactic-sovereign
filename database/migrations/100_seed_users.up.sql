
-- test-user@provider.com
INSERT INTO user_service_schema.api_user ("id", "email", "password")
  VALUES (
    '0463ed3d-bfc9-4c10-b6ee-c223bbca0fab',
    'test-user@provider.com',
    'strong-password'
  );

INSERT INTO user_service_schema.api_keys ("id", "key", "api_user")
  VALUES (
    'a5eff7a9-9bd6-4f51-9b42-a7ca5ffd3f5e',
    '3e8d49a3-9220-4ea0-88eb-299520c6ab85',
    '0463ed3d-bfc9-4c10-b6ee-c223bbca0fab'
  );

-- another-test-user@another-provider.com
INSERT INTO user_service_schema.api_user ("id", "email", "password")
  VALUES (
    '4f26321f-d0ea-46a3-83dd-6aa1c6053aaf',
    'another-test-user@another-provider.com',
    'super-strong-password'
  );

INSERT INTO user_service_schema.api_keys ("id", "key", "api_user")
  VALUES (
    'fd8136c4-c584-4bbf-a390-53d5c2548fb8',
    '2da3e9ec-7299-473a-be0f-d722d870f51a',
    '4f26321f-d0ea-46a3-83dd-6aa1c6053aaf'
  );
INSERT INTO user_service_schema.api_keys ("id", "key", "api_user")
  VALUES (
    '2e791bfe-7e35-465d-8269-1bbd7b4e86c5',
    '7ba9182e-f4a6-4eec-b216-d2a9f5179fc2',
    '4f26321f-d0ea-46a3-83dd-6aa1c6053aaf'
  );

-- better-test-user@mail-client.org
INSERT INTO user_service_schema.api_user ("id", "email", "password")
  VALUES (
    '00b265e6-6638-4b1b-aeac-5898c7307eb8',
    'better-test-user@mail-client.org',
    'weakpassword'
  );

INSERT INTO user_service_schema.api_keys ("id", "key", "api_user")
  VALUES (
    '42698272-5b8f-42db-a43c-8108eaad66e1',
    'e9c3ce0d-d6d6-45cb-ad93-c407d429469f',
    '00b265e6-6638-4b1b-aeac-5898c7307eb8'
  );
INSERT INTO user_service_schema.api_keys ("id", "key", "api_user", "enabled")
  VALUES (
    'f79d5502-4c57-41de-8237-893c6e1983f0',
    '586e86f0-b981-4d16-90ee-99114ae36d19',
    '00b265e6-6638-4b1b-aeac-5898c7307eb8',
    FALSE
  );
INSERT INTO user_service_schema.api_keys ("id", "key", "api_user", "enabled")
  VALUES (
    '7087e425-7c3e-40b3-9736-b2f77fefc0fb',
    '627a4de2-263e-4a2d-a24f-07bea157aabd',
    '00b265e6-6638-4b1b-aeac-5898c7307eb8',
    FALSE
  );
