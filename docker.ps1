docker rm redflowerdonate
docker build -t 117503445/redflowerdonate .
docker run --name redflowerdonate --restart=always -d -p 8000:8000 117503445/redflowerdonate 
