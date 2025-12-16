package api

import (
	"github.com/edvirons/ssp/ims/internal/auth"
	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// mountWorkOrderRoutes registers all work order-related routes.
func (s *Server) mountWorkOrderRoutes(r chi.Router, wo *handlers.WorkOrderHandler, woops *handlers.WorkOrderOpsHandler, woUpdate *handlers.WorkOrderUpdateHandler, woRework *handlers.WorkOrderReworkHandler, woBulk *handlers.WorkOrderBulkHandler) {
	// Work Orders - read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermWorkOrderRead, s.logger))
		r.Get("/work-orders/{id}", wo.GetByID)
		r.Get("/work-orders", wo.List)
	})

	// Work Orders - create operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermWorkOrderCreate, s.logger))
		r.Post("/work-orders", wo.Create)
	})

	// Work Orders - update operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermWorkOrderUpdate, s.logger))
		r.Patch("/work-orders/{id}/status", wo.UpdateStatus)
		r.Patch("/work-orders/{id}", woUpdate.Update)
	})

	// Work Order Rework/Rejection operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermWorkOrderReview, s.logger))
		r.Post("/work-orders/{id}/reject", woRework.Reject)
		r.Get("/work-orders/{id}/rework-history", woRework.GetReworkHistory)
	})

	// Work Order Bulk operations (stricter rate limiting)
	r.Group(func(r chi.Router) {
		r.Use(s.bulkRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermWorkOrderUpdate, s.logger))
		r.Post("/work-orders/bulk/status", woBulk.BulkStatusUpdate)
		r.Post("/work-orders/bulk/assignment", woBulk.BulkAssignment)
	})

	// Work Order Bulk approval operations
	r.Group(func(r chi.Router) {
		r.Use(s.bulkRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermWorkOrderApproval, s.logger))
		r.Post("/work-orders/bulk/approval", woBulk.BulkApproval)
	})

	// Work Order Scheduling
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermWorkOrderSchedule, s.logger))
		r.Post("/work-orders/{id}/schedule", woops.Schedule)
		r.Get("/work-orders/{id}/schedules", woops.Schedules)
	})

	// Work Order Deliverables
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermWorkOrderDeliverable, s.logger))
		r.Post("/work-orders/{id}/deliverables", woops.AddDeliverable)
		r.Get("/work-orders/{id}/deliverables", woops.Deliverables)
		r.Patch("/work-orders/{id}/deliverables/{deliverableId}/submit", woops.SubmitDeliverable)
	})

	// Work Order Deliverable Review (requires review permission)
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermWorkOrderReview, s.logger))
		r.Patch("/work-orders/{id}/deliverables/{deliverableId}/review", woops.ReviewDeliverable)
	})

	// Work Order Approvals
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermWorkOrderApproval, s.logger))
		r.Post("/work-orders/{id}/approvals", woops.RequestApproval)
		r.Patch("/work-orders/{id}/approvals/{approvalId}/decide", woops.DecideApproval)
	})
}
