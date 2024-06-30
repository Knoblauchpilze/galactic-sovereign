
-- test-user@provider.com
INSERT INTO user_service_schema.acl ("id", "api_user", "resource")
  VALUES (
    'f3e3b687-6033-4cc2-8723-3e7a28f74c52',
    '0463ed3d-bfc9-4c10-b6ee-c223bbca0fab',
    'v1/users'
  );

INSERT INTO user_service_schema.acl_permission ("acl", "permission")
  VALUES ('f3e3b687-6033-4cc2-8723-3e7a28f74c52', 'GET');
INSERT INTO user_service_schema.acl_permission ("acl", "permission")
  VALUES ('f3e3b687-6033-4cc2-8723-3e7a28f74c52', 'POST');
INSERT INTO user_service_schema.acl_permission ("acl", "permission")
  VALUES ('f3e3b687-6033-4cc2-8723-3e7a28f74c52', 'DELETE');

INSERT INTO user_service_schema.user_limit ("id", "name", "api_user")
  VALUES (
    'bb81dd7b-ef98-4bc7-9681-77d5403ecab5',
    'api-usage',
    '0463ed3d-bfc9-4c10-b6ee-c223bbca0fab'
  );

INSERT INTO user_service_schema.limits ("id", "name", "value", "user_limit")
  VALUES (
    '0d1acf04-b806-4573-9dcd-3dbd43a12cea',
    'group',
    'admin',
    'bb81dd7b-ef98-4bc7-9681-77d5403ecab5'
  );

-- another-test-user@another-provider.com
INSERT INTO user_service_schema.acl ("id", "api_user", "resource")
  VALUES (
    '8382acc1-f3d8-4ed0-be8b-0ac57ce3fc8b',
    '4f26321f-d0ea-46a3-83dd-6aa1c6053aaf',
    'v1/users'
  );

INSERT INTO user_service_schema.acl_permission ("acl", "permission")
  VALUES ('8382acc1-f3d8-4ed0-be8b-0ac57ce3fc8b', 'GET');
INSERT INTO user_service_schema.acl_permission ("acl", "permission")
  VALUES ('8382acc1-f3d8-4ed0-be8b-0ac57ce3fc8b', 'POST');

-- better-test-user@mail-client.org
INSERT INTO user_service_schema.user_limit ("id", "name", "api_user")
  VALUES (
    '5f0e2a17-8d14-41ef-9550-c0d490917818',
    'api-usage',
    '00b265e6-6638-4b1b-aeac-5898c7307eb8'
  );

INSERT INTO user_service_schema.limits ("id", "name", "value", "user_limit")
  VALUES (
    'd1686d58-6dd3-4517-8831-b70dd1d9eeaa',
    'group',
    'admin',
    '5f0e2a17-8d14-41ef-9550-c0d490917818'
  );
INSERT INTO user_service_schema.limits ("id", "name", "value", "user_limit")
  VALUES (
    'e9fb7ea4-6845-43dc-a358-a24c258025b7',
    'currency',
    'euro',
    '5f0e2a17-8d14-41ef-9550-c0d490917818'
  );
