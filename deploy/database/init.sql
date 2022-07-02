-- when a product state change to available state, users that interest in that product will be notified.
---------------------- marketplace ------------------------------------------
create table user(

);

create table user_permission(
-- admin
-- employee for support
);

create table address(); -- user can have multiple addresses.
create table shipping_method(); -- normal, express, ...


create table store(); -- a user can have a store or multiple store?; each store can have multiple products and can available them.
create table product(
    -- product specification jsonb
    -- off-percent
);
create table product_available_subscription();
create table review(
    -- rate
    -- text
    -- user
    -- ....
);
create table votes(); -- a vote relates only to one review.

create table category(); -- relation product to N category
create table order(); -- order can have multiple products.

create table promotion_code(); -- apply on order, each code can apply on just only one order.

--------------------- support and tracking system ---------------------------
create table ticket();
create table ticket_message();
