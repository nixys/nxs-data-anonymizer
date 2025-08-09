# Benchmark Faker: Python vs Go vs NXS-Data-Anonymizer

Ce benchmark compare les performances de trois approches pour anonymiser des dumps SQL :
1. **nxs-data-anonymizer** avec configuration command
2. **Python** avec la librairie Faker
3. **Go** avec la librairie gofakeit

## Structure des données

La table `users` contient 6 types de champs :
- `name` : String (nom complet)
- `email` : String (adresse email)
- `phone` : String (numéro de téléphone)
- `birth_date` : Date
- `salary` : Decimal (salaire)
- `description` : Text (description longue)

## Installation

```bash
# Installer Task (si pas déjà installé)
brew install go-task/tap/go-task

# Se placer dans le dossier benchmark
cd benchmark

# Installer les dépendances
task install
```

## Utilisation

### Workflow complet (recommandé)

```bash
# Lance tout : docker, init DB, populate 1000 users, benchmark
task run
```

### Commandes individuelles

```bash
# Démarrer les conteneurs Docker
task docker-up

# Initialiser les bases (18 lignes par défaut)
task init-db

# Populer avec N utilisateurs
task populate -- 100      # 100 users
task populate -- 1000     # 1000 users
task populate -- 10000    # 10000 users

# Populer seulement PostgreSQL avec 5000 users
task populate -- 5000 postgres

# Populer seulement MySQL avec 5000 users
task populate -- 5000 mysql

# Lancer le benchmark
task benchmark

# Nettoyer tout
task clean
```

## Scripts

### `populate_db.py`
Génère et insère N utilisateurs dans les bases de données.

```bash
python3 scripts/populate_db.py [count] [postgres|mysql|both]
```

### `benchmark.sh`
Lance les benchmarks sur les 3 méthodes d'anonymisation.

### `anonymize_faker.py`
Parse un dump SQL et anonymise avec Python Faker.

### `anonymize_gofakeit.go`
Parse un dump SQL et anonymise avec Go gofakeit.

## Configuration NXS

Le fichier `nxs-anonymizer.conf` utilise des commandes avec templates pour les 6 types de champs.

## Résultats

Les résultats sont sauvegardés dans `benchmark-results.txt` avec :
- Méthode utilisée
- Base de données (PostgreSQL/MySQL)
- Temps d'exécution en secondes

## Performances attendues

Sur une machine moderne avec 1000 lignes :
- **Go (gofakeit)** : Le plus rapide (~0.1s)
- **NXS-Data-Anonymizer** : Intermédiaire (~0.5s avec commands)
- **Python (Faker)** : Le plus lent (~1s)

Les performances varient selon :
- Le nombre de lignes
- La complexité des patterns regex
- L'overhead des commandes externes (pour NXS)
- Le type de base de données