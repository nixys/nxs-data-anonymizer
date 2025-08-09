// Exemple de modification pour misc/template.go
package misc

import (
	"bytes"
	"strings"
	"time"
	ttemplate "text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/brianvoe/gofakeit/v6"
)

var (
	null = "::NULL::"
	drop = "::DROP::"
)

type TemplateData struct {
	TableName string
	Values    map[string][]byte
	Variables map[string]string
}

// addFakerFunctions ajoute les fonctions gofakeit aux templates
func addFakerFunctions(funcMap ttemplate.FuncMap, locale string) {
	// Configuration gofakeit avec seed basé sur le temps
	gofakeit.Seed(time.Now().UnixNano())
	
	// Fonctions de noms
	funcMap["fakerName"] = gofakeit.Name
	funcMap["fakerFirstName"] = gofakeit.FirstName
	funcMap["fakerLastName"] = gofakeit.LastName
	funcMap["fakerMiddleName"] = gofakeit.MiddleName
	funcMap["fakerNamePrefix"] = gofakeit.NamePrefix
	funcMap["fakerNameSuffix"] = gofakeit.NameSuffix
	
	// Fonctions de contact
	funcMap["fakerEmail"] = gofakeit.Email
	funcMap["fakerPhone"] = gofakeit.Phone
	funcMap["fakerPhoneFormatted"] = gofakeit.PhoneFormatted
	
	// Fonctions de dates
	funcMap["fakerDate"] = func() string {
		return gofakeit.Date().Format("2006-01-02")
	}
	funcMap["fakerDateRange"] = func(start, end string) string {
		startTime, _ := time.Parse("2006-01-02", start)
		endTime, _ := time.Parse("2006-01-02", end)
		return gofakeit.DateRange(startTime, endTime).Format("2006-01-02")
	}
	funcMap["fakerBirthDate"] = func() string {
		return gofakeit.DateRange(
			time.Now().AddDate(-80, 0, 0),
			time.Now().AddDate(-18, 0, 0),
		).Format("2006-01-02")
	}
	
	// Fonctions numériques
	funcMap["fakerFloat"] = func(min, max float64, precision int) float64 {
		return gofakeit.Float64Range(min, max)
	}
	funcMap["fakerInt"] = func(min, max int) int {
		return gofakeit.IntRange(min, max)
	}
	funcMap["fakerNumber"] = func(digits int) string {
		return gofakeit.Numerify(strings.Repeat("#", digits))
	}
	
	// Fonctions de texte
	funcMap["fakerText"] = func(length int) string {
		return gofakeit.LoremIpsumSentence(length)
	}
	funcMap["fakerParagraph"] = gofakeit.LoremIpsumParagraph
	funcMap["fakerSentence"] = func(wordCount int) string {
		return gofakeit.LoremIpsumSentence(wordCount)
	}
	
	// Fonctions d'adresse
	funcMap["fakerAddress"] = gofakeit.Address
	funcMap["fakerStreet"] = gofakeit.Street
	funcMap["fakerCity"] = gofakeit.City
	funcMap["fakerState"] = gofakeit.State
	funcMap["fakerZip"] = gofakeit.Zip
	funcMap["fakerCountry"] = gofakeit.Country
	
	// Fonctions spécifiques à la locale
	if locale == "fr_FR" || locale == "fr" {
		funcMap["fakerPhoneFR"] = func() string {
			// Format téléphone français
			return "+33" + gofakeit.Numerify("#########")
		}
		funcMap["fakerEmailFR"] = func() string {
			// Email avec domaines français courants
			domains := []string{"orange.fr", "free.fr", "gmail.com", "yahoo.fr", "hotmail.fr"}
			domain := domains[gofakeit.IntRange(0, len(domains)-1)]
			return gofakeit.Username() + "@" + domain
		}
	}
	
	// Fonctions avancées
	funcMap["fakerUUID"] = gofakeit.UUID
	funcMap["fakerUsername"] = gofakeit.Username
	funcMap["fakerPassword"] = func(length int) string {
		return gofakeit.Password(true, true, true, true, false, length)
	}
	
	// Fonction pour définir la locale dynamiquement
	funcMap["fakerSetLocale"] = func(newLocale string) string {
		// Note: gofakeit ne support pas toutes les locales nativement
		// mais on peut configurer certains comportements
		return ""
	}
}

// TemplateExec makes message from given template `tpl` and data `d`
func TemplateExec(tpl string, d any) ([]byte, bool, error) {
	return TemplateExecWithLocale(tpl, d, "fr_FR")
}

// TemplateExecWithLocale permet de spécifier une locale
func TemplateExecWithLocale(tpl string, d any, locale string) ([]byte, bool, error) {

	var b bytes.Buffer

	// See http://masterminds.github.io/sprig/ for details
	t, err := ttemplate.New("template").Funcs(func() ttemplate.FuncMap {

		// Get current sprig functions
		funcMap := sprig.TxtFuncMap()

		// Add existing additional functions
		funcMap["null"] = func() string {
			return null
		}
		funcMap["isNull"] = func(v string) bool {
			if v == null {
				return true
			}
			return false
		}
		funcMap["drop"] = func() string {
			return drop
		}

		// NOUVEAU : Ajouter les fonctions faker
		addFakerFunctions(funcMap, locale)

		return funcMap
	}()).Parse(tpl)
	if err != nil {
		return []byte{}, false, err
	}

	err = t.Execute(&b, d)
	if err != nil {
		return []byte{}, false, err
	}

	// Return empty line if buffer is nil
	if b.Bytes() == nil {
		return []byte{}, false, nil
	}

	// Return nil if buffer is NULL (with special key)
	if bytes.Compare(b.Bytes(), []byte(null)) == 0 {
		return nil, false, nil
	}

	// Return `drop` value if buffer is DROP (with special key)
	if bytes.Compare(b.Bytes(), []byte(drop)) == 0 {
		return nil, true, nil
	}

	// Return buffer content otherwise
	return b.Bytes(), false, nil
}