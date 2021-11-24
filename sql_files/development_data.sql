
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
	('Finances');

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
    );

