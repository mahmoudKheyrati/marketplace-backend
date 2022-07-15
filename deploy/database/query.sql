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
from "user"
where id = ?;
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
select id, name
from category
where parent is null;
-- get all categories with sub-categories
-- todo: implement later if needed.
select id, name, parent
from category
order by id, parent;

-- get subcategories by category_id
with recursive cte as (
    select id, name, parent
    from category
    where id = ?
    union
    select c2.id, c2.name, c2.parent
    from category c2
             join cte as c1 on c2.parent = c1.id
)
select id, name, parent
from cte;

-- get parents by category_id
with recursive cte as (
    select id, name, parent
    from category
    where id = ?
    union
    select c2.id, c2.name, c2.parent
    from category c2
             join cte as c1 on c2.id = c1.parent
)
select id, name, parent
from cte;


--get store by store_id
select id,
       name,
       description,
       avatar_url,
       owner,
       creator,
       created_at
from store where id = ?;
-- create store
insert into store(name, description, avatar_url, owner, creator)
values (?, ?, ?, ?, ?);

-- update store
update store
set name= ?,
    description = ?,
    avatar_url = ?
where id = ?
  and owner = ?;

-- get user all stores
select id, name, description, avatar_url, owner, creator, created_at
from store
where creator =? or owner = ?
  and deleted_at is null;

-- get all marketplace stores
select id, name, description, avatar_url, owner, creator, created_at
from store
where deleted_at is null;
-- delete store
update store
set deleted_at = now()
where id = ?
  and owner = ?;

-- add address for store
insert into store_address (country, province, city, street, postal_code)
values (?, ?, ?, ?, ?);
-- check user is the owner of store or not
select exists(select 1 from store where owner = ? and id = ?);
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

-- store add product
insert into store_product(product_id, store_id, off_percent, max_off_price, price, available_count, warranty_id)
values (?, ?, ?, ?, ?, ?, ?);
-- update store product
update store_product
set off_percent= ?,
    max_off_price= ?,
    price = ?,
    available_count = ?,
    warranty_id = ?
where store_id = ?
  and product_id = ?;
-- get all products of store
select id,
       category_id,
       name,
       brand,
       description,
       picture_url,
       specification,
       p.created_at as created_at
from product p
         join store_product sp on p.id = sp.product_id
where id = ?;


-- add product
insert into product (category_id, name, brand, description, picture_url, specification)
values (?, ?, ?, ?, ?, ?);

-- delete product
delete
from product
where id = ?;

-- search products by name
select id, category_id, name, brand, description, picture_url, specification
from product
where name ilike '%?%';
-- important
-- search products by brand
select id, category_id, name, brand, description, picture_url, specification
from product
where brand ilike '%?%';
-- important
-- search products by category
select id,
       category_id,
       name,
       brand,
       description,
       picture_url,
       specification,
       created_at
from product
where category_id in (
    with recursive cte as (
        select id, name, parent
        from category
        where name = ?
        union
        select c2.id, c2.name, c2.parent
        from category c2
                 join cte as c1 on c2.parent = c1.id
    )
    select id
    from cte
) order by category_id;


-- similar products
with cte as (
    select product.id,
           category_id,
           name,
           brand,
           description,
           picture_url,
           specification,
           product.created_at,
           sum(2*(case when v.up_vote = true then 1 else 0 end) + (case when v.down_vote = true then 1 else 0 end)) as sum
    from product
             left join review r on product.id = r.product_id left join votes v on r.id = v.review_id
    where category_id = ?
    group by product.id, category_id, name, brand, description, picture_url, specification, product.created_at
)select id,
        category_id,
        name,
        brand,
        description,
        picture_url,
        specification,
        created_at
from cte order by sum desc limit 10;

-- get frequently bought together products
with cte as (
    select p.id,
           category_id,
           name,
           brand,
           description,
           picture_url,
           specification,
           p.created_at,
           count(po.quantity) as co

    from product p
             left join product_order po on p.id = po.product_id
             left join "order" o on o.id = po.order_id

    where p.id = ? and o.is_paid = true and payed_price != 0 and pay_date is not null
    group by p.id, category_id, name, brand, description, picture_url, specification, p.created_at

)select id,
        category_id,
        name,
        brand,
        description,
        picture_url,
        specification,
        created_at
from cte order by co desc limit 10;


-- get distinct brand by categoryId
select distinct brand from product where category_id = ?;

-- filter by brand
select id, category_id, name, brand, description, picture_url, specification
from product
where category_id = ?
  and brand in (?);
-- get price range by categoryId
select min(sp.price) as min, max(sp.price) as max from store_product sp join product p on p.id = sp.product_id
where p.category_id = ?;
-- filter by price
select id, category_id, name, brand, description, picture_url, specification
from product p
         join store_product sp on p.id = sp.product_id
where p.category_id = ?
  and sp.price between ? and ?;

-- get category specifications distinct
select distinct (jsonb_object_keys(p.specification))
from product p
         join category c on p.category_id = c.id
where p.category_id = ?;

-- filter by specification
select id, category_id, name, brand, description, picture_url, specification
from product
where category_id = ?
  and specification::jsonb ->> ? = ?;
-- important



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
         left join review r on p.id = r.product_id
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

-- =========================== search ===================================
-- with search_result as (
--     select id,
--            category_id,
--            name,
--            brand,
--            description,
--            picture_url,
--            specification,
--            product.created_at as created_at
--     from product
--
--     where (name ilike '%$1%' or brand ilike '%$1%' or category_id in (
--         with recursive cte as (
--             select id, name, parent
--             from category
--             where name ilike '%$1%'
--             union
--             select c2.id, c2.name, c2.parent
--             from category c2
--                      join cte as c1 on c2.parent = c1.id
--         ) select cte.id from cte))
--       -- filter queries
--         except
--     select id,
--            category_id,
--            name,
--            brand,
--            description,
--            picture_url,
--            specification,
--            created_at
--     from product
--     where brand != $4
--         except
--     select id,
--            category_id,
--            name,
--            brand,
--            description,
--            picture_url,
--            specification,
--            created_at
--     from product where specification::jsonb ->> ? != ?
-- ) select * from search_result sr left join store_product sp on sr.id = sp.product_id
-- where price between $5 and $6
-- order by (case
--     when $7 = 'byDate' then sr.created_at
--     when $7 = 'byCheapest' then sp.price
--     when $7 = 'byExpensive' then sp.price
--     when $7 = 'byReviewCount' then 1
--     end) ;


-- =========================== filter ===================================



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

-- get all store-products that have a specific product
select product_id,
       store_id,
       off_percent,
       max_off_price,
       price,
       available_count,
       warranty_id,
       store_product.created_at,
       w.name as warranty_name,
       w.type as warranty_type,
       w.month as warranty_month
from store_product
left join warranty w on store_product.warranty_id = w.id
where product_id = ?
  and available_count > 0 and is_last_version = true;

-- get product warranty by warranty id
select id, name, type, month, created_at
from warranty
where id = ?;
-- get warranty by store_product id
select id, name, type, month, created_at
from warranty
where id in (select warranty_id from store_product where store_id = ? and product_id = ?)
limit 1;



-- add warranty
insert into warranty(name, type, month)
values (?, ?, ?);
-- delete warranty by warranty_id
delete
from warranty
where id = ?;
-- get warranty by warranty_id
select id, name, type, month, created_at
from warranty where id = ?;

-- get product category
with recursive cte as (
    select c.id, c.name, c.parent
    from product p
             join category c on p.category_id = c.id
    where p.id = ?
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
set rate=?,
    review_text = ?
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
--apply promotion code to order
update "order" set applied_promotion_code = ? where id = ? and user_id = ? ;
-- delete promotion code from oder
update "order" set applied_promotion_code = null where id= ? and user_id = ? ;


-- apply shipping method
update "order" set shipping_method_id = ? where id= ? and user_id = ? ;

-- create new order
insert into "order"(user_id )
values (?) returning id;
-- get order by order_id
select id,
       status,
       tracking_code,
       user_id,
       address_id,
       shipping_method_id,
       applied_promotion_code,
       payed_price,
       is_paid,
       pay_date,
       created_at,
       deleted_at
from "order" where user_id = ? and id = ? ;
-- check user owns the order
select exists(select 1 from order where id = ? and user_id= ? );
-- check if order not deleted
select deleted_at is not null from "order" where id= ? and user_id = ? ;
-- delete order
delete from "order" where id = ? and user_id= ? and deleted_at is null  ;
-- check user paid for order
select is_paid from "order" where id = ? and user_id=? and deleted_at is null;
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
select product_id, store_id, order_id, quantity, created_at
from product_order
where order_id = ?;
-- get all user orders
select id,
       status,
       tracking_code,
       user_id,
       address_id,
       shipping_method_id,
       applied_promotion_code,
       payed_price,
       is_paid,
       pay_date,
       created_at
from "order"
where user_id = ? and deleted_at is null;
-- filter orders by status
select id,
       status,
       tracking_code,
       user_id,
       address_id,
       shipping_method_id,
       applied_promotion_code,
       payed_price,
       is_paid,
       pay_date,
       created_at
from "order"
where status = ?;
-- user payed for the order
update "order"
set is_paid = true,
    status='confirmed',
    payed_price = ?,
    pay_date=now()
where id = ?
  and user_id = ?;



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

-- get all unfinished tickets
select id, user_id, employee_id, ticket_type_id, is_done, done_at, created_at
from ticket where is_done = false order by created_at desc ;

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


