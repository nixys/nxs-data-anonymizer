# Proposition d'intégration GoFakeit dans nxs-data-anonymizer

## Résumé

Intégrer la bibliothèque gofakeit/v6 directement dans nxs-data-anonymizer pour améliorer les performances et éliminer les appels de commandes externes.

## Avantages

- **Performance** : Élimination des appels exec.Command() (24s → ~0.5s estimé)
- **Simplicité** : Pas besoin de binaires externes 
- **Robustesse** : Moins de points de défaillance
- **Localisation** : Support natif des locales (fr_FR, en_US, etc.)

## Architecture proposée

### Approche 1: Extension des templates Sprig

**Avantages** : Simplicité, compatibilité totale avec l'existant
**Implémentation** : Ajouter les fonctions gofakeit aux templates sprig existants

```yaml
filters:
  public.users:
    columns:
      name:
        value: "{{ fakerName }} {{ fakerLastName }}"
      email:
        value: "{{ fakerEmail }}"
      phone:
        value: "{{ fakerPhone }}"
```

### Approche 2: Nouveau type "faker" + configuration globale

**Avantages** : Configuration centralisée, plus flexible
**Implémentation** : Nouveau type de filtre + configuration globale

```yaml
faker:
  enabled: true
  locale: fr_FR
  
filters:
  public.users:
    columns:
      name:
        type: faker
        value: "name"  # ou "{{ name }} {{ lastName }}"
```

### Approche 3: Configuration par fichier externe (RECOMMANDÉE)

**Avantages** : Réutilisabilité, séparation des concerns, plus maintenable

```yaml
faker: true
faker-file: "./faker-profiles/french-names.yml"

filters:
  public.users:
    columns:
      name:
        type: faker
        value: "fullName"
```

Le fichier `french-names.yml` :
```yaml
locale: fr_FR
profiles:
  fullName: "{{ fakerFirstName }} {{ fakerLastName }}"
  email: "{{ fakerEmail }}"
  phone: "+33{{ fakerNumber 9 }}"
  salary: "{{ fakerFloat 25000 75000 2 }}"
```

## Implémentation détaillée

### 1. Modification de misc/values.go

```go
const (
    ValueTypeUnknown  ValueType = "unknown" 
    ValueTypeTemplate ValueType = "template"
    ValueTypeCommand  ValueType = "command"
    ValueTypeFaker    ValueType = "faker"    // NOUVEAU
)
```

### 2. Modification de ctx/conf.go

```go
type confOpts struct {
    // ... existant
    Faker     *fakerConf                    `conf:"faker"`     // NOUVEAU
    FakerFile string                        `conf:"faker-file"` // NOUVEAU
}

type fakerConf struct {
    Enabled bool   `conf:"enabled" conf_extraopts:"default=false"`
    Locale  string `conf:"locale" conf_extraopts:"default=fr_FR"`
}
```

### 3. Modification de misc/template.go

```go
import (
    "github.com/brianvoe/gofakeit/v6"
)

// Dans la fonction TemplateExec, ajouter aux sprig functions :
func addFakerFunctions(funcMap ttemplate.FuncMap) {
    // Configuration gofakeit
    gofakeit.Seed(time.Now().UnixNano())
    
    // Fonctions de base
    funcMap["fakerName"] = gofakeit.Name
    funcMap["fakerFirstName"] = gofakeit.FirstName
    funcMap["fakerLastName"] = gofakeit.LastName
    funcMap["fakerEmail"] = gofakeit.Email
    funcMap["fakerPhone"] = gofakeit.Phone
    funcMap["fakerDate"] = func() string {
        return gofakeit.Date().Format("2006-01-02")
    }
    funcMap["fakerFloat"] = func(min, max float64, precision int) float64 {
        return gofakeit.Float64Range(min, max)
    }
    funcMap["fakerNumber"] = func(digits int) string {
        return gofakeit.Numerify(strings.Repeat("#", digits))
    }
    funcMap["fakerText"] = func(length int) string {
        return gofakeit.LoremIpsumSentence(length)
    }
    
    // Support de locale
    funcMap["fakerSetLocale"] = func(locale string) string {
        // Configuration locale si supportée
        return ""
    }
}
```

### 4. Modification de modules/filters/relfilter/filter.go

```go
// Dans execFilter, ajouter le cas ValueTypeFaker :
case misc.ValueTypeFaker:
    // Le type faker utilise les templates avec fonctions faker étendues
    v, d, err = misc.TemplateExec(f.v, td)
    if err != nil {
        return []byte{}, false, fmt.Errorf("filter: faker template: %w", err)
    }
```

## Fichiers de configuration d'exemple

### Configuration simple (Approche 1)
```yaml
filters:
  public.users:
    columns:
      name:
        value: "{{ fakerFirstName }} {{ fakerLastName }}"
      email:
        value: "{{ fakerEmail }}"
      phone:
        value: "+33{{ fakerNumber 9 }}"
      birth_date:
        value: "{{ fakerDate }}"
      salary:
        value: "{{ fakerFloat 25000 75000 2 }}"
      description:
        value: "{{ fakerText 50 }}"
```

### Configuration avec type faker (Approche 2)
```yaml
faker:
  enabled: true
  locale: fr_FR

filters:
  public.users:
    columns:
      name:
        type: faker
        value: "{{ fakerFirstName }} {{ fakerLastName }}"
      email:
        type: faker
        value: "{{ fakerEmail }}"
```

### Configuration avec fichier externe (Approche 3 - RECOMMANDÉE)
```yaml
faker: true
faker-file: "./profiles/french-users.yml"

filters:
  public.users:
    columns:
      name:
        type: faker
        value: "fullName"
      email:
        type: faker
        value: "email"
      phone:
        type: faker
        value: "frenchPhone"
```

## Performance attendue

**Avant** (commandes externes) :
- PostgreSQL + Go Faker : 24.5s pour 1000 users
- MySQL + Go Faker : 24.7s pour 1000 users

**Après** (intégration native) :
- PostgreSQL + Faker intégré : ~0.8s estimé (comparable au natif)
- MySQL + Faker intégré : ~0.8s estimé

**Gain** : ~97% de réduction du temps d'anonymisation

## Rétrocompatibilité

L'implémentation propose trois niveaux :

1. **Niveau 1** : Extension transparente des templates sprig (aucun breaking change)
2. **Niveau 2** : Nouveau type "faker" optionnel 
3. **Niveau 3** : Configuration externe optionnelle

## Migration

### Depuis les commandes externes
**Avant** :
```yaml
name:
  type: command
  value: './scripts/faker-go name'
```

**Après** :
```yaml
name:
  value: "{{ fakerName }}"
```

### Depuis les templates natifs
**Avant** :
```yaml
name:
  value: "{{ randAlpha 8 }} {{ randAlpha 10 }}"
```

**Après** :
```yaml
name:
  value: "{{ fakerFirstName }} {{ fakerLastName }}"
```

## Plan d'implémentation

1. **Phase 1** : Extension des templates sprig avec fonctions faker de base
2. **Phase 2** : Ajout du type "faker" optionnel  
3. **Phase 3** : Support des fichiers de configuration externes
4. **Phase 4** : Support avancé des locales et profils personnalisés

## Test de validation

Créer des configurations de test pour valider :
- Performance équivalente au natif
- Génération de données réalistes
- Compatibilité avec les configurations existantes
- Support des locales françaises

## Conclusion

Cette intégration permettrait d'obtenir les performances du système natif tout en gardant la richesse des données générées par gofakeit, éliminant le compromis actuel performance vs qualité des données.