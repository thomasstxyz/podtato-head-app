package podtatoserver

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/podtato-head/podtato-head-app/pkg/assets"
	"github.com/podtato-head/podtato-head-app/pkg/handlers"
	metrics "github.com/podtato-head/podtato-head-app/pkg/metrics"
	"github.com/podtato-head/podtato-head-app/pkg/services"
	"github.com/podtato-head/podtato-head-app/pkg/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/pterm/pterm"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

var serviceMap = map[string]string{
	"leftArm":  "http://podtato-left-arm:8080",
	"rightArm": "http://podtato-right-arm:8080",
	"leftLeg":  "http://podtato-left-leg:8080",
	"rightLeg": "http://podtato-right-leg:8080",
}

const (
	assetsPrefix           = "/assets"
	externalServicesPrefix = "/parts"
)

type TemplateData struct {
	Version  string
	Hostname string
	Daytime  string
}

func serveMain(w http.ResponseWriter, r *http.Request) {
	homeTemplate, err := template.ParseFS(assets.Assets, "html/podtato-home.html")
	if err != nil {
		log.Fatalf("failed to parse file: %v", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("failed to get hostname: %v", err)
	}

	err = homeTemplate.Execute(w, TemplateData{Version: version.ServiceVersion(), Hostname: hostname, Daytime: getDayTime()})
	if err != nil {
		log.Fatalf("failed to execute template: %v", err)
	}
}

func getDayTime() string {
	hour := time.Now().Hour()
	if hour < 12 {
		return "morning"
	} else if hour < 18 {
		return "afternoon"
	} else {
		return "evening"
	}
}
func Run(component string, port string) {
	router := mux.NewRouter()

	router.Use(metrics.MetricsHandler)
	router.Path("/metrics").Handler(promhttp.Handler())

	switch component {
	case "all":
		router.Path("/").HandlerFunc(serveMain)

		// serve CSS and images
		router.PathPrefix(assetsPrefix).
			Handler(http.StripPrefix(assetsPrefix, http.FileServer(http.FS(assets.Assets))))

		router.Path(fmt.Sprintf("%s/{partName}/{imagePath}", externalServicesPrefix)).
			HandlerFunc(services.HandleLocalService)
		router.Path(fmt.Sprintf("/images/{partName}/{imageName}")).HandlerFunc(handlers.PartHandler)

		pterm.DefaultCenter.Println("Listening on port " + port + " in monolith mode")

	case "entry":
		router.PathPrefix(assetsPrefix).
			Handler(http.StripPrefix(assetsPrefix, http.FileServer(http.FS(assets.Assets))))

		router.Path("/").HandlerFunc(serveMain)
		router.Path(fmt.Sprintf("%s/{partName}/{imagePath}", externalServicesPrefix)).
			HandlerFunc(services.HandleExternalService)

	default:
		router.PathPrefix(assetsPrefix).
			Handler(http.StripPrefix(assetsPrefix, http.FileServer(http.FS(assets.Assets))))

		router.Path(fmt.Sprintf("%s/%s/{imagePath}", externalServicesPrefix, component)).
			HandlerFunc(services.HandleExternalService)

		router.Path(fmt.Sprintf("/images/%s/{imageName}", component)).HandlerFunc(handlers.PartHandler)

		pterm.DefaultCenter.Println("Listening on port " + port + " for " + component + " service")
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		log.Fatal(err)
	}
	log.Printf("exiting gracefully")

}
