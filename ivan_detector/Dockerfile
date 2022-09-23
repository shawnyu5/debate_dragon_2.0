FROM python:3.10.5-alpine3.16

WORKDIR /app
COPY . .
RUN apk update
RUN apk add make automake gcc g++ subversion python3-dev # Install some dependcies needed by numpy
RUN pip3 install -r requirements.txt -v

CMD flask run
