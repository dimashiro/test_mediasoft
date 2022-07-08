INSERT INTO departments VALUES ('d6ca8a87-bf57-4754-bc9f-7423ae2a4d05', 'Software department', 'd6ca8a87_bf57_4754_bc9f_7423ae2a4d05');
INSERT INTO departments VALUES ('db908e71-195f-4ea7-a726-1b477e932e49', 'iOS department', 'd6ca8a87_bf57_4754_bc9f_7423ae2a4d05.db908e71_195f_4ea7_a726_1b477e932e49');
INSERT INTO departments VALUES ('0e884d7a-6084-42ab-8bf8-0f4756195a12', 'Web department', 'd6ca8a87_bf57_4754_bc9f_7423ae2a4d05.0e884d7a_6084_42ab_8bf8_0f4756195a12');
INSERT INTO departments VALUES ('a45d1ce4-9670-4518-9197-3fa8d467f5df', 'Backend department', 'd6ca8a87_bf57_4754_bc9f_7423ae2a4d05.0e884d7a_6084_42ab_8bf8_0f4756195a12.a45d1ce4_9670_4518_9197_3fa8d467f5df');
INSERT INTO departments VALUES ('6d792313-10e6-4386-a8ce-2885c10d42e0', 'Frontend department', 'd6ca8a87_bf57_4754_bc9f_7423ae2a4d05.0e884d7a_6084_42ab_8bf8_0f4756195a12.6d792313_10e6_4386_a8ce_2885c10d42e0');

INSERT INTO employees VALUES ('fd73b60e-802b-439c-b44d-b6191ac9a368', 'John', 'Doe', 1984);
INSERT INTO employees VALUES ('2e6bcd61-50e9-405b-b1ac-59d66a337cf5', 'Darwin', 'Effertz', 1990);
INSERT INTO employees VALUES ('264264e5-29be-4908-b303-19ab5d11dda5', 'Meaghan', 'Stanton', 1999);
INSERT INTO employees VALUES ('3b02d89c-a297-42a3-9ccd-19f943713ef1', 'Cade', 'Schiller', 2000);
INSERT INTO employees VALUES ('3d6f72d1-2b98-4b2f-9a89-a9e2df42cdc9', 'Rubie', 'Kunze', 1985);

INSERT INTO employee_department VALUES ('fd73b60e-802b-439c-b44d-b6191ac9a368','d6ca8a87-bf57-4754-bc9f-7423ae2a4d05');
INSERT INTO employee_department VALUES ('fd73b60e-802b-439c-b44d-b6191ac9a368','0e884d7a-6084-42ab-8bf8-0f4756195a12');
INSERT INTO employee_department VALUES ('fd73b60e-802b-439c-b44d-b6191ac9a368','a45d1ce4-9670-4518-9197-3fa8d467f5df');
INSERT INTO employee_department VALUES ('2e6bcd61-50e9-405b-b1ac-59d66a337cf5','db908e71-195f-4ea7-a726-1b477e932e49');
INSERT INTO employee_department VALUES ('264264e5-29be-4908-b303-19ab5d11dda5','6d792313-10e6-4386-a8ce-2885c10d42e0');
INSERT INTO employee_department VALUES ('3b02d89c-a297-42a3-9ccd-19f943713ef1','6d792313-10e6-4386-a8ce-2885c10d42e0');
INSERT INTO employee_department VALUES ('3d6f72d1-2b98-4b2f-9a89-a9e2df42cdc9','6d792313-10e6-4386-a8ce-2885c10d42e0');