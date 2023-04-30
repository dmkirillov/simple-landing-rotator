package main

import (
	"crypto/rand"
	"math/big"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func randomChoiceWithWeights(links []string, probabilities []float64) string {
	totalWeight := 0.0
	for _, weight := range probabilities {
		totalWeight += weight
	}

	r, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to generate random number")
	}

	weightSum := 0.0
	for i, weight := range probabilities {
		weightSum += weight / totalWeight
		if float64(r.Int64()) < weightSum*1000000 {
			return links[i]
		}
	}
	return ""
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	links := []string{"http://yandex.ru/", "https://www.google.com/", "https://www.youtube.com/"}
	probabilities := []float64{0.2, 0.3, 0.5}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	r := gin.Default()

	r.GET("/redirect", func(c *gin.Context) {
		redirectUrl := randomChoiceWithWeights(links, probabilities)
		if redirectUrl == "" {
			c.String(500, "Failed to redirect")
			return
		}

		logger.Info().Str("redirect_url", redirectUrl).Msg("Redirecting")
		c.Redirect(302, redirectUrl)
	})

	r.POST("/links", func(c *gin.Context) {
		var link struct {
			URL         string  `json:"url"`
			Probability float64 `json:"probability"`
		}
		if err := c.BindJSON(&link); err != nil {
			c.String(400, "Invalid request body")
			return
		}

		links = append(links, link.URL)
		probabilities = append(probabilities, link.Probability)
		logger.Info().Str("url", link.URL).Msg("Link added")
		c.Status(204)
	})

	r.DELETE("/links/:url", func(c *gin.Context) {
		urlToDelete := c.Param("url")

		index := -1
		for i, url := range links {
			if url == urlToDelete {
				index = i
				break
			}
		}

		if index == -1 {
			c.String(404, "Link not found")
			return
		}

		links = append(links[:index], links[index+1:]...)
		probabilities = append(probabilities[:index], probabilities[index+1:]...)
		logger.Info().Str("url", urlToDelete).Msg("Link removed")
		c.Status(204)
	})

	err := r.Run(":8081")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
