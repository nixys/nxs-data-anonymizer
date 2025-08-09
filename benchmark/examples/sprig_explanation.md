# Templates Sprig - Explication concrète

## Qu'est-ce que Sprig ?

**Sprig** est une bibliothèque de fonctions prêtes à l'emploi pour les templates Go. C'est comme une boîte à outils avec plein de fonctions utiles.

### Syntaxe de base

```
{{ nomFonction }}                    # Appel simple
{{ nomFonction "parametre" }}        # Avec un paramètre
{{ nomFonction 123 456 }}           # Avec plusieurs paramètres
{{ fonction1 | fonction2 }}         # Pipeline (chaînage)
```

## Dans nxs-data-anonymizer actuellement

### Ce qui existe déjà (fonctions Sprig natives) :

```yaml
# Génération de texte aléatoire
name: "{{ randAlpha 8 }} {{ randAlpha 10 }}"
# Résultat: "KjHgFdSa QwErTyUiOp"

# Génération de nombres
phone: "+33{{ randNumeric 9 }}"
# Résultat: "+33123456789"

# Génération d'email
email: "{{ randAlphaNum 8 }}@example.com"
# Résultat: "abc123de@example.com"

# Manipulation de dates
birth_date: "{{ dateModify \"-20y\" now | date \"2006-01-02\" }}"
# Résultat: "2004-08-09" (20 ans avant aujourd'hui)

# Nombres dans une plage
salary: "{{ randInt 25000 75000 }}.{{ randInt 10 99 }}"
# Résultat: "45678.42"
```

### Fonctions Sprig disponibles

```yaml
# DATES
"{{ now }}"                          # Date actuelle
"{{ date \"2006-01-02\" now }}"      # Format date
"{{ dateModify \"+1y\" now }}"       # Ajouter 1 an

# NOMBRES
"{{ randInt 1 100 }}"                # Entier entre 1 et 100
"{{ add 5 3 }}"                      # Addition: 8
"{{ mul 4 7 }}"                      # Multiplication: 28

# TEXTE
"{{ randAlpha 10 }}"                 # 10 lettres aléatoires
"{{ randAlphaNum 8 }}"               # 8 caractères alphanumériques
"{{ randNumeric 6 }}"                # 6 chiffres
"{{ upper \"hello\" }}"              # "HELLO"
"{{ lower \"WORLD\" }}"              # "world"

# LISTES
"{{ list \"a\" \"b\" \"c\" | join \",\" }}"  # "a,b,c"

# CONDITIONS
"{{ if gt 5 3 }}oui{{ else }}non{{ end }}"   # "oui"
```

## Comment gofakeit s'intègre

Au lieu d'avoir des données "random" basiques, on aura des données **réalistes** :

### Avant (Sprig natif) :
```yaml
name: "{{ randAlpha 8 }} {{ randAlpha 10 }}"
# Résultat: "KjHgFdSa QwErTyUiOp" (pas réaliste)
```

### Après (avec gofakeit intégré) :
```yaml
name: "{{ fakerFirstName }} {{ fakerLastName }}"
# Résultat: "Jean Dupont" (réaliste !)
```

### Comparaison complète

```yaml
# AVANT - Données aléatoires
filters:
  public.users:
    columns:
      name: "{{ randAlpha 8 }} {{ randAlpha 10 }}"
      email: "{{ randAlphaNum 8 }}@example.com"
      phone: "+33{{ randNumeric 9 }}"
      # Résultats: "KjHgFdSa QwErTyUiOp", "abc123de@example.com", "+33123456789"

# APRÈS - Données réalistes
filters:
  public.users:
    columns:
      name: "{{ fakerFirstName }} {{ fakerLastName }}"
      email: "{{ fakerEmail }}"
      phone: "{{ fakerPhoneFR }}"
      # Résultats: "Jean Dupont", "jean.dupont@gmail.com", "+33678123456"
```