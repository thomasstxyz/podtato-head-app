package podtatoserver

import (
	"encoding/json"
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
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	assetsPrefix = "/assets"
)

type PodTatoServer struct {
	Component string
	Port      string
}

type FrontEndComponents struct {
}

type TemplateData struct {
	Version         string
	Hostname        string
	Daytime         string
	Components      FrontEndComponents
	LeftArm         string
	LeftArmVersion  string
	RightArm        string
	RightArmVersion string
	LeftLeg         string
	LeftLegVersion  string
	RightLeg        string
	RightLegVersion string
	Hat             string
	HatVersion      string
}

func (p PodTatoServer) frontendHandler(w http.ResponseWriter, r *http.Request) {

	homeTemplate, err := template.ParseFS(assets.Assets, "html/podtato-home.html")
	if err != nil {
		log.Fatalf("failed to parse file: %v", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("failed to get hostname: %v", err)
	}

	tpl := TemplateData{
		LeftArm:  p.fetchImage("left-arm"),
		RightArm: p.fetchImage("right-arm"),
		LeftLeg:  p.fetchImage("left-leg"),
		RightLeg: p.fetchImage("right-leg"),
		Hat:      p.fetchImage("hat"),
		Hostname: hostname,
		Daytime:  getDayTime(),
		Version:  version.ServiceVersion(),
	}

	err = homeTemplate.Execute(w, tpl)
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

func (p PodTatoServer) Serve() error {
	// Add metrics
	router := mux.NewRouter()
	router.Use(metrics.MetricsHandler)
	router.Path("/metrics").Handler(promhttp.Handler())

	switch p.Component {
	case "all":
		router.Path("/").HandlerFunc(p.frontendHandler)

		// serve CSS and images
		router.PathPrefix(assetsPrefix).
			Handler(http.StripPrefix(assetsPrefix, http.FileServer(http.FS(assets.Assets))))

		router.Path("/images/{partName}/{partName}").HandlerFunc(handlers.PartHandler)

		pterm.DefaultCenter.Println("Listening on port " + p.Port + " in monolith mode")

	case "frontend":
		router.Path("/").HandlerFunc(p.frontendHandler)

		// serve CSS and images
		router.PathPrefix(assetsPrefix).
			Handler(http.StripPrefix(assetsPrefix, http.FileServer(http.FS(assets.Assets))))

		pterm.DefaultCenter.Println("Listening on port " + p.Port + " in frontend mode")

	default:
		router.PathPrefix(assetsPrefix).
			Handler(http.StripPrefix(assetsPrefix, http.FileServer(http.FS(assets.Assets))))

		fmt.Println(p.Component)
		router.Path(fmt.Sprintf("/images/%s/{partName}", p.Component)).HandlerFunc(handlers.PartHandler)

		pterm.DefaultCenter.Println("Listening on port " + p.Port + " for " + p.Component + " service")
	}

	// Start server
	if err := http.ListenAndServe(fmt.Sprintf(":%s", p.Port), router); err != nil {
		return err
	}
	return nil
}

func (p PodTatoServer) fetchImage(component string) string {
	var serviceDiscoverer services.ServiceMap
	var err error
	if p.Component == "all" {
		serviceDiscoverer, err = services.NewLocalServiceDiscoverer(p.Port)
		if err != nil {
			log.Printf("failed to get service discoverer: %v", err)
			return ""

		}
	} else {
		serviceDiscoverer, err = services.ProvideServiceDiscoverer()
		if err != nil {
			log.Printf("failed to get service discoverer: %v", err)
			return ""
		}
	}
	fmt.Println(serviceDiscoverer)
	rootURL, err := serviceDiscoverer.GetServiceAddress(component)
	if err != nil {
		log.Printf("failed to discover address for service %s", component)
		return ""
	}

	response, err := http.Get(fmt.Sprintf("%s/images/%s/%s", rootURL, component, component))
	if err != nil {
		log.Printf("failed to reach dependency service: %v", err)
		return ""
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("failed to read body of dependency service response: %v", err)
		return ""
	}
	defer response.Body.Close()

	part := handlers.PartResponse{}
	err = json.Unmarshal(body, &part)
	if err != nil {
		log.Printf("failed to unmarshal body of dependency service response: %v", err)
		return ""
	}
	return part.Image
}
