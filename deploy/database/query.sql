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
insert into category( id,  name, parent)
values (1, 'electronics', null),
       (2, 'laptop', 1),
       (3, 'phone', 1),
       (4, 'monitor' , 1),
       (5, 'shoes', null),
       (6, 'men-shoes', 5),
       (7, 'boot', 6),
       (8, 'sneakers', 6),
       (9, 'women-shoes', 5),
       (10, 'boot', 9),
       (11, 'sneakers', 9),
       (12, 'flat', 9);

-- get main categories
select name from category where parent is null;
-- get sub-categories by category id
select name from category where parent = ? ;
-- get all categories with sub-categories
-- todo: implement later if needed.





