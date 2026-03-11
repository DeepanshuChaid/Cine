package moviecontroller 

import (
  "context"
  "log"
  "time"
  "github.com/gin-gonic/gin"  
  "github.com/DeepanshuChaid/Cine/tree/main/cine/models"
)

func GetAllMovies() gin.HandlerFunc {

  return func(c *gin.Context) {

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    rows, err := database.Pool.Query(ctx,
      "SELECT id, imdbid, title, posterpath, youtubeid, adminreview FROM movies")

    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{
        "error": err.Error(),
      })
      return
    }

    defer rows.Close()

    var movies []models.Movie

    for rows.Next() {
      var movie models.Movie

      err := rows.Scan(
        &movie.ID,
        &movie.ImdbID,
        &movie.Title,
        &movie.PosterPath,
        &movie.YoutubeID,
        &movie.AdminReview,
      )

      if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
          "error": err.Error(),
        })
        return
      }

      movies = append(movies, movie)
    }

    c.JSON(http.StatusOK, movies)
  }
}
