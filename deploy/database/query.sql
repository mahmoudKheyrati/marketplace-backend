-- user permissions
insert into user_permission(name, is_system_admin, system_read_access, system_edit_access_product,
                            system_create_access_product, system_create_new_store, store_create_new_product,
                            store_read_new_product, store_edit_new_product, is_employee)
values ('normal-user', false, false, false, false, false, false, false, false, false),
       ('store-admin', false, false, false, false, false, true, true, true, false),
       ('marketplace-admin', true, true, true, true, true, false, false, false, true);
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
insert into "user" (email, password, phone_number, first_name, last_name, avatar_url, national_id, permission_name)
values (?, ?, ?, ?, ?, '', '', 'normal-user');
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


-- shipping method
insert into shipping_method (name, expected_arrival_working_days, base_cost)
values ('ordinary', 7, 15),
       ('express', 3, 35),
       ('special', 1, 50);
-- get all shipping method
select name, expected_arrival_working_days, base_cost, created_at
from shipping_method;

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
-- send notification of product availability to user
begin;
delete
from product_available_subscription
where user_id = ?
  and product_id = ?;
insert into notification(user_id, product_id)
values (?, ?);
commit;
-- get all products that user subscribed on
select product_id, user_id, created_at, is_notification_sent
from product_available_subscription
where user_id = ?;

-- create review
insert into review (product_id, store_id, user_id, rate, review_text)
values (?, ?, ?, ?, ?);
-- todo: check if user by a product from store.
-- update review
update review
set rate=? and review_text = ?
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
       sum (case when v.up_vote=true then 1 else 0 end) as up_votes,
       sum (case when v.down_vote=true then 1 else 0 end) as down_votes
from review r
         join votes v on r.id = v.review_id
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
order by up_votes* 2 + down_votes desc;











