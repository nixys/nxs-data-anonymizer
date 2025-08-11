-- Table avec types de champs pour test anonymisation
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),           -- Type 1: String (nom)
    email VARCHAR(150),          -- Type 2: String (email)
    phone VARCHAR(20),           -- Type 3: String (téléphone)
    mobile VARCHAR(20),          -- Type 3b: String (mobile FR)
    birth_date DATE,             -- Type 4: Date
    salary DECIMAL(10,2),        -- Type 5: Numérique
    description TEXT,            -- Type 6: Texte long
    address VARCHAR(200),        -- Type 7: Adresse complète
    city VARCHAR(100),           -- Type 8: Ville
    postal_code VARCHAR(10),     -- Type 9: Code postal
    iban VARCHAR(34),            -- Type 10: IBAN
    ssn VARCHAR(15),             -- Type 11: Numéro sécurité sociale
    vat_number VARCHAR(15),      -- Type 12: Numéro TVA
    ip_address VARCHAR(45)       -- Type 13: Adresse IP
);

-- Insertion de 18 lignes de données avec tous les champs
INSERT INTO users (name, email, phone, mobile, birth_date, salary, description, address, city, postal_code, iban, ssn, vat_number, ip_address) VALUES
('Jean Dupont', 'jean.dupont@example.com', '+33123456789', '+33612345678', '1985-03-15', 45000.50, 'Développeur senior avec 10 ans d''expérience en Java et Python', '15 rue de la Paix', 'Paris', '75001', 'FR7612345678901234567890123', '1850315013456', 'FR12123456789', '192.168.1.10'),
('Marie Martin', 'marie.martin@company.fr', '+33156789012', '+33698765432', '1990-07-22', 38000.00, 'Chef de projet certifiée PMP, spécialisée en méthodes agiles', '42 avenue Victor Hugo', 'Lyon', '69001', 'FR7698765432109876543210987', '2900722043987', 'FR34987654321', '10.0.0.25'),
('Pierre Bernard', 'p.bernard@email.com', '+33145678901', '+33611223344', '1978-11-08', 52000.75, 'Architecte solutions cloud AWS et Azure', '8 boulevard Saint-Michel', 'Marseille', '13001', 'FR7611223344556677889900112', '1781108091234', 'FR56112233445', '172.16.0.5'),
('Sophie Dubois', 'sophie.d@gmail.com', '+33134567890', '+33655443322', '1992-05-30', 35000.25, 'Data analyst passionnée par le machine learning', '23 place de la République', 'Toulouse', '31000', 'FR7655443322998877665544332', '2920530134567', 'FR78554433221', '192.168.0.100'),
('Lucas Moreau', 'l.moreau@hotmail.fr', '+33167890123', '+33677889900', '1988-09-12', 42000.00, 'DevOps engineer expert en Kubernetes et Docker', '67 rue du Commerce', 'Nice', '06000', 'FR7677889900112233445566778', '1880912067890', 'FR90778899001', '10.1.1.50'),
('Emma Leroy', 'emma.leroy@pro.com', '+33178901234', '+33644556677', '1995-01-25', 33000.50, 'UX/UI designer créative avec portfolio impressionnant', '91 avenue des Champs', 'Nantes', '44000', 'FR7644556677223344556677889', '2950125178901', 'FR12445566772', '172.20.0.15'),
('Thomas Roux', 'thomas.r@yahoo.fr', '+33189012345', '+33622334455', '1982-06-18', 48000.00, 'Manager d''équipe technique de 15 personnes', '34 rue Nationale', 'Strasbourg', '67000', 'FR7622334455667788990011223', '1820618189012', 'FR34223344556', '192.168.2.75'),
('Camille Fournier', 'c.fournier@outlook.com', '+33190123456', '+33699887766', '1993-12-03', 36500.75, 'Développeuse full-stack React et Node.js', '56 place Wilson', 'Bordeaux', '33000', 'FR7699887766445566778899001', '2931203190123', 'FR56998877664', '10.2.0.30'),
('Alexandre Simon', 'alex.simon@tech.fr', '+33101234567', '+33611002233', '1987-04-27', 44000.25, 'Expert cybersécurité certifié CISSP', '78 boulevard Voltaire', 'Lille', '59000', 'FR7611002233556677889900112', '1870427101234', 'FR78110022335', '172.25.0.20'),
('Julie Laurent', 'j.laurent@business.com', '+33112345678', '+33666778899', '1991-08-14', 37000.00, 'Business analyst avec expertise en finance', '12 rue Lafayette', 'Rennes', '35000', 'FR7666778899334455667788990', '2910814112345', 'FR90667788992', '192.168.3.40'),
('Nicolas Michel', 'n.michel@startup.io', '+33123456780', '+33633445566', '1986-02-09', 46000.50, 'Lead développeur blockchain et smart contracts', '89 avenue Foch', 'Reims', '51100', 'FR7633445566778899001122334', '1860209123456', 'FR12334455667', '10.3.0.60'),
('Charlotte David', 'charlotte.d@corp.fr', '+33134567801', '+33688990011', '1994-10-21', 34500.00, 'Scrum master certifiée avec 5 ans d''expérience', '45 place Bellecour', 'Le Havre', '76600', 'FR7688990011223344556677889', '2941021134567', 'FR34889900112', '172.30.0.10'),
('Maxime Bertrand', 'max.bertrand@mail.com', '+33145678012', '+33655667788', '1989-07-05', 41000.75, 'Ingénieur système Linux et Windows Server', '67 rue de Rivoli', 'Saint-Étienne', '42000', 'FR7655667788990011223344556', '1890705145678', 'FR56556677889', '192.168.4.85'),
('Léa Robert', 'lea.robert@group.eu', '+33156789023', '+33622778899', '1996-03-28', 32000.25, 'Développeuse mobile iOS et Android', '23 cours Lafayette', 'Toulon', '83000', 'FR7622778899445566778899001', '2960328156789', 'FR78227788990', '10.4.0.95'),
('Antoine Richard', 'a.richard@enterprise.fr', '+33167890134', '+33699001122', '1983-11-16', 50000.00, 'Directeur technique avec vision stratégique', '34 boulevard Haussmann', 'Angers', '49000', 'FR7699001122667788990011223', '1831116167890', 'FR90990011223', '172.35.0.45'),
('Manon Petit', 'manon.p@consulting.com', '+33178901245', '+33611223355', '1992-09-07', 38500.50, 'Consultante SAP avec certifications multiples', '56 avenue Marceau', 'Grenoble', '38000', 'FR7611223355889900112233445', '2920907178901', 'FR12112233556', '192.168.5.120'),
('Hugo Durand', 'h.durand@digital.fr', '+33189012356', '+33677445566', '1990-01-13', 40000.00, 'Ingénieur QA automatisation avec Selenium', '78 rue Saint-Honoré', 'Dijon', '21000', 'FR7677445566112233445566778', '1900113189012', 'FR34774455667', '10.5.0.200'),
('Chloé Lefevre', 'chloe.l@agency.com', '+33190123467', '+33644889900', '1997-06-02', 31000.75, 'Junior développeuse motivée et polyvalente', '90 place Stanislas', 'Nîmes', '30000', 'FR7644889900556677889900112', '2970602190123', 'FR56448899001', '172.40.0.150');