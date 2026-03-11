,package models

type Genre struct {
  Genreid   string `json:"genreid" validate:"required"`
  Genrename string `json:"genrename" validate:"required,min=2,max=255"`
}

type Ranking struct {
  Rankingvalue int    `json:"rankingvalue" validate:"required"`
  Rankingname  string `json:"rankingname" validate:"oneof=Excellent Good Okay Bad Terrible"`
}

type Movie struct {
  ID          string  `json:"id"`
  Imdbid      string  `json:"imdbid" validate:"required"`
  Title       string  `json:"title" validate:"required,min=2,max=255"`
  Posterpath  string  `json:"posterpath" validate:"required,url"`
  Youtubeid   string  `json:"youtubeid" validate:"required"`
  Genre       []Genre `json:"genre" validate:"required"`
  Adminreview string  `json:"adminreview" validate:"required"`
  Ranking     Ranking `json:"ranking" validate:"required"`
}