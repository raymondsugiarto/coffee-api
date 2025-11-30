-- Organization initialization
INSERT INTO organization (id, name, code, origin, created_at, updated_at)
 VALUES ('16f31f95-b356-4e96-b0df-c7f5052beb95', '1', '1', 'sekian', NOW(), NOW());    

-- User base entries
INSERT INTO "user" (id, organization_id, user_type, created_at, updated_at)
 VALUES 
    ('e55be757-5f65-4423-92b7-5b573cba374b', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'ADMIN', NOW(), NOW()), -- adm
    ('3596df3b-13aa-473d-b4d7-2a744b504af3', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'ADMIN', NOW(), NOW()), -- e1
    ('adc0d06f-cc21-4419-9936-d249e8a2494e', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'ADMIN', NOW(), NOW()) -- e2
;   

-- Admin users
INSERT INTO admin (id, organization_id, user_id, admin_type, phone_number, email, first_name, last_name, created_at, updated_at)
VALUES 
   ('7c167f13-0587-4a35-9298-dded451ca929', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'e55be757-5f65-4423-92b7-5b573cba374b', 'ADMIN', '08179340556', 'adm@gmail.com', 'Admin', '', NOW(), NOW()),
   ('1ce4485c-bb03-4224-bd7c-d4c0f852db43', '16f31f95-b356-4e96-b0df-c7f5052beb95', '3596df3b-13aa-473d-b4d7-2a744b504af3', 'EMPLOYEE', '081123456789', 'e1@gmail.com', 'Mita', '', NOW(), NOW()),
   ('b916d6a4-f325-465e-8c91-308f03d01405', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'adc0d06f-cc21-4419-9936-d249e8a2494e', 'EMPLOYEE', '081123456780', 'e2@gmail.com', 'Berry', '', NOW(), NOW())
;

-- Company entries
INSERT INTO company (id, organization_id, phone_number, name, created_at, updated_at)
VALUES 
   ('a93150c9-eb99-4c62-8fc1-c414c8a0f78d', '16f31f95-b356-4e96-b0df-c7f5052beb95', '1234567890', 'Sekian',  NOW(), NOW()),
   ('3b80cdb0-5a52-42e5-9a9b-38f7c3e9164f', '16f31f95-b356-4e96-b0df-c7f5052beb95', '1234567890', 'Mawaru',  NOW(), NOW())
;

-- Company entries
INSERT INTO admin_company (id, organization_id, company_id, admin_id, created_at, updated_at)
VALUES 
   ('9438e093-e944-4a02-8a46-07cac69e7be2', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'a93150c9-eb99-4c62-8fc1-c414c8a0f78d', '7c167f13-0587-4a35-9298-dded451ca929',  NOW(), NOW()),
   ('16cdad6a-f3fd-4b7e-9a2e-764ad468d080', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'a93150c9-eb99-4c62-8fc1-c414c8a0f78d', '1ce4485c-bb03-4224-bd7c-d4c0f852db43',  NOW(), NOW()),
   ('45ec3144-92d7-4641-a970-4b3e11f5a4ae', '16f31f95-b356-4e96-b0df-c7f5052beb95', '3b80cdb0-5a52-42e5-9a9b-38f7c3e9164f', 'b916d6a4-f325-465e-8c91-308f03d01405',  NOW(), NOW())
;


-- User credentials
INSERT INTO user_credential (id, organization_id, user_id, username, password, created_at, updated_at)
 VALUES 
   ('1918b887f-6bbd-4c85-a944-e04a6bac5a6c', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'e55be757-5f65-4423-92b7-5b573cba374b', 'admin', '$2a$07$rXsQxQHRwxwHYNTzHKTl.eilofdCZ9Ci0TTJmLdV6I7rxsYn/O74.', NOW(), NOW()), -- pass admin
   ('3eb49372-e394-4ee1-8871-00fa7e5e9a77', '16f31f95-b356-4e96-b0df-c7f5052beb95', '3596df3b-13aa-473d-b4d7-2a744b504af3', '081123456789', '$2a$07$rXsQxQHRwxwHYNTzHKTl.eilofdCZ9Ci0TTJmLdV6I7rxsYn/O74.', NOW(), NOW()),
   ('99ce7f88-dd72-4426-88c9-1775d3c30ce1', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'adc0d06f-cc21-4419-9936-d249e8a2494e', '081123456780', '$2a$07$rXsQxQHRwxwHYNTzHKTl.eilofdCZ9Ci0TTJmLdV6I7rxsYn/O74.', NOW(), NOW())
;

INSERT INTO item (id, organization_id, code, name, price, created_at, updated_at)
 VALUES 
   ('390aed07-5681-4a1e-8819-8a6f5aff3ee4', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'KAB', 'Kopi Aren Blend', 8000, NOW(), NOW()),
   ('a0ec3af5-dc96-4a66-ad43-0f27f69d00de', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'CMB', 'Caramel Macchiato Blend', 9000, NOW(), NOW()),
   ('816e96f9-27f1-4725-993c-053a51293362', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'AMB', 'Americano Blend', 9000, NOW(), NOW()),
   ('92bf21ac-7759-4e03-972b-1209406ebd83', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'MAM', 'Mawaru Matcha', 15000, NOW(), NOW()),
   ('7f99a6c7-99b9-435c-ba2a-58ea77760e8d', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'CMM', 'Caramel Macchiato Matcha', 17000, NOW(), NOW()),
   ('9733b9e5-97c1-4135-829c-058fc282239c', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'SMM', 'Strawberry Macchiato Matcha', 18000, NOW(), NOW()),
   ('963bcf80-b3a7-40fe-88c4-ff8c23a1d81d', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'MMM', 'Mango Macchiato Matcha', 18000, NOW(), NOW()),
   ('65b5a7a6-6eb1-4f3f-8ef8-797c676feae2', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'COK', 'Cokelat', 9000, NOW(), NOW()),
   ('51d39c73-528e-43ff-9fb9-2bb4b659d3ba', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'LYC', 'Lychee Tea', 8000, NOW(), NOW())
;

INSERT INTO item_company (id, organization_id, item_id, company_id, created_at, updated_at)
VALUES 
   -- SEKIAN
   ('e1bd2514-7b1d-4eb4-9703-3b0f7b4d56a2', '16f31f95-b356-4e96-b0df-c7f5052beb95', '390aed07-5681-4a1e-8819-8a6f5aff3ee4','a93150c9-eb99-4c62-8fc1-c414c8a0f78d', NOW(), NOW()),
   ('b1817dd2-b9e2-43a0-985e-e5b9b82a47c5', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'a0ec3af5-dc96-4a66-ad43-0f27f69d00de','a93150c9-eb99-4c62-8fc1-c414c8a0f78d', NOW(), NOW()),
   ('224e6173-1a50-40c7-89fc-21ce53087fef', '16f31f95-b356-4e96-b0df-c7f5052beb95', '816e96f9-27f1-4725-993c-053a51293362','a93150c9-eb99-4c62-8fc1-c414c8a0f78d', NOW(), NOW()),
   ('ea03e693-c420-4747-bd97-630e79de76fc', '16f31f95-b356-4e96-b0df-c7f5052beb95', '92bf21ac-7759-4e03-972b-1209406ebd83','a93150c9-eb99-4c62-8fc1-c414c8a0f78d', NOW(), NOW()),
   ('62d8d8df-da07-497e-a17d-fbe11512c745', '16f31f95-b356-4e96-b0df-c7f5052beb95', '7f99a6c7-99b9-435c-ba2a-58ea77760e8d','a93150c9-eb99-4c62-8fc1-c414c8a0f78d', NOW(), NOW()),
   ('66f9a024-fb8d-4305-af68-358c77fe9918', '16f31f95-b356-4e96-b0df-c7f5052beb95', '65b5a7a6-6eb1-4f3f-8ef8-797c676feae2','a93150c9-eb99-4c62-8fc1-c414c8a0f78d', NOW(), NOW()),
   ('7e9dbaa9-4fc7-43ac-8e46-162af3714a3f', '16f31f95-b356-4e96-b0df-c7f5052beb95', '51d39c73-528e-43ff-9fb9-2bb4b659d3ba','a93150c9-eb99-4c62-8fc1-c414c8a0f78d', NOW(), NOW()),
   

   -- MAWARU
   ('81b261b1-e82c-4b3b-902d-254ef78ecdfb', '16f31f95-b356-4e96-b0df-c7f5052beb95', '92bf21ac-7759-4e03-972b-1209406ebd83','3b80cdb0-5a52-42e5-9a9b-38f7c3e9164f', NOW(), NOW()),
   ('6d34a537-f30d-4a41-9af4-8c105dd08f97', '16f31f95-b356-4e96-b0df-c7f5052beb95', '7f99a6c7-99b9-435c-ba2a-58ea77760e8d','3b80cdb0-5a52-42e5-9a9b-38f7c3e9164f', NOW(), NOW()),
   ('284481f4-0be5-47b0-81e6-25fdf10e98b9', '16f31f95-b356-4e96-b0df-c7f5052beb95', '9733b9e5-97c1-4135-829c-058fc282239c','3b80cdb0-5a52-42e5-9a9b-38f7c3e9164f', NOW(), NOW()),
   ('4193dd10-6658-4243-8114-c34e34130aba', '16f31f95-b356-4e96-b0df-c7f5052beb95', '963bcf80-b3a7-40fe-88c4-ff8c23a1d81d','3b80cdb0-5a52-42e5-9a9b-38f7c3e9164f', NOW(), NOW()),
   ('a542bb86-90ec-45d6-bd8c-460d0b8835c2', '16f31f95-b356-4e96-b0df-c7f5052beb95', '390aed07-5681-4a1e-8819-8a6f5aff3ee4','3b80cdb0-5a52-42e5-9a9b-38f7c3e9164f', NOW(), NOW()),
   ('1a0a20d2-d105-4d1c-ad6b-3943541416a4', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'a0ec3af5-dc96-4a66-ad43-0f27f69d00de','3b80cdb0-5a52-42e5-9a9b-38f7c3e9164f', NOW(), NOW()),
   ('19a74187-2531-43e4-aea9-a5be4b6b881b', '16f31f95-b356-4e96-b0df-c7f5052beb95', '65b5a7a6-6eb1-4f3f-8ef8-797c676feae2','3b80cdb0-5a52-42e5-9a9b-38f7c3e9164f', NOW(), NOW()),
   ('2910144a-8a25-461c-8c10-8e954bc7bf51', '16f31f95-b356-4e96-b0df-c7f5052beb95', '51d39c73-528e-43ff-9fb9-2bb4b659d3ba','3b80cdb0-5a52-42e5-9a9b-38f7c3e9164f', NOW(), NOW())

;



INSERT INTO "user" (id, organization_id, user_type, created_at, updated_at)
 VALUES 
    ('e55be757-5f65-4423-92b7-5b573cba374b', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'ADMIN', NOW(), NOW()), -- adm
    ('3596df3b-13aa-473d-b4d7-2a744b504af3', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'ADMIN', NOW(), NOW()), -- e1
    ('adc0d06f-cc21-4419-9936-d249e8a2494e', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'ADMIN', NOW(), NOW()) -- e2
;  
INSERT INTO admin (id, organization_id, user_id, admin_type, phone_number, email, first_name, last_name, created_at, updated_at)
VALUES 
   ('7c167f13-0587-4a35-9298-dded451ca929', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'e55be757-5f65-4423-92b7-5b573cba374b', 'ADMIN', '08179340556', 'adm@gmail.com', 'Admin', '', NOW(), NOW()),
   ('1ce4485c-bb03-4224-bd7c-d4c0f852db43', '16f31f95-b356-4e96-b0df-c7f5052beb95', '3596df3b-13aa-473d-b4d7-2a744b504af3', 'EMPLOYEE', '081123456789', 'e1@gmail.com', 'Mita', '', NOW(), NOW()),
   ('b916d6a4-f325-465e-8c91-308f03d01405', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'adc0d06f-cc21-4419-9936-d249e8a2494e', 'EMPLOYEE', '081123456780', 'e2@gmail.com', 'Berry', '', NOW(), NOW())
;
INSERT INTO admin_company (id, organization_id, company_id, admin_id, created_at, updated_at)
VALUES 
   ('9438e093-e944-4a02-8a46-07cac69e7be2', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'a93150c9-eb99-4c62-8fc1-c414c8a0f78d', '7c167f13-0587-4a35-9298-dded451ca929',  NOW(), NOW()),
   ('16cdad6a-f3fd-4b7e-9a2e-764ad468d080', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'a93150c9-eb99-4c62-8fc1-c414c8a0f78d', '1ce4485c-bb03-4224-bd7c-d4c0f852db43',  NOW(), NOW()),
   ('45ec3144-92d7-4641-a970-4b3e11f5a4ae', '16f31f95-b356-4e96-b0df-c7f5052beb95', '3b80cdb0-5a52-42e5-9a9b-38f7c3e9164f', 'b916d6a4-f325-465e-8c91-308f03d01405',  NOW(), NOW())
;
INSERT INTO user_credential (id, organization_id, user_id, username, password, created_at, updated_at)
 VALUES 
   ('1918b887f-6bbd-4c85-a944-e04a6bac5a6c', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'e55be757-5f65-4423-92b7-5b573cba374b', 'admin', '$2a$07$rXsQxQHRwxwHYNTzHKTl.eilofdCZ9Ci0TTJmLdV6I7rxsYn/O74.', NOW(), NOW()), -- pass admin
   ('3eb49372-e394-4ee1-8871-00fa7e5e9a77', '16f31f95-b356-4e96-b0df-c7f5052beb95', '3596df3b-13aa-473d-b4d7-2a744b504af3', '081123456789', '$2a$07$rXsQxQHRwxwHYNTzHKTl.eilofdCZ9Ci0TTJmLdV6I7rxsYn/O74.', NOW(), NOW()),
   ('99ce7f88-dd72-4426-88c9-1775d3c30ce1', '16f31f95-b356-4e96-b0df-c7f5052beb95', 'adc0d06f-cc21-4419-9936-d249e8a2494e', '081123456780', '$2a$07$rXsQxQHRwxwHYNTzHKTl.eilofdCZ9Ci0TTJmLdV6I7rxsYn/O74.', NOW(), NOW())
;



9f5bb25f-ca8c-4a12-851e-81d3adef8f7a
79169b7d-ffaf-46b4-9475-23ae206ffe65
a5971e36-6891-4983-997f-8f515301dc0f
835de52c-a861-4920-927f-c2737966ff20
aaf51601-4c89-4547-9b66-f901f19c4e26
331a3830-6e22-4bf3-aa09-3f8340a871e2
b063bd67-b066-4f71-b1a3-b220b523c3c5
1ee8a3d9-6236-4ec7-af98-4bb38debd2e3




-- ===


-- Customer/Employee entries
INSERT INTO customer (id, user_id, organization_id, company_id, phone_number, email, first_name, last_name, 
                     annual_income, employer_percentage, employer_amount, customer_percentage, customer_amount,
                     identity_card_file, customer_photo, created_at, updated_at)
 VALUES 
    ('190e7153c-7923-4db2-81ce-754e446d8048', 'bdc0ff16-6633-4acf-b736-78193b4b0b47', '1', NULL, '1234567890', 'employee@gmail.com', 'employee', 'employee', 
     '75000000', 5.0000, 3750000.0000, 2.5000, 1875000.0000, 'id_card_001.jpg', 'photo_001.jpg', NOW(), NOW()),
    
    ('54827bb1-873b-49b2-8871-02d1749bd47b', '3596df3b-13aa-473d-b4d7-2a744b504af3', '1', 'a93150c9-eb99-4c62-8fc1-c414c8a0f78d', 
     '12345678901', 'employee1@gmail.com', 'E1', '1', '60000000', 4.0000, 2400000.0000, 2.0000, 1200000.0000, 
     'id_card_002.jpg', 'photo_002.jpg', NOW(), NOW()),
    
    ('e5e67bc7-d2f4-4d98-8ccc-9e11ec61a27e', 'adc0d06f-cc21-4419-9936-d249e8a2494e', '1', 'a93150c9-eb99-4c62-8fc1-c414c8a0f78d', 
     '12345678902', 'employee2@gmail.com', 'E2', '2', '55000000', 4.5000, 2475000.0000, 2.2500, 1237500.0000,
     'id_card_003.jpg', 'photo_003.jpg', NOW(), NOW());

-- Roles
INSERT INTO role (id, organization_id, code, name, created_at, updated_at)
 VALUES 
    ('b6c6a0b7-c8a1-40de-840e-0af344607a10', '1', 'ADMIN', 'ADMIN', NOW(), NOW()),
    ('d331bf1d-67a7-4948-8203-29251732f947', '1', 'COMPANY', 'COMPANY', NOW(), NOW()),
    ('e4bb9af3-b884-4988-83c3-9dcecbd1f467', '1', 'EMPLOYEE', 'EMPLOYEE', NOW(), NOW());

-- User role assignments
INSERT INTO user_has_role (id, user_id, role_id, created_at, updated_at)
 VALUES 
    ('11b2f59cc-e908-4942-a564-2e33835090df', 'e55be757-5f65-4423-92b7-5b573cba374b', 'b6c6a0b7-c8a1-40de-840e-0af344607a10', NOW(), NOW()),
    ('1047c67e-0c12-4887-83d8-29c812171b8f', 'c108752a-8833-4c87-8bda-fcace2979b19', 'd331bf1d-67a7-4948-8203-29251732f947', NOW(), NOW());
