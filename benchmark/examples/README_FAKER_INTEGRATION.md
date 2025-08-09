# Intégration GoFakeit dans nxs-data-anonymizer

## Vue d'ensemble

Cette proposition détaille l'intégration native de la bibliothèque GoFakeit dans nxs-data-anonymizer pour améliorer drastiquement les performances tout en gardant la qualité des données générées.

## Problème actuel

- **Performance médiocre** : 24+ secondes pour anonymiser 1000 utilisateurs avec gofakeit via commandes externes
- **Complexité** : Nécessité de binaires externes et scripts wrapper
- **Latence** : Chaque valeur nécessite un appel `exec.Command()`

## Solution proposée

Intégration native de gofakeit dans le moteur de templates Sprig existant.

### Gain de performance attendu

| Méthode | Temps actuel | Temps estimé après intégration | Amélioration |
|---------|--------------|--------------------------------|--------------|
| PostgreSQL + Faker | 24.5s | ~0.8s | **97% plus rapide** |
| MySQL + Faker | 24.7s | ~0.8s | **97% plus rapide** |

## Architecture des fichiers

### Fichiers de modifications du code source

- `examples/misc_values_modified.go` - Ajout du type `ValueTypeFaker`
- `examples/misc_template_modified.go` - Extension des templates avec fonctions gofakeit
- `examples/ctx_conf_modified.go` - Support de la configuration faker

### Configurations d'exemple

#### Approche 1 : Extension transparente des templates (RECOMMANDÉE)
```yaml
# examples/nxs-faker-integrated.conf
filters:
  public.users:
    columns:
      name:
        value: "{{ fakerFirstName }} {{ fakerLastName }}"
      email:
        value: "{{ fakerEmailFR }}"
```

**Avantages** :
- Aucun breaking change
- Compatible avec toutes les configurations existantes
- Ajout transparent de nouvelles fonctions

#### Approche 2 : Type "faker" dédié
```yaml
# examples/nxs-faker-type.conf
faker:
  enabled: true
  locale: fr_FR

filters:
  public.users:
    columns:
      name:
        type: faker
        value: "{{ fakerFirstName }} {{ fakerLastName }}"
```

**Avantages** :
- Configuration centralisée
- Support explicite des locales
- Plus de contrôle sur le comportement

#### Approche 3 : Profils externes
```yaml
# examples/nxs-faker-profile.conf
faker: true
faker-file: "./examples/french-users-profile.yml"

filters:
  public.users:
    columns:
      name:
        type: faker
        value: "fullName"
```

**Avantages** :
- Réutilisabilité des profils
- Séparation des concerns
- Maintenance centralisée des patterns

## Fonctions faker ajoutées

### Fonctions de base
- `fakerName` - Nom complet
- `fakerFirstName` - Prénom
- `fakerLastName` - Nom de famille
- `fakerEmail` - Adresse email
- `fakerPhone` - Numéro de téléphone

### Fonctions localisées (fr_FR)
- `fakerPhoneFR` - Téléphone français (+33...)
- `fakerEmailFR` - Email avec domaines français

### Fonctions de dates
- `fakerDate` - Date aléatoire
- `fakerBirthDate` - Date de naissance réaliste
- `fakerDateRange "2020-01-01" "2025-01-01"` - Date dans une plage

### Fonctions numériques
- `fakerFloat 25000 75000 2` - Nombre décimal avec précision
- `fakerInt 1 100` - Entier dans une plage
- `fakerNumber 9` - Nombre avec N chiffres

### Fonctions de texte
- `fakerText 50` - Texte Lorem Ipsum
- `fakerSentence 12` - Phrase avec N mots
- `fakerParagraph` - Paragraphe complet

## Migration depuis les configurations existantes

### Depuis commandes externes
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

### Depuis templates natifs
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

## Implémentation technique

### 1. Ajout du type ValueTypeFaker

```go
// misc/values.go
const (
    ValueTypeUnknown  ValueType = "unknown"
    ValueTypeTemplate ValueType = "template"
    ValueTypeCommand  ValueType = "command"
    ValueTypeFaker    ValueType = "faker"    // NOUVEAU
)
```

### 2. Extension des templates Sprig

```go
// misc/template.go
import "github.com/brianvoe/gofakeit/v6"

func addFakerFunctions(funcMap ttemplate.FuncMap, locale string) {
    funcMap["fakerName"] = gofakeit.Name
    funcMap["fakerEmail"] = gofakeit.Email
    // ... autres fonctions
}
```

### 3. Configuration faker

```go
// ctx/conf.go
type confOpts struct {
    // ... existant
    Faker     bool       `conf:"faker"`
    FakerFile string     `conf:"faker-file"`
}
```

## Tests de validation

Le script `examples/test_faker_integration.sh` permet de :
- Tester les performances de chaque approche
- Valider la qualité des données générées
- Comparer avec les méthodes existantes

```bash
cd benchmark/examples
chmod +x test_faker_integration.sh
./test_faker_integration.sh
```

## Rétrocompatibilité

- **Approche 1** : Aucun breaking change, ajout transparent
- **Approche 2** : Optionnelle, n'affecte pas les configurations existantes
- **Approche 3** : Optionnelle, pour les nouveaux projets

## Recommandation d'implémentation

1. **Phase 1** : Implémentation de l'Approche 1 (extension Sprig)
   - Gain immédiat en performance
   - Aucun breaking change
   - Support des cas d'usage actuels

2. **Phase 2** : Ajout de l'Approche 2 (type faker)
   - Configuration plus explicite
   - Support des locales

3. **Phase 3** : Support des profils externes
   - Réutilisabilité avancée
   - Patterns d'entreprise

## Prochaines étapes

1. Validation de l'approche avec l'équipe nxs-data-anonymizer
2. Implémentation de la Phase 1 (extension Sprig)
3. Tests de performance sur données réelles
4. Documentation utilisateur
5. Déploiement progressif