# prunusavium-api

https://en.wikipedia.org/wiki/Prunus_avium

+ deploy cloud functions

```
gcloud functions deploy v1 --entry-point V1 --runtime go111 --memory 128m \
       --trigger-http \
       --region asia-northeast1
```
