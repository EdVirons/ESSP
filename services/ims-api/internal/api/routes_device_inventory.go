package api

import (
	"github.com/edvirons/ssp/ims/internal/auth"
	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// mountDeviceInventoryRoutes registers device inventory routes for school contacts and admins.
func (s *Server) mountDeviceInventoryRoutes(r chi.Router, inv *handlers.DeviceInventoryHandler) {
	// School Inventory - read operations (for school contacts)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermDeviceInventory, s.logger))
		r.Get("/schools/{schoolId}/inventory", inv.GetSchoolInventory)
	})

	// Locations - read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermLocationRead, s.logger))
		r.Get("/schools/{schoolId}/locations", inv.ListLocations)
	})

	// Locations - write operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermLocationWrite, s.logger))
		r.Post("/schools/{schoolId}/locations", inv.CreateLocation)
		r.Put("/schools/{schoolId}/locations/{id}", inv.UpdateLocation)
		r.Delete("/schools/{schoolId}/locations/{id}", inv.DeleteLocation)
	})

	// Device Assignments - read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermAssignmentRead, s.logger))
		r.Get("/devices/{deviceId}/assignments", inv.GetDeviceAssignments)
		r.Get("/locations/{locationId}/devices", inv.GetLocationDevices)
	})

	// Device Assignments - write operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermAssignmentWrite, s.logger))
		r.Post("/devices/{deviceId}/assign", inv.AssignDevice)
		r.Delete("/devices/{deviceId}/assign", inv.UnassignDevice)
	})

	// Device Groups - read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermGroupRead, s.logger))
		r.Get("/schools/{schoolId}/groups", inv.ListGroups)
		r.Get("/groups/{id}", inv.GetGroup)
		r.Get("/groups/{id}/devices", inv.GetGroupDevices)
	})

	// Device Groups - write operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermGroupWrite, s.logger))
		r.Post("/groups", inv.CreateGroup)
		r.Post("/groups/{id}/members", inv.AddGroupMembers)
		r.Delete("/groups/{id}/members", inv.RemoveGroupMembers)
	})

	// Device Registration - for school contacts to add new devices
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermDeviceCreate, s.logger))
		r.Post("/schools/{schoolId}/devices", inv.RegisterDevice)
	})
}
