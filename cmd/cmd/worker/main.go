package main
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/2spy/vinted-discord-bot/internal/scrapers/vinted"
	"github.com/2spy/vinted-discord-bot/pkg/models"
)
func main() {
	webhook := os.Getenv("DISCORD_WEBHOOK_URL")
	query := os.Getenv("SEARCH_QUERY")
	maxPriceText := os.Getenv("MAX_PRICE")
	if webhook == "" || query == "" || maxPriceText == "" {
		log.Fatal("Il manque DISCORD_WEBHOOK_URL, SEARCH_QUERY ou MAX_PRICE dans les variables d'environnement.")
	}
	maxPrice, err := strconv.ParseFloat(maxPriceText, 64)
	if err != nil {
		log.Fatal("MAX_PRICE doit être un nombre.")
	}
	scraper := vinted.NewVintedScraper()
	job := models.ScrapeJob{
		Query:   query,
		MaxPrice: int(maxPrice),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	items, err := scraper.Search(ctx, job)
	if err != nil {
		log.Fatalf("Recherche Vinted impossible : %v", err)
	}
	client := &http.Client{Timeout: 15 * time.Second}
	sent := 0
	for _, item := range items {
		if item.Price > maxPrice {
			continue
		}
		message := fmt.Sprintf("**%s**\nPrix : %.2f %s\n%s",
			item.Title, item.Price, item.Currency, item.URL)
		payload, err := json.Marshal(map[string]string{"content": message})
		if err != nil {
			log.Printf("Annonce ignorée : %v", err)
			continue
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhook, bytes.NewReader(payload))
		if err != nil {
			log.Printf("Création du message impossible : %v", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Envoi Discord impossible : %v", err)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			log.Printf("Discord a refusé le message (HTTP %d)", resp.StatusCode)
			continue
		}
		sent++
		time.Sleep(1200 * time.Millisecond)
	}
	log.Printf("Recherche terminée : %d annonces envoyées sur Discord.", sent)
}
