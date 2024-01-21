# Скрипт находит все файлы необходимые для заполнения XML документа используя дату в имени CSV файлов.
# Поиск идет по файлам в папках 'Amangeldy', 'Ayraqty', 'Zharkum'.
# Скрипт подписывает XML документ и генерирует тело SOAP сообщения через докер-контейнер.
# Результат сохраняет в файл в папке './xml_data' в следующем формате 'output_2023.07.07.xml'

import os, sys, re
import time
import subprocess
from datetime import datetime

# Пути к директориям
destination_folder = 'd:\\isun_log\\xml_data'
print(destination_folder)
source_folders = ['amangeldy', 'ayraqty', 'zharkum']
file_extension = '.csv'  # Расширение файлов для поиска и копирования

def search_files(directory, extension):
    result = []
    for root, dirs, files in os.walk(directory):
        for file in files:
            if file.endswith(extension):
                result.append(file)
                # print(file)
    return result

def exec_cmd_sing_xml(date):
    print(date)

    file = f'output_{date}.xml'
    file_name = os.path.basename(file)
    destination_path = os.path.join(destination_folder, file_name)
    print(destination_path)

    if date:
        if not os.path.exists(destination_path):
            print('destination_path: '+destination_path)
            # Выполнение команды Bash
            command = f'docker run --rm -e DATE={date} -e DIR="/app/xml_data" -v d:\isun_log:/app --name=go_isun_container go_isun'
            result = subprocess.run(command, shell=True, capture_output=True, text=True)
            #print(result)
        else:
            print(f'File already exist: {destination_path}')


def get_date_from_filename(file):
    print('get_date_from_filename: ' + file)
    # Используем регулярное выражение для поиска даты в имени файла
    match = re.search(r'\d{2}\.\d{2}\.\d{4}', file)
    print(match)
    if match:
        date_str = match.group(0)
        # Преобразуем строку даты в объект datetime
        #date = datetime.strptime(date_str, "%d.%m.%Y")
        print("Извлеченная дата:", date_str)
        return date_str
    else:
        print("Дата не найдена в имени файла.")
        print(file)
        return None

file_extension = '.csv'  # Расширение файлов для поиска

# Find files in first folder
files1 = search_files(source_folders[0], file_extension)
files2 = search_files(source_folders[1], file_extension)
files3 = search_files(source_folders[2], file_extension)

# Проверка наличия одинаковых файлов во всех трех списках
common_files = set(files1) & set(files2) & set(files3)
print(common_files)

for file in common_files:
    date = get_date_from_filename(file)
    print(date)
    exec_cmd_sing_xml(date)

