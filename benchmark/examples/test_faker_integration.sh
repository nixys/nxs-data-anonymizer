#!/bin/bash

# Script de test pour valider l'intégration faker

echo "=== Test d'intégration GoFakeit dans nxs-data-anonymizer ==="
echo ""

# Configuration des bases de données (réutilise les variables existantes)
source ../scripts/benchmark.sh

# Fonction de test de performance
test_performance() {
    local config_file=$1
    local test_name=$2
    
    echo "Test: $test_name"
    echo "Configuration: $config_file"
    
    # Test PostgreSQL
    echo "  PostgreSQL..."
    time_start=$(date +%s.%N)
    pg_dump -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB \
        --no-owner --no-privileges \
        | ./nxs-data-anonymizer -t pgsql -c "$config_file" > /dev/null
    time_end=$(date +%s.%N)
    pg_time=$(echo "$time_end - $time_start" | bc)
    echo "    Temps: ${pg_time}s"
    
    # Test MySQL
    echo "  MySQL..."
    time_start=$(date +%s.%N)
    mysqldump -h $MYSQL_HOST -P $MYSQL_PORT -u$MYSQL_USER $MYSQL_DB \
        | ./nxs-data-anonymizer -t mysql -c "$config_file" > /dev/null
    time_end=$(date +%s.%N)
    mysql_time=$(echo "$time_end - $time_start" | bc)
    echo "    Temps: ${mysql_time}s"
    
    echo ""
}

# Test de validation des données générées
test_data_quality() {
    local config_file=$1
    local test_name=$2
    
    echo "Validation des données: $test_name"
    
    # Générer un échantillon de données
    pg_dump -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB \
        --no-owner --no-privileges --where="1 limit 5" \
        | ./nxs-data-anonymizer -t pgsql -c "$config_file" \
        > "/tmp/faker_test_${test_name}.sql"
    
    echo "  Échantillon généré dans /tmp/faker_test_${test_name}.sql"
    
    # Vérifier que les données ne sont pas vides
    if grep -q "INSERT" "/tmp/faker_test_${test_name}.sql"; then
        echo "  ✓ Données générées trouvées"
    else
        echo "  ✗ Aucune donnée générée"
        return 1
    fi
    
    # Vérifier la diversité des données (pas de répétition exacte)
    unique_count=$(grep "INSERT" "/tmp/faker_test_${test_name}.sql" | sort | uniq | wc -l)
    total_count=$(grep "INSERT" "/tmp/faker_test_${test_name}.sql" | wc -l)
    
    if [ "$unique_count" -eq "$total_count" ]; then
        echo "  ✓ Données uniques générées"
    else
        echo "  ! Attention: $((total_count - unique_count)) doublons détectés"
    fi
    
    echo ""
}

# Vérifier la présence du binaire nxs-data-anonymizer
if [ ! -f "./nxs-data-anonymizer" ]; then
    echo "Erreur: binaire nxs-data-anonymizer non trouvé"
    echo "Veuillez copier le binaire dans le dossier benchmark/"
    exit 1
fi

echo "1. Test de performance avec l'approche 1 (extension sprig)"
test_performance "examples/nxs-faker-integrated.conf" "sprig_extension"

echo "2. Test de performance avec l'approche 2 (type faker)"
if [ -f "examples/nxs-faker-type.conf" ]; then
    test_performance "examples/nxs-faker-type.conf" "faker_type"
else
    echo "  Configuration type faker non implémentée, ignoré"
fi

echo "3. Test de performance avec l'approche 3 (profil externe)"
if [ -f "examples/nxs-faker-profile.conf" ]; then
    test_performance "examples/nxs-faker-profile.conf" "faker_profile"
else
    echo "  Configuration profil faker non implémentée, ignoré"
fi

echo "4. Comparaison avec les méthodes existantes"
echo "  Configuration native actuelle..."
test_performance "../nxs-native.conf" "native_current"

echo "  Configuration commande externe actuelle..."
test_performance "../nxs-go-faker.conf" "external_current"

echo "5. Tests de validation des données"
test_data_quality "examples/nxs-faker-integrated.conf" "sprig_extension"

echo "=== Résumé des tests d'intégration ==="
echo "Les fichiers de test sont disponibles dans /tmp/faker_test_*.sql"
echo "Comparez les performances et la qualité des données générées."
echo ""
echo "Performance attendue après intégration:"
echo "  - Faker intégré: ~0.5-1s (comparable au natif)"
echo "  - Commandes externes: ~24s (actuel)"
echo "  - Natif NXS: ~0.4s (référence)"