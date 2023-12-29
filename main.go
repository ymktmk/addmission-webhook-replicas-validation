package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	var (
		serverCert = flag.String("server-cert", "./server.crt", "Server certificate")
		serverKey  = flag.String("server-key", "./server.key", "Server key")
	)

	e := echo.New()
	
	e.POST("/replicas/validate", validateReplicas)

	s := http.Server{
		Addr:           ":8080",
		Handler:        e,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		TLSConfig:      &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	if err := s.ListenAndServeTLS(*serverCert, *serverKey); err != nil {
		e.Logger.Fatal(err)
	}

}

func validateReplicas(e echo.Context) error {
	var admissionReview admissionv1.AdmissionReview

	if err := e.Bind(&admissionReview); err != nil {
		log.Println(err)
		return err
	}

	var deployment appsv1.Deployment
	if err := json.Unmarshal(admissionReview.Request.Object.Raw, &deployment); err != nil {
		log.Println(err)
		return err
	}

	if deployment.Spec.Replicas != nil && *deployment.Spec.Replicas > 1 {
		admissionReview.Response.Allowed = true
		return e.JSON(http.StatusOK, admissionReview)
	} else {
		admissionReview.Response.Result = &metav1.Status{
			Code:    http.StatusForbidden,
			Message: "Replicas must be 1 or less",
		}
		admissionReview.Response.Allowed = false
		return e.JSON(http.StatusForbidden, admissionReview)
	}
}
