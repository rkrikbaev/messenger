# Use an official Python runtime as a parent image
FROM python:3.9

LABEL Auth: Krikbayev Rustam 
LABEL Email: "rkrikbaev@gmail.com"
ENV REFRESHED_AT 2023-02-02

# Install any needed packages
RUN python -m pip install --upgrade pip
RUN pip install selenium

EXPOSE 80

# Copy the current directory contents into the container at /app
RUN mkdir app
WORKDIR /app

COPY chromedriver .
COPY graber.py .

ENTRYPOINT [ "python", "graber.py" ]