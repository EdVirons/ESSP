-- +goose Up
-- SSOT reference tables for Kenya counties and sub-counties

CREATE TABLE IF NOT EXISTS ssot_counties (
    code TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ssot_sub_counties (
    code TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    county_code TEXT NOT NULL REFERENCES ssot_counties(code),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(county_code, name)
);

CREATE INDEX IF NOT EXISTS idx_ssot_sub_counties_county ON ssot_sub_counties(county_code);
CREATE INDEX IF NOT EXISTS idx_ssot_counties_name ON ssot_counties(name);
CREATE INDEX IF NOT EXISTS idx_ssot_sub_counties_name ON ssot_sub_counties(name);

-- Seed Kenya counties (47 counties)
INSERT INTO ssot_counties (code, name) VALUES
('001', 'Mombasa'),
('002', 'Kwale'),
('003', 'Kilifi'),
('004', 'Tana River'),
('005', 'Lamu'),
('006', 'Taita Taveta'),
('007', 'Garissa'),
('008', 'Wajir'),
('009', 'Mandera'),
('010', 'Marsabit'),
('011', 'Isiolo'),
('012', 'Meru'),
('013', 'Tharaka Nithi'),
('014', 'Embu'),
('015', 'Kitui'),
('016', 'Machakos'),
('017', 'Makueni'),
('018', 'Nyandarua'),
('019', 'Nyeri'),
('020', 'Kirinyaga'),
('021', 'Muranga'),
('022', 'Kiambu'),
('023', 'Turkana'),
('024', 'West Pokot'),
('025', 'Samburu'),
('026', 'Trans Nzoia'),
('027', 'Uasin Gishu'),
('028', 'Elgeyo Marakwet'),
('029', 'Nandi'),
('030', 'Baringo'),
('031', 'Laikipia'),
('032', 'Nakuru'),
('033', 'Narok'),
('034', 'Kajiado'),
('035', 'Kericho'),
('036', 'Bomet'),
('037', 'Kakamega'),
('038', 'Vihiga'),
('039', 'Bungoma'),
('040', 'Busia'),
('041', 'Siaya'),
('042', 'Kisumu'),
('043', 'Homa Bay'),
('044', 'Migori'),
('045', 'Kisii'),
('046', 'Nyamira'),
('047', 'Nairobi')
ON CONFLICT (code) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS ssot_sub_counties;
DROP TABLE IF EXISTS ssot_counties;
