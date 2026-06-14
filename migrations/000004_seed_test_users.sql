-- +goose Up
INSERT INTO users (id, telegram_id, username, name, role, status, noise_level, karma, created_at, updated_at) VALUES
('f0000000-0000-0000-0000-000000000001', 10001, 'tsoy_legend', 'Виктор Цой', 'musician', 'active', 'light', 100, now() - INTERVAL '15 days', now() - INTERVAL '15 days'),
('f0000000-0000-0000-0000-000000000002', 10002, 'zemfira_live', 'Земфира', 'musician', 'active', 'light', 95, now() - INTERVAL '14 days', now() - INTERVAL '14 days'),
('f0000000-0000-0000-0000-000000000003', 10003, 'basta_nagano', 'Баста', 'musician', 'active', 'hard', 80, now() - INTERVAL '13 days', now() - INTERVAL '13 days'),
('f0000000-0000-0000-0000-000000000004', 10004, 'splean_official', 'Сплин', 'musician', 'active', 'medium', 90, now() - INTERVAL '12 days', now() - INTERVAL '12 days'),
('f0000000-0000-0000-0000-000000000005', 10005, 'bi2_band', 'Би-2', 'musician', 'active', 'hard', 85, now() - INTERVAL '11 days', now() - INTERVAL '11 days'),
('f0000000-0000-0000-0000-000000000006', 10006, 'ddt_shevchuk', 'ДДТ', 'musician', 'active', 'medium', 100, now() - INTERVAL '10 days', now() - INTERVAL '10 days'),
('f0000000-0000-0000-0000-000000000007', 10007, 'nautilus_pomp', 'Наутилус Помпилиус', 'musician', 'banned', 'light', 40, now() - INTERVAL '9 days', now() - INTERVAL '9 days'),
('f0000000-0000-0000-0000-000000000008', 10008, 'kino_group', 'Группа Кино', 'musician', 'active', 'hard', 95, now() - INTERVAL '8 days', now() - INTERVAL '8 days'),
('f0000000-0000-0000-0000-000000000009', 10009, 'alisa_army', 'Алиса', 'musician', 'active', 'hard', 70, now() - INTERVAL '7 days', now() - INTERVAL '7 days'),
('f0000000-0000-0000-0000-000000000010', 10010, 'aquarium_bg', 'Аквариум', 'musician', 'active', 'light', 100, now() - INTERVAL '6 days', now() - INTERVAL '6 days'),
('f0000000-0000-0000-0000-000000000011', 10011, 'korol_i_shut', 'Король и Шут', 'musician', 'active', 'hard', 75, now() - INTERVAL '5 days', now() - INTERVAL '5 days'),
('f0000000-0000-0000-0000-000000000012', 10012, 'grazhdanskaya_oborona', 'Гражданская Оборона', 'musician', 'banned', 'medium', 20, now() - INTERVAL '4 days', now() - INTERVAL '4 days'),
('f0000000-0000-0000-0000-000000000013', 10013, 'mumiytroll_band', 'Мумий Тролль', 'musician', 'pending', 'light', 100, now() - INTERVAL '3 days', now() - INTERVAL '3 days'),
('f0000000-0000-0000-0000-000000000014', 10014, 'leningrad_shnur', 'Ленинград', 'musician', 'active', 'hard', 60, now() - INTERVAL '2 days', now() - INTERVAL '2 days'),
('f0000000-0000-0000-0000-000000000015', 10015, 'piknik_official', 'Пикник', 'musician', 'pending', 'medium', 100, now() - INTERVAL '1 day', now() - INTERVAL '1 day')
ON CONFLICT (telegram_id) DO NOTHING;

INSERT INTO user_stats (user_id, total_bookings, successful_checkins, no_shows) VALUES
('f0000000-0000-0000-0000-000000000001', 20, 20, 0),
('f0000000-0000-0000-0000-000000000002', 15, 14, 1),
('f0000000-0000-0000-0000-000000000003', 30, 25, 5),
('f0000000-0000-0000-0000-000000000004', 12, 11, 1),
('f0000000-0000-0000-0000-000000000005', 25, 22, 3),
('f0000000-0000-0000-0000-000000000006', 18, 18, 0),
('f0000000-0000-0000-0000-000000000007', 10, 5, 5),
('f0000000-0000-0000-0000-000000000008', 40, 39, 1),
('f0000000-0000-0000-0000-000000000009', 14, 11, 3),
('f0000000-0000-0000-0000-000000000010', 50, 50, 0),
('f0000000-0000-0000-0000-000000000011', 22, 17, 5),
('f0000000-0000-0000-0000-000000000012', 8, 2, 6),
('f0000000-0000-0000-0000-000000000013', 0, 0, 0),
('f0000000-0000-0000-0000-000000000014', 35, 27, 8),
('f0000000-0000-0000-0000-000000000015', 0, 0, 0)
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM users WHERE telegram_id IN (
    10001, 10002, 10003, 10004, 10005, 10006, 10007, 10008, 10009, 10010, 10011, 10012, 10013, 10014, 10015
);
