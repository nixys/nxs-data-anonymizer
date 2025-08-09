#!/bin/bash

# Script de benchmark unifié pour les 4 configurations
# Usage: ./test_benchmark.sh

# Couleurs pour l'affichage
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration bases de données
POSTGRES_HOST="localhost"
POSTGRES_PORT="5432"
POSTGRES_DB="testdb"
POSTGRES_USER="postgres"
export PGPASSWORD="postgres"

MYSQL_HOST="127.0.0.1"
MYSQL_PORT="3306"
MYSQL_DB="testdb"
MYSQL_USER="root"
export MYSQL_PWD="root"

echo -e "${BLUE}==================================================${NC}"
echo -e "${BLUE}         BENCHMARK NXS-DATA-ANONYMIZER${NC}"
echo -e "${BLUE}==================================================${NC}"
echo ""

# Vérifier les binaires
echo -e "${YELLOW}Vérification des binaires...${NC}"
if [ ! -f "./nxs-data-anonymizer" ]; then
    echo -e "${RED}❌ Binaire nxs-data-anonymizer non trouvé${NC}"
    echo "   Compilez avec: task build-native"
    exit 1
fi

if [ ! -f "./nxs-data-anonymizer-faker" ]; then
    echo -e "${RED}❌ Binaire nxs-data-anonymizer-faker non trouvé${NC}"
    echo "   Compilez avec: go build -o benchmark/nxs-data-anonymizer-faker"
    exit 1
fi

echo -e "${GREEN}✓ Binaires trouvés${NC}"
echo ""

# Fonction de test
run_test() {
    local db_type=$1
    local method=$2
    local config=$3
    local binary=$4
    local description=$5
    
    echo -e "${YELLOW}Test: $description${NC}"
    
    # Mesure du temps
    time_start=$(date +%s.%N)
    
    if [ "$db_type" = "postgres" ]; then
        pg_dump -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d "$POSTGRES_DB" \
            --no-owner --no-privileges -t users 2>/dev/null | \
            head -500 | \
            $binary -t pgsql -c $config > /dev/null 2>&1
    else
        mysqldump -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u"$MYSQL_USER" "$MYSQL_DB" \
            --tables users --extended-insert=false 2>/dev/null | \
            head -500 | \
            $binary -t mysql -c $config > /dev/null 2>&1
    fi
    
    time_end=$(date +%s.%N)
    elapsed=$(echo "$time_end - $time_start" | bc)
    
    echo -e "   ${GREEN}✓ Temps: ${elapsed}s${NC}"
    echo ""
    
    # Stocker le résultat
    echo "$description: ${elapsed}s" >> benchmark-results.txt
}

# Nettoyer les résultats précédents
rm -f benchmark-results.txt
echo "=== Résultats du benchmark ===" > benchmark-results.txt
echo "Date: $(date)" >> benchmark-results.txt
echo "" >> benchmark-results.txt

# Lancer les 4 tests
echo -e "${BLUE}=== TESTS DE PERFORMANCE ===${NC}"
echo ""

run_test "postgres" "faker" "postgres-faker.conf" "./nxs-data-anonymizer-faker" "PostgreSQL avec Faker intégré"
run_test "postgres" "native" "postgres-native.conf" "./nxs-data-anonymizer" "PostgreSQL avec fonctions natives"
run_test "mysql" "faker" "mysql-faker.conf" "./nxs-data-anonymizer-faker" "MySQL avec Faker intégré"
run_test "mysql" "native" "mysql-native.conf" "./nxs-data-anonymizer" "MySQL avec fonctions natives"

echo -e "${BLUE}==================================================${NC}"
echo -e "${GREEN}📊 RÉSUMÉ DES RÉSULTATS${NC}"
echo -e "${BLUE}==================================================${NC}"
echo ""
cat benchmark-results.txt | grep -E "(PostgreSQL|MySQL)"
echo ""
echo -e "${GREEN}✅ Benchmark terminé!${NC}"
echo ""