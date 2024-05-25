package jobs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-scheduler/models"
	"go-scheduler/utils/scheduler"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

type ExecutionJSON struct {
	Command     string   `json:"command"`
	HostIds     []string `json:"hostIds"`
	Platform    string   `json:"platform"`
	Interpreter string   `json:"interpreter"`
	Timeout     float64  `json:"timeout"`
}

var ZeroAptHost = os.Getenv("BACKEND_INTERNAL_HOST")
var BatutaApiHost = os.Getenv("BATUTA_URL")

func GetJob(id string) (models.JobDetails, error) {
	var jsonResponse models.JobDetails
	url := fmt.Sprintf("%s/api/v1/batuta/go/job/details-api/%s", ZeroAptHost, id)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := client.Get(url)

	if err != nil {
		return jsonResponse, err
	}

	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return jsonResponse, err
	}

	err = json.Unmarshal(body, &jsonResponse)
	if err == nil {
		return jsonResponse, nil
	}

	log.Error("Body isn't a JSON")
	return jsonResponse, err
}

func ExecCommand(data models.JobDetails) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/rport/platforms/run", BatutaApiHost)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	tempJob := data.Job

	job := models.Job{
		ID:            tempJob.ID,
		Command:       tempJob.Command,
		CreationDate:  tempJob.CreationDate,
		CustomerID:    tempJob.CustomerID,
		ExecutionDate: tempJob.ExecutionDate,
		HostId:        tempJob.HostId,
		Interpreter:   tempJob.Interpreter,
		LastAttemptAt: tempJob.LastAttemptAt,
		Platform:      tempJob.Platform,
		RelatedJobId:  tempJob.RelatedJobId,
		ScenarioId:    tempJob.ScenarioId,
		Status:        tempJob.Status,
		TaskType:      tempJob.TaskType,
		Timeout:       tempJob.Timeout,
	}

	// Crear un JSON
	executionJson := ExecutionJSON{
		Command:     job.Command,
		HostIds:     []string{job.HostId},
		Platform:    job.Platform,
		Interpreter: job.Interpreter,
		Timeout:     job.Timeout,
	}

	// Serializar el JSON
	jsonBytes, err := json.Marshal(executionJson)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonBytes)))

	if err != nil {
		fmt.Println("Error al crear la solicitud POST:", err)
		return nil, err
	}

	req.Header.Set("X-Api-Token", data.CustomerToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ClientToken", "Nf45EpkYt0lbfCmb44FEZRUattmGXb4O")
	req.Header.Set("SoarId", data.SoardId)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error al enviar la solicitud POST:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Leer la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error al leer la respuesta:", err)
		return nil, err
	}

	// Intentar deserializar el cuerpo como JSON
	var jsonResponse map[string]interface{}
	err = json.Unmarshal(body, &jsonResponse)
	if err == nil {
		return jsonResponse, nil
	}

	// El cuerpo no es un JSON
	log.Info("El cuerpo no es un JSON", string(body))
	return nil, err
}

func UpdateJobDetails(data models.JobResults) bool {
	var url string

	if data.Report != "" {
		url = fmt.Sprintf("%s/api/v1/batuta/go/job/initial-receiver/%s", ZeroAptHost, data.ID)
	} else {
		url = fmt.Sprintf("%s/api/v1/batuta/go/job/report-receiver/%s", ZeroAptHost, data.ID)
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Serializar el JSON
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return false
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonBytes)))

	if err != nil {
		fmt.Println("Error al crear la solicitud PATCH:", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error al enviar la solicitud PATCH:", err)
		return false
	}

	defer resp.Body.Close()

	jsonByte, _ := io.ReadAll(resp.Body)

	fmt.Println("Respuesta:", string(jsonByte))
	// Ver si todo salio okis
	return resp.StatusCode == 200
}

func InitialJob(jobID string) {
	now := time.Now()
	execution_date := now.Format(time.RFC3339)

	fmt.Printf("✅ Task #%s executed %s\n", jobID, execution_date)

	//? hacer peticion GET DETAILS JOB - ZeroAPT
	data, err := GetJob(jobID)

	if err != nil {
		log.Errorf("Error al hacer la peticion GET DETAILS JOB: %s", err.Error())
	}

	//? hacer peticion POST EXECUTION JOB - BATUTA
	batutaResp, err := ExecCommand(data)

	if err != nil {
		log.Errorf("Error al hacer la peticion POST EXECUTION JOB: %s", err.Error())
	}

	//? Hacer peticion POST JOB RESULTS - ZeroAPT
	jsonBytes, err := json.Marshal(batutaResp)
	if err != nil {
		log.Errorf("Error al serializar el JSON (Pero TODO Ok): %s", err.Error())
	}

	jsonString := string(jsonBytes)

	if strings.Contains(jsonString, "RPORT_HOST_DISCONNECTED") {
		updateJobResult := models.JobResults{
			ID:            jobID,
			LastAttemptAt: execution_date,
			Batuta:        batutaResp,
			Report:        "Cualquier_cosa",
		}

		//? hacer peticion POST JOB RESULTS - ZeroAPT
		if ok := UpdateJobDetails(updateJobResult); ok {
			log.Infof("Job details updated successfully - [RPORT_HOST_DISCONNECTED] %s", updateJobResult.ID)
		}

	} else {

		newUUID := uuid.New()
		uuidString := newUUID.String()

		freezeTime := 4

		if data.Freeze_time != 0 {
			freezeTime = data.Freeze_time
		}

		reportTask := &models.Task{
			ID:            uuidString,
			ExecutionTime: time.Now().Add(time.Duration(freezeTime * int(time.Minute))),
			Action: func() {
				//?Funcion para Mandar el reporte
				ReportJob(uuidString)

			},
		}

		fmt.Println("Adding report task")
		_, err = scheduler.AddTask(reportTask)
		if err != nil {
			log.Errorf("Failed to add report task: %v", err)
		}

		updateJobResult := models.JobResults{
			ID:            jobID,
			LastAttemptAt: execution_date,
			Batuta:        batutaResp,
			Report:        uuidString,
		}

		//? hacer peticion POST JOB RESULTS - ZeroAPT
		if ok := UpdateJobDetails(updateJobResult); ok {
			log.Infof("Job details updated successfully %s", updateJobResult.ID)
		}

		scheduler.TaskChannel <- reportTask
	}
}

func ReportJob(jobID string) {
	now := time.Now()
	execution_date := now.Format(time.RFC3339)

	fmt.Printf("✅ Report Task #%s executed at :%s\n", jobID, time.Now().String())

	//? hacer peticion GET DETAILS JOB - ZeroAPT
	data, err := GetJob(jobID)

	if err != nil {
		log.Errorf("Error al hacer la peticion GET DETAILS JOB: %s", err.Error())
	}

	//? hacer peticion POST EXECUTION JOB - BATUTA
	batutaResp, err := ExecCommand(data)

	if err != nil {
		log.Errorf("Error al hacer la peticion POST EXECUTION JOB: %s", err.Error())
	}

	//? hacer peticion POST EXECUTION JOB - BATUTA
	jsonBytes, err := json.Marshal(batutaResp)
	if err != nil {
		log.Errorf("Error al serializar el JSON (Pero TODO Ok): %s", err.Error())
	}

	jsonString := string(jsonBytes)

	if strings.Contains(jsonString, "RPORT_HOST_DISCONNECTED") {
		//? hacer peticion POST JOB RESULTS - ZeroAPT
		updateJobResult := models.JobResults{
			ID:            jobID,
			LastAttemptAt: execution_date,
			Batuta:        batutaResp,
		}

		if ok := UpdateJobDetails(updateJobResult); ok {
			log.Infof("Job Report details updated successfully - [RPORT_HOST_DISCONNECTED] %s", updateJobResult.ID)
		}

	} else {
		updateJobResult := models.JobResults{
			ID:            jobID,
			LastAttemptAt: execution_date,
			Batuta:        batutaResp,
		}

		if ok := UpdateJobDetails(updateJobResult); ok {
			//? hacer peticion POST JOB RESULTS - ZeroAPT
			log.Infof("Job Report details updated successfully %s", updateJobResult.ID)
		}
	}
}
