#!/usr/bin/env python3
"""
Script pour comparer visuellement les données avant/après anonymisation
"""

import re
import sys
from pathlib import Path
# from tabulate import tabulate

def extract_insert_values(content):
    """Extrait les valeurs depuis les INSERT SQL ou COPY PostgreSQL"""
    values = []
    
    # Essayer d'abord le format PostgreSQL COPY FROM stdin
    if "COPY public.users" in content:
        lines = content.split('\n')
        in_data_section = False
        
        for line in lines:
            if line.startswith("COPY public.users"):
                in_data_section = True
                continue
            elif line.strip() == '\\.' or line.startswith('--'):
                in_data_section = False
                continue
            
            if in_data_section and line.strip() and not line.startswith('SET') and not line.startswith('--'):
                # Parser les données séparées par des tabulations
                parts = line.split('\t')
                if len(parts) >= 6:  # id + 6 champs
                    clean_parts = []
                    for p in parts[1:7]:  # Ignorer l'ID, prendre les 6 champs
                        if len(p) > 30:
                            p = p[:27] + "..."
                        clean_parts.append(p)
                    values.append(clean_parts)
                    if len(values) >= 10:  # Limiter à 10 lignes
                        break
        
        return values
    
    # Sinon essayer le format MySQL INSERT INTO ... VALUES
    insert_pattern = re.compile(r"INSERT INTO.*?users.*?VALUES\s*(.+);", re.IGNORECASE | re.DOTALL)
    matches = insert_pattern.findall(content)
    
    if not matches:
        return values
    
    # Parser chaque ensemble de valeurs
    for match in matches:
        # Extraire chaque groupe de valeurs entre parenthèses
        value_groups = re.findall(r'\(([^)]+)\)', match)
        
        for group in value_groups[:10]:  # Limiter à 10 lignes
            # Parser les valeurs individuelles
            parts = []
            current_value = ""
            in_string = False
            escape_next = False
            
            for char in group:
                if escape_next:
                    current_value += char
                    escape_next = False
                elif char == '\\':
                    escape_next = True
                    current_value += char
                elif char == "'" and not escape_next:
                    in_string = not in_string
                    current_value += char
                elif char == ',' and not in_string:
                    parts.append(current_value.strip())
                    current_value = ""
                else:
                    current_value += char
            
            if current_value:
                parts.append(current_value.strip())
            
            if len(parts) >= 6:  # On attend au moins 6 champs
                # Nettoyer les valeurs
                clean_parts = []
                for p in parts[:6]:  # Prendre les 6 premiers champs
                    p = p.strip().strip("'").strip('"')
                    if len(p) > 30:  # Tronquer les descriptions longues
                        p = p[:27] + "..."
                    clean_parts.append(p)
                values.append(clean_parts)
    
    return values

def compare_files(original_file, anonymized_file, method_name):
    """Compare les fichiers original et anonymisé"""
    
    print(f"\n=== Comparaison {method_name} ===")
    
    # Lire les fichiers
    try:
        with open(original_file, 'r') as f:
            original_content = f.read()
    except FileNotFoundError:
        print(f"❌ Fichier original non trouvé: {original_file}")
        return
    
    try:
        with open(anonymized_file, 'r') as f:
            anon_content = f.read()
    except FileNotFoundError:
        print(f"❌ Fichier anonymisé non trouvé: {anonymized_file}")
        return
    
    # Si c'est un fichier original avec format tabulaire
    if "original" in str(original_file):
        # Parser le format tabulaire PostgreSQL/MySQL
        lines = original_content.strip().split('\n')
        original_values = []
        
        # Détecter le type de séparateur
        if lines and '\t' in lines[0]:
            # Format MySQL avec tabulations
            for i, line in enumerate(lines[1:]):  # Skip header
                if line.strip():
                    parts = line.split('\t')
                    if len(parts) >= 7:  # id + 6 champs
                        # Nettoyer les données
                        clean_parts = []
                        for p in parts[1:7]:  # Ignorer l'ID, prendre les 6 champs
                            # Gérer les descriptions sur plusieurs lignes
                            if '\\n' in p:
                                p = p.replace('\\n', ' ')
                            if len(p) > 30:
                                p = p[:27] + "..."
                            clean_parts.append(p.strip())
                        original_values.append(clean_parts)
                        
                        if len(original_values) >= 5:  # Limiter à 5 lignes
                            break
        else:
            # Format PostgreSQL avec pipes
            data_started = False
            for line in lines:
                # Détecter le début des données
                if not data_started and '|' in line and not line.strip().startswith('-'):
                    # Skip la ligne d'en-tête avec les noms de colonnes
                    if 'id' in line.lower() and 'name' in line.lower():
                        continue
                    data_started = True
                
                if data_started and line.strip() and '|' in line:
                    # Ignorer les lignes de séparation
                    if line.strip().startswith('+') or line.strip().startswith('-'):
                        continue
                        
                    parts = [p.strip() for p in line.split('|')]
                    if len(parts) >= 7:  # id + 6 champs
                        # Nettoyer les données
                        clean_parts = []
                        for i, p in enumerate(parts[1:7]):  # Ignorer l'ID, prendre les 6 champs
                            # Gérer les descriptions sur plusieurs lignes
                            if '\n' in p:
                                p = p.replace('\n', ' ')
                            if len(p) > 30:
                                p = p[:27] + "..."
                            clean_parts.append(p.strip())
                        original_values.append(clean_parts)
                        
                        if len(original_values) >= 5:  # Limiter à 5 lignes
                            break
    else:
        original_values = extract_insert_values(original_content)[:5]
    
    anon_values = extract_insert_values(anon_content)[:5]
    
    if not original_values:
        print("⚠️  Pas de données originales trouvées")
        return
    
    if not anon_values:
        print("⚠️  Pas de données anonymisées trouvées")
        return
    
    # Afficher la comparaison
    headers = ["Champ", "Original", "Anonymisé"]
    field_names = ["Nom", "Email", "Téléphone", "Date naissance", "Salaire", "Description"]
    
    print(f"\n📊 Exemple de la première ligne:")
    
    if original_values and anon_values:
        comparison_data = []
        for i, field in enumerate(field_names):
            if i < len(original_values[0]) and i < len(anon_values[0]):
                orig = original_values[0][i]
                anon = anon_values[0][i]
                
                # Vérifier si anonymisé
                is_different = orig != anon
                status = "✅" if is_different else "❌"
                
                comparison_data.append([
                    f"{field} {status}",
                    orig[:30] if len(orig) > 30 else orig,
                    anon[:30] if len(anon) > 30 else anon
                ])
        
        # Affichage simple sans tabulate
        for row in comparison_data:
            print(f"  {row[0]:<15} | {row[1]:<30} | {row[2]:<30}")
    
    # Statistiques
    print(f"\n📈 Statistiques:")
    print(f"  - Lignes originales analysées: {len(original_values)}")
    print(f"  - Lignes anonymisées trouvées: {len(anon_values)}")

def main():
    print("=" * 60)
    print("     ANALYSE DE L'ANONYMISATION DES DONNÉES")
    print("=" * 60)
    
    results_dir = Path("/Users/chris/Documents/GitHub/external/nxs-data-anonymizer/benchmark/verification_results")
    
    if not results_dir.exists():
        print("❌ Le dossier verification_results n'existe pas.")
        print("   Lancez d'abord: ./scripts/verify_anonymization.sh")
        return
    
    # PostgreSQL
    print("\n" + "="*30 + " POSTGRESQL " + "="*30)
    
    postgres_original = results_dir / "postgres_original.txt"
    
    compare_files(
        postgres_original,
        results_dir / "postgres_nxs_native.txt",
        "NXS-Native"
    )
    
    compare_files(
        postgres_original,
        results_dir / "postgres_nxs_python.txt",
        "NXS-Python-Faker"
    )
    
    compare_files(
        postgres_original,
        results_dir / "postgres_nxs_go.txt",
        "NXS-Go-Faker"
    )
    
    # MySQL
    print("\n" + "="*30 + " MYSQL " + "="*30)
    
    mysql_original = results_dir / "mysql_original.txt"
    
    compare_files(
        mysql_original,
        results_dir / "mysql_nxs_native.txt",
        "NXS-Native"
    )
    
    compare_files(
        mysql_original,
        results_dir / "mysql_nxs_python.txt",
        "NXS-Python-Faker"
    )
    
    compare_files(
        mysql_original,
        results_dir / "mysql_nxs_go.txt",
        "NXS-Go-Faker"
    )

if __name__ == "__main__":
    main()