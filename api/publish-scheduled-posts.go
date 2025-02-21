package handler

import (
	"net/http"
	"pixi/api/controllers"
)

func PublishScheduledPostsHandler(w http.ResponseWriter, r *http.Request) {
	// Call the function to publish scheduled posts
	controllers.PublishScheduledPosts()

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Scheduled posts published"))
}
