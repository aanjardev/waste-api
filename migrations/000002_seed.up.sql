-- Seed data untuk testing

INSERT INTO households (id, owner_name, address) VALUES
('11111111-1111-1111-1111-111111111111', 'Budi Santoso', 'Jl. Merdeka No. 10, Jakarta'),
('22222222-2222-2222-2222-222222222222', 'Siti Aminah', 'Jl. Sudirman No. 25, Bandung'),
('33333333-3333-3333-3333-333333333333', 'Andi Wijaya', 'Jl. Diponegoro No. 5, Surabaya');

INSERT INTO waste_pickups (id, household_id, type, status, safety_check) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111', 'organic', 'pending', FALSE),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '22222222-2222-2222-2222-222222222222', 'plastic', 'scheduled', FALSE);

--