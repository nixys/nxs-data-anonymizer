#!/bin/bash

echo "Attente du démarrage de PostgreSQL..."
until docker exec benchmark-postgres pg_isready -U postgres > /dev/null 2>&1; do
  echo -n "."
  sleep 1
done
echo " PostgreSQL prêt!"

echo "Attente du démarrage de MySQL..."
TIMEOUT=60
count=0
until docker exec benchmark-mysql mysqladmin ping -h localhost -uroot -proot --silent > /dev/null 2>&1; do
  echo -n "."
  sleep 2
  count=$((count + 2))
  if [ $count -ge $TIMEOUT ]; then
    echo " Timeout atteint pour MySQL!"
    docker logs benchmark-mysql --tail=20
    exit 1
  fi
done
echo " MySQL prêt!"

# Attendre encore 5s pour que MySQL soit complètement initialisé
sleep 5

echo "Toutes les bases de données sont prêtes!"