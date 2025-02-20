// utils/scheduler.go
package scheduler

import (
	"context"
	"log"
	"pixi/api/controllers" // Adjusted to match the module name
	"time"
)

func StartScheduler(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			controllers.PublishScheduledPosts()
		case <-ctx.Done():
			log.Println("Scheduler stopped")
			return
		}
	}
}
