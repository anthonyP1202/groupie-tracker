package groupietracker

import (
	"net/http"
	"sync"
	"time"
)

// Game représente une partie du Petit Bac
type Game struct {
	ID           string
	Letter       string
	Categories   []string
	Players      map[string]Player
	Timer        *time.Timer
	TimerTimeout time.Duration
	Mutex        sync.Mutex
}

// Player représente un joueur
type Player struct {
	Name      string
	Responses map[string]string // map[category]word
	Score     int
}

// NouvellePartie crée une nouvelle partie de Petit Bac
func NouvellePartie() *Game {
	return &Game{
		ID:           "uniqueID", // Générer un ID unique
		Categories:   []string{"Artiste", "Album", "Groupe de musique", "Instrument de musique", "Featuring"},
		Players:      make(map[string]Player),
		TimerTimeout: 60 * time.Second, // Temps de réponse par défaut
	}
}

// Endpoint pour démarrer une nouvelle partie
func startGameHandler(w http.ResponseWriter, r *http.Request) {
	// Envoyer game aux clients via WebSocket
}

// Endpoint pour rejoindre une partie
func joinGameHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID de la partie et le nom du joueur à partir de la requête
	// Ajouter le joueur à la partie en cours
}

// Handler pour la soumission des réponses des joueurs
func submitResponseHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les données de réponse du joueur à partir de la requête WebSocket
	// Vérifier si la réponse est valide et unique pour chaque catégorie
	// Mettre à jour les scores des joueurs en conséquence
	// Si tous les joueurs ont soumis leurs réponses ou si le temps est écoulé, passer au tour suivant
}
