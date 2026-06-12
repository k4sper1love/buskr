-- +goose Up
INSERT INTO locations (id, name, description, coords, max_noise, status) VALUES
('a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d', 'Арбат (Жибек Жолы)', 'Пешеходная зона на Жибек Жолы. Самое популярное место с высокой проходимостью.', ST_SetSRID(ST_MakePoint(76.9423, 43.2621), 4326), 'hard', 'active'),
('b2c3d4e5-f67a-8b9c-0d1e-2f3a4b5c6d7e', 'Панфилова (у КБТУ)', 'Красивая пешеходная зона с отличной акустикой напротив исторического сквера.', ST_SetSRID(ST_MakePoint(76.9431, 43.2562), 4326), 'medium', 'active'),
('c3d4e5f6-7a8b-9c0d-1e2f-3a4b5c6d7e8f', 'Терренкур (река)', 'Тихая и уютная прогулочная зона у реки Малая Алматинка. Идеально для акустики.', ST_SetSRID(ST_MakePoint(76.9602, 43.2215), 4326), 'light', 'active'),
('d4e5f67a-8b9c-0d1e-2f3a-4b5c6d7e8f9a', 'Центральный Парк (вход)', 'Площадка прямо перед главным входом в парк культуры. Много семейной аудитории.', ST_SetSRID(ST_MakePoint(76.9712, 43.2618), 4326), 'medium', 'active'),
('e5f67a8b-9c0d-1e2f-3a4b-5c6d7e8f9a0b', 'Театр Лермонтова', 'Просторная мощеная площадь перед театром. Отлично подходит для больших групп.', ST_SetSRID(ST_MakePoint(76.9465, 43.2428), 4326), 'hard', 'active')
ON CONFLICT (id) DO NOTHING;

-- +goose Down
DELETE FROM locations WHERE id IN (
  'a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d',
  'b2c3d4e5-f67a-8b9c-0d1e-2f3a4b5c6d7e',
  'c3d4e5f6-7a8b-9c0d-1e2f-3a4b5c6d7e8f',
  'd4e5f67a-8b9c-0d1e-2f3a-4b5c6d7e8f9a',
  'e5f67a8b-9c0d-1e2f-3a4b-5c6d7e8f9a0b'
);
