-- when a product state change to available state, users that interest in that product will be notified.
---------------------- marketplace ------------------------------------------
drop table if exists user_permission cascade ;
create table user_permission
(
    id                           bigserial primary key,
    name                         text unique,
    is_system_admin              bool,
    system_read_access           bool,
    system_edit_access_product   bool,
    system_create_access_product bool,
    system_create_new_store      bool,
    is_employee                  bool,
    created_at                   timestamptz default now()
);

drop table if exists "user" cascade ;
create table "user" -- login with email and password
(
    id           bigserial primary key,
    email        text not null unique,
    password     text not null,
    phone_number text not null unique,
    first_name   text,
    last_name    text,
    avatar_url   text,
    national_id  text,
    permission   bigint references user_permission (id) on delete set null on update cascade,
    created_at   timestamptz default now()
);

drop table if exists address;
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

drop table if exists shipping_method;
create table shipping_method
(
    id                            serial primary key,
    name                          text not null,
    expected_arrival_working_days int,
    base_cost                     int  not null,
    created_at                    timestamptz default now()
);


drop table if exists category;
create table category
(
    id         bigserial primary key,
    name       text unique not null,
    parent     bigint references category (id),
    created_at timestamptz default now()
);


drop table if exists store cascade ;
create table store
(
    id          bigserial primary key,
    name        text,
    description text,
    avatar_url  text,
    owner       bigint references "user" (id),
    creator     bigint references "user" (id),
    created_at  timestamptz default now()
); -- a user can have a store or multiple store?; each store can have multiple products and can available them.

drop table if exists store_address;
create table store_address
(
    store_id    bigint references "store" (id),
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


drop table if exists product;
create table product
(
    id            bigserial primary key,
    category_id   bigint references category (id),
    name          text not null,
    brand         text,
    description   text,
    picture_url   text,
    specification jsonb
);

drop table if exists warranty;
create table warranty
(
    id         bigserial primary key,
    name       text not null,
    type       text not null,
    month      int,
    created_at timestamptz default now()
);
drop table if exists store_product;
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

    constraint unique (product_id, store_id)
);

drop table if exists product_category;
create table product_category
(
    category_id bigint references category (id) on delete set null on update cascade,
    product_id  bigint references product (id) on delete set null on update cascade,
    constraint unique (category_id, product_id)
);

drop table if exists product_available_subscription;
create table product_available_subscription
(
    product_id           bigint not null references product (id) on delete cascade on update cascade,
    user_id              bigint not null references "user" (id) on delete cascade on update cascade,
    created_at           timestamptz     default now(),
    is_notification_sent bool   not null default false,
    constraint unique (product_id, user_id)
);
drop table if exists review;
create table review
(
    id          bigserial primary key,
    product_id  bigint           not null references product (id) on delete cascade on update cascade,
    store_id    bigint           not null references store (id) on delete cascade on update cascade,
    user_id     bigint           not null references "user" (id) on delete cascade on update cascade,
    rate        double precision not null check ( value between 0 and 5),
    review_text text,
    created_at  timestamptz default now()
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

drop table if exists promotion_code;
create table promotion_code
(
    id            text primary key,
    percent       double precision default 0.0,
    max_off_price double precision default 0.0,
    created_at    timestamptz      default now()
);

drop table if exists "order";
create table "order"
(
    order_id               bigserial primary key,
    status                 text   not null default 'confirmed' check ( value in ('confirmed', 'is-packing', 'packed', 'shipped')),
    tracking_code          text            default 'not-set',
    user_id                bigint not null references "user" (id),
    address_id             bigint not null references address (id) on update cascade,
    store_product_id       bigint not null references store_product on update cascade,
    shipping_method_id     int    not null references shipping_method (id),
    applied_promotion_code text references promotion_code (id),
    is_paid                bool   not null default false,
    created_at             timestamptz     default now()
);

--------------------- support and tracking system ---------------------------
drop table if exists ticket_type;
create table ticket_type
(
    id              bigserial primary key,
    name            text unique,
    description     text,
    is_last_version bool        default true,
    created_at      timestamptz default now()
);
drop table if exists ticket;
create table ticket
(
    id             bigserial primary key,
    user_id        bigint                 not null references "user" (id) on update cascade,
    employee_id    bigint      default -1 not null references "user" (id) on update cascade,
    ticket_type_id bigint                 not null references ticket_type (id) on update cascade,
    is_done        bool,
    done_at        timestamptz,
    created_at     timestamptz default now()

);
drop table if exists ticket_message;
create table ticket_message
(
    ticket_id    bigint not null references ticket (id) on update cascade,
    sender_id    bigint not null references "user" (id),
    message_text text,
    status       text check ( value in ('sent', 'received', 'seen')),
    created_at   timestamptz default now()
);
