-- Create databases for different services
CREATE DATABASE theraclosure_auth;
CREATE DATABASE theraclosure_users;
CREATE DATABASE theraclosure_payments;
CREATE DATABASE theraclosure_core;
CREATE DATABASE theraclosure_geolocation;

-- Create users table in auth database
\c theraclosure_auth;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'THERAPIST',
    subscription_status VARCHAR(20) NOT NULL DEFAULT 'trialing',
    stripe_customer_id VARCHAR(255),
    cognito_id VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT true,
    email_verified BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_cognito_id ON users(cognito_id);
CREATE INDEX idx_users_stripe_customer_id ON users(stripe_customer_id);
CREATE INDEX idx_users_role ON users(role);

-- Create users table in users database (for user management service)
\c theraclosure_users;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(50),
    zip_code VARCHAR(20),
    license_number VARCHAR(100),
    license_state VARCHAR(50),
    license_expiration DATE,
    specializations TEXT[],
    certifications TEXT[],
    practice_name VARCHAR(255),
    practice_type VARCHAR(50),
    years_experience INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_user_profiles_user_id ON user_profiles(user_id);

-- Create enrollment data table
CREATE TABLE IF NOT EXISTS enrollment_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    step VARCHAR(50) NOT NULL,
    data JSONB,
    is_completed BOOLEAN DEFAULT false,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_enrollment_data_user_id ON enrollment_data(user_id);
CREATE INDEX idx_enrollment_data_step ON enrollment_data(step);

-- Create payments tables
\c theraclosure_payments;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    stripe_subscription_id VARCHAR(255) UNIQUE NOT NULL,
    stripe_customer_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    current_period_start TIMESTAMP WITH TIME ZONE,
    current_period_end TIMESTAMP WITH TIME ZONE,
    price_id VARCHAR(255),
    product_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_stripe_subscription_id ON subscriptions(stripe_subscription_id);

CREATE TABLE IF NOT EXISTS payment_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    stripe_event_id VARCHAR(255) UNIQUE NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    subscription_id UUID,
    data JSONB,
    processed BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create templates table in core database
\c theraclosure_core;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    category VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_templates_category ON templates(category);
CREATE INDEX idx_templates_created_by ON templates(created_by);

CREATE TABLE IF NOT EXISTS support_tickets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    subject VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    priority VARCHAR(20) DEFAULT 'medium',
    status VARCHAR(20) DEFAULT 'open',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_support_tickets_user_id ON support_tickets(user_id);
CREATE INDEX idx_support_tickets_status ON support_tickets(status);

-- Insert default admin user
\c theraclosure_auth;
INSERT INTO users (
    id,
    email,
    password_hash,
    first_name,
    last_name,
    role,
    subscription_status,
    is_active,
    email_verified
) VALUES (
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'admin@theraclosure.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', -- password: admin123
    'Admin',
    'User',
    'ADMIN',
    'active',
    true,
    true
) ON CONFLICT (email) DO NOTHING;

-- Insert default templates
\c theraclosure_core;
INSERT INTO templates (title, content, category, created_by) VALUES
(
    'Client Closure Notification',
    'Dear [CLIENT_NAME],

I hope this letter finds you well. I am writing to inform you that I will be closing my practice effective [CLOSURE_DATE].

This decision has not been made lightly, and I want to ensure that your therapeutic journey continues with minimal disruption.

Transition Options:
1. I will be providing referrals to qualified therapists in our area
2. Your records will be maintained as required by law
3. I am available for consultation during the transition period

Next Steps:
- Please schedule a final session to discuss your transition plan
- We will work together to identify the best therapeutic match for your needs
- All necessary documentation will be prepared for transfer

Thank you for allowing me to be part of your therapeutic journey.

Warm regards,

[THERAPIST_NAME]
[LICENSE_NUMBER]',
    'closure',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'
),
(
    'Referral Request Template',
    'Dear [COLLEAGUE_NAME],

I am referring [CLIENT_NAME] to you as I am closing my practice on [CLOSURE_DATE].

Client Background:
- Current concerns: [CONCERNS]
- Treatment duration: [DURATION]
- Progress made: [PROGRESS]
- Recommended approach: [APPROACH]

The client has consented to this referral and is motivated to continue their therapeutic work.

Please let me know if you need any additional information or if you would like to schedule a consultation.

Best regards,

[THERAPIST_NAME]
[CONTACT_INFO]',
    'referral',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'
);

-- Create geolocation database schema
\c theraclosure_geolocation;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";

-- Countries table
CREATE TABLE IF NOT EXISTS countries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(3) UNIQUE NOT NULL,
    code2 VARCHAR(2) UNIQUE NOT NULL,
    region VARCHAR(50),
    currency VARCHAR(10),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- States table
CREATE TABLE IF NOT EXISTS states (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_id UUID NOT NULL REFERENCES countries(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(10),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Cities table
CREATE TABLE IF NOT EXISTS cities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    state_id UUID NOT NULL REFERENCES states(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    zip_code VARCHAR(20),
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for better performance
CREATE INDEX idx_countries_code ON countries(code);
CREATE INDEX idx_countries_code2 ON countries(code2);
CREATE INDEX idx_countries_active ON countries(active);
CREATE INDEX idx_countries_deleted_at ON countries(deleted_at);

CREATE INDEX idx_states_country_id ON states(country_id);
CREATE INDEX idx_states_active ON states(active);
CREATE INDEX idx_states_deleted_at ON states(deleted_at);

CREATE INDEX idx_cities_state_id ON cities(state_id);
CREATE INDEX idx_cities_zip_code ON cities(zip_code);
CREATE INDEX idx_cities_active ON cities(active);
CREATE INDEX idx_cities_deleted_at ON cities(deleted_at);
CREATE INDEX idx_cities_location ON cities(latitude, longitude);
ON CONFLICT DO NOTHING;