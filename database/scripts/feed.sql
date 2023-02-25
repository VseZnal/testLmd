INSERT INTO warehouse (
    id,
    name,
    active
)
VALUES
    (111, 'main', true),
    (222, 'moscow', true),
    (333, 'spb', true),
    (444, 'ekb', false);

INSERT INTO products (
    name,
    size,
    quantity,
    warehouse_id
)
VALUES
('test name1', 'test seze1', 1, 111),
('test name2', 'test seze2', 2, 111),
('test name3', 'test seze3', 3, 111),
('test name4', 'test seze4', 1, 222),
('test name5', 'test seze5', 2, 222),
('test name6', 'test seze6', 3, 222),
('test name7', 'test seze7', 1, 333),
('test name8', 'test seze8', 2, 333),
('test name9', 'test seze9', 3, 333),
('test name10', 'test seze10', 4, 444);
