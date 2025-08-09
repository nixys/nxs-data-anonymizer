# NXS-DATA-ANONYMIZER - Benchmark

Benchmark comparatif entre l'anonymization avec **faker intégré** vs **fonctions natives** pour PostgreSQL et MySQL.

## 🎯 Objectif

Comparer les performances et la qualité des données anonymisées entre :
- **Faker intégré** : Données réalistes françaises (noms, emails, téléphones)
- **Fonctions natives** : Données aléatoires alphanumériques

## 📊 Configurations testées

| Configuration | Base | Méthode | Description |
|---------------|------|---------|-------------|
| `postgres-faker.conf` | PostgreSQL | Faker | Données réalistes françaises |
| `postgres-native.conf` | PostgreSQL | Native | Données aléatoires |
| `mysql-faker.conf` | MySQL | Faker | Données réalistes françaises |
| `mysql-native.conf` | MySQL | Native | Données aléatoires |

## 🔧 Prérequis

- Docker (PostgreSQL + MySQL)
- Go (compilation des binaires)
- Python 3 (population des bases)

## 🚀 Utilisation

### Setup complet
```bash
task setup        # Démarre Docker, initialise les BDD, popule 1000 users
```

### Tests rapides
```bash
task test-pg-faker     # Test PostgreSQL + Faker
task test-pg-native    # Test PostgreSQL + Native  
task test-mysql-faker  # Test MySQL + Faker
task test-mysql-native # Test MySQL + Native
```

### Benchmark complet
```bash
task benchmark    # Mesure les performances des 4 configs
task verify       # Vérifie l'anonymisation (avant/après)
```

### Workflow complet
```bash
task run-all      # Setup + Benchmark + Vérification
```

## 📈 Résultats attendus

**Performance :** 
- PostgreSQL généralement plus rapide que MySQL
- Différence mineure entre faker et native

**Qualité des données :**
- **Faker :** `Marc Delahaye` → `Marie Dubois`, `test@example.com` → `marie.dubois@hotmail.fr`
- **Native :** `Marc Delahaye` → `XyZ9kL AbC3mN`, `test@example.com` → `Kj8mP@example.com`

## 🔍 Détails techniques

### Binaires utilisés
- `nxs-data-anonymizer` : Binaire natif standard
- `nxs-data-anonymizer-faker` : Binaire avec faker intégré (template.go modifié)

### Différence MySQL vs PostgreSQL
- **PostgreSQL :** Fonctionne avec dumps partiels (data-only)
- **MySQL :** Nécessite dumps complets avec headers (`--extended-insert=false`)

### Modifications apportées
- **Template modifié :** `misc/template.go` avec intégration `gofakeit`
- **Fonctions ajoutées :** `fakerEmailFR`, `fakerPhoneFR`, `fakerBirthDate`, etc.

## 🧹 Nettoyage

```bash
task clean        # Arrête Docker + supprime tous les fichiers temporaires
```

## 🛠️ Structure

```
benchmark/
├── postgres-faker.conf      # Config PostgreSQL + Faker
├── postgres-native.conf     # Config PostgreSQL + Native  
├── mysql-faker.conf         # Config MySQL + Faker
├── mysql-native.conf        # Config MySQL + Native
├── scripts/
│   ├── test_benchmark.sh    # Mesure des performances
│   ├── verify_anonymization.sh # Vérification avant/après
│   ├── populate_db.py       # Population des BDD
│   └── wait-for-db.sh       # Attente Docker
├── sql/init.sql             # Schéma de la table users
└── docker-compose.yml       # PostgreSQL + MySQL
```