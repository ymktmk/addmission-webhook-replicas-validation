package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	e := echo.New()
	
	e.POST("/replicas/validate", validateReplicas)

	s := http.Server{
		Addr:           ":8080",
		Handler:        e,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
	}

	if err := s.ListenAndServe(); err != nil {
		e.Logger.Fatal(err)
	}

}

func validateReplicas(e echo.Context) error {
	var admissionRequest admissionv1.AdmissionReview

	if err := e.Bind(&admissionRequest); err != nil {
		log.Println(err)
		return err
	}

	var deployment appsv1.Deployment
	if err := json.Unmarshal(admissionRequest.Request.Object.Raw, &deployment); err != nil {
		log.Println(err)
		return err
	}

	if deployment.Spec.Replicas != nil && *deployment.Spec.Replicas > 1 {
		admissionResponse := admissionv1.AdmissionReview{
			TypeMeta: admissionRequest.TypeMeta,
			Response: &admissionv1.AdmissionResponse{
				UID:     admissionRequest.Request.UID,
				Allowed: false,
				Result: &metav1.Status{
					Message: "Replicas must be 1 or less",
				},
			},
		}
		return e.JSON(http.StatusForbidden, admissionResponse)
	}


	admissionResponse := admissionv1.AdmissionReview{
		TypeMeta: admissionRequest.TypeMeta,
		Response: &admissionv1.AdmissionResponse{
			UID:     admissionRequest.Request.UID,
			Allowed: true,
		},
	}

	return e.JSON(http.StatusOK, admissionResponse)
}
