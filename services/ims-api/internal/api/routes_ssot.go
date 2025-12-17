package api

import (
	"github.com/edvirons/ssp/ims/internal/auth"
	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// mountSSOTRoutes registers all SSOT (Single Source of Truth) related routes.
func (s *Server) mountSSOTRoutes(r chi.Router, sync *handlers.SSOTSyncHandler, ssotList *handlers.SSOTListHandler, wh *handlers.SSOTWebhookHandler) {
	// SSOT Sync (admin/system operations)
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermSSOTSync, s.logger))
		r.Post("/ssot/sync/schools", sync.SyncSchools)
		r.Post("/ssot/sync/devices", sync.SyncDevices)
		r.Post("/ssot/sync/parts", sync.SyncParts)
		// HR SSOT sync routes
		r.Post("/ssot/sync/people", sync.SyncPeople)
		r.Post("/ssot/sync/teams", sync.SyncTeams)
		r.Post("/ssot/sync/org-units", sync.SyncOrgUnits)
		r.Post("/ssot/sync/team-memberships", sync.SyncTeamMemberships)
	})

	// SSOT Webhooks (admin/system operations)
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermSSOTWebhook, s.logger))
		r.Post("/ssot/events/schools", wh.Schools)
		r.Post("/ssot/events/devices", wh.Devices)
		r.Post("/ssot/events/parts", wh.Parts)
	})

	// SSOT List (browse snapshot data)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermSSOTRead, s.logger))
		r.Get("/ssot/schools", ssotList.ListSchools)
		r.Get("/ssot/schools/counties", ssotList.ListCounties)
		r.Get("/ssot/schools/sub-counties", ssotList.ListSubCounties)
		r.Get("/ssot/devices", ssotList.ListDevices)
		r.Get("/ssot/devices/stats", ssotList.GetDeviceStats)
		r.Get("/ssot/device-models", ssotList.ListDeviceModels)
		r.Get("/ssot/device-models/makes", ssotList.GetDeviceMakes)
		r.Get("/ssot/parts", ssotList.ListParts)
		r.Get("/ssot/status", ssotList.GetSyncStatus)

		// HR SSOT endpoints
		r.Get("/ssot/people", ssotList.ListPeople)
		r.Get("/ssot/people/{personId}", ssotList.GetPerson)
		r.Get("/ssot/teams", ssotList.ListTeams)
		r.Get("/ssot/teams/{teamId}", ssotList.GetTeam)
		r.Get("/ssot/org-units", ssotList.ListOrgUnits)
		r.Get("/ssot/org-units/tree", ssotList.GetOrgTree)
		r.Get("/ssot/org-units/{orgUnitId}", ssotList.GetOrgUnit)
		r.Get("/ssot/team-memberships", ssotList.ListTeamMemberships)
		r.Get("/ssot/team-memberships/{membershipId}", ssotList.GetTeamMembership)
	})

	// SSOT snapshot lookup helpers (debug/internal)
	s.mountSSOTLookupRoutes(r)
}
