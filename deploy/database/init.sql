-- when a product state change to available state, users that interest in that product will be notified.
---------------------- marketplace ------------------------------------------
create table "user" -- login with email and password
(
    id           bigserial primary key,
    email        text   not null unique,
    password     text,
    phone_number text   not null unique,
    first_name   text,
    last_name    text,
    avatar_url   text,
    national_id  text,
    permission   bigint references user_permission (id) on delete set null on update set null,
    created_at   timestamptz
);

create table user_permission
(
    id                           bigserial primary key,
    name                         text unique,
    is_system_admin              bool,
    system_read_access           bool,
    system_edit_access_product   bool,
    system_create_access_product bool,
    system_create_new_store      bool,
    is_employee                  bool
);

create table address
(
); -- user can have multiple addresses.
create table shipping_method
(
); -- normal, express, ...


create table store
(
    id         bigserial primary key,
    name       text,
    avatar_url text,
    owner      bigint references "user" (id),
    creator    bigint references "user" (id),
    created_at timestamptz
); -- a user can have a store or multiple store?; each store can have multiple products and can available them.

create table product
(
    -- product specification jsonb
    -- off-percent
);
create table product_available_subscription
(
);
create table review
(
    -- rate
    -- text
    -- user
    -- ....
);
create table votes
(
); -- a vote relates only to one review.

create table category
(
); -- relation product to N category
create table order
(
); -- order can have multiple products.

create table promotion_code
(
);
-- apply on order, each code can apply on just only one order.

--------------------- support and tracking system ---------------------------
create table ticket
(
);
create table ticket_message
(
);
