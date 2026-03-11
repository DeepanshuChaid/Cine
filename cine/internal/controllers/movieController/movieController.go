package moviecontroller

import (
	"context"
	// "log"
	"time"

	"github.com/DeepanshuChaid/Cine/tree/main/cine/internal/database"
	"github.com/DeepanshuChaid/Cine/tree/main/cine/internal/models"
	"github.com/gin-gonic/gin"
)

func GetAllMovies() gin.HandlerFunc {
  return func(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
    defer cancel()

    var movies []models.Movie

    rows, err := database.Pool.Query(ctx,
      "SELECT id, imdbid, title, posterpath, youtubeid, adminreview FROM movies")

    if err != nil {
      c.JSON(500, gin.H{"error": "Failed to fetch movies", "details": err.Error()})
      return
    }
    defer rows.Close()

    for rows.Next() {
      var movie models.Movie
    
      err = rows.Scan(
        &movie.ID,
        &movie.Imdbid,
        &movie.Title,
        &movie.Posterpath,
        &movie.Youtubeid,
        &movie.Adminreview
      )
      if err != nil {
        c.JSON(500, gin.H{"error": "Failed to scan movie row", "details": err.Error()})
        return
      }

      movies = append(movies, movie)
    }

    if err := rows.Err(); err != nil {
      c.JSON(500, gin.H{
        "error": "Error iterating rows",
        "details": err.Error(),
      })
      return
    }

    c.JSON(200, gin.H{
      "movies": movies,
    })
    
  }
}
