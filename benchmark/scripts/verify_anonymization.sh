#!/bin/bash

# Couleurs pour l'affichage
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
POSTGRES_HOST="localhost"
POSTGRES_PORT="5432"
POSTGRES_DB="testdb"
POSTGRES_USER="postgres"
POSTGRES_PASS="postgres"

MYSQL_HOST="127.0.0.1"
MYSQL_PORT="3306"
MYSQL_DB="testdb"
MYSQL_USER="root"
MYSQL_PASS="root"

# Export pour éviter les warnings mysqldump
export MYSQL_PWD="$MYSQL_PASS"

echo -e "${BLUE}=== VÉRIFICATION DE L'ANONYMISATION ===${NC}"

# Purger les anciens résultats
echo -e "${YELLOW}Nettoyage des anciens résultats...${NC}"
rm -rf verification_results/
mkdir -p verification_results

echo ""


# Fonction pour extraire et afficher les données
extract_data() {
    local db_type=$1
    local method=$2
    local output_file=$3
    
    echo -e "${YELLOW}Extraction avec $method sur $db_type...${NC}"
    
    if [ "$db_type" = "postgres" ]; then
        if [ "$method" = "original" ]; then
            PGPASSWORD=$POSTGRES_PASS psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB \
                -c "SELECT * FROM users LIMIT 10" > "$output_file"
        else
            PGPASSWORD=$POSTGRES_PASS pg_dump -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER \
                -d $POSTGRES_DB --data-only --table=users | head -100 > "$output_file.tmp"
            
            if [ "$method" = "nxs-native" ]; then
                cat "$output_file.tmp" | ./nxs-data-anonymizer -t pgsql -c nxs-native.conf > "$output_file"
            elif [ "$method" = "nxs-go" ]; then
                cat "$output_file.tmp" | ./nxs-data-anonymizer -t pgsql -c nxs-go-faker.conf > "$output_file"
            fi
            rm -f "$output_file.tmp"
        fi
    elif [ "$db_type" = "mysql" ]; then
        if [ "$method" = "original" ]; then
            mysql -h $MYSQL_HOST -P $MYSQL_PORT -u$MYSQL_USER $MYSQL_DB \
                -e "SELECT * FROM users LIMIT 10" > "$output_file"
        else
            mysqldump -h $MYSQL_HOST -P $MYSQL_PORT -u$MYSQL_USER $MYSQL_DB \
                --tables users --where="1 limit 10" > "$output_file.tmp" 2>/dev/null
            
            if [ "$method" = "nxs-native" ]; then
                cat "$output_file.tmp" | ./nxs-data-anonymizer -t mysql -c nxs-native-mysql.conf > "$output_file"
            elif [ "$method" = "nxs-go" ]; then
                cat "$output_file.tmp" | ./nxs-data-anonymizer -t mysql -c nxs-go-faker-mysql.conf > "$output_file"
            fi
            rm -f "$output_file.tmp"
        fi
    fi
}

# Compiler les scripts Go si nécessaire
if [ ! -f "./scripts/faker-go" ]; then
    echo "Compilation du script Go..."
    cd scripts
    go build -o faker-go faker-go.go
    cd ..
fi

# PostgreSQL
echo -e "${BLUE}=== POSTGRESQL ===${NC}"
echo ""

# Données originales
extract_data "postgres" "original" "verification_results/postgres_original.txt"

# Anonymisation NXS Native
extract_data "postgres" "nxs-native" "verification_results/postgres_nxs_native.txt"


# Anonymisation NXS + Go Faker
extract_data "postgres" "nxs-go" "verification_results/postgres_nxs_go.txt"


echo ""
echo -e "${BLUE}=== MYSQL ===${NC}"
echo ""

# Données originales
extract_data "mysql" "original" "verification_results/mysql_original.txt"

# Anonymisation NXS Native
extract_data "mysql" "nxs-native" "verification_results/mysql_nxs_native.txt"

# Anonymisation NXS + Go Faker
extract_data "mysql" "nxs-go" "verification_results/mysql_nxs_go.txt"


echo ""
echo -e "${GREEN}=== COMPARAISON DES RÉSULTATS ===${NC}"
echo ""

echo -e "${YELLOW}PostgreSQL - Données originales (3 premières lignes):${NC}"
head -5 verification_results/postgres_original.txt

echo ""
echo -e "${YELLOW}PostgreSQL - Après anonymisation NXS Native (extrait INSERT):${NC}"
grep -m 1 "INSERT INTO" verification_results/postgres_nxs_native.txt | cut -c1-200 || echo "Pas de données INSERT trouvées"

echo ""

echo ""
echo -e "${YELLOW}PostgreSQL - Après anonymisation NXS-Go (extrait INSERT):${NC}"
grep -m 1 "INSERT INTO" verification_results/postgres_nxs_go.txt | cut -c1-200 || echo "Pas de données INSERT trouvées"

echo ""
echo -e "${GREEN}✓ Résultats sauvegardés dans verification_results/${NC}"
echo "  - postgres_original.txt : Données originales PostgreSQL"
echo "  - postgres_nxs_[native|go].txt : Données anonymisées PostgreSQL"
echo "  - mysql_original.txt : Données originales MySQL"
echo "  - mysql_nxs_[native|go].txt : Données anonymisées MySQL"