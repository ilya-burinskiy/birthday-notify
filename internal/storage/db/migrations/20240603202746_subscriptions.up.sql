CREATE TABLE "subscriptions" (
    "id" bigserial PRIMARY KEY,
    "subscribed_user_id" bigint references "users"("id") NOT NULL,
    "subscribing_user_id" bigint references "users"("id") NOT NULL,
    UNIQUE ("subscribed_user_id", "subscribing_user_id")
);
