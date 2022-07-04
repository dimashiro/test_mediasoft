CREATE EXTENSION IF NOT EXISTS ltree;
CREATE TABLE IF NOT EXISTS departments (
    department_id UUID,
    department_name text NOT NULL,
    department_path ltree,

    PRIMARY KEY (department_id)
);
create index department_path_idx on departments using gist (department_path);

CREATE TABLE IF NOT EXISTS employees (
    employee_id UUID,
    employee_name text NOT NULL,
    employee_surname text NOT NULL,
    employee_birthyear integer,

    PRIMARY KEY (employee_id)
);

CREATE TABLE IF NOT EXISTS employee_department (
  employee_id   UUID REFERENCES employees (employee_id) ON UPDATE CASCADE ON DELETE CASCADE,
  department_id UUID REFERENCES departments (department_id) ON UPDATE CASCADE,
  CONSTRAINT employee_department_pkey PRIMARY KEY (employee_id, department_id)
);