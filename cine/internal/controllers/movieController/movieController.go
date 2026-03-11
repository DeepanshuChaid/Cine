package moviecontroller

import (
	"context"
	// "log"
	"time"
  "github.com/jackc/pgx/v5"

	"github.com/DeepanshuChaid/Cine/tree/main/cine/internal/database"
	"github.com/DeepanshuChaid/Cine/tree/main/cine/internal/models"
	"github.com/gin-gonic/gin"
)

//GET ALL MOVIES
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
        &movie.Adminreview,
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


// GET SINGLE MOVIE  
func GetMovie() gin.HandlerFunc {
  return func(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
    defer cancel()

    movieId := c.Param("id")

    if movieId == "" {
      c.JSON(400, gin.H{
        "error": "Movie ID is required",
      })
      return
    }

    var movie models.Movie

    err := database.Pool.QueryRow(ctx, "SELECT id, imdbid, title, posterpath, youtubeid, adminreview FROM movies WHERE id = $1", movieId).Scan(
        &movie.ID,
        &movie.Imdbid,
        &movie.Title,
        &movie.Posterpath,
        &movie.Youtubeid,
        &movie.Adminreview,
      )
  
    if err != nil {
        if err == pgx.ErrNoRows {
            c.JSON(404, gin.H{
                "error": "Movie not found",
            })
            return
        }

        c.JSON(500, gin.H{
            "error": "Failed to fetch movie",
            "details": err.Error(),
        })
        return
    }


    c.JSON(200, gin.H{
      "message": "Successfully fetched Movie",
      "movie": movie,
    })
  }
}