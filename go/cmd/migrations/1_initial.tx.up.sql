CREATE SCHEMA partman;
CREATE EXTENSION pg_partman WITH SCHEMA partman;

CREATE OR REPLACE FUNCTION public.set_created_dttm()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_dttm := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION public.set_updated_dttm()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_dttm := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS public.account (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	legal_id VARCHAR(32) NOT NULL UNIQUE,
    created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL
);

CREATE TRIGGER set_created_dttm_trigger
BEFORE INSERT ON public.account
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.account
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

CREATE TABLE IF NOT EXISTS public.power_region (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL
);

CREATE TRIGGER set_created_dttm_trigger
BEFORE INSERT ON public.power_region
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.power_region
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

CREATE TABLE IF NOT EXISTS public.premise_type (
	code VARCHAR(64) PRIMARY KEY,
	name VARCHAR(64) UNIQUE NOT NULL,
	description VARCHAR(300) NOT NULL,
	created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL
);

CREATE TRIGGER set_created_dttm_trigger
BEFORE INSERT ON public.premise_type
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.premise_type
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

CREATE TABLE IF NOT EXISTS public.premise (
	id UUID PRIMARY KEY,
	power_region_id UUID NOT NULL,
	premise_type_code VARCHAR(64) NOT NULL,
	name VARCHAR(200),
	code VARCHAR(128) NOT NULL,
	customer_name VARCHAR(200) NOT NULL,
	created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL,
    address_line_1 VARCHAR(200),
    address_line_2 VARCHAR(200),
    city VARCHAR(200),
    state VARCHAR(200),
    zip VARCHAR(200),
    country VARCHAR(200),
    CONSTRAINT fk_power_region_id
        FOREIGN KEY(power_region_id)
        REFERENCES public.power_region(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_premise_type_code
        FOREIGN KEY(premise_type_code)
        REFERENCES public.premise_type(code),
	CONSTRAINT unique_code_power_region 
        UNIQUE (code, power_region_id)
);

CREATE TRIGGER set_created_dttm_trigger
BEFORE INSERT ON public.premise
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.premise
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

CREATE TABLE IF NOT EXISTS public.meter (
	id UUID PRIMARY KEY,
	premise_id UUID,
	power_region_id UUID NOT NULL,
	name VARCHAR(64) NOT NULL,
	type VARCHAR(32),
	load_profile VARCHAR(32),
	cycle_code VARCHAR(32),
	created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL,
	is_active BOOLEAN NOT NULL DEFAULT TRUE,
    CONSTRAINT fk_premise_id
        FOREIGN KEY(premise_id)
        REFERENCES public.premise(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_power_region_id
        FOREIGN KEY(power_region_id)
        REFERENCES public.power_region(id)
        ON DELETE CASCADE,
    CONSTRAINT unique_power_region_name
        UNIQUE (name, power_region_id)
);

CREATE TRIGGER set_created_dttm_trigger
BEFORE INSERT ON public.meter
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.meter
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

CREATE TABLE IF NOT EXISTS public.asset (
	id UUID PRIMARY KEY,
	name VARCHAR(200) NOT NULL,
	asset_code VARCHAR(200) NOT NULL,
	account_id UUID NOT NULL,
	meter_id UUID,
	premise_id UUID NOT NULL,
	created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL,
    CONSTRAINT fk_account_id
        FOREIGN KEY(account_id)
        REFERENCES public.account(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_premise_id
        FOREIGN KEY(premise_id)
        REFERENCES public.premise(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_meter_id
        FOREIGN KEY(meter_id)
        REFERENCES public.meter(id)
        ON DELETE CASCADE
);

CREATE TRIGGER asset_set_created_dttm_trigger
BEFORE INSERT ON public.asset
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER asset_set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.asset
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

CREATE TABLE IF NOT EXISTS public.premise_account_status (
	code VARCHAR(64) PRIMARY KEY,
	name VARCHAR(32) NOT NULL,
    is_active BOOLEAN NOT NULL,
	created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL
);

CREATE TRIGGER set_created_dttm_trigger
BEFORE INSERT ON public.premise_account_status
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.premise_account_status
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

CREATE TABLE IF NOT EXISTS public.premise_account_junction (
	account_id UUID NOT NULL,
	premise_id UUID NOT NULL,
	premise_account_status_code VARCHAR(64) NOT NULL,
	min_start_dt DATE NOT NULL,
	max_end_dt DATE,
	created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL,
    CONSTRAINT fk_account_id
        FOREIGN KEY(account_id)
        REFERENCES public.account(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_premise_account_status_code
        FOREIGN KEY(premise_account_status_code)
        REFERENCES public.premise_account_status(code),
    CONSTRAINT fk_premise_id
        FOREIGN KEY(premise_id)
        REFERENCES public.premise(id)
        ON DELETE CASCADE,
    CONSTRAINT pk_premise_account_junction
        PRIMARY KEY (account_id, premise_id)
);

CREATE TRIGGER set_created_dttm_trigger
BEFORE INSERT ON public.premise_account_junction
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.premise_account_junction
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

CREATE TABLE IF NOT EXISTS public.premise_account_history (
	account_id UUID NOT NULL,
	premise_id UUID NOT NULL,
	estimated_start_dt DATE NOT NULL,
	estimated_end_dt DATE,
	start_dt DATE,
	end_dt DATE,
	created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL,
    CONSTRAINT fk_account_id
        FOREIGN KEY(account_id)
        REFERENCES public.account(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_premise_id
        FOREIGN KEY(premise_id)
        REFERENCES public.premise(id)
        ON DELETE CASCADE,
    CONSTRAINT pk_premise_account_history
        PRIMARY KEY (account_id, premise_id)
);

CREATE TRIGGER set_created_dttm_trigger
BEFORE INSERT ON public.premise_account_history
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.premise_account_history
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

CREATE TABLE IF NOT EXISTS public.tdsp (
	id UUID PRIMARY KEY,
	legal_entity_name TEXT NOT NULL UNIQUE,
	name VARCHAR(200) NOT NULL UNIQUE,
	legal_id VARCHAR(32) NOT NULL UNIQUE,
	abbreviation VARCHAR(64) NOT NULL UNIQUE,
	premise_code_validation_expression TEXT,
	created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL
);

CREATE TRIGGER set_created_dttm_trigger
BEFORE INSERT ON public.tdsp
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.tdsp
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

CREATE TABLE IF NOT EXISTS public.tdsp_power_region_junction (
	tdsp_id UUID NOT NULL,
	power_region_id UUID NOT NULL,
	created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL,
    CONSTRAINT fk_tdsp_id
        FOREIGN KEY(tdsp_id)
        REFERENCES public.tdsp(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_power_region_id
        FOREIGN KEY(power_region_id)
        REFERENCES public.power_region(id)
        ON DELETE CASCADE,
    CONSTRAINT pk_tdsp_power_region_junction
        PRIMARY KEY (tdsp_id, power_region_id)
);

CREATE TRIGGER set_created_dttm_trigger
BEFORE INSERT ON public.tdsp_power_region_junction
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.tdsp_power_region_junction
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

CREATE TABLE IF NOT EXISTS public.usage_transaction_purpose (
	code VARCHAR(64) PRIMARY KEY,
	name VARCHAR(64) UNIQUE NOT NULL,
    is_cancel BOOLEAN NOT NULL DEFAULT FALSE,
	description VARCHAR(300) NOT NULL,
	created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL
);

CREATE TRIGGER set_created_dttm_trigger
BEFORE INSERT ON public.usage_transaction_purpose
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.usage_transaction_purpose
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

CREATE TABLE IF NOT EXISTS public.power_region_usage_transaction_purpose (
	code VARCHAR(64) PRIMARY KEY,
    usage_transaction_purpose_code VARCHAR(64) NOT NULL,
    power_region_id UUID NOT NULL,
	name VARCHAR(64) UNIQUE NOT NULL,
    is_cancel BOOLEAN NOT NULL DEFAULT FALSE,
	description VARCHAR(300) NOT NULL,
	created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL
);

CREATE TRIGGER set_created_dttm_trigger
BEFORE INSERT ON public.power_region_usage_transaction_purpose
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.power_region_usage_transaction_purpose
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

CREATE TABLE IF NOT EXISTS public.transaction_type (
	code VARCHAR(64) PRIMARY KEY,
	name VARCHAR(64) UNIQUE NOT NULL,
	description VARCHAR(300) NOT NULL,
	created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL
);

CREATE TRIGGER set_created_dttm_trigger
BEFORE INSERT ON public.transaction_type
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.transaction_type
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

CREATE TABLE IF NOT EXISTS public.transaction_sub_type (
	code VARCHAR(64) PRIMARY KEY,
    transaction_type_code VARCHAR(64) NOT NULL,
	name VARCHAR(64) UNIQUE NOT NULL,
	description VARCHAR(300) NOT NULL,
	created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL,
    CONSTRAINT fk_transaction_type_code
        FOREIGN KEY(transaction_type_code)
        REFERENCES public.transaction_type(code)
        ON DELETE CASCADE
);

CREATE TRIGGER set_created_dttm_trigger
BEFORE INSERT ON public.transaction_sub_type
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.transaction_sub_type
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

CREATE TABLE IF NOT EXISTS public.power_region_transaction_type (
    transaction_type_code VARCHAR(64) NOT NULL,
    power_region_id UUID NOT NULL,
	name VARCHAR(64) NOT NULL,
	code VARCHAR(64) NOT NULL,
	description VARCHAR(300) NOT NULL,
	created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL,
    CONSTRAINT fk_transaction_type_code
        FOREIGN KEY(transaction_type_code)
        REFERENCES public.transaction_type(code)
        ON DELETE CASCADE,
    CONSTRAINT fk_power_region_id
        FOREIGN KEY(power_region_id)
        REFERENCES public.power_region(id)
        ON DELETE CASCADE,
    CONSTRAINT pk_power_region_transaction_type
        PRIMARY KEY (name, power_region_id),
    CONSTRAINT unique_power_region_transaction_type_code_power_region 
        UNIQUE (code, power_region_id)
);

CREATE TRIGGER set_created_dttm_trigger
BEFORE INSERT ON public.power_region_transaction_type
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.power_region_transaction_type
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

CREATE TABLE IF NOT EXISTS public.power_region_usage_transaction_product_transfer_detail_type (
    power_region_id UUID NOT NULL,
	code VARCHAR(64) NOT NULL,
    name VARCHAR(64) not null,
	description VARCHAR(300) NOT NULL,
    is_interval BOOLEAN NOT NULL,
    is_summary BOOLEAN NOT NULL,
    is_meter BOOLEAN NOT NULL,
	created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL,
    CONSTRAINT fk_power_region_id
        FOREIGN KEY(power_region_id)
        REFERENCES public.power_region(id)
        ON DELETE CASCADE,
    CONSTRAINT pk_power_region_usage_transaction_product_transfer_detail_type
        PRIMARY KEY (code, power_region_id),
    CONSTRAINT unique_power_region_usage_transaction_product_transfer_detail_type_name
        UNIQUE (name, power_region_id)
);

CREATE TRIGGER set_created_dttm_trigger
BEFORE INSERT ON public.power_region_usage_transaction_product_transfer_detail_type
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.power_region_usage_transaction_product_transfer_detail_type
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

CREATE TABLE IF NOT EXISTS public.power_region_transaction_sub_type (
    transaction_sub_type_code VARCHAR(64) NOT NULL,
    power_region_id UUID NOT NULL,
	name VARCHAR(64) NOT NULL,
	code VARCHAR(64) NOT NULL,
	description VARCHAR(300) NOT NULL,
	created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL,
    CONSTRAINT fk_transaction_sub_type_code
        FOREIGN KEY(transaction_sub_type_code)
        REFERENCES public.transaction_sub_type(code)
        ON DELETE CASCADE,
    CONSTRAINT fk_power_region_id
        FOREIGN KEY(power_region_id)
        REFERENCES public.power_region(id)
        ON DELETE CASCADE,
    CONSTRAINT pk_power_region_transaction_sub_type
        PRIMARY KEY (name, power_region_id),
    CONSTRAINT unique_power_region_transaction_sub_type_code_power_region 
        UNIQUE (code, power_region_id)
);

CREATE TRIGGER set_created_dttm_trigger
BEFORE INSERT ON public.power_region_transaction_sub_type
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.power_region_transaction_sub_type
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

CREATE TABLE IF NOT EXISTS public.usage_transaction (
	id UUID PRIMARY KEY,
	transaction_id VARCHAR(32) NOT NULL UNIQUE,
    transaction_type_code VARCHAR(64) NOT NULL,
    transaction_sub_type_code VARCHAR(64) NOT NULL,
    transaction_dt DATE NOT NULL,
    service_period_start_dt DATE NOT NULL,
    service_period_end_dt DATE NOT NULL,
    is_final BOOLEAN NOT NULL DEFAULT FALSE,
    is_canceled BOOLEAN NOT NULL DEFAULT FALSE,
    usage_transaction_purpose_code VARCHAR(64) NOT NULL,
	premise_id UUID NOT NULL,
    tdsp_id UUID NOT NULL,
    power_region_id UUID NOT NULL,
	created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL,
    CONSTRAINT fk_premise_id
        FOREIGN KEY(premise_id)
        REFERENCES public.premise(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_power_region_id
        FOREIGN KEY(power_region_id)
        REFERENCES public.power_region(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_tdsp_id
        FOREIGN KEY(tdsp_id)
        REFERENCES public.tdsp(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_usage_transaction_purpose_code
        FOREIGN KEY(usage_transaction_purpose_code)
        REFERENCES public.usage_transaction_purpose(code)
        ON DELETE CASCADE,
    CONSTRAINT fk_transaction_type_code
        FOREIGN KEY(transaction_type_code)
        REFERENCES public.transaction_type(code)
        ON DELETE CASCADE,
    CONSTRAINT fk_transaction_sub_type_code
        FOREIGN KEY(transaction_sub_type_code)
        REFERENCES public.transaction_sub_type(code)
        ON DELETE CASCADE
);

CREATE TRIGGER set_created_dttm_trigger
BEFORE INSERT ON public.usage_transaction
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.usage_transaction
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

CREATE TABLE IF NOT EXISTS public.usage_transaction_detail (
	start_dttm TIMESTAMP NOT NULL,
	end_dttm TIMESTAMP NOT NULL,
    service_period_start_dt DATE NOT NULL,
    service_period_end_dt DATE NOT NULL,
    usage_transaction_id UUID NOT NULL,
    is_canceled BOOLEAN NOT NULL DEFAULT FALSE,
	premise_id UUID NOT NULL,
    power_region_id UUID NOT NULL,
	meter_id UUID,
	meter_name VARCHAR(64) NOT NULL,
	consumption DECIMAL(10,5),
	generation DECIMAL(10,5),
	created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL,
    CONSTRAINT fk_premise_id
        FOREIGN KEY(premise_id)
        REFERENCES public.premise(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_power_region_id
        FOREIGN KEY(power_region_id)
        REFERENCES public.power_region(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_usage_transaction_id
        FOREIGN KEY(usage_transaction_id)
        REFERENCES public.usage_transaction(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_meter_id
        FOREIGN KEY(meter_id)
        REFERENCES public.meter(id)
        ON DELETE CASCADE,
    CONSTRAINT pk_usage_transaction_detail
        PRIMARY KEY (start_dttm,usage_transaction_id, meter_name)
) PARTITION BY RANGE (start_dttm);

CREATE TRIGGER set_created_dttm_trigger
BEFORE INSERT ON public.usage_transaction_detail
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.usage_transaction_detail
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

SELECT partman.create_parent( 
 p_parent_table => 'public.usage_transaction_detail',
 p_control      => 'start_dttm',
 p_type         => 'range',
 p_interval     => '1 month',
 p_start_partition      => to_char(date_trunc('year', current_date) - interval '3 years', 'YYYY-MM-DD'),
 p_premake      => 12);

UPDATE partman.part_config
SET infinite_time_partitions = true,
retention = '36 month',
retention_keep_table=false
WHERE parent_table = 'public.usage_transaction_detail';

CREATE TABLE IF NOT EXISTS public.meter_usage_15_minute (
	start_dttm TIMESTAMP NOT NULL,
	end_dttm TIMESTAMP NOT NULL,
    service_period_start_dt DATE NOT NULL,
    service_period_end_dt DATE NOT NULL,
    is_canceled BOOLEAN NOT NULL DEFAULT FALSE,
	premise_id UUID NOT NULL,
	meter_id UUID NOT NULL,
	account_id UUID NOT NULL,
	consumption DECIMAL(10,5),
	generation DECIMAL(10,5),
	created_dttm timestamp NOT NULL,
	updated_dttm timestamp NOT NULL,
    CONSTRAINT fk_premise_id
        FOREIGN KEY(premise_id)
        REFERENCES public.premise(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_account_id
        FOREIGN KEY(account_id)
        REFERENCES public.account(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_meter_id
        FOREIGN KEY(meter_id)
        REFERENCES public.meter(id)
        ON DELETE CASCADE,
    CONSTRAINT pk_meter_usage_15_minute
        PRIMARY KEY (start_dttm, meter_id)
) PARTITION BY RANGE (start_dttm);

CREATE TRIGGER set_created_dttm_trigger
BEFORE INSERT ON public.meter_usage_15_minute
FOR EACH ROW
EXECUTE FUNCTION public.set_created_dttm();

CREATE TRIGGER set_updated_dttm_trigger
BEFORE INSERT OR UPDATE  ON public.meter_usage_15_minute
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_dttm();

SELECT partman.create_parent( 
 p_parent_table => 'public.meter_usage_15_minute',
 p_control      => 'start_dttm',
 p_type         => 'range',
 p_interval     => '1 month',
 p_start_partition      => to_char(date_trunc('year', current_date) - interval '3 years', 'YYYY-MM-DD'),
 p_premake      => 12);

UPDATE partman.part_config
SET infinite_time_partitions = true,
retention = '36 month',
retention_keep_table=false
WHERE parent_table = 'public.meter_usage_15_minute';