package moviecontroller

import (
	"context"
	// "log"
	"time"

  "github.com/go-playground/validator/v10"
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


// CREATE MOVIE
func CreateMovie() gin.HandlerFunc {
  return func(c *gin.Context) {

    ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
    defer cancel()

    var movie models.Movie

    if err := c.ShouldBindJSON(&movie); err != nil {
      c.JSON(400, gin.H{
        "error":   "Invalid request body",
        "details": err.Error(),
      })
      return
    }

    validate := validator.New(validator.WithRequiredStructEnabled())

    if err := validate.Struct(movie); err != nil {
      c.JSON(400, gin.H{
        "error":   "Validation failed",
        "details": err.Error(),
      })
      return
    }

    tx, err := database.Pool.Begin(ctx)
    if err != nil {
      c.JSON(500, gin.H{
        "error":   "Failed to begin transaction",
        "details": err.Error(),
      })
      return
    }
    defer tx.Rollback(ctx)

    // Insert movie
    err = tx.QueryRow(ctx,
      `INSERT INTO movies (imdbid, title, posterpath, youtubeid, adminreview)
       VALUES ($1,$2,$3,$4,$5)
       RETURNING id`,
      movie.Imdbid,
      movie.Title,
      movie.Posterpath,
      movie.Youtubeid,
      movie.Adminreview,
    ).Scan(&movie.ID)

    if err != nil {
      c.JSON(500, gin.H{
        "error":   "Failed to insert movie",
        "details": err.Error(),
      })
      return
    }

    // Insert ranking
    _, err = tx.Exec(ctx,
      `INSERT INTO rankings (movie_id, rankingvalue, rankingname)
       VALUES ($1,$2,$3)`,
      movie.ID,
      movie.Ranking.Rankingvalue,
      movie.Ranking.Rankingname,
    )

    if err != nil {
      c.JSON(500, gin.H{
        "error":   "Failed to insert ranking",
        "details": err.Error(),
      })
      return
    }

    // Insert genres and movie_genres relation
    for _, g := range movie.Genre {

      _, err = tx.Exec(ctx,
        `INSERT INTO genres (genreid, genrename)
         VALUES ($1,$2)
         ON CONFLICT (genreid) DO NOTHING`,
        g.Genreid,
        g.Genrename,
      )

      if err != nil {
        c.JSON(500, gin.H{
          "error":   "Failed to insert genre",
          "details": err.Error(),
        })
        return
      }

      _, err = tx.Exec(ctx,
        `INSERT INTO movie_genres (movie_id, genre_id)
         VALUES ($1,$2)`,
        movie.ID,
        g.Genreid,
      )

      if err != nil {
        c.JSON(500, gin.H{
          "error":   "Failed to insert movie genre",
          "details": err.Error(),
        })
        return
      }
    }

    // Commit transaction
    if err := tx.Commit(ctx); err != nil {
      c.JSON(500, gin.H{"error": "Transaction failed"})
      return
    }

    c.JSON(201, gin.H{
      "message": "Movie created successfully",
      "movie":   movie,
    })
  }
}
