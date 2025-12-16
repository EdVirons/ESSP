package api

import (
	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/go-chi/chi/v5"
)

// mountEdTechRoutes registers all EdTech profile routes.
func (s *Server) mountEdTechRoutes(r chi.Router, h *handlers.EdTechProfilesHandler) {
	r.Route("/edtech-profiles", func(r chi.Router) {
		// Get form options (no special permission needed)
		r.Get("/options", h.GetOptions)

		// Get profile by school ID
		r.Get("/school/{schoolId}", h.GetBySchoolID)

		// Create or update profile
		r.Post("/", h.SaveProfile)

		// Profile-specific operations
		r.Route("/{id}", func(r chi.Router) {
			// Generate AI analysis
			r.Post("/generate-ai", h.GenerateAI)

			// Submit follow-up responses
			r.Post("/submit-followup", h.SubmitFollowUp)

			// Mark as complete
			r.Post("/complete", h.Complete)
		})

		// Get version history
		r.Get("/school/{schoolId}/history", h.GetHistory)
	})
}
