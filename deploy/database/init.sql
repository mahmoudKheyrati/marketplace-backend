-- when a product state change to available state, users that interest in that product will be notified.
---------------------- marketplace ------------------------------------------
drop table if exists user_permission cascade;
create table user_permission
(
    name                         text primary key,
    is_system_admin              bool        default false,
    system_read_access           bool        default false,
    system_edit_access_product   bool        default false,
    system_create_access_product bool        default false,
    system_create_new_store      bool        default false,
    store_create_new_product     bool        default false,
    store_read_new_product       bool        default false,
    store_edit_new_product       bool        default false,
    is_employee                  bool        default false,
    created_at                   timestamptz default now()
);

insert into user_permission(name, is_system_admin, system_read_access, system_edit_access_product,
                            system_create_access_product, system_create_new_store, store_create_new_product,
                            store_read_new_product, store_edit_new_product, is_employee)
values ('normal-user', false, false, false, false, false, false, false, false, false),
       ('store-admin', false, false, false, false, false, true, true, true, false),
       ('marketplace-admin', true, true, true, true, true, false, false, false, true);



drop table if exists "user" cascade;
create table "user" -- login with email and password
(
    id              bigserial primary key,
    email           text not null unique,
    password        text not null,
    phone_number    text not null unique,
    first_name      text        default '',
    last_name       text        default '',
    avatar_url      text        default '',
    national_id     text        default '',
    permission_name text references user_permission (name) on delete set null on update cascade,
    created_at      timestamptz default now()
);

drop table if exists address cascade;
create table address
(
    id                bigserial primary key,
    user_id           bigint references "user" (id) on delete cascade on update cascade,
    country           text,
    province          text,
    city              text,
    street            text,
    postal_code       text,
    home_phone_number text,
    is_last_version   bool        default true,
    created_at        timestamptz default now()
);

drop table if exists shipping_method cascade;
create table shipping_method
(
    name                          text primary key,
    expected_arrival_working_days int,
    base_cost                     int not null,
    created_at                    timestamptz default now()
);
-- shipping method
insert into shipping_method (name, expected_arrival_working_days, base_cost)
values ('ordinary', 7, 15),
       ('express', 3, 35),
       ('special', 1, 50);


drop table if exists category cascade;
create table category
(
    id         bigserial primary key,
    name       text not null,
    parent     bigint references category (id),
    created_at timestamptz default now()
);
-- category
insert into category(id, name, parent)
values (1, 'electronics', null),
       (2, 'laptop', 1),
       (3, 'phone', 1),
       (4, 'monitor', 1),
       (5, 'shoes', null),
       (6, 'men-shoes', 5),
       (7, 'boot', 6),
       (8, 'sneakers', 6),
       (9, 'women-shoes', 5),
       (10, 'boot', 9),
       (11, 'sneakers', 9),
       (12, 'flat', 9);


drop table if exists store cascade;
create table store
(
    id          bigserial primary key,
    name        text,
    description text,
    avatar_url  text,
    owner       bigint references "user" (id),
    creator     bigint references "user" (id),
    created_at  timestamptz default now(),
    deleted_at  timestamptz
); -- a user can have a store or multiple store?; each store can have multiple products and can available them.

drop table if exists store_address;
create table store_address
(
    id          bigserial primary key,
    store_id    bigint references "store" (id) on update cascade on delete set null,
    country     text,
    province    text,
    city        text,
    street      text,
    postal_code text,
    created_at  timestamptz default now()
);

drop table if exists store_category;
create table store_category
(
    category_id bigint references category (id) on delete set null on update cascade,
    store_id    bigint references store (id) on delete set null on update cascade,
    unique (category_id, store_id)
);


drop table if exists product cascade;
create table product
(
    id            bigserial primary key,
    category_id   bigint references category (id),
    name          text not null,
    brand         text,
    description   text,
    picture_url   text,
    specification jsonb,
    created_at    timestamptz default now()

);

drop table if exists warranty cascade;
create table warranty
(
    id         bigserial primary key,
    name       text not null,
    type       text not null,
    month      int,
    created_at timestamptz default now()
);
drop table if exists store_product cascade;
create table store_product
(
    product_id      bigint           not null references product (id),
    store_id        bigint           not null references store (id),
    off_percent     double precision          default 0.0,
    max_off_price   double precision          default 0.0,
    price           double precision not null default 0.0,
    available_count int              not null default 0,
    warranty_id     bigint references warranty (id) on delete set null on update cascade,
    created_at      timestamptz               default now(),
    is_last_version bool                      default true,
    unique (product_id, store_id)
);

-- drop table if exists product_category;
-- create table product_category
-- (
--     category_id bigint references category (id) on delete set null on update cascade,
--     product_id  bigint references product (id) on delete set null on update cascade,
--     unique (category_id, product_id)
-- );

drop table if exists product_available_subscription;
create table product_available_subscription
(
    id                   bigserial primary key,
    product_id           bigint not null references product (id) on delete cascade on update cascade,
    user_id              bigint not null references "user" (id) on delete cascade on update cascade,
    created_at           timestamptz     default now(),
    is_notification_sent bool   not null default false,
    available_status     bool   not null default false
);


drop table if exists review cascade;
create table review
(
    id          bigserial primary key,
    product_id  bigint           not null references product (id) on delete cascade on update cascade,
    store_id    bigint           not null references store (id) on delete cascade on update cascade,
    user_id     bigint           not null references "user" (id) on delete cascade on update cascade,
    rate        double precision not null check ( rate between 0 and 5),
    review_text text,
    created_at  timestamptz default now(),
    deleted_at  timestamptz
);

drop table if exists votes;
create table votes
(
    review_id bigint not null references review (id) on delete cascade on update cascade,
    user_id   bigint not null references "user" (id) on delete cascade on update cascade,
    up_vote   bool,
    down_vote bool,
    check ( up_vote != votes.down_vote
        )
);

drop table if exists promotion_code cascade;
create table promotion_code
(
    id            text primary key,
    percent       double precision default 0.0,
    max_off_price double precision default 0.0,
    created_at    timestamptz      default now(),
    deleted_at    timestamptz
);

drop table if exists "order" cascade;
create table "order"
(
    order_id               bigserial primary key,
    status                 text   not null default 'ongoing' check ( status in
                                                                     ('ongoing', 'cancelled-by-user', 'confirmed',
                                                                      'is-packing', 'packed', 'shipped')),
    tracking_code          text            default 'not-set',
    user_id                bigint not null references "user" (id),
    address_id             bigint references address (id) on update cascade,
    shipping_method_id     text references shipping_method (name),
    applied_promotion_code text references promotion_code (id),
    is_paid                bool   not null default false,
    pay_date               timestamptz,
    created_at             timestamptz     default now()
);

drop table if exists product_order;
create table product_order
(
    product_id bigint not null references product (id),
    store_id   bigint not null references store (id),
    order_id   bigint references "order" (order_id),
    quantity   bigint      default 1,
    created_at timestamptz default now(),
    unique (product_id, store_id, order_id)
);

drop table if exists payment;
create table payment
(
    order_id    bigint,
    user_id     bigint,
    total_price bigint,
    pay_at      timestamptz default now(),
    unique (order_id, user_id)
);

--------------------- support and tracking system ---------------------------
drop table if exists ticket_type cascade;
create table ticket_type
(
    id              bigserial primary key,
    name            text,
    description     text,
    is_last_version bool        default true,
    created_at      timestamptz default now()
);
drop table if exists ticket cascade;
create table ticket
(
    id             bigserial primary key,
    user_id        bigint                 not null references "user" (id) on update cascade,
    employee_id    bigint      default -1 not null references "user" (id) on update cascade,
    ticket_type_id bigint                 not null references ticket_type (id) on update cascade,
    is_done        bool        default false,
    done_at        timestamptz,
    created_at     timestamptz default now()

);
drop table if exists ticket_message cascade;
create table ticket_message
(
    id           bigserial primary key,
    ticket_id    bigint not null references ticket (id) on update cascade,
    sender_id    bigint not null references "user" (id),
    message_text text   not null,
    status       text        default 'sent' check ( status in ('sent', 'received', 'seen')),
    created_at   timestamptz default now()
);
