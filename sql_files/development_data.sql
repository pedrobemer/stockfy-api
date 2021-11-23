
\connect stockfy;

-- Insert user with admin privileges
INSERT INTO
	public.users ("uid", "username", email, "type")
VALUES
	('TestAdminID', 'Test Name Admin', 'test_admin@email.com', 'admin');