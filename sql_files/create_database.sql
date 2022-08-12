-- Create Database for the stockfy and connect to it
CREATE DATABASE stockfy;
\connect stockfy;

-- Create functions extension to generate the UUID values
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create functions to enable timestamp trigger
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create Sectors table
CREATE TABLE public.sectors (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	created_at timestamp without time zone NOT NULL DEFAULT now(),
	updated_at timestamp without time zone NOT NULL DEFAULT now(),
	"name" text NOT NULL,
	CONSTRAINT sectors_pk PRIMARY KEY (id)
);
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON public.sectors
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- Create Asset Types table
CREATE TABLE public.asset_types (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	created_at timestamp without time zone NOT NULL DEFAULT now(),
	updated_at timestamp without time zone NOT NULL DEFAULT now(),
	"type" text NOT NULL,
    "name" text NOT NULL,
	country text NOT NULL,
	CONSTRAINT asset_types_pk PRIMARY KEY (id)
);
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON public.asset_types
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- Create Assets table
CREATE TABLE public.assets (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	created_at timestamp without time zone NOT NULL DEFAULT now(),
	updated_at timestamp without time zone NOT NULL DEFAULT now(),
	symbol text NOT NULL,
    fullname text NOT NULL,
	preference text NOT NULL,
    asset_type_id uuid NOT NULL,
    sector_id uuid NOT NULL,
	CONSTRAINT assets_pk PRIMARY KEY (id),
    CONSTRAINT assets_types_fk FOREIGN KEY (asset_type_id) REFERENCES public.asset_types(id),
    CONSTRAINT sectors_fk FOREIGN KEY (sector_id) REFERENCES public.sectors(id)
);
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON public.assets
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- Creat Users table
CREATE TABLE public.users (
	"uid" text NOT NULL,
	created_at timestamp without time zone NOT NULL DEFAULT now(),
	updated_at timestamp without time zone NOT NULL DEFAULT now(),
	username text NOT NULL,
    email text NOT NULL,
	"type" text NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY ("uid"),
	UNIQUE("uid")
);
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON public.users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- Create Asset Users table
CREATE TABLE public.asset_users (
    asset_id uuid NOT NULL,
    user_uid text NOT NULL,
	created_at timestamp without time zone NOT NULL DEFAULT now(),
	updated_at timestamp without time zone NOT NULL DEFAULT now(),
	CONSTRAINT asset_users_pk PRIMARY KEY (asset_id,user_uid),
    CONSTRAINT asset_users_assets_fk FOREIGN KEY (asset_id) REFERENCES public.assets(id) ON DELETE CASCADE,
    CONSTRAINT asset_users_users_fk FOREIGN KEY (user_uid) REFERENCES public.users("uid") ON DELETE CASCADE,
	UNIQUE(asset_id, user_uid)
);
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON public.asset_users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- Create Brokerages table
CREATE TABLE public.brokerages (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	created_at timestamp without time zone NOT NULL DEFAULT now(),
	updated_at timestamp without time zone NOT NULL DEFAULT now(),
	"name" text NOT NULL,
    fullname text NOT NULL,
	country text NOT NULL,
	CONSTRAINT brokerages_pk PRIMARY KEY (id)
);
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON public.brokerages
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- Create Orders table
CREATE TABLE public.orders (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now(),
	asset_id uuid NOT NULL,
	user_uid text NOT NULL,
	brokerage_id uuid NOT NULL,
	quantity float8 NOT NULL,
	price float8 NOT NULL,
	currency text NOT NULL,
	order_type text NOT NULL,
	"date" date NOT NULL,
	CONSTRAINT orders_pk PRIMARY KEY (id),
	CONSTRAINT orders_brokerage_fk FOREIGN KEY (brokerage_id) REFERENCES public.brokerages(id),
	CONSTRAINT orders_asset_fk FOREIGN KEY (asset_id) REFERENCES public.assets(id) ON DELETE CASCADE,
	CONSTRAINT orders_user_fk FOREIGN KEY (user_uid) REFERENCES public.users("uid") ON DELETE CASCADE
);
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON public.orders
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- Create Earnings table
CREATE TABLE public.earnings (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now(),
	asset_id uuid NOT NULL,
	user_uid text NOT NULL,
	"type" text NOT NULL,
	earning float8 NOT NULL,
	"date" date NOT NULL,
	currency text NOT NULL,
	CONSTRAINT earnings_pk PRIMARY KEY (id),
	CONSTRAINT earnings_asset_fk FOREIGN KEY (asset_id) REFERENCES public.assets(id) ON DELETE CASCADE,
	CONSTRAINT earnings_user_fk FOREIGN KEY (user_uid) REFERENCES public.users("uid") ON DELETE CASCADE
);
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON public.earnings
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


-- Populate database with important datas regarding the asset types
INSERT INTO
	public.asset_types ("type", "name", country)
VALUES
	('ETF', 'ETFs Brasil', 'BR'),
	('ETF', 'ETFs EUA', 'US'),
	('STOCK', 'Ações Brasil', 'BR'),
	('STOCK', 'Ações EUA', 'US'),
	('REIT', 'REITs', 'US'),
	('FII', 'Fundos Imobiliários', 'BR');

-- -- Populate database with initial Brokerage Firms information
INSERT INTO
	public.brokerages ("name", "fullname", country)
VALUES
	('Clear', 'Clear Corretora', 'BR'),
	('Rico', 'Rico Corretora - Grupo XP', 'BR'),
	('Passfolio', 'Passfolio Securities', 'US'),
	('Avenue', 'Avenue Securities', 'US');