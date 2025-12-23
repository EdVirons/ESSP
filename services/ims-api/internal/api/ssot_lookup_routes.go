package api

import (
	"net/http"

	"github.com/edvirons/ssp/ims/internal/lookups"
	"github.com/edvirons/ssp/shared/pkg/httpx"
	"github.com/go-chi/chi/v5"
)

func (s *Server) mountSSOTLookupRoutes(r chi.Router) {
	r.Route("/ssot", func(r chi.Router) {
		r.Get("/school/{id}", s.handleSchoolLookup)
		r.Get("/school/{id}/primary-contact", s.handlePrimaryContact)
		r.Get("/device/{id}", s.handleDeviceLookup)
		r.Get("/device/serial/{serial}", s.handleDeviceBySerialLookup)
		r.Get("/part/{id}", s.handlePartLookup)
		r.Get("/part/puk/{puk}", s.handlePartByPUKLookup)
	})
}

func (s *Server) lookupStore() *lookups.Store {
	return lookups.New(s.pg.RawPool())
}

func (s *Server) handleSchoolLookup(w http.ResponseWriter, r *http.Request) {
	tenant := httpx.TenantID(r)
	if tenant == "" {
		httpx.Error(w, 400, "X-Tenant-Id required")
		return
	}
	id := chi.URLParam(r, "id")
	res, err := s.lookupStore().SchoolByID(r.Context(), tenant, id)
	if err != nil {
		httpx.Error(w, statusFromLookupErr(err), err.Error())
		return
	}
	httpx.WriteJSON(w, 200, res)
}

func (s *Server) handlePrimaryContact(w http.ResponseWriter, r *http.Request) {
	tenant := httpx.TenantID(r)
	if tenant == "" {
		httpx.Error(w, 400, "X-Tenant-Id required")
		return
	}
	id := chi.URLParam(r, "id")
	res, err := s.lookupStore().PrimaryContactBySchoolID(r.Context(), tenant, id)
	if err != nil {
		httpx.Error(w, statusFromLookupErr(err), err.Error())
		return
	}
	httpx.WriteJSON(w, 200, res)
}

func (s *Server) handleDeviceLookup(w http.ResponseWriter, r *http.Request) {
	tenant := httpx.TenantID(r)
	if tenant == "" {
		httpx.Error(w, 400, "X-Tenant-Id required")
		return
	}
	id := chi.URLParam(r, "id")
	res, err := s.lookupStore().DeviceByID(r.Context(), tenant, id)
	if err != nil {
		httpx.Error(w, statusFromLookupErr(err), err.Error())
		return
	}
	httpx.WriteJSON(w, 200, res)
}

func (s *Server) handleDeviceBySerialLookup(w http.ResponseWriter, r *http.Request) {
	tenant := httpx.TenantID(r)
	if tenant == "" {
		httpx.Error(w, 400, "X-Tenant-Id required")
		return
	}
	serial := chi.URLParam(r, "serial")
	res, err := s.lookupStore().DeviceBySerial(r.Context(), tenant, serial)
	if err != nil {
		httpx.Error(w, statusFromLookupErr(err), err.Error())
		return
	}
	httpx.WriteJSON(w, 200, res)
}

func (s *Server) handlePartLookup(w http.ResponseWriter, r *http.Request) {
	tenant := httpx.TenantID(r)
	if tenant == "" {
		httpx.Error(w, 400, "X-Tenant-Id required")
		return
	}
	id := chi.URLParam(r, "id")
	res, err := s.lookupStore().PartByID(r.Context(), tenant, id)
	if err != nil {
		httpx.Error(w, statusFromLookupErr(err), err.Error())
		return
	}
	httpx.WriteJSON(w, 200, res)
}

func (s *Server) handlePartByPUKLookup(w http.ResponseWriter, r *http.Request) {
	tenant := httpx.TenantID(r)
	if tenant == "" {
		httpx.Error(w, 400, "X-Tenant-Id required")
		return
	}
	puk := chi.URLParam(r, "puk")
	res, err := s.lookupStore().PartByPUK(r.Context(), tenant, puk)
	if err != nil {
		httpx.Error(w, statusFromLookupErr(err), err.Error())
		return
	}
	httpx.WriteJSON(w, 200, res)
}

func statusFromLookupErr(err error) int {
	if err == lookups.ErrNotFound {
		return 404
	}
	if err == lookups.ErrSnapshotMissing {
		return 424
	} // Failed Dependency
	return 500
}
