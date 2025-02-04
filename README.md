# snapmatch-ai

This project utilizes Google Cloud services such as Gemini AI, Vertex AI, Cloud Storage, and BigQuery.

## Prerequisites

Before running `snapmatch-ai`, you need to have the following:

*   A Google Cloud Platform (GCP) account.
*   The Google Cloud SDK (gcloud) installed and configured.
*   Appropriate permissions to access the GCP services used by the project (Cloud Storage, BigQuery, Vertex AI).
*   Docker
*   Go 1.23 or higher installed.

   ## Setup

1.  **Clone the repository:**

  ```bash
  git clone https://github.com/trapajim/snapmatch-ai.git
  cd snapmatch-ai
  ```
2.  **Set up your environment variables:**

Create a `.env` file in the root of the project with the following content:

```
GEMINI_API_KEY="" // api key for gemini
STORAGE_BUCKET="" // storage bucket for uploaded images
JOBS_STORAGE_BUCKET="" // vertex ai batch job storage
PROJECT_ID="" // gcp project id 
LOCATION="" // gcp location
DATASET_ID="" bigquery dataset
TABLE_ID="" bigquery table id
BQ_VERTEX_CONN="" // big query vertex connection
BQ_MULTI_MODAL_MODEL="" //big query multi modal model 
BQ_TEXT_MODEL="" // big query text model
```

3. Run The project

The project can be starte with docker
 ```bash
   docker compose up
 ```
