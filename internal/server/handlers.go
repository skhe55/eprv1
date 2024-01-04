package server

import (
	httpConverter "erpv1/internal/converter/delivery/http"
	repoConverter "erpv1/internal/converter/repository"
	usecaseConverter "erpv1/internal/converter/usecase"
	"erpv1/pkg/inmemdb"
	"net/http"
)

func (s *Server) MapHandlers(mux *http.ServeMux) {
	inmemDbInstance := inmemdb.NewInmemDB()

	converterRepo := repoConverter.NewConverterInmemRepository(inmemDbInstance)
	converterUsecase := usecaseConverter.NewConverterUseCase(converterRepo)

	converterHandlers := httpConverter.NewConverterHandlers(converterUsecase)

	httpConverter.MapConverterRoutes("api/v1/converter", mux, converterHandlers)
}
