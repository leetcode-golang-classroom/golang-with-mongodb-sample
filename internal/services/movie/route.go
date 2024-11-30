package movie

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoute(router *gin.RouterGroup) {
	router.POST("/", h.CreateMovie)
	router.PUT("/:id", h.UpdateMovie)
	router.DELETE("/:id", h.DeleteMovie)
	router.DELETE("/", h.DeleteAllMovies)
	router.GET("/", h.ListAllMovies)
	router.GET("/one/:name", h.FindMovieByName)
	router.GET("/all/:name", h.FindAllMoviesByName)
	router.POST("/multiple", h.CreateMovies)
}

func (h *Handler) CreateMovie(ctx *gin.Context) {
	var movie Movie
	if err := ctx.ShouldBindJSON(&movie); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.store.CreateMovie(ctx, movie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert movie"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Movie inserted successfully!"})
}

func (h *Handler) UpdateMovie(ctx *gin.Context) {
	movieID := ctx.Param("id")
	var movie Movie
	if err := ctx.ShouldBindJSON(&movie); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.store.UpdateMovie(ctx, movieID, movie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Movie updated successfully!"})
}

func (h *Handler) DeleteMovie(ctx *gin.Context) {
	movieID := ctx.Param("id")
	err := h.store.DeleteMovie(ctx, movieID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Movie deleted successfully!"})
}

func (h *Handler) DeleteAllMovies(ctx *gin.Context) {
	err := h.store.DeleteAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete all movies"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "All movies deleted successfully"})
}

func (h *Handler) FindMovieByName(ctx *gin.Context) {
	movieName := ctx.Param("name")
	movie, err := h.store.Find(ctx, movieName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server internal error"})
		return
	}
	if movie.Movie == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}
	ctx.JSON(http.StatusOK, movie)
}

func (h *Handler) FindAllMoviesByName(ctx *gin.Context) {
	movieName := ctx.Param("name")
	movies, err := h.store.FindAll(ctx, movieName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server internal error"})
		return
	}
	if len(movies) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No movies found"})
		return
	}
	ctx.JSON(http.StatusOK, movies)
}

func (h *Handler) ListAllMovies(ctx *gin.Context) {
	movies, err := h.store.ListAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server internal error"})
		return
	}
	ctx.JSON(http.StatusOK, movies)
}

func (h *Handler) CreateMovies(ctx *gin.Context) {
	var movies []Movie
	if err := ctx.ShouldBindJSON(&movies); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.store.CreateMovies(ctx, movies)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert movies"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Movies inserted successfully!"})
}
