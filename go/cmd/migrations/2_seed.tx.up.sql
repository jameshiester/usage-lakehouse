-- SQL Seed Data Migration
-- This file contains all the seed data that was previously in 2_seed.go

-- Generate UUIDs for consistent references
DO $$
DECLARE
    -- Power Region UUIDs
    ercot_id UUID := gen_random_uuid();
    pjm_id UUID := gen_random_uuid();
    miso_id UUID := gen_random_uuid();
    caiso_id UUID := gen_random_uuid();
    nyiso_id UUID := gen_random_uuid();
    isone_id UUID := gen_random_uuid();
    spp_id UUID := gen_random_uuid();
    
    -- TDSP UUIDs
    oncor_id UUID := gen_random_uuid();
    aep_north_id UUID := gen_random_uuid();
    aep_central_id UUID := gen_random_uuid();
    centerpoint_id UUID := gen_random_uuid();
    tnmp_id UUID := gen_random_uuid();
    
    -- Account UUIDs
    faber_id UUID := gen_random_uuid();
    acme_id UUID := gen_random_uuid();
    
    -- Premise and Meter UUIDs
    premise1_id UUID := gen_random_uuid();
    meter1_id UUID := gen_random_uuid();
    
    -- Start time for account junctions
    start_time TIMESTAMP := (CURRENT_DATE - INTERVAL '4 years')::TIMESTAMP;
BEGIN
    -- Insert premise account statuses
    INSERT INTO public.premise_account_status (code, name, is_active, created_dttm, updated_dttm) VALUES
        ('ACTIVE', 'Active', true, NOW(), NOW()),
        ('ENROLL_REQUESTED', 'Enrollment Requested', false, NOW(), NOW()),
        ('ENROLL_REJECTED', 'Enrollment Requested', false, NOW(), NOW()),
        ('PENDING_ENROLLMENT', 'Pending Enrollment', false, NOW(), NOW()),
        ('DELETE_REQUESTED', 'Delete Requested', true, NOW(), NOW()),
        ('PENDING DELETE', 'Pending Delete', true, NOW(), NOW()),
        ('DELETED', 'Deleted', false, NOW(), NOW());
    
    -- Insert power regions
    INSERT INTO public.power_region (id, name, created_dttm, updated_dttm) VALUES
        (ercot_id, 'ERCOT', NOW(), NOW()),
        (pjm_id, 'PJM', NOW(), NOW()),
        (miso_id, 'MISO', NOW(), NOW()),
        (caiso_id, 'CAISO', NOW(), NOW()),
        (nyiso_id, 'NYISO', NOW(), NOW()),
        (isone_id, 'ISONE', NOW(), NOW()),
        (spp_id, 'SPP', NOW(), NOW());
    
    -- Insert TDSPS
    INSERT INTO public.tdsp (id, legal_entity_name, name, legal_id, abbreviation, premise_code_validation_expression, created_dttm, updated_dttm) VALUES
        (oncor_id, 'Oncor Electric Delivery Company LLC', 'Oncor Electric Delivery Company LLC', '1039940674000', 'ONCOR', '^(1044372|1017699)\d{10}$', NOW(), NOW()),
        (aep_central_id, 'AEP Texas Central Company', 'AEP Texas Central Company', '007924772', 'AEP-C', '^1000288\d{15}$', NOW(), NOW()),
        (aep_north_id, 'AEP Texas North Company', 'AEP Texas North Company', '007923311', 'AEP-N', '^1000078\d{15}$', NOW(), NOW()),
        (centerpoint_id, 'CenterPoint Energy Houston Electric LLC', 'CenterPoint Energy Houston Electric LLC', '957877905', 'CENTERPOINT', '^10089\d{17}$', NOW(), NOW()),
        (tnmp_id, 'Texas-New Mexico Power Company', 'Texas-New Mexico Power Company', '007929441', 'TNMP', '^10168\d{17}$', NOW(), NOW());
    
    -- Insert TDSP-Power Region mappings
    INSERT INTO public.tdsp_power_region_junction (tdsp_id, power_region_id, created_dttm, updated_dttm) VALUES
        (oncor_id, ercot_id, NOW(), NOW()),
        (aep_north_id, ercot_id, NOW(), NOW()),
        (aep_central_id, ercot_id, NOW(), NOW()),
        (centerpoint_id, ercot_id, NOW(), NOW()),
        (tnmp_id, ercot_id, NOW(), NOW());
    
    -- Insert transaction types
    INSERT INTO public.transaction_type (code, name, description, created_dttm, updated_dttm) VALUES
        ('USAGE', 'Usage', 'Historic or monthly usage transaction', NOW(), NOW()),
        ('ENROLLMENT', 'Enrollment', 'Enrollment requests and responses', NOW(), NOW());
    
    -- Insert power region transaction types
    INSERT INTO public.power_region_transaction_type (code, name, power_region_id, transaction_type_code, description, created_dttm, updated_dttm) VALUES
        ('867', 'Usage', ercot_id, 'USAGE', 'Historic or monthly usage transaction', NOW(), NOW());
    
    -- Insert transaction sub types
    INSERT INTO public.transaction_sub_type (code, name, transaction_type_code, description, created_dttm, updated_dttm) VALUES
        ('UH', 'Historic Usage', 'USAGE', 'Historic usage transaction', NOW(), NOW()),
        ('UM', 'Monthly Usage', 'USAGE', 'Monthly usage transaction', NOW(), NOW());
    
    -- Insert power region transaction sub types
    INSERT INTO public.power_region_transaction_sub_type (code, name, power_region_id, transaction_sub_type_code, description, created_dttm, updated_dttm) VALUES
        ('2', 'Historic Usage', ercot_id, 'UH', 'Historic Usage transmitted from ERCOT to CR', NOW(), NOW()),
        ('3', 'Monthly Usage', ercot_id, 'UM', 'Monthly or Final Usage transmitted from ERCOT to CR', NOW(), NOW());
    
    -- Insert usage transaction purposes
    INSERT INTO public.usage_transaction_purpose (code, name, is_cancel, description, created_dttm, updated_dttm) VALUES
        ('N', 'New', false, 'New usage transaction', NOW(), NOW()),
        ('C', 'Cancel', true, 'Canceled usage transaction', NOW(), NOW()),
        ('R', 'Replace', true, 'Used when the TDSP cancels and sends a replacement transaction for corrected data', NOW(), NOW());
    
    -- Insert power region usage transaction purposes
    INSERT INTO public.power_region_usage_transaction_purpose (code, name, power_region_id, usage_transaction_purpose_code, description, created_dttm, updated_dttm) VALUES
        ('00', 'Original', ercot_id, 'N', 'Conveys original readings for the account being reported.', NOW(), NOW()),
        ('01', 'Cancellation', ercot_id, 'C', 'Readings previously reported for the account are to be ignored. This would cancel the entire period of usage for the period.', NOW(), NOW()),
        ('02', 'Replace', ercot_id, 'R', 'Used when the TDSP cancels and sends a replacement transaction for corrected data.', NOW(), NOW());
    
    -- Insert premise types
    INSERT INTO public.premise_type (code, name, description, created_dttm, updated_dttm) VALUES
        ('RESIDENTIAL', 'RESIDENTIAL', 'Residential premise', NOW(), NOW()),
        ('COMMERCIAL', 'COMMERCIAL', 'Commercial premise', NOW(), NOW());
    
    -- Insert accounts
    INSERT INTO public.account (id, legal_id, name, created_dttm, updated_dttm) VALUES
        (acme_id, '1234567', 'Acme Corporation', NOW(), NOW()),
        (faber_id, '7654321', 'Faber LLC', NOW(), NOW());
    
    -- Insert premises
    INSERT INTO public.premise (id, power_region_id, premise_type_code, name, code, customer_name, address_line_1, city, state, zip, country, created_dttm, updated_dttm) VALUES
        (premise1_id, ercot_id, 'RESIDENTIAL', '10443720008808467', '10443720008808467', 'John Smith', '04914 BAYONNE DR', 'ROWLETT', 'TX', '750881851', 'USA', NOW(), NOW());
    
    -- Insert meters
    INSERT INTO public.meter (id, premise_id, power_region_id, name, type, load_profile, cycle_code, created_dttm, updated_dttm) VALUES
        (meter1_id, premise1_id, ercot_id, 'LG12345', 'IDR', '', '13', NOW(), NOW());
    
    -- Insert premise account junctions
    INSERT INTO public.premise_account_junction (account_id, premise_id, premise_account_status_code, min_start_dt, created_dttm, updated_dttm) VALUES
        (acme_id, premise1_id, 'ACTIVE', start_time, NOW(), NOW());
    
    -- Insert premise account history
    INSERT INTO public.premise_account_history (account_id, premise_id, estimated_start_dt, start_dt, created_dttm, updated_dttm) VALUES
        (acme_id, premise1_id, start_time, start_time, NOW(), NOW());
    
    -- Insert power region usage transaction product transfer detail types
    INSERT INTO public.power_region_usage_transaction_product_transfer_detail_type (code, name, power_region_id, is_interval, is_meter, is_summary, description, created_dttm, updated_dttm) VALUES
        ('PL', 'Non-Interval Detail', ercot_id, false, true, false, 'Non-Interval Detail', NOW(), NOW()),
        ('SU', 'Non-Interval Usage Summary', ercot_id, false, true, true, 'Non-Interval Usage Summary', NOW(), NOW()),
        ('BD', 'Unmetered Services Detail', ercot_id, false, false, false, 'Unmetered Services Detail', NOW(), NOW()),
        ('BO', 'Interval Summary', ercot_id, true, true, true, 'Interval Summary', NOW(), NOW()),
        ('IA', 'Net Interval Usage Summary', ercot_id, true, true, true, 'Net Interval Usage Summary', NOW(), NOW()),
        ('PM', 'Interval Detail', ercot_id, true, true, false, 'Interval Detail', NOW(), NOW()),
        ('PP', 'Net Interval Usage Summary Across Meters', ercot_id, true, true, true, 'Net Interval Usage Summary Across Meters', NOW(), NOW());
END $$; 