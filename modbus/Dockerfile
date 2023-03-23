# Use an official Python runtime as a parent image
FROM python:3.9

LABEL Auth: Krikbayev Rustam 
LABEL Email: "rkrikbaev@gmail.com"
ENV REFRESHED_AT 2023-02-02

# Install any needed packages
RUN python -m pip install --upgrade pip
RUN pip install pyModbusTCP
RUN pip install selenium

EXPOSE 5020

# Copy the current directory contents into the container at /app
RUN mkdir app
WORKDIR /app

COPY ./data_graber.py
COPY ./modbus_slave.py .
COPY ./chromedriver .

# ENTRYPOINT [ "python", "modbus_slave.py" ]