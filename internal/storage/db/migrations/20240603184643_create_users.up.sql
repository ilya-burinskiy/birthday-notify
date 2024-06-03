CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "email" varchar(256) UNIQUE NOT NULL,
    "encrypted_password" bytea NOT NULL,
    "birthdate" date NOT NULL
);
