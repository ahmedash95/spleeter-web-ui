FROM golang:1.11

RUN apt update && apt install python3 python3-pip ffmpeg -y && pip3 install spleeter 