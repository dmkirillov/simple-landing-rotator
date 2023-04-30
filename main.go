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

	err := r.Run(":8081")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
