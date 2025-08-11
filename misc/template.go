package misc

import (
	"bytes"
	"fmt"
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

// TemplateExec makes message from given template `tpl` and data `d`
func TemplateExec(tpl string, d any) ([]byte, bool, error) {

	var b bytes.Buffer

	// See http://masterminds.github.io/sprig/ for details
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

		// Integration gofakeit pour données réalistes
		gofakeit.Seed(time.Now().UnixNano())
		
		// Noms
		t["fakerName"] = gofakeit.Name
		t["fakerFirstName"] = gofakeit.FirstName
		t["fakerLastName"] = gofakeit.LastName
		
		// Contact
		t["fakerEmail"] = gofakeit.Email
		t["fakerPhone"] = gofakeit.Phone
		
		// Dates
		t["fakerDate"] = func() string {
			return gofakeit.Date().Format("2006-01-02")
		}
		t["fakerBirthDate"] = func() string {
			return gofakeit.DateRange(
				time.Now().AddDate(-80, 0, 0),
				time.Now().AddDate(-18, 0, 0),
			).Format("2006-01-02")
		}
		
		// Nombres
		t["fakerFloat"] = func(min, max float64, precision int) string {
			value := gofakeit.Float64Range(min, max)
			return fmt.Sprintf("%."+fmt.Sprint(precision)+"f", value)
		}
		t["fakerInt"] = gofakeit.IntRange
		
		// Texte
		t["fakerSentence"] = func(words int) string {
			return gofakeit.LoremIpsumSentence(words)
		}
		
		// Français
		t["fakerPhoneFR"] = func() string {
			return "+33" + gofakeit.Numerify("#########")
		}
		t["fakerMobileFR"] = func() string {
			return "+336" + gofakeit.Numerify("########")
		}
		t["fakerEmailFR"] = func() string {
			domains := []string{"gmail.com", "hotmail.fr", "orange.fr", "free.fr"}
			return strings.ToLower(gofakeit.FirstName()) + "." + 
			       strings.ToLower(gofakeit.LastName()) + "@" + 
			       domains[gofakeit.IntRange(0, len(domains)-1)]
		}
		t["fakerAddress"] = func() string {
			return fmt.Sprintf("%d %s %s", gofakeit.IntRange(1, 999), gofakeit.Street(), gofakeit.StreetSuffix())
		}
		t["fakerCityFR"] = func() string {
			cities := []string{"Paris", "Lyon", "Marseille", "Toulouse", "Nice", "Nantes", "Strasbourg", "Montpellier", "Bordeaux", "Lille"}
			return cities[gofakeit.IntRange(0, len(cities)-1)]
		}
		t["fakerPostalCodeFR"] = func() string {
			return gofakeit.Numerify("#####")
		}
		t["fakerIBAN"] = func() string {
			return "FR" + gofakeit.Numerify("#########################")
		}
		t["fakerBIC"] = func() string {
			return gofakeit.LetterN(4) + "FR" + gofakeit.Numerify("##") + gofakeit.LetterN(3)
		}
		t["fakerSecuriteSocialeFR"] = func() string {
			sexe := gofakeit.IntRange(1, 2)
			annee := gofakeit.IntRange(20, 99)
			mois := gofakeit.IntRange(1, 12)
			dept := []string{"01", "75", "69", "13", "31", "44", "67", "34", "33", "59"}
			commune := gofakeit.Numerify("###")
			ordre := gofakeit.Numerify("###")
			return fmt.Sprintf("%d%02d%02d%s%s%s", sexe, annee, mois, dept[gofakeit.IntRange(0, len(dept)-1)], commune, ordre)
		}
		t["fakerVATFR"] = func() string {
			return "FR" + gofakeit.Numerify("###########")
		}
		t["fakerIPv4"] = gofakeit.IPv4Address
		t["fakerIPv6"] = gofakeit.IPv6Address
		
		// Cartes de crédit
		t["fakerCreditCard"] = gofakeit.CreditCardNumber
		t["fakerCreditCardType"] = gofakeit.CreditCardType
		t["fakerCreditCardCVV"] = gofakeit.CreditCardCvv
		t["fakerCreditCardExpiry"] = func() string {
			return gofakeit.CreditCardExp()
		}

		return t
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
