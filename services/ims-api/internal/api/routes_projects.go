package api

import (
	"github.com/edvirons/ssp/ims/internal/auth"
	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// mountProjectRoutes registers all project-related routes including phases, surveys, BOQ, team, activities, and work orders.
func (s *Server) mountProjectRoutes(r chi.Router, proj *handlers.ProjectsHandler, ph *handlers.PhasesHandler, surv *handlers.SurveysHandler, boq *handlers.BOQHandler, team *handlers.ProjectTeamHandler, activities *handlers.ProjectActivitiesHandler, projectWO *handlers.ProjectWorkOrdersHandler) {
	// Projects - read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermProjectRead, s.logger))
		r.Get("/projects/types", proj.GetProjectTypes)
		r.Get("/projects/counts", proj.GetProjectTypeCounts)
		r.Get("/projects/{id}", proj.GetByID)
		r.Get("/projects", proj.List)
	})

	// Projects - create operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermProjectCreate, s.logger))
		r.Post("/projects", proj.Create)
	})

	// Phases - read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermPhaseRead, s.logger))
		r.Get("/projects/{id}/phases", ph.List)
	})

	// Phases - create operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermPhaseCreate, s.logger))
		r.Post("/projects/{id}/phases", ph.Create)
	})

	// Phases - update operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermPhaseUpdate, s.logger))
		r.Patch("/phases/{phaseId}/status", ph.UpdateStatus)
	})

	// Surveys - read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermSurveyRead, s.logger))
		r.Get("/projects/{id}/surveys", surv.List)
		r.Get("/surveys/{surveyId}", surv.GetByID)
	})

	// Surveys - create operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermSurveyCreate, s.logger))
		r.Post("/projects/{id}/surveys", surv.Create)
		r.Post("/surveys/{surveyId}/rooms", surv.AddRoom)
		r.Post("/surveys/{surveyId}/photos", surv.AddPhoto)
	})

	// BOQ - read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermBOQRead, s.logger))
		r.Get("/projects/{id}/boq", boq.List)
	})

	// BOQ - create operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermBOQCreate, s.logger))
		r.Post("/projects/{id}/boq/items", boq.AddItem)
	})

	// Project Team - read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermProjectTeamRead, s.logger))
		r.Get("/projects/{id}/team", team.ListMembers)
		r.Get("/users/me/projects", team.ListMyProjects)
		r.Get("/phases/{phaseId}/assignments", team.ListPhaseAssignments)
	})

	// Project Team - update operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermProjectTeamUpdate, s.logger))
		r.Post("/projects/{id}/team", team.AddMember)
		r.Patch("/projects/{id}/team/{memberId}", team.UpdateMember)
		r.Delete("/projects/{id}/team/{memberId}", team.RemoveMember)
		r.Post("/phases/{phaseId}/assignments", team.AddPhaseAssignment)
		r.Delete("/phases/{phaseId}/assignments/{assignmentId}", team.RemovePhaseAssignment)
	})

	// Project Activities - read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermActivityRead, s.logger))
		r.Get("/projects/{id}/activities", activities.ListActivities)
		r.Get("/projects/{id}/attachments", activities.ListAttachments)
	})

	// Project Activities - create operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermActivityCreate, s.logger))
		r.Post("/projects/{id}/activities", activities.CreateActivity)
	})

	// Project Activities - update operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermActivityUpdate, s.logger))
		r.Patch("/activities/{activityId}", activities.UpdateActivity)
		r.Post("/activities/{activityId}/pin", activities.TogglePin)
	})

	// Project Activities - delete operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermActivityDelete, s.logger))
		r.Delete("/activities/{activityId}", activities.DeleteActivity)
		r.Delete("/attachments/{attachmentId}", activities.DeleteAttachment)
	})

	// Project Work Orders - read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermWorkOrderRead, s.logger))
		r.Get("/projects/{id}/work-orders", projectWO.ListByProject)
	})

	// Project Work Orders - create operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermWorkOrderCreate, s.logger))
		r.Post("/projects/{id}/work-orders", projectWO.CreateFromProject)
	})
}
