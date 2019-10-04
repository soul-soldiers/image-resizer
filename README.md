# Soul-Soldiers image resize
-----------------------------

Based on: https://github.com/didil/gcf-go-image-resizer

Modified to run on Google Cloud Run

Builds on Google Cloud Build

```
gcloud beta run deploy --image eu.gcr.io/soul-soldiers/image-resizer --platform managed --region europe-west1
```