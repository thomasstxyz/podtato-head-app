package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleLocalService(w http.ResponseWriter, r *http.Request) {
	service := mux.Vars(r)["partName"]
	imagePath := mux.Vars(r)["imagePath"]

	serviceDiscoverer, err := NewLocalServiceDiscoverer()
	if err != nil {
		log.Printf("failed to get service discoverer: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rootURL, err := serviceDiscoverer.GetServiceAddress(service)
	if err != nil {
		log.Printf("failed to discover address for service %s", service)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// call dependency service and get response
	response, err := http.Get(fmt.Sprintf("%s/images/%s/%s", rootURL, service, imagePath))
	if err != nil {
		log.Printf("failed to reach dependency service: %v", err)
		http.Error(w, err.Error(), http.StatusFailedDependency)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("failed to read body of dependency service response: %v", err)
		http.Error(w, err.Error(), http.StatusFailedDependency)
		return
	}
	defer response.Body.Close()

	// write dependency's response into our response
	w.Header().Set("Content-Type", "image/svg+xml")
	_, err = w.Write(body)
	if err != nil {
		log.Printf("failed to write body into our response: %v", err)
		http.Error(w, err.Error(), http.StatusFailedDependency)
		return
	}

}

func HandleExternalService(w http.ResponseWriter, r *http.Request) {
	service := mux.Vars(r)["partName"]
	imagePath := mux.Vars(r)["imagePath"]

	// discover address of dependency service
	serviceDiscoverer, err := ProvideServiceDiscoverer()
	if err != nil {
		log.Printf("failed to get service discoverer: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(service)
	rootURL, err := serviceDiscoverer.GetServiceAddress(service)
	if err != nil {
		log.Printf("failed to discover address for service %s", service)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// call dependency service and get response
	response, err := http.Get(fmt.Sprintf("%s/images/%s/%s", rootURL, service, imagePath))
	if err != nil {
		log.Printf("failed to reach dependency service: %v", err)
		http.Error(w, err.Error(), http.StatusFailedDependency)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("failed to read body of dependency service response: %v", err)
		http.Error(w, err.Error(), http.StatusFailedDependency)
		return
	}
	defer response.Body.Close()

	// write dependency's response into our response
	w.Header().Set("Content-Type", "image/svg+xml")
	_, err = w.Write(body)
	if err != nil {
		log.Printf("failed to write body into our response: %v", err)
		http.Error(w, err.Error(), http.StatusFailedDependency)
		return
	}
}
