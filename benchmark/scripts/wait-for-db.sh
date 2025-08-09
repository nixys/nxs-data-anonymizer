#!/bin/bash

echo "Attente du démarrage de PostgreSQL..."
until docker exec benchmark-postgres pg_isready -U postgres > /dev/null 2>&1; do
  echo -n "."
  sleep 1
done
echo " PostgreSQL prêt!"

echo "Attente du démarrage de MySQL..."
until docker exec benchmark-mysql mysqladmin ping -h localhost -uroot -proot --silent > /dev/null 2>&1; do
  echo -n "."
  sleep 1
done
echo " MySQL prêt!"

echo "Toutes les bases de données sont prêtes!"