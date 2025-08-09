# Implémentation simple de l'option 1

## Résumé ultra-simple

**Objectif** : Ajouter des fonctions `fakerXXX` aux templates existants de nxs-data-anonymizer.

**Résultat** : Au lieu de `{{ randAlpha 8 }}` qui donne "KjHgFdSa", on aura `{{ fakerFirstName }}` qui donne "Jean".

## Modification minimale nécessaire

### 1. Fichier `misc/template.go` - SEULE modification obligatoire

```go
// AVANT (lignes 27-46)
t, err := ttemplate.New("template").Funcs(func() ttemplate.FuncMap {
    // Get current sprig functions
    t := sprig.TxtFuncMap()

    // Add additional functions
    t["null"] = func() string {
        return null
    }
    t["isNull"] = func(v string) bool {
        if v == null {
            return true
        }
        return false
    }
    t["drop"] = func() string {
        return drop
    }

    return t
}())
```

```go
// APRÈS - Ajouter juste ces lignes
import (
    // ... imports existants
    "github.com/brianvoe/gofakeit/v6"  // NOUVEAU
    "time"                             // NOUVEAU
)

t, err := ttemplate.New("template").Funcs(func() ttemplate.FuncMap {
    // Get current sprig functions
    t := sprig.TxtFuncMap()

    // Add additional functions (existant)
    t["null"] = func() string {
        return null
    }
    t["isNull"] = func(v string) bool {
        if v == null {
            return true
        }
        return false
    }
    t["drop"] = func() string {
        return drop
    }

    // NOUVEAU : Ajouter les fonctions faker
    gofakeit.Seed(time.Now().UnixNano())
    
    t["fakerName"] = gofakeit.Name
    t["fakerFirstName"] = gofakeit.FirstName
    t["fakerLastName"] = gofakeit.LastName
    t["fakerEmail"] = gofakeit.Email
    t["fakerPhone"] = gofakeit.Phone
    t["fakerDate"] = func() string {
        return gofakeit.Date().Format("2006-01-02")
    }
    t["fakerBirthDate"] = func() string {
        return gofakeit.DateRange(
            time.Now().AddDate(-80, 0, 0),
            time.Now().AddDate(-18, 0, 0),
        ).Format("2006-01-02")
    }
    t["fakerFloat"] = func(min, max float64, precision int) float64 {
        return gofakeit.Float64Range(min, max)
    }
    t["fakerInt"] = gofakeit.IntRange
    t["fakerSentence"] = func(words int) string {
        return gofakeit.LoremIpsumSentence(words)
    }
    
    // Fonctions françaises
    t["fakerPhoneFR"] = func() string {
        return "+33" + gofakeit.Numerify("#########")
    }
    t["fakerEmailFR"] = func() string {
        domains := []string{"gmail.com", "hotmail.fr", "orange.fr", "free.fr", "yahoo.fr"}
        return gofakeit.Username() + "@" + domains[gofakeit.IntRange(0, len(domains)-1)]
    }

    return t
}())
```

## C'est tout !

Avec **SEULEMENT** cette modification, tu peux maintenant utiliser :

```yaml
# Au lieu de ça :
name: "{{ randAlpha 8 }} {{ randAlpha 10 }}"

# Tu peux écrire ça :
name: "{{ fakerFirstName }} {{ fakerLastName }}"
```

## Test immédiat

1. **Créer une config de test** :
```yaml
# test-faker.conf
filters:
  public.users:
    columns:
      name:
        value: "{{ fakerFirstName }} {{ fakerLastName }}"
      email:
        value: "{{ fakerEmailFR }}"
      phone:
        value: "{{ fakerPhoneFR }}"
```

2. **Tester** :
```bash
echo "INSERT INTO users VALUES (1, 'Test', 'test@test.com', '0123456789');" | ./nxs-data-anonymizer -t mysql -c test-faker.conf
```

3. **Résultat attendu** :
```sql
INSERT INTO users VALUES (1, 'Jean Dupont', 'jean@gmail.com', '+33678123456');
```

## Avantages de cette approche

- **Zéro breaking change** : Les configs existantes continuent de fonctionner
- **Adoption progressive** : Tu peux migrer colonne par colonne  
- **Performance maximale** : Plus d'appels externes, tout en Go natif
- **Simplicité** : Une seule modification de fichier

## Fonctions ajoutées

| Fonction | Exemple de résultat |
|----------|-------------------|
| `fakerName` | "Jean Dupont" |
| `fakerFirstName` | "Jean" |
| `fakerLastName` | "Dupont" |  
| `fakerEmail` | "jean@example.com" |
| `fakerEmailFR` | "jean@gmail.com" |
| `fakerPhone` | "123-456-7890" |
| `fakerPhoneFR` | "+33678123456" |
| `fakerDate` | "1985-03-15" |
| `fakerBirthDate` | "1975-06-22" (18-80 ans) |
| `fakerFloat 1000 9999 2` | 1234.56 |
| `fakerInt 18 65` | 42 |
| `fakerSentence 10` | "Lorem ipsum dolor sit..." |

## Migration progressive

```yaml
# Étape 1 : Migrer juste les noms
filters:
  public.users:
    columns:
      name:
        value: "{{ fakerFirstName }} {{ fakerLastName }}"  # NOUVEAU
      email:
        value: "{{ randAlphaNum 8 }}@example.com"          # ANCIEN
      phone:
        value: "+33{{ randNumeric 9 }}"                    # ANCIEN

# Étape 2 : Migrer les emails  
      email:
        value: "{{ fakerEmailFR }}"                        # NOUVEAU

# Étape 3 : etc...
```

**Résultat** : Performance native + données réalistes !