#!/bin/bash

echo "=== TEST MANUEL RAPIDE ==="
echo ""

echo "1. PostgreSQL Faker:"
echo "Avant:"
PGPASSWORD=postgres pg_dump -h localhost -U postgres -d testdb --data-only -t users | grep "^[0-9]" | head -1 | cut -f2-4
echo "Après:"
PGPASSWORD=postgres pg_dump -h localhost -U postgres -d testdb --data-only -t users | grep -A 4 "COPY" | head -5 | ./nxs-data-anonymizer-faker -t pgsql -c postgres-faker.conf | grep "^[0-9]" | head -1 | cut -f2-4
echo ""

echo "2. PostgreSQL Natif:"
echo "Avant:"
PGPASSWORD=postgres pg_dump -h localhost -U postgres -d testdb --data-only -t users | grep "^[0-9]" | head -1 | cut -f2-4
echo "Après:"
PGPASSWORD=postgres pg_dump -h localhost -U postgres -d testdb --data-only -t users | grep -A 4 "COPY" | head -5 | ./nxs-data-anonymizer -t pgsql -c postgres-native.conf | grep "^[0-9]" | head -1 | cut -f2-4
echo ""

echo "3. MySQL Faker:"
echo "Avant:"
mysqldump -h 127.0.0.1 -uroot -proot testdb users --no-create-info 2>/dev/null | grep INSERT | head -1 | sed 's/.*VALUES (//' | sed 's/);$//' | cut -d, -f2-4
echo "Après:"
mysqldump -h 127.0.0.1 -uroot -proot testdb users --no-create-info 2>/dev/null | grep INSERT | head -1 | ./nxs-data-anonymizer-faker -t mysql -c mysql-faker.conf | grep INSERT | sed 's/.*VALUES (//' | sed 's/);$//' | cut -d, -f2-4
echo ""

echo "4. MySQL Natif:"
echo "Avant:"
mysqldump -h 127.0.0.1 -uroot -proot testdb users --no-create-info 2>/dev/null | grep INSERT | head -1 | sed 's/.*VALUES (//' | sed 's/);$//' | cut -d, -f2-4
echo "Après:"
mysqldump -h 127.0.0.1 -uroot -proot testdb users --no-create-info 2>/dev/null | grep INSERT | head -1 | ./nxs-data-anonymizer -t mysql -c mysql-native.conf | grep INSERT | sed 's/.*VALUES (//' | sed 's/);$//' | cut -d, -f2-4