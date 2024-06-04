CREATE TABLE "notify_settings" (
    "id" bigserial PRIMARY KEY,
    "user_id" bigint references "users"("id"),
    "days_before_notify" int NOT NULL DEFAULT 1,
    UNIQUE ("user_id")
);
