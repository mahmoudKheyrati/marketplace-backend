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

