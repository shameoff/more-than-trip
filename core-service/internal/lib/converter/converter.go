package converter

import (
	"fmt"
	"net/url"

	"github.com/shameoff/more-than-trip/core/internal/config"
)

// ConvertPostgresUriToDSN converts Postgres URI with sslmode=disable to DSN
// is not used now but may be helpful when work with gorm
func ConvertPostgresUriToDSN(dbUri string) (string, error) {
	parsedURL, err := url.Parse(dbUri)
	if err != nil {
		return "", err
	}

	password, is_set := parsedURL.User.Password()
	if is_set {
		password = fmt.Sprintf("password=%s", password)
	} else {
		password = ""
	}

	sslmode := parsedURL.Query().Get("sslmode")
	if sslmode != "" {
		sslmode = fmt.Sprintf("sslmode=%s", sslmode)
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s %s dbname=%s %s",
		parsedURL.Hostname(), parsedURL.Port(), parsedURL.User.Username(),
		password, parsedURL.Path[1:], sslmode)
	return dsn, nil
}

// ConvertConfigToDSN converts DatabaseConfig to DSN
func ConvertDatabaseConfigToDSN(c config.DatabaseConfig) (string, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s %s",
		c.Hostname, c.Port, c.Username, c.Password, c.Database, "sslmode=disable")

	return dsn, nil
}
