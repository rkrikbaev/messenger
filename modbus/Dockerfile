# Use an official Python runtime as a parent image
FROM python:3.9

LABEL Auth: Krikbayev Rustam 
LABEL Email: "rkrikbaev@gmail.com"
ENV REFRESHED_AT 2023-02-02

# Install any needed packages
RUN python -m pip install --upgrade pip
RUN pip install pyModbusTCP

# Copy the current directory contents into the container at /app
RUN mkdir app
WORKDIR /app

COPY ./modbus_slave.py .
COPY ./config.py .

ENTRYPOINT [ "python", "modbus_slave.py" ]