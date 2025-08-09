-- Table avec 6 types de champs différents
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),           -- Type 1: String (nom)
    email VARCHAR(150),          -- Type 2: String (email)
    phone VARCHAR(20),           -- Type 3: String (téléphone)
    birth_date DATE,             -- Type 4: Date
    salary DECIMAL(10,2),        -- Type 5: Numérique
    description TEXT             -- Type 6: Texte long
);

-- Insertion de 18 lignes de données
INSERT INTO users (name, email, phone, birth_date, salary, description) VALUES
('Jean Dupont', 'jean.dupont@example.com', '+33612345678', '1985-03-15', 45000.50, 'Développeur senior avec 10 ans d''expérience en Java et Python'),
('Marie Martin', 'marie.martin@company.fr', '+33698765432', '1990-07-22', 38000.00, 'Chef de projet certifiée PMP, spécialisée en méthodes agiles'),
('Pierre Bernard', 'p.bernard@email.com', '+33611223344', '1978-11-08', 52000.75, 'Architecte solutions cloud AWS et Azure'),
('Sophie Dubois', 'sophie.d@gmail.com', '+33655443322', '1992-05-30', 35000.25, 'Data analyst passionnée par le machine learning'),
('Lucas Moreau', 'l.moreau@hotmail.fr', '+33677889900', '1988-09-12', 42000.00, 'DevOps engineer expert en Kubernetes et Docker'),
('Emma Leroy', 'emma.leroy@pro.com', '+33644556677', '1995-01-25', 33000.50, 'UX/UI designer créative avec portfolio impressionnant'),
('Thomas Roux', 'thomas.r@yahoo.fr', '+33622334455', '1982-06-18', 48000.00, 'Manager d''équipe technique de 15 personnes'),
('Camille Fournier', 'c.fournier@outlook.com', '+33699887766', '1993-12-03', 36500.75, 'Développeuse full-stack React et Node.js'),
('Alexandre Simon', 'alex.simon@tech.fr', '+33611002233', '1987-04-27', 44000.25, 'Expert cybersécurité certifié CISSP'),
('Julie Laurent', 'j.laurent@business.com', '+33666778899', '1991-08-14', 37000.00, 'Business analyst avec expertise en finance'),
('Nicolas Michel', 'n.michel@startup.io', '+33633445566', '1986-02-09', 46000.50, 'Lead développeur blockchain et smart contracts'),
('Charlotte David', 'charlotte.d@corp.fr', '+33688990011', '1994-10-21', 34500.00, 'Scrum master certifiée avec 5 ans d''expérience'),
('Maxime Bertrand', 'max.bertrand@mail.com', '+33655667788', '1989-07-05', 41000.75, 'Ingénieur système Linux et Windows Server'),
('Léa Robert', 'lea.robert@group.eu', '+33622778899', '1996-03-28', 32000.25, 'Développeuse mobile iOS et Android'),
('Antoine Richard', 'a.richard@enterprise.fr', '+33699001122', '1983-11-16', 50000.00, 'Directeur technique avec vision stratégique'),
('Manon Petit', 'manon.p@consulting.com', '+33611223355', '1992-09-07', 38500.50, 'Consultante SAP avec certifications multiples'),
('Hugo Durand', 'h.durand@digital.fr', '+33677445566', '1990-01-13', 40000.00, 'Ingénieur QA automatisation avec Selenium'),
('Chloé Lefevre', 'chloe.l@agency.com', '+33644889900', '1997-06-02', 31000.75, 'Junior développeuse motivée et polyvalente');