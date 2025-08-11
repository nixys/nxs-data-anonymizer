#!/usr/bin/env python3
"""
Script pour populer les bases de données avec des données de test
Usage: python populate_db.py [nombre_users] [postgres|mysql|both]
"""

import sys
import argparse
from faker import Faker
import psycopg2
import mysql.connector
from datetime import datetime
import time

fake = Faker('fr_FR')

def generate_user_data(count):
    """Génère N utilisateurs avec tous les champs"""
    users = []
    for i in range(count):
        # Générer un mobile français (06 ou 07)
        mobile_prefix = fake.random.choice(['06', '07'])
        mobile = f"+33{mobile_prefix}{fake.numerify('########')}"
        
        # Générer IBAN français
        iban = f"FR{fake.numerify('##')} {fake.numerify('#### #### #### #### #### ###')}"
        
        # Générer numéro sécurité sociale français
        sexe = fake.random.choice([1, 2])
        annee = fake.random.randint(50, 99)
        mois = fake.random.randint(1, 12)
        dept = fake.random.randint(1, 95)
        commune = fake.random.randint(1, 999)
        cle = fake.random.randint(1, 97)
        ssn = f"{sexe}{annee:02d}{mois:02d}{dept:02d}{commune:03d}{cle:02d}"
        
        # Générer numéro TVA français
        vat_key = fake.random.randint(10, 99)
        vat_siren = fake.numerify('#########')
        vat = f"FR{vat_key}{vat_siren}"
        
        user = (
            fake.name(),                           # name
            fake.email(),                          # email  
            fake.phone_number(),                   # phone
            mobile,                                # mobile
            fake.date_between('-50y', '-20y'),    # birth_date
            round(fake.random.uniform(25000, 75000), 2),  # salary
            fake.text(max_nb_chars=200),          # description
            fake.address().replace('\n', ', '),   # address
            fake.city(),                          # city
            fake.postcode(),                      # postal_code
            iban,                                 # iban
            ssn,                                  # ssn
            vat,                                  # vat_number
            fake.ipv4()                          # ip_address
        )
        users.append(user)
    return users

def populate_postgres(users, host='localhost', port=5432, db='testdb', user='postgres', password='postgres'):
    """Insère les données dans PostgreSQL"""
    try:
        conn = psycopg2.connect(
            host=host,
            port=port,
            database=db,
            user=user,
            password=password
        )
        cursor = conn.cursor()
        
        # Vider la table existante
        cursor.execute("TRUNCATE TABLE users RESTART IDENTITY")
        
        # Insertion en masse
        insert_query = """
            INSERT INTO users (name, email, phone, mobile, birth_date, salary, description, address, city, postal_code, iban, ssn, vat_number, ip_address)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
        """
        
        start = time.time()
        cursor.executemany(insert_query, users)
        conn.commit()
        elapsed = time.time() - start
        
        cursor.execute("SELECT COUNT(*) FROM users")
        count = cursor.fetchone()[0]
        
        print(f"✓ PostgreSQL: {count} lignes insérées en {elapsed:.2f}s")
        print(f"  Vitesse: {count/elapsed:.0f} lignes/seconde")
        
        cursor.close()
        conn.close()
        return True
        
    except Exception as e:
        print(f"✗ Erreur PostgreSQL: {e}")
        return False

def populate_mysql(users, host='localhost', port=3306, db='testdb', user='root', password='root'):
    """Insère les données dans MySQL"""
    try:
        conn = mysql.connector.connect(
            host=host,
            port=port,
            database=db,
            user=user,
            password=password,
            autocommit=True
        )
        cursor = conn.cursor()
        
        # Vider la table existante
        cursor.execute("TRUNCATE TABLE users")
        
        # Insertion en masse
        insert_query = """
            INSERT INTO users (name, email, phone, mobile, birth_date, salary, description, address, city, postal_code, iban, ssn, vat_number, ip_address)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
        """
        
        start = time.time()
        cursor.executemany(insert_query, users)
        conn.commit()
        elapsed = time.time() - start
        
        cursor.execute("SELECT COUNT(*) FROM users")
        count = cursor.fetchone()[0]
        
        print(f"✓ MySQL: {count} lignes insérées en {elapsed:.2f}s")
        print(f"  Vitesse: {count/elapsed:.0f} lignes/seconde")
        
        cursor.close()
        conn.close()
        return True
        
    except Exception as e:
        print(f"✗ Erreur MySQL: {e}")
        return False

def main():
    parser = argparse.ArgumentParser(description='Populer les bases de données avec des données de test')
    parser.add_argument('count', type=int, nargs='?', default=1000,
                       help='Nombre d\'utilisateurs à créer (défaut: 1000)')
    parser.add_argument('database', choices=['postgres', 'mysql', 'both'], 
                       nargs='?', default='both',
                       help='Base de données à populer (défaut: both)')
    parser.add_argument('--host', default='localhost', help='Hôte de la base de données')
    parser.add_argument('--pg-port', type=int, default=5432, help='Port PostgreSQL')
    parser.add_argument('--mysql-port', type=int, default=3306, help='Port MySQL')
    
    args = parser.parse_args()
    
    print(f"\n=== Génération de {args.count} utilisateurs ===")
    start = time.time()
    users = generate_user_data(args.count)
    elapsed = time.time() - start
    print(f"✓ Données générées en {elapsed:.2f}s\n")
    
    print("=== Insertion dans les bases de données ===")
    
    if args.database in ['postgres', 'both']:
        populate_postgres(users, host=args.host, port=args.pg_port)
    
    if args.database in ['mysql', 'both']:
        populate_mysql(users, host=args.host, port=args.mysql_port)
    
    print("\n✓ Population terminée!")

if __name__ == "__main__":
    main()