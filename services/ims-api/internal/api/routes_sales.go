package api

import (
	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/go-chi/chi/v5"
)

// mountSalesRoutes registers all sales-related routes.
func (s *Server) mountSalesRoutes(r chi.Router, demo *handlers.DemoPipelineHandler, pres *handlers.PresentationsHandler, metrics *handlers.SalesMetricsHandler) {
	// Demo Pipeline routes
	r.Route("/demo-pipeline", func(r chi.Router) {
		// Lead management
		r.Get("/leads", demo.ListLeads)
		r.Post("/leads", demo.CreateLead)
		r.Get("/leads/{id}", demo.GetLead)
		r.Put("/leads/{id}", demo.UpdateLead)
		r.Delete("/leads/{id}", demo.DeleteLead)

		// Stage management
		r.Put("/leads/{id}/stage", demo.UpdateLeadStage)

		// Activities
		r.Get("/leads/{id}/activities", demo.ListActivities)
		r.Post("/leads/{id}/notes", demo.AddNote)

		// Demo scheduling
		r.Post("/leads/{id}/schedule-demo", demo.ScheduleDemo)

		// Summary
		r.Get("/summary", demo.GetPipelineSummary)
		r.Get("/activities", demo.GetRecentActivities)
	})

	// Presentations routes
	r.Route("/presentations", func(r chi.Router) {
		r.Get("/", pres.List)
		r.Post("/", pres.Create)
		r.Get("/types", pres.GetTypes)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", pres.GetByID)
			r.Put("/", pres.Update)
			r.Delete("/", pres.Delete)
			r.Get("/download", pres.DownloadURL)
			r.Post("/view", pres.RecordView)
		})
	})

	// Sales Metrics routes
	r.Route("/sales", func(r chi.Router) {
		r.Get("/dashboard", metrics.GetDashboard)
		r.Get("/metrics", metrics.GetMetricsSummary)
		r.Get("/pipeline-stages", metrics.GetPipelineStages)
		r.Post("/metrics/increment", metrics.IncrementMetric)
	})
}
