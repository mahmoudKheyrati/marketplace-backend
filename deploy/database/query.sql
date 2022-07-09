--  login
select id,
       email,
       password,
       phone_number,
       first_name,
       last_name,
       avatar_url,
       national_id,
       permission_name,
       created_at
from "user"
where email = ?
  and password = ?;
-- signup
insert into "user" (email, password, phone_number, first_name, last_name, permission_name)
values (?, ?, ?, ?, ?, 'normal-user');
-- get user profile
select id,
       email,
       phone_number,
       first_name,
       last_name,
       avatar_url,
       national_id,
       permission_name,
       created_at
from "user" where id = ?;
-- add address for user
insert into address(user_id, country, province, city, street, postal_code, home_phone_number)
values (?, ?, ?, ?, ?, ?, ?);
-- get user addresses
select id,
       user_id,
       country,
       province,
       city,
       street,
       postal_code,
       home_phone_number,
       is_last_version,
       created_at
from address
where user_id = ?
  and is_last_version = true;
-- update user address
begin;

update address
set is_last_version = false
where id = ?;

insert into address(user_id, country, province, city, street, postal_code, home_phone_number)
values (?, ?, ?, ?, ?, ?, ?);
commit;



-- get all shipping method
select name, expected_arrival_working_days, base_cost, created_at
from shipping_method;


-- get main categories
select name
from category
where parent is null;
-- get sub-categories by category id
select name
from category
where parent = ?;
-- get all categories with sub-categories
-- todo: implement later if needed.

-- create store
insert into store(name, description, avatar_url, owner, creator)
values (?, ?, ?, ?, ?);

-- update store
update store
set name= ?,
    description = ?,
    avatar_url = ?,
    owner = ?
where id = ?;

-- get user all stores
select id, name, description, avatar_url, owner, creator, created_at
from store
where creator = ?;

-- get all marketplace stores
select id, name, description, avatar_url, owner, creator, created_at
from store;
-- delete store
delete
from store
where id = ?
  and creator = ?;


-- add address for store

insert into store_address (country, province, city, street, postal_code)
values (?, ?, ?, ?, ?);
-- get store addresses
select store_id, country, province, city, street, postal_code, created_at
from store_address
where store_id = ?;
-- update store address
update store_address
set country     =?,
    province=?,
    city= ?,
    street= ?,
    postal_code = ?
where store_id = ?;

-- insert store category
insert into store_category(category_id, store_id)
values (?, ?);

-- delete store category
delete
from store_category
where category_id =?
  and store_id = ?;

-- add product
insert into product (category_id, name, brand, description, picture_url, specification)
values (?, ?, ?, ?, ?, ?);

-- delete product
delete
from product
where id = ?;

-- select products by name
select id, category_id, name, brand, description, picture_url, specification
from product
where name ilike '%?%';
-- important
-- select products by brand
select id, category_id, name, brand, description, picture_url, specification
from product
where brand ilike '%?%';
-- important
-- select products by category
-- todo: add recursive query to get sub categories
select *
from product
where category_id in (
    select id
    from category
    where category.name = ?
);

-- filter by brand
select id, category_id, name, brand, description, picture_url, specification
from product
where category_id = ?
  and brand in (?);
-- filter by price
select id, category_id, name, brand, description, picture_url, specification
from product p
         join store_product sp on p.id = sp.product_id
where p.category_id = ?
  and sp.price between ? and ?;

-- filter by specification
select id, category_id, name, brand, description, picture_url, specification
from product
where category_id = ?
  and specification::jsonb ->> ? = ?;
-- important

-- get category specifications distinct
select distinct (jsonb_object_keys(p.specification))
from product p
         join category c on p.category_id = c.id
where p.category_id = ?;

-- sort by price cheapest to most expensive
select p.id,
       p.category_id,
       p.name,
       p.brand,
       p.description,
       p.picture_url,
       p.specification,
       min(sp.price) as price
from product p
         join store_product sp on p.id = sp.product_id
where p.category_id = ?
group by p.id, p.category_id, p.name, p.brand, p.description, p.picture_url, p.specification
order by min(sp.price);

-- sort by price most expensive to cheapest
select p.id,
       p.category_id,
       p.name,
       p.brand,
       p.description,
       p.picture_url,
       p.specification,
       min(sp.price) as price
from product p
         join store_product sp on p.id = sp.product_id
where p.category_id = ?
group by p.id, p.category_id, p.name, p.brand, p.description, p.picture_url, p.specification
order by min(sp.price) desc;
-- sort by rating highest to lowest
select p.id,
       p.category_id,
       p.name,
       p.brand,
       p.description,
       p.picture_url,
       p.specification,
       avg(r.rate) as rate
from product p
         join review r on p.id = r.product_id
where p.category_id = ?
group by p.id, p.category_id, p.name, p.brand, p.description, p.picture_url, p.specification
order by avg(r.rate) desc;
-- sort by recently added
select id,
       category_id,
       name,
       brand,
       description,
       picture_url,
       specification
from product
where category_id = ?
order by created_at desc;
-- get product by product_id
select id,
       category_id,
       name,
       brand,
       description,
       picture_url,
       specification
from product
where id = ?;

-- get all store-products that have this product
select product_id,
       store_id,
       off_percent,
       max_off_price,
       price,
       available_count,
       warranty_id,
       created_at,
       is_last_version
from store_product
where product_id = ?
  and available_count > 0;

-- get product warranty by warranty id
select id, name, type, month, created_at
from warranty
where id = ?;

-- add warranty
insert into warranty(name, type, month)
values (?, ?, ?);
-- delete warranty by warranty_id
delete
from warranty
where id = ?;

-- get product category
with recursive cte as (
    select c.id, c.name, c.parent
    from product p
             join category c on p.category_id = c.id
    where p.id = 3
    union
    select c.id, c.name, c.parent
    from category c
             join cte ct on ct.parent = c.id
)
select id, name
from cte;
-- subscribe to product availability
insert into product_available_subscription(product_id, user_id)
values (?, ?);
-- get all user notifications
select id, product_id, user_id, created_at, is_notification_sent, available_status
from product_available_subscription
where user_id = ?
  and is_notification_sent = false
  and available_status = true;
-- make notification seen
update product_available_subscription
set is_notification_sent = true
where user_id = ?
  and id = ?;
-- get all products that user subscribed on and not available
select product_id, user_id, created_at, is_notification_sent
from product_available_subscription
where user_id = ?
  and available_status = false
  and is_notification_sent = false;

-- create review
insert into review (product_id, store_id, user_id, rate, review_text)
values (?, ?, ?, ?, ?);
-- todo: check if user by a product from store.
-- update review
update review
set rate=?, review_text = ?
where id = ?
  and user_id=?
  and deleted_at is null;
-- get user all reviews
select id, product_id, store_id, user_id, rate, review_text, created_at
from review
where user_id = ?
  and deleted_at is null;
-- delete review
update review
set deleted_at = now()
where id = ?
  and user_id = ?
  and deleted_at is null;

-- create votes
insert into votes (review_id, user_id, up_vote, down_vote)
values (?, ?, ?, ?);
-- delete vote
delete
from votes
where user_id = ?
  and review_id = ?;

-- sort reviews by votes count
with cte as (
    select r.id,
           r.product_id,
           r.store_id,
           r.user_id,
           r.rate,
           r.review_text,
           r.created_at,
           r.deleted_at,
           sum(case when v.up_vote = true then 1 else 0 end)   as up_votes,
           sum(case when v.down_vote = true then 1 else 0 end) as down_votes
    from review r
             left join votes v on r.id = v.review_id
    group by r.id, r.product_id, r.store_id, r.user_id, r.rate, r.review_text, r.created_at
)
select id,
       product_id,
       store_id,
       user_id,
       rate,
       review_text,
       created_at,
       up_votes,
       down_votes
from cte
where product_id = ?
  and deleted_at is null
order by up_votes * 2 + down_votes desc;


-- sort reviews by created date
select r.id,
       r.product_id,
       r.store_id,
       r.user_id,
       r.rate,
       r.review_text,
       r.created_at,
       sum(case when v.up_vote = true then 1 else 0 end)   as up_votes,
       sum(case when v.down_vote = true then 1 else 0 end) as down_votes
from review r
         left join votes v on r.id = v.review_id
where product_id = ?
  and deleted_at is null
group by r.id, r.product_id, r.store_id, r.user_id, r.rate, r.review_text, r.created_at
order by created_at desc;

-- add promotion_code
insert into promotion_code(id, percent, max_off_price)
values (?, ?, ?);
-- delete promotion_code
update promotion_code
set deleted_at = now()
where id = ?;

-- create new order
insert into "order"(user_id, address_id, product_id, store_id, shipping_method_id, applied_promotion_code)
values (?, ?, ?, ?, ?);

-- add product to order
insert into product_order(product_id, store_id, order_id)
values (?, ?, ?);
-- delete product from order
delete
from product_order
where order_id = ?
  and product_id = ?
  and store_id = ?;
-- update product order quantity
update product_order
set quantity = ?
where order_id = ?
  and product_id = ?
  and store_id = ?;
-- get all product_ids in the order
select product_id, store_id
from product_order
where order_id = ?;
-- get all user orders
select order_id,
       status,
       tracking_code,
       user_id,
       address_id,
       product_id,
       store_id,
       shipping_method_id,
       applied_promotion_code,
       is_paid,
       pay_date,
       created_at
from "order"
where user_id = ?;
-- filter orders by status
select order_id,
       status,
       tracking_code,
       user_id,
       address_id,
       product_id,
       store_id,
       shipping_method_id,
       applied_promotion_code,
       is_paid,
       pay_date,
       created_at
from "order"
where status = ?;
-- user payed for the order
begin;
update "order"
set is_paid = true,
    status='confirmed',
    pay_date=now()
where order_id = ?
  and user_id = ?;

insert into payment(order_id, user_id, total_price)
values (?, ?, ?);
commit;


-- crate ticket-types
-- todo: insert some types
insert into ticket_type(name, description)
values ('', '');
-- update ticket-types
begin;
update ticket_type
set is_last_version = false
where id = ?;
insert into ticket_type(name, description)
values (?, ?);
commit;
-- delete ticket_type
update ticket_type
set is_last_version = false
where id = ?;
-- get all ticket types
select id, name, description
from ticket_type
where is_last_version = true;

-- create new ticket
insert into ticket(user_id, ticket_type_id)
values (?, ?);
-- get user all tickets
select id, user_id, employee_id, ticket_type_id, is_done, done_at, created_at
from ticket
where user_id = ?
order by created_at desc;
-- filter user tickets by ticket_type
select id, user_id, employee_id, ticket_type_id, is_done, done_at, created_at
from ticket
where user_id = ?;

-- make chat done!
update ticket
set is_done = true,
    done_at = now()
where id = ?
  and employee_id = ?;

-- add ticket message
insert into ticket_message(ticket_id, sender_id, message_text)
values (?, ?, ?);
-- set message status to received
update ticket_message
set status = 'received'
where id = ?
  and ticket_id = ?
  and sender_id = ?;

-- set message status to seen
update ticket_message
set status = 'seen'
where id = ?
  and ticket_id = ?
  and sender_id = ?;
-- get ticket messages
select id, ticket_id, sender_id, message_text, status, created_at
from ticket_message
where ticket_id = ?
order by created_at desc
limit 5 offset ?;


