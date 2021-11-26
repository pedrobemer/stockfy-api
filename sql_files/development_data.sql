
\connect stockfy;

-- Insert user with admin privileges
INSERT INTO
	public.users ("uid", "username", email, "type")
VALUES
	('TestAdminID', 'Test Name Admin', 'test_admin@email.com', 'admin'),
    ('TestNoAdminID', 'Test Name NoAdmin', 'test_noadmin@email.com', 'normal');

INSERT INTO
	public.sectors ("name")
VALUES
	('Finances'),
    ('Technology');

-- Insert asset
INSERT INTO
	public.assets (symbol, fullname, preference, asset_type_id, sector_id)
VALUES
	('ITUB4', 'Itau Unibanco Holding S.A', 'PN', (
        SELECT
            "at".id
        FROM asset_types as "at"
        WHERE "at"."type" = 'STOCK' and "at".country = 'BR'
        ),(
        SELECT
            s.id
        FROM sectors as s
        WHERE s."name" = 'Finances'
        )
    ),
    ('AAPL', 'Apple Inc', '', (
        SELECT
            "at".id
        FROM asset_types as "at"
        WHERE "at"."type" = 'STOCK' and "at".country = 'US'
        ),(
        SELECT
            s.id
        FROM sectors as s
        WHERE s."name" = 'Technology'
        )
    );

-- Insert Asset User
INSERT INTO
    public.asset_users (asset_id, user_uid)
VALUES
    ((
     SELECT
        a.id
     FROM assets as a
     WHERE a.symbol = 'ITUB4'
    ), (
     SELECT
        u.uid
     FROM users as u
     WHERE u.uid = 'TestAdminID'
    ));

-- Insert Orders
INSERT INTO
    public.orders (asset_id, user_uid, brokerage_id, quantity, price, currency,
        order_type, "date")
VALUES
    (
     (SELECT a.id FROM assets as a WHERE a.symbol = 'ITUB4'), (
     SELECT u.uid FROM users as u WHERE u.uid = 'TestAdminID'), (
     SELECT b.id FROM brokerages as b WHERE b.name = 'Avenue'), 5, 20.39, 'BRL',
     'buy', '2021-10-05'
    ),
    (
     (SELECT a.id FROM assets as a WHERE a.symbol = 'ITUB4'), (
     SELECT u.uid FROM users as u WHERE u.uid = 'TestAdminID'), (
     SELECT b.id FROM brokerages as b WHERE b.name = 'Avenue'), 4, 29.39, 'BRL',
     'buy', '2019-12-06'
    ),
    (
     (SELECT a.id FROM assets as a WHERE a.symbol = 'ITUB4'), (
     SELECT u.uid FROM users as u WHERE u.uid = 'TestAdminID'), (
     SELECT b.id FROM brokerages as b WHERE b.name = 'Avenue'), -2, 19.1, 'BRL',
     'sell', '2020-04-01'
    ),
    (
     (SELECT a.id FROM assets as a WHERE a.symbol = 'ITUB4'), (
     SELECT u.uid FROM users as u WHERE u.uid = 'TestAdminID'), (
     SELECT b.id FROM brokerages as b WHERE b.name = 'Avenue'), 2, 20.05, 'BRL',
     'buy', '2020-04-20'
    ),
    (
     (SELECT a.id FROM assets as a WHERE a.symbol = 'ITUB4'), (
     SELECT u.uid FROM users as u WHERE u.uid = 'TestAdminID'), (
     SELECT b.id FROM brokerages as b WHERE b.name = 'Avenue'), 8, 22.20, 'BRL',
     'buy', '2021-08-10'
    ),
    (
     (SELECT a.id FROM assets as a WHERE a.symbol = 'ITUB4'), (
     SELECT u.uid FROM users as u WHERE u.uid = 'TestAdminID'), (
     SELECT b.id FROM brokerages as b WHERE b.name = 'Avenue'), 12, 25.58, 'BRL',
     'buy', '2021-09-13'
    );

