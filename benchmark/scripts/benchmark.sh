#!/bin/bash

# Couleurs pour l'affichage
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
RESULTS_FILE="benchmark-results.txt"
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

# Purger les anciens résultats
echo -e "${YELLOW}Nettoyage des anciens résultats...${NC}"
rm -f $RESULTS_FILE
rm -rf verification_results/
mkdir -p verification_results

echo -e "${BLUE}=== BENCHMARK COMPARATIF FAKER ===${NC}"
echo "Date: $(date)" | tee $RESULTS_FILE
echo "" | tee -a $RESULTS_FILE

# Fonction pour mesurer le temps d'exécution
benchmark_command() {
    local cmd="$1"
    local name="$2"
    local db_type="$3"
    
    echo -e "${YELLOW}Test: $name ($db_type)${NC}"
    
    # Mesurer le temps
    local start=$(date +%s.%N)
    eval "$cmd" 2>/tmp/bench_err.log
    local end=$(date +%s.%N)
    
    # Calculer la durée
    local duration=$(echo "$end - $start" | bc)
    
    echo -e "${GREEN}✓ $name ($db_type): ${duration}s${NC}"
    echo "$name,$db_type,$duration" >> $RESULTS_FILE
    
    # Afficher les erreurs s'il y en a
    if [ -s /tmp/bench_err.log ]; then
        cat /tmp/bench_err.log
    fi
}

# Compiler les scripts Go si nécessaire
echo -e "${BLUE}Compilation des scripts Go...${NC}"
cd scripts
go build -o faker-go faker-go.go
cd ..

# Rendre les scripts exécutables
chmod +x scripts/faker-go

echo "" | tee -a $RESULTS_FILE
echo "=== BENCHMARK POSTGRESQL ===" | tee -a $RESULTS_FILE
echo "" | tee -a $RESULTS_FILE

# Test 1: NXS natif avec PostgreSQL
benchmark_command \
    "PGPASSWORD=$POSTGRES_PASS pg_dump -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB | ./nxs-data-anonymizer -t pgsql -c nxs-native.conf > /dev/null" \
    "NXS-Native" \
    "PostgreSQL"


# Test 3: NXS + Go Faker avec PostgreSQL
benchmark_command \
    "PGPASSWORD=$POSTGRES_PASS pg_dump -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB | ./nxs-data-anonymizer -t pgsql -c nxs-go-faker.conf > /dev/null" \
    "NXS-Go-Faker" \
    "PostgreSQL"

echo "" | tee -a $RESULTS_FILE
echo "=== BENCHMARK MYSQL ===" | tee -a $RESULTS_FILE
echo "" | tee -a $RESULTS_FILE

# Test 3: NXS natif avec MySQL
benchmark_command \
    "mysqldump -h $MYSQL_HOST -P $MYSQL_PORT -u$MYSQL_USER $MYSQL_DB | ./nxs-data-anonymizer -t mysql -c nxs-native-mysql.conf > /dev/null" \
    "NXS-Native" \
    "MySQL"

# Test 4: NXS + Go Faker avec MySQL
benchmark_command \
    "mysqldump -h $MYSQL_HOST -P $MYSQL_PORT -u$MYSQL_USER $MYSQL_DB | ./nxs-data-anonymizer -t mysql -c nxs-go-faker-mysql.conf > /dev/null" \
    "NXS-Go-Faker" \
    "MySQL"

echo "" | tee -a $RESULTS_FILE
echo -e "${BLUE}=== RÉSUMÉ DES RÉSULTATS ===${NC}" | tee -a $RESULTS_FILE
echo "" | tee -a $RESULTS_FILE

# Analyser et afficher les résultats
echo "Méthode,Base de données,Temps (secondes)" | tee -a $RESULTS_FILE
tail -n 6 $RESULTS_FILE | grep -E "NXS|Python|Go"

echo ""
echo -e "${GREEN}Benchmark terminé! Résultats sauvegardés dans: $RESULTS_FILE${NC}"