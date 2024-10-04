package tests

import (
	"context"
	"os"
	"testing"
	"trocup-transaction/config"
)

func TestMain(m *testing.M) {
	// Configure l'environnement pour utiliser une base de données de test
	os.Setenv("MONGODB_DBNAME", "test_db")
	os.Setenv("MONGODB_PASSWORD", "8VOTPXVSGtdukkyO")
	// Initialise la connexion à MongoDB
	config.InitMongo()

	// Exécute les tests
	code := m.Run()

	// Nettoie la base de données après les tests
	config.CleanUpTestDatabase("test_db")

	// Ferme la connexion après les tests
	if err := config.Client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}

	os.Exit(code)
}
